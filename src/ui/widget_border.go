package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Border struct {
	Text   string
	Fg     color.NRGBA
	Bg     color.NRGBA
	Width  unit.Dp
	Radius unit.Dp

	Is_inset_zero bool // 内边距设置为0
}

func (self *Border) Default(style *Style) {
	self.Fg = Ternary(self.Fg == (color.NRGBA{}), style.Palette.Fg_1, self.Fg)
	self.Bg = Ternary(self.Bg == (color.NRGBA{}), style.Palette.Bg_1, self.Bg)
	self.Radius = Ternary(self.Radius == 0, unit.Dp(4), self.Radius)
	self.Width = Ternary(self.Width == 0, unit.Dp(1), self.Width)
}

func (self Border) Layout(gtx C, style *Style, context W) D {
	self.Default(style)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// return widget.Border{Width: self.Width, CornerRadius: self.Radius, Color: self.Fg}.Layout(gtx, func(gtx C) D {
			draw_rectangle_line(gtx, gtx.Constraints.Min, self.Fg, int(self.Radius), float32(self.Width))
			return D{Size: gtx.Constraints.Min}
			// })
		}),

		layout.Expanded(func(gtx C) D {
			trans := op.Offset(image.Pt(gtx.Dp(12), -gtx.Dp(10))).Push(gtx.Ops)
			defer trans.Pop()

			return layout.Stack{}.Layout(gtx,
				// 背景矩形
				layout.Expanded(func(gtx C) D {
					if self.Text == "" {
						return D{}
					}
					pt := gtx.Constraints.Min
					draw_rectangle(gtx, pt, self.Bg, 0)
					return D{}
				}),
				// 前景文字
				layout.Stacked(func(gtx C) D {
					if self.Text == "" {
						return D{}
					}
					return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
						l := material.Label(style.Theme, unit.Sp(16), self.Text)
						l.Color = self.Fg
						return l.Layout(gtx)
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
			}.Layout(gtx, context)
		}),
	)
}
