package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// 横向占位
func FlexerX() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Max.Y = 0
		return D{Size: gtx.Constraints.Max}
	}
}

// 垂直占位
func FlexerY() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Max.X = 0
		return D{Size: gtx.Constraints.Max}
	}
}

// 占位具体空间
func Spacer(width, height int) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: width, Y: height}
		return layout.Dimensions{Size: size}
	}
}

// 绘制矩形
func draw_rectangle(gtx C, area image.Point, c color.NRGBA, radius int) {
	if radius == 0 {
		paint.FillShape(gtx.Ops, c, clip.Rect{Max: area}.Op())
	} else {
		paint.FillShape(gtx.Ops, c, clip.RRect{Rect: image.Rectangle{Max: area},
			SE: radius,
			SW: radius,
			NE: radius,
			NW: radius,
		}.Op(gtx.Ops))
	}
}

// 绘制矩形边框
func draw_rectangle_line(gtx C, area image.Point, c color.NRGBA, radius int, width float32) {
	if radius == 0 {
		paint.FillShape(gtx.Ops, c, clip.Stroke{Path: clip.Rect{Max: area}.Path(), Width: width}.Op())
	} else {
		paint.FillShape(gtx.Ops, c, clip.Stroke{Path: clip.RRect{Rect: image.Rectangle{Max: area},
			SE: radius,
			SW: radius,
			NE: radius,
			NW: radius,
		}.Path(gtx.Ops), Width: width}.Op())
	}

}

// 绘制任意位置矩形
func draw_rectangle_2(gtx C, x int, y int, w int, h int, c color.NRGBA, radius int) {
	defer op.Offset(image.Pt(x, y)).Push(gtx.Ops).Pop()
	if radius == 0 {
		paint.FillShape(gtx.Ops, c, clip.Rect{Max: image.Pt(w, h)}.Op())
	} else {
		paint.FillShape(gtx.Ops, c, clip.RRect{Rect: image.Rectangle{Max: image.Pt(w, h)},
			SE: radius,
			SW: radius,
			NE: radius,
			NW: radius,
		}.Op(gtx.Ops))
	}
}

// 绘制任意位置矩形
func draw_rectangle_line_2(gtx C, x int, y int, w int, h int, c color.NRGBA, radius int, width float32) {
	offset := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	if radius == 0 {
		paint.FillShape(gtx.Ops, c, clip.Stroke{Path: clip.Rect{Max: image.Pt(w, h)}.Path(), Width: width}.Op())
	} else {
		paint.FillShape(gtx.Ops, c, clip.Stroke{Path: clip.RRect{Rect: image.Rectangle{Max: image.Pt(w, h)},
			SE: radius,
			SW: radius,
			NE: radius,
			NW: radius,
		}.Path(gtx.Ops), Width: width}.Op())
	}
	offset.Pop()
}

// 绘制坐标系
func draw_coordinate_system(gtx C) {
	var axisLength int = 1e6
	var halfWidth int = 2

	paint.FillShape(gtx.Ops, COLOR_RED, clip.Rect{
		Min: image.Pt(-axisLength, -halfWidth),
		Max: image.Pt(axisLength, halfWidth),
	}.Op())

	paint.FillShape(gtx.Ops, COLOR_GREEN, clip.Rect{
		Min: image.Pt(-halfWidth, -axisLength),
		Max: image.Pt(halfWidth, axisLength),
	}.Op())
}
