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
}

func NewCustomDeco(style *Style) *CustomDeco {
	deco := widget.Decorations{}
	self := CustomDeco{
		Deco:  &deco,
		Style: style,
		min:   deco.Clickable(system.ActionMinimize),
		max:   deco.Clickable(system.ActionMaximize),
		close: deco.Clickable(system.ActionClose),
	}
	return &self
}

func (self *CustomDeco) Layout(gtx C) D {
	paint.FillShape(gtx.Ops, self.Style.Palette.Deco_bg, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, int(titleBarHeight))}.Op())
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return self.Deco.LayoutMove(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(material.Label(self.Style.Theme, unit.Sp(30), "").Layout),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return self.actionButton(
				gtx,
				self.min,
				self.Style.Palette.Decoicon_fg_idle,
				self.Style.Palette.Decoicon_fg_hover,
				self.Style.Palette.Decoicon_fg_press,
				self.Style.Palette.Decoicon_bg_idle,
				self.Style.Palette.Decoicon_bg_hover,
				self.Style.Palette.Decoicon_bg_press,
				IconMinimizeWindow)
		}),
		layout.Rigid(func(gtx C) D {
			if self.Deco.Maximized {
				return self.actionButton(
					gtx,
					self.max,
					self.Style.Palette.Decoicon_fg_idle,
					self.Style.Palette.Decoicon_fg_hover,
					self.Style.Palette.Decoicon_fg_press,
					self.Style.Palette.Decoicon_bg_idle,
					self.Style.Palette.Decoicon_bg_hover,
					self.Style.Palette.Decoicon_bg_press,
					IconMaximizedWindow)
			} else {
				return self.actionButton(
					gtx,
					self.max,
					self.Style.Palette.Decoicon_fg_idle,
					self.Style.Palette.Decoicon_fg_hover,
					self.Style.Palette.Decoicon_fg_press,
					self.Style.Palette.Decoicon_bg_idle,
					self.Style.Palette.Decoicon_bg_hover,
					self.Style.Palette.Decoicon_bg_press,
					IconMaximizeWindow)
			}
		}),
		layout.Rigid(func(gtx C) D {
			return self.actionButton(
				gtx,
				self.close,
				self.Style.Palette.Decoicon_fg_idle_close,
				self.Style.Palette.Decoicon_fg_hover_close,
				self.Style.Palette.Decoicon_fg_press_close,
				self.Style.Palette.Decoicon_bg_idle_close,
				self.Style.Palette.Decoicon_bg_hover_close,
				self.Style.Palette.Decoicon_bg_press_close,
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
