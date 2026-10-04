package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DropDown struct {
	Header W
	Panel  W

	Text      string // 经典边框文本,需要显示
	Is_fill   bool   // 是否填充绘制背景
	Is_border bool   // 是否绘制边框
	Is_expand bool   // 是否展开

	Bg color.NRGBA // 背景颜色
	Fg color.NRGBA // 边框颜色

	height_tar float64
	height_cur float64

	clickable_header *widget.Clickable
}

// 下拉栏布局, 布局中只能包含绘制逻辑
//   - header 顶部布局
//   - panel 下拉内容布局
func (self DropDown) Layout(gtx C, style *Style, clickable_header *widget.Clickable, header W, panel W) D {
	self.height_cur = 10
	self.clickable_header = clickable_header

	if self.clickable_header.Clicked(gtx) {
		self.Is_expand = !self.Is_expand
	}

	btn_header := material.Button(style.Theme, self.clickable_header, "")
	btn_header.Background.A = 0

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					if self.Is_fill {
						draw_rectangle(gtx, gtx.Constraints.Min, COLOR_GREEN, gtx.Dp(4))
					}
					draw_rectangle_line(gtx, gtx.Constraints.Min, COLOR_GREEN, gtx.Dp(4), float32(gtx.Dp(1)))
					return D{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(func(gtx C) D { return Border{Text: "fdsfs"}.Layout(gtx, style, header) }),
				layout.Expanded(btn_header.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {

			macro := op.Record(gtx.Ops)
			p := panel(gtx)
			record := macro.Stop()

			self.height_tar = Ternary(self.Is_expand, float64(p.Size.Y), 0)

			clip := clip.RRect{Rect: image.Rect(0, 0, p.Size.X, int(100))}.Push(gtx.Ops)

			draw_rectangle(gtx, p.Size, COLOR_GREEN, gtx.Dp(8))
			draw_rectangle_line(gtx, p.Size, COLOR_ORANGE, gtx.Dp(8), float32(gtx.Dp(1)))
			record.Add(gtx.Ops)

			clip.Pop()
			return D{Size: p.Size}
		}),
	)
}
