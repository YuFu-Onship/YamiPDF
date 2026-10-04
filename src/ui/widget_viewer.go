package ui

import (
	rumia "app/src"
	"app/src/tool"
	"fmt"
	"image"
	"log"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// 渲染着色器
type modeShader int

const (
	mode_shader_oriangle modeShader = iota
	mode_shader_gameboy
	mode_shader_warm
	mode_shader_gray
	mode_shader_binary
	mode_shader_end
)

func select_mode_shader(img *image.RGBA, mode modeShader) *image.RGBA {
	switch mode {
	case mode_shader_oriangle:
		return img
	case mode_shader_gameboy:
		return tool.Shader_gameboy(img)
	case mode_shader_warm:
		return tool.Shader_warm(img)
	case mode_shader_gray:
		return tool.Shader_gray(img)
	case mode_shader_binary:
		return tool.Shader_binary(img)
	default:
		return img
	}
}

// 布局模式
type modeLayout int

const (
	mode_layout_single        modeLayout = iota // 单页
	mode_layout_double                          // 双页
	mode_layout_single_double                   // 单页+双页
	mode_layout_end
)

// 填充模式
type modeFill int

const (
	mode_fill_fit   modeFill = iota // 整页
	mode_fill_width                 // 填充宽度
	mode_fill_end
)

// 阅读顺序模式
type modeOrder int

const (
	mode_order_left  modeOrder = iota // 从左往右
	mode_order_right                  // 从右往左
	mode_order_end
)

// 渲染模式
type modeRender int

const (
	mode_render_page modeRender = iota
	mode_render_tile
)

// 接口
type ViewInterface interface {
	Render(filePath string, page int) (*image.RGBA, error)
}

// 页元素信息
type PageInfo struct {
	Num    int             // 当前页索引
	Size   image.Rectangle // 当前页尺寸
	Widget *widget.Image   // 当前页图像
	img    *image.RGBA
	DPI    float64 // 当前页的渲染dpi

	mode_shader modeShader

	box_w, box_h float64    // 当前页面元素尺寸,高度始终为base_height
	box          [4]float64 // 包围框信息,x,y,w,h

	zoom_img       *image.RGBA
	zoom_widget    *widget.Image // 放大显示的瓦片
	zoom_rect      *types.Rectangle
	zoom_rect_last *types.Rectangle
	zoom_dpi       float64          // 瓦片dpi
	zoom_size      image.Rectangle  // 瓦片尺寸
	zoom_ratio     *types.Rectangle //百分比值,左上,右下
}

// 骨架信息
type SkeletonInfo struct {
	Num  int
	Pos  image.Point
	Size image.Rectangle
}

func (self *PageInfo) Destory() {
	if self == nil {
		return
	}

	// 清空页
	if self.img != nil {
		self.img.Pix = nil
		self.img = nil
	}
	self.Widget = nil

	// 清空瓦片
	if self.zoom_img != nil {
		self.zoom_img.Pix = nil
		self.zoom_img = nil
	}
	self.zoom_widget = nil
}

// pdf渲染器 ---------------------------------------------------------
type PDFViewer struct {
	motion *Motion // 动作检测
	style  *Style  // 样式主题

	buffer_page     map[int]*PageInfo        // 页缓冲区
	buffer_size     map[int]*image.Rectangle // 尺寸缓冲区
	buffer_skeleton []int                    // 骨架缓冲区
	buffer_box_size map[int][2]float64       // 元素包围框 w,h (世界坐标系)
	buffer_box_pos  map[int][2]float64       // 元素包围框 x,y (世界坐标系)

	offsetX_tar float64
	offsetX_cur float64
	offsetY_tar float64
	offsetY_cur float64
	scale_tar   float64
	scale_cur   float64

	rotateArc_tar float64
	rotateArc_cur float64

	index_focus  int // 当前聚焦页
	index_center int // 当前居中页

	clickable_reset   widget.Clickable // 按钮组件
	clickable_next    widget.Clickable
	clickable_last    widget.Clickable
	clickable_rotate  widget.Clickable
	clickable_shader  widget.Clickable
	clickable_layout  widget.Clickable
	clickable_order   widget.Clickable
	clickable_release widget.Clickable

	doc_id       string         // 文档id
	doc_scale    float64        // 文档控件尺寸/文件实际尺寸
	doc_info     map[string]any // 文档信息
	doc_page_num int            // 文档页数

	viewer_area image.Point // 渲染区域

	pad         float64 // 页间距
	page_radius int     // 固定的元素圆角半径
	base_height float64 // 固定的元素高

	bridge_trunk TrunkBridge // 外部接口
	bridge_ui    UIBridge    // 内部ui接口

	dpi float64 // 全局dpi

	debouncer tool.Debouncer // 抖动检测

	mode_render modeRender // 渲染模式
	mode_shader modeShader
	mode_layout modeLayout
	mode_fill   modeFill
	mode_order  modeOrder

	is_update_page bool // 更新页面
}

// 新的pdf渲染窗口
func NewPDFView(doc_id string, bridge_trunk TrunkBridge, bridge_ui UIBridge) *PDFViewer {
	self := PDFViewer{
		motion: NewMotion(),
		pad:    8,

		buffer_page:     map[int]*PageInfo{},
		buffer_size:     make(map[int]*image.Rectangle),
		buffer_skeleton: []int{},
		buffer_box_size: map[int][2]float64{},
		buffer_box_pos:  map[int][2]float64{},

		offsetX_tar: 0,
		offsetX_cur: 0,
		offsetY_tar: 0,
		offsetY_cur: 0,
		scale_tar:   1.0,
		scale_cur:   1.0,

		rotateArc_tar: 0.0,
		rotateArc_cur: 0.0,

		index_focus: 0,

		base_height: 1000.0,
		page_radius: 8,

		doc_id:    doc_id,
		doc_scale: 1,

		bridge_trunk: bridge_trunk,
		bridge_ui:    bridge_ui,

		mode_render: mode_render_page,
		debouncer:   *tool.NewDebouncer(time.Millisecond * 100),

		is_update_page: true,

		clickable_reset:   widget.Clickable{},
		clickable_next:    widget.Clickable{},
		clickable_last:    widget.Clickable{},
		clickable_rotate:  widget.Clickable{},
		clickable_shader:  widget.Clickable{},
		clickable_layout:  widget.Clickable{},
		clickable_order:   widget.Clickable{},
		clickable_release: widget.Clickable{},
	}

	self.doc_info = self.bridge_trunk.Get_info(doc_id)
	self.doc_page_num = self.doc_info["page_num"].(int)

	self.mode_shader = mode_shader_oriangle
	self.mode_layout = mode_layout_double
	self.mode_order = mode_order_right

	for i := range self.doc_info["page_num"].(int) {
		size, _ := self.bridge_trunk.Get_size(self.doc_id, i)
		self.buffer_size[i] = &size
	}

	return &self
}

// 界面渲染
func (self *PDFViewer) Layout(gtx C, style *Style) D {
	self.style = style

	// self.offsetY_tar = math.Min(0, self.offsetY_tar)
	// self.offsetY_cur = math.Min(0, self.offsetY_cur)

	// 按钮
	btn_reset := material.Button(self.style.Theme, &self.clickable_reset, "重置")
	btn_reset.Background = self.style.Palette.Bg_3
	btn_reset.Color = self.style.Palette.Fg_3

	if self.clickable_reset.Clicked(gtx) {
		self.apply_center(self.index_focus)
		// p := rand.N(self.doc_page_num)
		// self.apply_center(p)
		self.is_update_page = true
	}

	btn_next := material.Button(self.style.Theme, &self.clickable_next, "下一页")
	btn_next.Background = self.style.Palette.Bg_2
	btn_next.Color = self.style.Palette.Fg_2
	if self.clickable_next.Clicked(gtx) {
		// self.index_focus += 1
	}

	btn_last := material.Button(self.style.Theme, &self.clickable_last, "上一页")
	btn_last.Background = self.style.Palette.Bg_1
	btn_last.Color = self.style.Palette.Fg_1
	if self.clickable_last.Clicked(gtx) {
		// self.index_focus -= 1

	}

	btn_rotate := material.Button(self.style.Theme, &self.clickable_rotate, "旋转")
	btn_rotate.Background = self.style.Palette.Bg_3
	btn_rotate.Color = self.style.Palette.Fg_3
	if self.clickable_rotate.Clicked(gtx) {
		self.rotateArc_tar += math.Pi * 0.5
	}

	btn_shader := material.Button(self.style.Theme, &self.clickable_shader, "着色器")
	btn_shader.Background = self.style.Palette.Bg_3
	btn_shader.Color = self.style.Palette.Fg_3
	if self.clickable_shader.Clicked(gtx) {
		self.mode_shader = (self.mode_shader + 1) % mode_shader_end
		self.is_update_page = true
	}

	btn_layout := material.Button(self.style.Theme, &self.clickable_layout, fmt.Sprintf("布局%v", self.mode_layout))
	btn_layout.Background = self.style.Palette.Bg_3
	btn_layout.Color = self.style.Palette.Fg_3
	if self.clickable_layout.Clicked(gtx) {
		self.mode_layout = (self.mode_layout + 1) % mode_layout_end
		self.is_update_page = true
		// self.apply_center(self.index_focus)
		self.apply_center_pos(self.index_focus)
	}

	btn_order := material.Button(self.style.Theme, &self.clickable_order, "顺序")
	btn_order.Background = self.style.Palette.Bg_3
	btn_order.Color = self.style.Palette.Fg_3
	if self.clickable_order.Clicked(gtx) {
		self.mode_order = (self.mode_order + 1) % mode_order_end
	}

	// 	btn_release := material.Button(self.style.Theme, &self.clickable_release, "释放")
	// 	if self.clickable_release.Clicked(gtx) {
	// 	}

	// 界面部分
	self.dpi = self.get_dpi()

	// 更新缓冲区
	if self.viewer_area.Y != 0 {
		self.update_buffer()
	}

	// viewer 绘制布局
	dims := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, btn_reset.Layout),
				layout.Flexed(1, btn_next.Layout),
				layout.Flexed(1, btn_last.Layout),
				layout.Flexed(1, btn_rotate.Layout),
				layout.Flexed(1, btn_shader.Layout),
				layout.Flexed(1, btn_layout.Layout),
				layout.Flexed(1, btn_order.Layout),
			)
		}),
		layout.Flexed(1, self.layout_viewer),
	)

	// 一些参数的限定值
	self.scale_tar = math.Round(self.scale_tar*100) / 100

	// 动画效果
	if self.scale_cur != self.scale_tar || self.offsetY_cur != self.offsetY_tar || self.offsetX_cur != self.offsetX_tar {
		style.AniSys.Ease64Group(&self.offsetY_cur, []EaseItem{
			{Cur: &self.scale_cur, Tar: self.scale_tar, Precision: 0.01},
			{Cur: &self.offsetY_cur, Tar: self.offsetY_tar, Precision: 0.1},
			{Cur: &self.offsetX_cur, Tar: self.offsetX_tar, Precision: 0.1},
		}, 0.4)
	}

	style.AniSys.Ease64(&self.rotateArc_cur, 0.0005, &self.rotateArc_cur, self.rotateArc_tar, 0.3)

	return dims
}

// ui部分 ------------------------------------------------------------------------------------------------------

// 内容渲染区域
func (self *PDFViewer) layout_viewer(gtx C) D {
	self.viewer_area = gtx.Constraints.Max

	self.mode_render = mode_render_page
	if self.base_height*self.scale_tar*0.8 > float64(self.viewer_area.Y) {
		self.mode_render = mode_render_tile
	}

	// if self.bridge_ui.IsWinSizeChange() {
	// self.center_aligned(0, false)
	// }

	draw_rectangle_2(gtx, 0, self.viewer_area.Y/2-2, self.viewer_area.X, 4, COLOR_ORANGE, 0)
	l := layout.Stack{}.Layout(gtx,
		// 内容渲染
		layout.Stacked(func(gtx C) D {

			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			defer op.Offset(image.Pt(int(self.offsetX_cur), int(self.offsetY_cur))).Push(gtx.Ops).Pop()
			defer op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(float32(self.scale_cur), float32(self.scale_cur)))).Push(gtx.Ops).Pop()
			pdf_dims := func(gtx C) D {
				gtx.Constraints.Max = image.Pt(gtx.Dp(1e6), gtx.Dp(1e6))
				gtx.Constraints.Min = image.Pt(0, 0)

				// 绘制骨架
				for _, num := range self.buffer_skeleton {
					if info, exist := self.buffer_page[num]; exist {
						if info.Widget != nil {
							continue
						}
					}

					// 局部变换参数计算
					size := self.buffer_size[num]
					sin_a, cos_a := math.Abs(math.Sin(self.rotateArc_cur)), math.Abs(math.Cos(self.rotateArc_cur))
					pdf_w, pdf_h := float64(size.Dx()), float64(size.Dy())
					local_scale_coe := self.base_height / (pdf_w*sin_a + pdf_h*cos_a)
					min_x, min_y, _, _ := self.get_page_box(num)

					// 世界坐标位移
					global_offset := op.Offset(image.Pt(
						int(self.buffer_box_pos[num][0]),
						int(self.buffer_box_pos[num][1]),
					)).Push(gtx.Ops)

					self.trans_local(
						gtx,
						-min_x, -min_y,
						self.rotateArc_cur,
						local_scale_coe,
						pdf_w,
						pdf_h,
						func() {
							draw_rectangle(gtx, image.Pt(int(pdf_w), int(pdf_h)), COLOR_WHITE, 0)
						})

					global_offset.Pop()
				}

				// 绘制内容
				for num, inst := range self.buffer_page {
					size := self.buffer_size[inst.Num]

					// 局部变换参数计算
					sin_a, cos_a := math.Abs(math.Sin(self.rotateArc_cur)), math.Abs(math.Cos(self.rotateArc_cur))
					pdf_w, pdf_h := float64(size.Dx()), float64(size.Dy())
					local_scale_coe := self.base_height / (pdf_w*sin_a + pdf_h*cos_a)

					min_x, min_y, box_w, box_h := self.get_page_box(inst.Num)
					inst.box_w, inst.box_h = box_w, box_h

					// self.buffer_box[num] = [4]float64{0, 0, box_w, box_h}

					// 世界坐标位移
					self.calc_box_pos_global(num)
					global_offset := op.Offset(image.Pt(
						int(self.buffer_box_pos[num][0]),
						int(self.buffer_box_pos[num][1]),
					)).Push(gtx.Ops)

					self.trans_local(
						gtx,
						-min_x, -min_y,
						self.rotateArc_cur,
						local_scale_coe,
						pdf_w,
						pdf_h,
						func() {
							self.draw_full(gtx, inst)
							self.draw_tile(gtx, inst)
						})

					global_offset.Pop()
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
				scroll_delta := (action.Scroll * float64(self.viewer_area.Y) * 0.0007) / self.scale_cur
				self.offsetY_tar -= scroll_delta

				self.debouncer.Do(func() { self.is_update_page = true })

			case ActionZoom:
				if action.Scroll != 0 {
					oldScale := self.scale_tar
					zoomFactor := 1.3
					if action.Scroll > 0 {
						self.scale_tar /= zoomFactor
					} else {
						self.scale_tar *= zoomFactor
					}

					self.scale_tar = math.Max(0.15, math.Min(10, self.scale_tar))
					self.scale_tar = math.Round(self.scale_tar*100) / 100

					// 设置最小缩放倍率

					// 变化公式：NewOffset = MousePos - (MousePos - OldOffset) * (NewScale / OldScale)
					scaleRatio := self.scale_tar / oldScale

					// size := self.buffer_page[self.index_focus].Size
					// size := self.buffer_size[self.index_focus]
					// cont_w := float64(size.Dx()) * self.base_height / float64(size.Dy())

					cw, _ := self.calc_row_page_size(self.index_focus)

					if cw*self.scale_tar < float64(self.viewer_area.X)*0.8 {
						self.offsetX_tar = (float64(self.viewer_area.X)-cw*self.scale_tar)*0.5 + cw*0.5*self.scale_tar
					} else {
						self.offsetX_tar = float64(action.Pos.X) - (float64(action.Pos.X)-self.offsetX_tar)*scaleRatio
					}

					self.offsetY_tar = float64(action.Pos.Y) - (float64(action.Pos.Y)-self.offsetY_tar)*scaleRatio
				}

				self.debouncer.Do(func() { self.is_update_page = true })

			case ActionMove:
			default:
			}
			return D{}
		}),
	)

	var s int64
	for _, v := range self.buffer_page {
		if v.img != nil {
			s += int64(len(v.img.Pix))
		}
	}

	return l
}

// 局部变换, 整页与瓦片绘制, widget坐标系层级
func (self *PDFViewer) trans_local(gtx C, x, y, arc, scale, pdf_w, pdf_h float64, context func()) {
	local_offset := op.Offset(image.Pt(int(x), int(y)))
	local_rotate := f32.Affine2D{}.Rotate(f32.Pt(0, 0), float32(arc))
	local_clip := clip.RRect{
		Rect: image.Rect(0, 0, int(pdf_w*scale), int(pdf_h*scale)),
		SW:   self.page_radius,
		SE:   self.page_radius,
		NW:   self.page_radius,
		NE:   self.page_radius,
	}
	local_scale := f32.Affine2D{}.Scale(
		f32.Pt(0, 0),
		f32.Pt(float32(scale), float32(scale)),
	)

	// 局部变换
	apply_offset := local_offset.Push(gtx.Ops)
	apply_rotate := op.Affine(local_rotate).Push(gtx.Ops)
	apply_clip := local_clip.Push(gtx.Ops)
	apply_scale := op.Affine(local_scale).Push(gtx.Ops)

	context()

	apply_scale.Pop()
	apply_clip.Pop()
	apply_rotate.Pop()
	apply_offset.Pop()
}

// 绘制整页
func (self *PDFViewer) draw_full(gtx C, inst *PageInfo) {
	// 整页绘制
	if inst.Widget != nil {
		// img -> pdf 坐标系
		size := self.buffer_size[inst.Num]
		scale_coe := float32(size.Dy()) / float32(inst.Widget.Src.Size().Y)
		inst.Widget.Scale = 1.0 / gtx.Metric.PxPerDp
		scale := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0.0, 0.0), f32.Pt(scale_coe, scale_coe))).Push(gtx.Ops)
		inst.Widget.Layout(gtx)
		scale.Pop()
	}
}

// 绘制瓦片
func (self *PDFViewer) draw_tile(gtx C, inst *PageInfo) {
	if inst.zoom_widget != nil {
		ratio_1 := (inst.zoom_ratio.UR.X - inst.zoom_ratio.LL.X) / float64(inst.zoom_widget.Src.Size().X)
		tile_ratio := ratio_1

		x_1 := inst.zoom_ratio.LL.X
		y_1 := inst.zoom_ratio.LL.Y

		tile_affine := f32.Affine2D{}.
			Scale(f32.Pt(0, 0), f32.Pt(float32(tile_ratio), float32(tile_ratio))).
			Offset(f32.Pt(float32(x_1), float32(y_1)))

		tile_transform := op.Affine(tile_affine).Push(gtx.Ops)
		inst.zoom_widget.Scale = 1.0 / gtx.Metric.PxPerDp
		inst.zoom_widget.Layout(gtx)
		tile_transform.Pop()
		// s := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(float32(tile_ratio), float32(tile_ratio)))
		// as := op.Affine(s).Push(gtx.Ops)
		// as.Pop()
	}
}

// 得到指定页元素尺寸包围框信息,widget坐标系
//   - x,y,w,h (局部坐标系)
func (self *PDFViewer) get_page_box(page int) (float64, float64, float64, float64) {
	size := self.buffer_size[page]
	arc := self.rotateArc_cur

	// 局部变换参数计算
	sin_a, cos_a := math.Abs(math.Sin(self.rotateArc_cur)), math.Abs(math.Cos(self.rotateArc_cur))

	pdf_w, pdf_h := float64(size.Dx()), float64(size.Dy())

	ele_w := pdf_w*cos_a + pdf_h*sin_a
	ele_h := pdf_w*sin_a + pdf_h*cos_a

	local_scale_coe := self.base_height / ele_h

	// 计算位移量
	w, h := pdf_w*local_scale_coe, pdf_h*local_scale_coe
	x0, y0 := tool.Rotate_point(w, 0, 0, 0, arc)
	x1, y1 := tool.Rotate_point(w, h, 0, 0, arc)
	x2, y2 := tool.Rotate_point(0, h, 0, 0, arc)
	x3, y3 := 0.0, 0.0
	min_x, min_y := min(x0, x1, x2, x3), min(y0, y1, y2, y3)

	box_w, box_h := ele_w*local_scale_coe, ele_h*local_scale_coe
	return min_x, min_y, box_w, box_h
}

// 更新缓冲区
func (self *PDFViewer) update_buffer() {
	self.dpi = self.get_dpi()
	pages := self.calc_page_list()

	self.index_focus = self.index_center

	page := self.calc_center_count(pages)
	self.index_center = page

	// 尺寸缓冲区
	size_index_min := max(page-30, 0)
	size_index_max := min(page+30, self.doc_page_num-1)

	if self.buffer_size[size_index_min] == nil || self.buffer_size[size_index_max] == nil {
		size_index_start := max(page-60, 0)
		size_index_end := min(page+60, self.doc_page_num-1)

		for i := size_index_start; i <= size_index_end; i++ {
			r, _ := self.bridge_trunk.Get_size(self.doc_id, i)
			self.buffer_size[i] = &r
		}
	}

	// 骨架缓冲区
	self.buffer_skeleton = []int{}
	for i := range len(pages) * 3 {
		if last := page - i; last >= 0 {
			self.buffer_skeleton = append(self.buffer_skeleton, last)
		}
		if next := page + i; next < self.doc_page_num {
			self.buffer_skeleton = append(self.buffer_skeleton, next)
		}
	}

	// 包围盒缓冲区
	for _, page := range self.buffer_skeleton {
		box_w, box_h := self.calc_box_size(page)
		self.buffer_box_size[page] = [2]float64{box_w, box_h}
		x, y := self.calc_box_pos_global(page)
		self.buffer_box_pos[page] = [2]float64{x, y}
	}

	// 渲染缓冲区
	if self.is_update_page {
		self.is_update_page = false

		for _, page := range pages {
			if self.buffer_page[page] == nil && self.buffer_size != nil {
				info := PageInfo{
					Num:    page,
					Widget: nil,
					Size:   *self.buffer_size[page],
					DPI:    self.dpi,
				}
				self.buffer_page[page] = &info
			}
		}

		if self.mode_render == mode_render_tile {
			self.add_tile(pages)
		} else if self.mode_render == mode_render_page {
			self.is_update_pages(pages)
		}
	}
}

// 更新某一目标页附近的缓冲区
func (self *PDFViewer) update_buffer_target(page int) {
	if page < 0 || page >= self.doc_page_num {
		return
	}
	self.apply_center(page)
	self.update_buffer()
}

// 得到当前视口横轴中心的页号
func (self *PDFViewer) get_focus_page() int {
	_, y := self.from_screen_to_world(0, float64(self.viewer_area.Y)*0.5)
	p := int(math.Max(0, y/self.base_height))
	return p
}

// 增加新的页面
func (self *PDFViewer) is_update_pages(pages []int) {
	// 当前视口的页

	// 清除页缓冲区中的无用页
	keep_pages := make(map[int]bool)
	for _, p := range pages {
		keep_pages[p] = true
	}

	for k, v := range self.buffer_page {
		if !keep_pages[k] {
			if v != nil {
				v.Destory()
				self.bridge_ui.WindowRefresh()
			}
			delete(self.buffer_page, k)
		}
	}

	// 更新页缓冲区中的新加页信息
	// for _, page := range pages {
	// 	if self.buffer_page[page] == nil && self.buffer_size != nil {
	// 		info := PageInfo{
	// 			Num:    page,
	// 			Widget: nil,
	// 			Size:   *self.buffer_size[page],
	// 			DPI:    self.dpi,
	// 		}
	// 		self.buffer_page[page] = &info
	// 	}
	// }

	// 渲染页缓冲区中的图像
	full_tasks := []*rumia.RenderFullTask{}
	for _, p := range pages {
		if self.buffer_page[p].Widget == nil ||
			self.buffer_page[p].DPI != self.dpi ||
			self.buffer_page[p].mode_shader != self.mode_shader {

			full_tasks = append(full_tasks, &rumia.RenderFullTask{
				Page:    p,
				DPI:     self.dpi,
				Call_fn: func(img *image.RGBA) { self.callback_apply_full(p, self.dpi, self.mode_shader, img) },
			})
		}
	}

	self.bridge_trunk.Add_render_page_list(self.doc_id, full_tasks)
}

// update tiles
func (self *PDFViewer) update_tile(pages []int) {

}

// 增加瓦片
func (self *PDFViewer) add_tile(pages []int) {
	tasks := self.get_tile_rect(pages)
	for _, p := range self.buffer_page {
		if p.zoom_dpi != self.dpi || p.zoom_widget == nil {
			self.bridge_trunk.Add_render_tile_list(self.doc_id, tasks)
		}
	}
}

// 得到当前视口的瓦片信息
func (self *PDFViewer) get_tile_rect(pages []int) []*rumia.RenderTileTask {
	tasks := []*rumia.RenderTileTask{}
	// for i := self.index_focus - 1; i <= self.index_focus+1; i += 1 {
	for _, i := range pages {
		if i >= 0 && i < self.doc_page_num {
			size := self.buffer_size[i]

			// 有效显示区域坐标
			box_w, box_h := self.buffer_box_size[i][0], self.buffer_box_size[i][1]

			// 在视口坐标系上的位移坐标
			box_x, box_y := self.from_world_to_screen(self.buffer_box_pos[i][0], self.buffer_box_pos[i][1])

			// 视口坐标系, 显示区域
			ll_x := math.Max(0, box_x)
			ll_y := math.Max(0, box_y)
			ur_x := math.Min(float64(self.viewer_area.X), box_x+box_w*self.scale_tar)
			ur_y := math.Min(float64(self.viewer_area.Y), box_y+(box_h+self.pad)*self.scale_tar+box_h*self.scale_tar)

			if ur_y <= 0 || ll_y >= float64(self.viewer_area.Y) {
				continue
			}

			// 十字参数
			// (   )(h_1)(   )
			// (w_1)(w,h)(w_2)
			// (   )(h_2)(   )

			widget_w := box_w * self.scale_tar
			widget_h := box_h * self.scale_tar

			h_1 := ll_y - box_y
			h_2 := ur_y - ll_y
			h_3 := widget_h - h_1 - h_2

			w_1 := ll_x - box_x
			w_2 := ur_x - ll_x
			w_3 := widget_w - w_1 - w_2

			var x_nor, y_nor, w_nor, h_nor float64

			state := self.to_arc_state(self.rotateArc_tar)
			switch state {
			case 0:
				x_nor = w_1 / widget_w
				y_nor = h_1 / widget_h
				w_nor = w_2 / widget_w
				h_nor = h_2 / widget_h
			case 1:
				x_nor = h_1 / widget_h
				y_nor = w_3 / widget_w
				w_nor = h_2 / widget_h
				h_nor = w_2 / widget_w
			case 2:
				x_nor = w_3 / widget_w
				y_nor = h_3 / widget_h
				w_nor = w_2 / widget_w
				h_nor = h_2 / widget_h
			case 3:
				x_nor = h_3 / widget_h
				y_nor = w_1 / widget_w
				w_nor = h_2 / widget_h
				h_nor = w_2 / widget_w
			}

			px := float64(size.Dx()) * x_nor
			py := float64(size.Dy()) * y_nor
			pw := float64(size.Dx()) * w_nor
			ph := float64(size.Dy()) * h_nor

			px = max(0, px)
			py = max(0, py)
			pw = min(float64(size.Dx())-px, pw)

			// pdf坐标,y轴反转
			py_pdf := float64(size.Dy()) * (1.0 - (y_nor + h_nor))

			// pdf页的百分比, 显示部分的百分比
			rect_pdf := types.NewRectangle(px, py_pdf, px+pw, py_pdf+ph)
			rect_page := types.NewRectangle(px, py, px+pw, py+ph)

			task := rumia.RenderTileTask{
				Page:       i,
				Dpi:        self.dpi,
				Rect:       rect_pdf,
				Rect_ratio: rect_page,
				Call_fn:    func(img *image.RGBA) { self.callback_apply_tile(i, self.dpi, rect_pdf, rect_page, img) },
			}
			tasks = append(tasks, &task)
		}
	}

	return tasks
}

// 得到当前视口的整页数,从中心页向上和向下螺旋排布
func (self *PDFViewer) get_page_list() []int {
	min_count := 0
	max_count := self.doc_page_num - 1
	_, y := self.from_screen_to_world(0, float64(self.viewer_area.Y)*0.5)
	center_page := int(y / (self.base_height + self.pad))
	p := max(min_count, min(max_count, center_page))

	pages := []int{p}
	b := math.Ceil(float64(self.viewer_area.Y)/((self.base_height+self.pad)*self.scale_tar)) + 1
	for i := range int(b * 0.5) {
		if np := p + (i + 1); np < self.doc_page_num-1 {
			_, nsy, _, _ := self.get_page_pos_on_viewer_tar(np)
			if nsy < float64(self.viewer_area.Y) {
				pages = append(pages, np)
			}
		}

		if lp := p - (i + 1); lp >= 0 {
			_, _, _, ley := self.get_page_pos_on_viewer_tar(lp)
			if ley >= 0 {
				pages = append(pages, lp)
			}
		}
	}
	return pages
}

// 计算当前视口内可见的页面列表,自上而下
func (self *PDFViewer) calc_page_list() []int {
	pages := []int{}
	if self.doc_page_num == 0 {
		return pages
	}

	_, lly := self.from_screen_to_world(0, 0)
	_, ury := self.from_screen_to_world(float64(self.viewer_area.X), float64(self.viewer_area.Y))

	var c1, c2 float64
	switch self.mode_layout {
	default:
		c1 = math.Floor(lly / (self.base_height + self.pad))
		c2 = math.Ceil(ury / (self.base_height + self.pad))
	case mode_layout_double:
		c1 = math.Floor(lly/(self.base_height+self.pad)) * 2
		c2 = math.Ceil(ury/(self.base_height+self.pad))*2 - 1
	case mode_layout_single_double:
		c1 = math.Floor(lly/(self.base_height+self.pad))*2 - 1
		c2 = math.Ceil(ury/(self.base_height+self.pad))*2 - 2
	}

	for i := c1; i <= c2; i++ {
		if i >= 0 && i < float64(self.doc_page_num) {
			pages = append(pages, int(i))
		}
	}

	return pages
}

// 计算一组页中的中间数
func (self *PDFViewer) calc_center_count(pages []int) int {
	if len(pages) == 0 {
		_, cy := self.from_screen_to_world(0.0, float64(self.viewer_area.Y))
		if cy <= 0 {
			return 0
		} else {
			return self.doc_page_num - 1
		}
	}

	mid_idx := len(pages) / 2
	_, viewer_cy := self.from_screen_to_world(0.0, float64(self.viewer_area.Y/2))
	minDist := math.MaxFloat64

	best_p := pages[mid_idx]

	for i := -3; i <= 3; i++ {
		idx := mid_idx + i
		if idx < 0 || idx >= len(pages) {
			continue
		}

		if page_num := pages[idx]; page_num >= 0 && page_num < self.doc_page_num {
			_, lly := self.calc_box_pos_global(page_num)
			ury := lly + self.base_height + self.pad

			if dist := math.Abs((lly+ury)*0.5 - viewer_cy); dist < minDist {
				minDist = dist
				best_p = page_num
			}
		}
	}

	return best_p
}

// 得到一页在视口坐标系中的坐标(根据目标值算)
func (self *PDFViewer) get_page_pos_on_viewer_tar(page int) (float64, float64, float64, float64) {
	img_w := float64(self.buffer_size[page].Dx())
	img_h := float64(self.buffer_size[page].Dy())
	ratio := self.base_height / img_h

	start_x := math.Max(0.0, self.offsetX_tar-img_w*ratio*self.scale_tar*0.5)
	start_y := math.Max(0.0, self.offsetY_tar+(self.base_height+self.pad)*float64(page)*self.scale_tar)

	end_x := math.Min(float64(self.viewer_area.X), self.offsetX_tar-img_w*ratio*self.scale_tar*0.5+img_w*ratio*self.scale_tar)
	end_y := math.Min(float64(self.viewer_area.Y), self.offsetY_tar+(self.base_height+self.pad)*float64(page+1)*self.scale_tar-self.pad*self.scale_tar)

	return start_x, start_y, end_x, end_y
}

// 得到一页在视口坐标系中的坐标(根据当前值算)
func (self *PDFViewer) get_page_pos_on_viewer_cur(page int) (float64, float64, float64, float64) {
	img_w := float64(self.buffer_size[page].Dx())
	img_h := float64(self.buffer_size[page].Dy())
	ratio := self.base_height / img_h

	start_x := math.Max(0.0, self.offsetX_cur-img_w*ratio*self.scale_cur*0.5)
	start_y := math.Max(0.0, self.offsetY_cur+(self.base_height+self.pad)*float64(page)*self.scale_cur)

	end_x := math.Min(float64(self.viewer_area.X), self.offsetX_cur-img_w*ratio*self.scale_cur*0.5+img_w*ratio*self.scale_cur)
	end_y := math.Min(float64(self.viewer_area.Y), self.offsetY_cur+(self.base_height+self.pad)*float64(page+1)*self.scale_cur-self.pad*self.scale_cur)

	return start_x, start_y, end_x, end_y
}

// 得到dpi
func (self *PDFViewer) get_dpi() float64 {
	scale := math.Round(self.scale_tar*10) / 10
	dpi := 72.0 * self.base_height * scale / float64(self.buffer_size[self.index_focus].Dy())
	dpi = math.Round(math.Min(1000, dpi))
	return dpi
}

// 渲染回调,为控件应用渲染图
func (self *PDFViewer) callback_apply_full(page int, dpi float64, mode_shader modeShader, img *image.RGBA) {
	go func() {
		if dpi != self.dpi {
			img.Pix = []uint8{}
			img = nil
			return
		}
		if info, exist := self.buffer_page[page]; exist {
			img := select_mode_shader(img, self.mode_shader)
			// img := tool.Shader_warm(img)
			info.img = img
			info.DPI = dpi
			info.mode_shader = mode_shader
			if info.Widget != nil {
				info.Widget.Src = paint.NewImageOp(img)
			} else {
				info.Widget = &widget.Image{
					Src: paint.NewImageOp(info.img),
					Fit: widget.Unscaled,
				}
				info.Widget.Src.Filter = paint.FilterNearest
			}
			self.bridge_ui.WindowRefresh()
		} else {
			log.Println("page info not exist")
		}
	}()
}

// 渲染回调,应用瓦片图
func (self *PDFViewer) callback_apply_tile(page int, dpi float64, rect *types.Rectangle, rect_ratio *types.Rectangle, img *image.RGBA) {
	go func() {
		if info, exist := self.buffer_page[page]; exist {
			img := select_mode_shader(img, self.mode_shader)
			// img := tool.Shader_binary(img)

			info.zoom_img = img
			info.zoom_dpi = dpi
			info.zoom_rect = rect
			info.zoom_ratio = rect_ratio

			imgop := paint.NewImageOp(info.zoom_img)
			imgop.Filter = paint.FilterNearest

			if info.zoom_widget != nil {
				info.zoom_widget.Src = paint.NewImageOp(info.zoom_img)
			} else {
				info.zoom_widget = &widget.Image{
					Src: paint.NewImageOp(info.zoom_img),
					Fit: widget.Unscaled,
				}
			}

			self.bridge_ui.WindowRefresh()
		} else {
			log.Println("not exist")
		}
	}()
}

// 得到指定页的包围盒尺寸
func (self *PDFViewer) get_box_size(page int) (float64, float64) {
	if page < 0 || page >= self.doc_page_num {
		return 0, 0
	}

	s, exist := self.buffer_box_size[page]
	if exist {
		return s[0], s[1]
	}
	return self.calc_box_size(page)
}

// 计算指定页的包围框尺寸
//   - 前提:尺寸缓冲区存在
func (self *PDFViewer) calc_box_size(page int) (float64, float64) {
	size, exist := self.buffer_size[page]
	if !exist {
		s, _ := self.bridge_trunk.Get_size(self.doc_id, page)
		size = &s
	}

	// 局部变换参数计算
	sin_a, cos_a := math.Abs(math.Sin(self.rotateArc_cur)), math.Abs(math.Cos(self.rotateArc_cur))
	pdf_w, pdf_h := float64(size.Dx()), float64(size.Dy())

	ele_w := pdf_w*cos_a + pdf_h*sin_a
	ele_h := pdf_w*sin_a + pdf_h*cos_a
	local_scale_coe := self.base_height / ele_h

	box_w := ele_w * local_scale_coe
	box_h := ele_h * local_scale_coe

	return box_w, box_h
}

// 得到指定页在控件坐标系内的偏移
func (self *PDFViewer) calc_box_pos_global(page int) (float64, float64) {
	w := self.buffer_box_size[page][0]
	h := self.buffer_box_size[page][1]
	var x, y float64 = 0.0, 0.0

	// 是否从右往左
	isRTL := self.mode_order == mode_order_right

	switch self.mode_layout {

	default: // 单页排版
		x = -w * 0.5
		y = (h + self.pad) * float64(page)

	case mode_layout_double: // 双页排版
		l := page % 2
		if l == 0 {
			// 左侧页
			w1 := self.buffer_box_size[page][0]
			w2 := Ternary(page+1 >= self.doc_page_num, 0, self.buffer_box_size[page+1][0])

			if !isRTL {
				x = -(w1 + w2) * 0.5
			} else {
				x = -(w1+w2)*0.5 + w2
			}

			y = (h + self.pad) * float64(page/2)

		} else {
			// 右侧页
			w1 := Ternary(page-1 < 0, 0, self.buffer_box_size[page-1][0])
			w2 := self.buffer_box_size[page][0]

			if !isRTL {
				x = -(w1+w2)*0.5 + w1
			} else {
				x = -(w1 + w2) * 0.5
			}

			y = (h + self.pad) * float64(page/2)
		}

	case mode_layout_single_double: // 单页+双页 (通常第0页单独居中)
		if page == 0 {
			x = -w * 0.5
			y = (h + self.pad) * float64(page)
			break
		}

		l := (page - 1) % 2
		if l == 0 {
			// 右侧页
			w1 := self.buffer_box_size[page][0]
			w2 := Ternary(page+1 < self.doc_page_num, self.buffer_box_size[page+1][0], 0)

			if !isRTL {
				x = -(w1 + w2) * 0.5
			} else {
				x = -(w1+w2)*0.5 + w2
			}
			y = (h + self.pad) * float64((page-1)/2+1)

		} else {
			// 左侧页
			w1 := self.buffer_box_size[page-1][0]
			w2 := self.buffer_box_size[page][0]

			if !isRTL {
				x = -(w1+w2)*0.5 + w1
			} else {
				x = -(w1 + w2) * 0.5
			}
			y = (h + self.pad) * float64((page-1)/2+1)
		}
	}

	self.buffer_box_pos[page] = [2]float64{x, y}
	return x, y
}

// 缩放变换

// 其他 -----------------------------------------------------------------------------------------------------

// 计算某一页在当前排版下所在行的尺寸
func (self *PDFViewer) calc_row_page_size(page int) (float64, float64) {
	switch self.mode_layout {
	default:
		return self.get_box_size(page)
	case mode_layout_double:
		if r := page % 2; r == 0 {
			w1, h1 := self.get_box_size(page)
			w2, _ := self.get_box_size(page + 1)
			return w1 + w2, h1
		} else {
			w1, h1 := self.get_box_size(page)
			w2, _ := self.get_box_size(page - 1)
			return w1 + w2, h1
		}
	case mode_layout_single_double:
		if page == 0 {
			return self.get_box_size(0)
		}
		if r := page % 2; r == 0 {
			w1, h1 := self.get_box_size(page)
			w2, _ := self.get_box_size(page - 1)
			return w1 + w2, h1
		} else {
			w1, h1 := self.get_box_size(page)
			w2, _ := self.get_box_size(page + 1)
			return w1 + w2, h1
		}
	}

}

// 跳转到某页并居中,返回Scale,OffsetX,OffsetY
func (self *PDFViewer) apply_center(page int) (float64, float64, float64) {
	// // 得到目标行的尺寸
	// row_w, row_h := self.calc_row_page_size(page)
	// area_w, area_h := float64(self.viewer_area.X), float64(self.viewer_area.Y)
	//
	// switch self.mode_layout {
	// case mode_layout_double:
	//
	//	self.scale_tar = math.Round(math.Min(area_w/row_w, area_h/row_h)*100) / 100
	//	self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
	//	self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page/2)*(self.base_height+self.pad)*self.scale_tar - 0.5)
	//
	// case mode_layout_single_double:
	//
	//	self.scale_tar = math.Round(math.Min(area_w/row_w, area_h/row_h)*100) / 100
	//	self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
	//	self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page/2+1)*(self.base_height+self.pad)*self.scale_tar - 0.5)
	//
	// case mode_layout_fit:
	//
	//	self.scale_tar = math.Max(area_w/row_w, area_h/row_h)
	//	self.scale_tar = math.Round(self.scale_tar*100) / 100
	//	self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
	//	self.offsetY_tar = math.Round(-float64(page)*(self.base_height+self.pad)*self.scale_tar - 0.5)
	//
	// default:
	//		self.scale_tar = math.Min(area_w/row_w, area_h/row_h)
	//		self.scale_tar = math.Round(self.scale_tar*100) / 100
	//		self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
	//		self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page)*(self.base_height+self.pad)*self.scale_tar - 0.5)
	//	}
	//
	// return self.scale_tar, self.offsetX_tar, self.offsetY_tar

	scale := self.apply_center_scale(page)
	x, y := self.apply_center_pos(page)
	return scale, x, y
}

func (self *PDFViewer) apply_center_scale(page int) float64 {
	row_w, row_h := self.calc_row_page_size(page)
	area_w, area_h := float64(self.viewer_area.X), float64(self.viewer_area.Y)

	self.scale_tar = math.Round(math.Min(area_w/row_w, area_h/row_h)*100) / 100
	return self.scale_tar
}

func (self *PDFViewer) apply_center_pos(page int) (float64, float64) {
	log.Println(page)
	row_w, _ := self.calc_row_page_size(page)
	area_w, area_h := float64(self.viewer_area.X), float64(self.viewer_area.Y)

	switch self.mode_layout {
	case mode_layout_double:
		self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
		self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page/2)*(self.base_height+self.pad)*self.scale_tar - 0.5)

	case mode_layout_single_double:
		self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
		self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page/2+1)*(self.base_height+self.pad)*self.scale_tar - 0.5)

	// case mode_layout_fit:
	// 	self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
	// 	self.offsetY_tar = math.Round(-float64(page)*(self.base_height+self.pad)*self.scale_tar - 0.5)

	default:
		self.offsetX_tar = math.Round((area_w-row_w*self.scale_tar)*.5 + row_w*self.scale_tar*0.5)
		self.offsetY_tar = math.Round((area_h-self.base_height*self.scale_tar)*.5 - float64(page)*(self.base_height+self.pad)*self.scale_tar - 0.5)
	}
	return self.offsetX_tar, self.offsetY_tar
}

// 世界坐标转屏幕坐标
func (self *PDFViewer) from_world_to_screen(x float64, y float64) (float64, float64) {
	sx := x*self.scale_tar + self.offsetX_tar
	sy := y*self.scale_tar + self.offsetY_tar
	return sx, sy
}

// 屏幕坐标转世界坐标
func (self *PDFViewer) from_screen_to_world(x, y float64) (float64, float64) {
	wx := (x - self.offsetX_tar) / self.scale_tar
	wy := (y - self.offsetY_tar) / self.scale_tar
	return wx, wy
}

// 弧度转状态
func (self *PDFViewer) to_arc_state(arc float64) int {
	rad := math.Mod(arc, 2*math.Pi)
	if rad < 0 {
		rad += 2 * math.Pi
	}
	step := int(math.Round(rad/(math.Pi*0.5))) % 4
	return step
}
