package ui

import (
	"image"
	"log"
	"math"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gen2brain/go-fitz"
)

// 接口
type ViewInterface interface {
	Render(filePath string, page int) (*image.RGBA, error)
}

// 页元素信息
type PageInfo struct {
	Num  int
	Size image.Rectangle
	// Scale  float64
	Widget *widget.Image
}

type ImgInfo struct {
	Num   int
	imgop *paint.ImageOp
}

// pdf渲染器 ---------------------------------------------------------
type PDFViewer struct {
	motion *Motion
	style  *Style
	mu     sync.Mutex

	// 页元素与缓冲区
	max_buffer       int
	page_buffer      map[int]*PageInfo
	thumbnail_buffer map[int]*ImgInfo
	pad              int

	// 缓冲区
	size_buffer map[int]*image.Rectangle

	// 位移与缩放
	offsetX_tar float64
	offsetX_cur float64
	offsetY_tar float64
	offsetY_cur float64
	scale_tar   float64
	scale_cur   float64

	focus_index int

	// 按钮组件
	reset_btn widget.Clickable

	// 文档信息
	doc_id    string
	doc_scale float64 // 文档控件尺寸/文件实际尺寸

	// 渲染组件
	viewer_area     image.Point
	viewer_children []layout.FlexChild
	base_height     float64
	viewer_radius   int

	// 逐步加载
	height_list []float64

	// 接口
	bridge_trunk TrunkBridge
	bridge_ui    UIBridge

	// pdf_doc
	doc *fitz.Document
}

func NewPDFView(doc_id string, bridge_trunk TrunkBridge, bridge_ui UIBridge) *PDFViewer {
	self := PDFViewer{
		motion:           NewMotion(),
		page_buffer:      map[int]*PageInfo{},
		thumbnail_buffer: map[int]*ImgInfo{},
		pad:              8,

		size_buffer: make(map[int]*image.Rectangle),

		offsetX_tar: 0,
		offsetX_cur: 0,
		offsetY_tar: 0,
		offsetY_cur: 0,
		scale_tar:   1.0,
		scale_cur:   1.0,

		focus_index: 0,
		max_buffer:  5,

		base_height:   1000.0,
		viewer_radius: 8,

		doc_id:    doc_id,
		doc_scale: 1,

		bridge_trunk: bridge_trunk,
		bridge_ui:    bridge_ui,
	}

	for i := range 100 {
		size, _ := self.bridge_trunk.Get_size(self.doc_id, i)
		self.size_buffer[i] = &size
	}

	return &self
}

// 界面渲染
func (self *PDFViewer) Layout(gtx C, style *Style) D {
	self.style = style

	// 按钮
	resetBtn := material.Button(self.style.Theme, &self.reset_btn, "reset")
	if self.reset_btn.Clicked(gtx) {
		self.center_aligned(false)
	}

	self.add_thumbnail_buffer()

	// viewer 布局
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(resetBtn.Layout),
		layout.Flexed(1, self.layout_viewer),
	)

	// 更新页面
	self.add_new_page()

	// 动画效果
	if self.scale_cur != self.scale_tar || self.offsetY_cur != self.offsetY_tar || self.offsetX_cur != self.offsetX_tar {
		style.AniSys.Ease64Group(&self.offsetY_cur, []EaseItem{
			{Cur: &self.scale_cur, Tar: self.scale_tar, Precision: 0.01},
			{Cur: &self.offsetY_cur, Tar: self.offsetY_tar, Precision: 1},
			{Cur: &self.offsetX_cur, Tar: self.offsetX_tar, Precision: 1},
		}, 0.4)
	}

	return dims
}

// ui部分 ------------------------------------------------------------------------------------------------------

// 内容渲染区域
func (self *PDFViewer) layout_viewer(gtx C) D {
	self.viewer_area = gtx.Constraints.Max

	if self.bridge_ui.IsWinSizeChange() {
		self.center_aligned(true)
	}

	l := layout.Stack{}.Layout(gtx,

		// 内容渲染
		layout.Stacked(func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			defer op.Offset(image.Pt(int(self.offsetX_cur), int(self.offsetY_cur))).Push(gtx.Ops).Pop()
			defer op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(float32(self.scale_cur), float32(self.scale_cur)))).Push(gtx.Ops).Pop()

			pdf_dims := func(gtx C) D {
				gtx.Constraints.Max = image.Pt(gtx.Dp(1e6), gtx.Dp(1e6))
				gtx.Constraints.Min = image.Pt(0, 0)

				for index, inst := range self.page_buffer {
					size := inst.Size
					ratio := self.base_height / float64(size.Dy())
					cont_h := self.base_height
					cont_w := float64(size.Dx()) * ratio

					offset_x := -int(cont_w * 0.5)
					offset_y := index * (int(self.base_height) + self.pad)

					trans := op.Offset(image.Pt(offset_x, offset_y)).Push(gtx.Ops)
					clip := clip.RRect{Rect: image.Rect(0, 0, int(cont_w), int(cont_h)), SW: self.viewer_radius, SE: self.viewer_radius, NW: self.viewer_radius, NE: self.viewer_radius}.Push(gtx.Ops)

					if inst.Widget == nil {
						draw_rectangle(gtx, image.Pt(int(cont_w), int(cont_h)), COLOR_WHITE, 0)
					} else {
						s := float32(float64(size.Dy()) / float64(inst.Widget.Src.Size().Y))
						self.doc_scale = float64(s)
						scale := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0.0, 0.0), f32.Pt(s*float32(ratio), s*float32(ratio)))).Push(gtx.Ops)
						inst.Widget.Scale = 1.0 / gtx.Metric.PxPerDp
						inst.Widget.Layout(gtx)
						scale.Pop()
					}

					clip.Pop()
					trans.Pop()
				}

				return D{}
			}(gtx)
			draw_coordinate_system(gtx)
			return pdf_dims
		}),

		// 事件检测
		layout.Expanded(func(gtx C) D {
			action := self.motion.Update(gtx, self.viewer_area)
			switch action.Kind {
			case ActionScroll:
				scroll_delta := action.Scroll / self.scale_cur
				self.offsetY_tar -= scroll_delta

			case ActionZoom:
				if action.Scroll != 0 {
					oldScale := self.scale_tar
					zoomFactor := 1.3
					if action.Scroll > 0 {
						self.scale_tar /= zoomFactor
					} else {
						self.scale_tar *= zoomFactor
					}

					// 设置最小缩放倍率
					self.scale_tar = math.Max(self.scale_tar, 0.1)

					// 变化公式：NewOffset = MousePos - (MousePos - OldOffset) * (NewScale / OldScale)
					scaleRatio := self.scale_tar / oldScale

					size := self.page_buffer[self.focus_index].Size
					cont_w := float64(size.Dx()) * self.base_height / float64(size.Dy())
					if cont_w*self.scale_tar < float64(self.viewer_area.X)*0.8 {
						self.offsetX_tar = (float64(self.viewer_area.X)-cont_w*self.scale_tar)*0.5 + cont_w*0.5*self.scale_tar
					} else {
						self.offsetX_tar = float64(action.Pos.X) - (float64(action.Pos.X)-self.offsetX_tar)*scaleRatio
					}

					self.offsetY_tar = float64(action.Pos.Y) - (float64(action.Pos.Y)-self.offsetY_tar)*scaleRatio
				}

			case ActionMove:
			default:
			}
			return D{}
		}),
	)

	return l
}

// 逻辑部分 ---------------------------------------------------------------------------------------------------

// 增加新的页面
func (self *PDFViewer) add_new_page() {
	// 增加内容
	_, y := self.screen_to_world(0, float64(self.viewer_area.Y)*0.5)
	p := int(math.Max(0, y/self.base_height))
	self.focus_index = p

	is_have := false
	for i, _ := range self.page_buffer {
		if i == p {
			is_have = true
			break
		}
	}

	if !is_have {
		page := p
		self.page_buffer[p] = &PageInfo{
			Num: p,
			// Scale:  1,
			Widget: nil,
		}

		go func(page int) {
			size, _ := self.bridge_trunk.Get_size(self.doc_id, p)

			self.page_buffer[p].Size = size
			self.bridge_trunk.Add_render_page(self.doc_id, page, func(img *image.RGBA) {
				imgop := paint.NewImageOp(img)
				imgop.Filter = paint.FilterNearest
				self.page_buffer[p].Widget = &widget.Image{
					Src: imgop,
					Fit: widget.Unscaled,
				}
				self.bridge_ui.WindowRefresh()
			})
		}(page)
	}
}

// 检测某一页在某个列表中是否已经存在
func (self *PDFViewer) check_page_alive(page int, dict *map[int]*PageInfo) bool {
	for k, _ := range *dict {
		if k == page {
			return true
		}
	}
	return false
}

// 缓冲区
func (self *PDFViewer) add_thumbnail_buffer() {
	_, y := self.screen_to_world(0, float64(self.viewer_area.Y)*0.5)
	p := int(math.Max(0, y/self.base_height))

	load_page := []int{p}
	keepPages := make(map[int]bool)
	keepPages[p] = true

	// 判断当前缓冲区中的页
	for i := range 5 {
		np := p + (i + 1)
		lp := p - (i + 1)

		load_page = append(load_page, np)
		keepPages[np] = true

		if lp >= 0 {
			load_page = append(load_page, lp)
			keepPages[lp] = true
		}

	}

	for n := range self.page_buffer {
		if !keepPages[n] {
			delete(self.page_buffer, n)
		}
	}

	// for _, page_num := range load_page {
	// 	if _, exists := self.page_buffer[page_num]; !exists {
	// 		info := &PageInfo{
	// 			Num:    page_num,
	// 			Widget: nil,
	// 			Size:   image.Rect(0, 0, 0, 0),
	// 		}
	// 		self.page_buffer[page_num] = info

	// 		go func(page int, page_info *PageInfo) {
	// 			size, _ := self.bridge_trunk.Get_size(self.doc_id, page)
	// 			page_info.Size = size

	// 			self.bridge_trunk.Add_render_page(self.doc_id, page_num, func(img *image.RGBA) {

	// 				imgop := paint.NewImageOp(img)
	// 				if keepPages[page] {
	// 					page_info.Widget = &widget.Image{
	// 						Src: imgop,
	// 						Fit: widget.Unscaled,
	// 					}
	// 				}

	// 				self.bridge_ui.WindowRefresh()
	// 			})
	// 		}(page_num, info)
	// 	}
	// }

	for _, page := range load_page {
		if self.page_buffer[page] == nil {
			info := PageInfo{
				Num:    page,
				Widget: nil,
				Size:   image.Rect(0, 0, 1, 1),
			}

			self.mu.Lock()
			self.page_buffer[page] = &info
			self.mu.Unlock()
		}
	}

	ld := []int{}
	for _, p := range load_page {
		if self.page_buffer[p].Widget == nil {
			ld = append(ld, p)
		}
	}

	self.bridge_trunk.Add_render_page_list(self.doc_id, ld, func(page int, img *image.RGBA) {
		go func() {
			start := time.Now()
			if info, exist := self.page_buffer[page]; exist {

				size, _ := self.bridge_trunk.Get_size(self.doc_id, page)
				info.Size = size
				info.Widget = &widget.Image{
					Src: paint.NewImageOp(img),
					Fit: widget.Unscaled,
				}
				info.Widget.Src.Filter = paint.FilterNearest
			}
			self.bridge_ui.WindowRefresh()
			log.Println("应用耗时:%v", time.Since(start))
		}()
	})
}

// 添加页面
func (self *PDFViewer) add_page() {

}

// 缩放变换

// 其他 -----------------------------------------------------------------------------------------------------

// 当前页居中
//   - true 立即应用
//   - false 使用动画效果
func (self *PDFViewer) center_aligned(instant bool) {
	if len(self.page_buffer) == 0 {
		return
	}

	// 页面尺寸
	size := self.page_buffer[self.focus_index].Size
	ratio := self.base_height / float64(size.Dy())
	page_w := float64(size.Dx()) * ratio
	page_h := self.base_height

	// 显示区域大小
	area_w := float64(self.viewer_area.X)
	area_h := float64(self.viewer_area.Y)

	self.scale_tar = math.Min(area_w/page_w, area_h/page_h)
	self.offsetX_tar = math.Round((area_w-page_w*self.scale_tar)*.5 + page_w*self.scale_tar*0.5)
	self.offsetY_tar = math.Round((area_h-page_h*self.scale_tar)*.5 - float64(self.focus_index)*(page_h+float64(self.pad))*self.scale_tar - 0.5)

	if instant {
		self.scale_cur = self.scale_tar
		self.offsetX_cur = self.offsetX_tar
		self.offsetY_cur = self.offsetY_tar
	}
}

// 世界坐标转屏幕坐标
// func (self *PDFViewer) world_to_screen(x float64, y float64) (float64, float64) {
// 	scale := self.scale_cur
// 	offset_x := self.offsetX_cur
// 	offset_y := self.offsetY_cur
// 	pivot_x := float64(self.pointer_pos_cur.X)
// 	pivot_y := float64(self.pointer_pos_cur.Y)
// 	screen_x := scale*(x-pivot_x) + pivot_x + offset_x
// 	screen_y := scale*(y-pivot_y) + pivot_y + offset_y
// 	return screen_x, screen_y
// }

// 屏幕坐标转世界坐标
func (self *PDFViewer) screen_to_world(x, y float64) (float64, float64) {
	wx := (x - self.offsetX_cur) / self.scale_cur
	wy := (y - self.offsetY_cur) / self.scale_cur
	return wx, wy
}
