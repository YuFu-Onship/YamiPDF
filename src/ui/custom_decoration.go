package ui

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const (
	iconSize       = unit.Dp(20)
	titleBarHeight = unit.Dp(36)
	actionBtnW     = unit.Dp(40)
	winIconSize    = unit.Dp(20)
	winIconMargin  = unit.Dp(4)
	winIconStroke  = unit.Dp(2)
)

type CustomDeco struct {
	Deco  *widget.Decorations
	Style *Style
	min   *widget.Clickable
	max   *widget.Clickable
	close *widget.Clickable

	text string

	is_show_arrow  bool              // 是否显示返回按钮
	fn_back        func()            // 点击返回按钮后的回调
	clickable_back *widget.Clickable // 返回按钮
}

func NewCustomDeco(style *Style) *CustomDeco {
	deco := widget.Decorations{}
	self := CustomDeco{
		Deco:  &deco,
		Style: style,
		min:   deco.Clickable(system.ActionMinimize),
		max:   deco.Clickable(system.ActionMaximize),
		close: deco.Clickable(system.ActionClose),

		text: "",

		is_show_arrow:  true,
		clickable_back: &widget.Clickable{},
		fn_back:        func() {},
	}
	return &self
}

func (self *CustomDeco) Layout(gtx C) D {
	// paint.FillShape(gtx.Ops, self.Style.Palette_deco.Deco_bg, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, int(titleBarHeight))}.Op())
	// paint.FillShape(gtx.Ops, self.Style.Palette.Bg_1, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, int(titleBarHeight))}.Op())
	gtx.Constraints.Max.Y = gtx.Dp(titleBarHeight)
	gtx.Constraints.Min.Y = 0

	if self.clickable_back.Clicked(gtx) {
		self.fn_back()
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if self.is_show_arrow {
				return self.actionButton(
					gtx,
					self.clickable_back,
					self.Style.Palette_deco.Decoicon_fg_idle,
					self.Style.Palette_deco.Decoicon_fg_hover,
					self.Style.Palette_deco.Decoicon_fg_press,
					self.Style.Palette_deco.Decoicon_bg_idle,
					self.Style.Palette_deco.Decoicon_bg_hover,
					self.Style.Palette_deco.Decoicon_bg_press,
					IconBackArrow)
			} else {
				return D{}
			}
		}),

		layout.Rigid(Spacer(gtx.Dp(5), gtx.Dp(0))),

		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return self.Deco.LayoutMove(gtx, func(gtx layout.Context) layout.Dimensions {
				if self.text == "" {
					return D{Size: gtx.Constraints.Max}
				}

				macro := op.Record(gtx.Ops)
				l := material.Label(self.Style.Theme, unit.Sp(16), self.text)
				l.Color = self.Style.Palette_deco.Decoicon_fg_idle
				l.MaxLines = 1
				t := l.Layout(gtx)
				call := macro.Stop()

				textHeight := t.Size.Y
				textWidth := t.Size.X

				offsetX := (gtx.Constraints.Max.X - textWidth) / 2
				offsetY := (gtx.Constraints.Max.Y - textHeight) / 2

				offsetX = max(offsetX, 0)
				offsetY = max(offsetY, 0)

				defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()
				call.Add(gtx.Ops)

				return D{Size: gtx.Constraints.Max}
			})
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return self.actionButton(
				gtx,
				self.min,
				self.Style.Palette_deco.Decoicon_fg_idle,
				self.Style.Palette_deco.Decoicon_fg_hover,
				self.Style.Palette_deco.Decoicon_fg_press,
				self.Style.Palette_deco.Decoicon_bg_idle,
				self.Style.Palette_deco.Decoicon_bg_hover,
				self.Style.Palette_deco.Decoicon_bg_press,
				IconMinimizeWindow)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if self.Deco.Maximized {
				return self.actionButton(
					gtx,
					self.max,
					self.Style.Palette_deco.Decoicon_fg_idle,
					self.Style.Palette_deco.Decoicon_fg_hover,
					self.Style.Palette_deco.Decoicon_fg_press,
					self.Style.Palette_deco.Decoicon_bg_idle,
					self.Style.Palette_deco.Decoicon_bg_hover,
					self.Style.Palette_deco.Decoicon_bg_press,
					IconMaximizedWindow)
			} else {
				return self.actionButton(
					gtx,
					self.max,
					self.Style.Palette_deco.Decoicon_fg_idle,
					self.Style.Palette_deco.Decoicon_fg_hover,
					self.Style.Palette_deco.Decoicon_fg_press,
					self.Style.Palette_deco.Decoicon_bg_idle,
					self.Style.Palette_deco.Decoicon_bg_hover,
					self.Style.Palette_deco.Decoicon_bg_press,
					IconMaximizeWindow)
			}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return self.actionButton(
				gtx,
				self.close,
				self.Style.Palette_deco.Decoicon_fg_idle_close,
				self.Style.Palette_deco.Decoicon_fg_hover_close,
				self.Style.Palette_deco.Decoicon_fg_press_close,
				self.Style.Palette_deco.Decoicon_bg_idle_close,
				self.Style.Palette_deco.Decoicon_bg_hover_close,
				self.Style.Palette_deco.Decoicon_bg_press_close,
				IconCloseWindow)
		}),
	)
}

func (self *CustomDeco) Actions(gtx C, win *app.Window) {
	for {
		action := self.Deco.Update(gtx)
		if action == 0 {
			break
		}
		win.Perform(action)
	}
}

func (self *CustomDeco) actionButton(
	gtx C,
	cl *widget.Clickable,
	fg_idle color.NRGBA,
	fg_hover color.NRGBA,
	fg_press color.NRGBA,
	bg_idle color.NRGBA,
	bg_hover color.NRGBA,
	bg_press color.NRGBA,
	draw func(C) D,
) D {

	btnWPx := gtx.Dp(actionBtnW)
	btnHPx := gtx.Dp(titleBarHeight)
	iconPx := gtx.Dp(iconSize)
	size := image.Pt(btnWPx, btnHPx)

	var bg color.NRGBA
	var fg color.NRGBA

	switch {
	case cl.Pressed():
		bg = bg_press
		fg = fg_press
	case cl.Hovered():
		bg = bg_hover
		fg = fg_hover
	default:
		bg = bg_idle
		fg = fg_idle
	}

	// 	btn := material.Button(self.Style.Theme, cl, "")
	// 	btn.Background.A = 0
	// 	btn.CornerRadius = 0
	// 	btn.Inset = layout.Inset{}
	// 	if btn.Button.Hovered() {
	// 		btn.Background = bg
	// 	}

	// 波纹动画
	// 	return layout.Stack{}.Layout(gtx,
	// 		layout.Stacked(func(gtx C) D {
	// 			gtx.Constraints.Min = image.Pt(int(actionBtnW), int(titleBarHeight))
	// 			return btn.Layout(gtx)
	// 		}),
	// 		layout.Stacked(func(gtx C) D {
	// 			paint.ColorOp{Color: fg}.Add(gtx.Ops)
	// 			off := image.Pt((size.X-iconPx)/2, (size.Y-iconPx)/2)
	// 			offset := op.Offset(off).Push(gtx.Ops)
	// 			draw(gtx)
	// 			offset.Pop()
	//
	// 			return D{Size: size}
	// 		}),
	// 	)

	return cl.Layout(gtx, func(gtx C) D {
		paint.FillShape(gtx.Ops, bg, clip.Rect{Max: size}.Op())
		paint.ColorOp{Color: fg}.Add(gtx.Ops)
		off := image.Pt((size.X-iconPx)/2, (size.Y-iconPx)/2)
		offset := op.Offset(off).Push(gtx.Ops)
		draw(gtx)
		offset.Pop()

		return D{Size: size}
	})
}

// 设置显示文本
func (self *CustomDeco) API_set_text(text string) {
	self.text = text
}

// 启用返回按钮
//   - true 显示返回按钮
//   - false 不显示返回按钮
func (self *CustomDeco) API_set_back(state bool) {

}

// 标题栏图标 ---------------------------------------------------------------------

// 最小化图标
func IconMinimizeWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Dp(winIconMargin))
	width := float32(gtx.Dp(winIconStroke))
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f32.Point{X: margin, Y: size32 - margin})
	p.LineTo(f32.Point{X: size32 - 2*margin, Y: size32 - margin})
	st := clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// 最大化图标
func IconMaximizeWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(winIconSize)
	margin := gtx.Dp(winIconMargin)
	width := gtx.Dp(winIconStroke)
	r := clip.RRect{
		Rect: image.Rect(margin, margin, size-margin, size-margin),
	}
	st := clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: float32(width),
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	r.Rect.Max = image.Pt(size-margin, 2*margin)
	st = clip.Outline{
		Path: r.Path(gtx.Ops),
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// 非最大化图标
func IconMaximizedWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(winIconSize)
	margin := gtx.Dp(winIconMargin)
	width := gtx.Dp(winIconStroke)
	r := clip.RRect{
		Rect: image.Rect(margin, margin, size-2*margin, size-2*margin),
	}
	st := clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: float32(width),
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	r = clip.RRect{
		Rect: image.Rect(2*margin, 2*margin, size-margin, size-margin),
	}
	st = clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: float32(width),
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// 关闭图标
func IconCloseWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Dp(winIconMargin))
	width := float32(gtx.Dp(winIconStroke))
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f32.Point{X: margin, Y: margin})
	p.LineTo(f32.Point{X: size32 - margin, Y: size32 - margin})
	p.MoveTo(f32.Point{X: size32 - margin, Y: margin})
	p.LineTo(f32.Point{X: margin, Y: size32 - margin})
	st := clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// 返回图标
func IconBackArrow(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Dp(winIconMargin))
	width := float32(gtx.Dp(winIconStroke))

	// 计算中心与两端的位置
	centerY := size32 / 2.0
	leftX := margin
	rightX := size32 - margin

	// 计算箭头头部折角线的长度（这里取与箭身相平衡的高度）
	arrowHeadYOffset := (size32 - 2*margin) / 2.0

	var p clip.Path
	p.Begin(gtx.Ops)

	// 1. 绘制水平主线（从右向左）
	p.MoveTo(f32.Point{X: rightX, Y: centerY})
	p.LineTo(f32.Point{X: leftX, Y: centerY})

	// 2. 绘制左侧箭头的上翼
	p.LineTo(f32.Point{X: leftX + arrowHeadYOffset, Y: centerY - arrowHeadYOffset})

	// 3. 绘制左侧箭头的下翼
	p.MoveTo(f32.Point{X: leftX, Y: centerY})
	p.LineTo(f32.Point{X: leftX + arrowHeadYOffset, Y: centerY + arrowHeadYOffset})

	st := clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()

	return layout.Dimensions{Size: image.Pt(size, size)}
}
