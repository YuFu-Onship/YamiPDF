package ui

import (
	"image"
	"image/color"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

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

// 绘制坐标系
func draw_coordinate_system(gtx C) {
	const axisLength = 1e6
	const halfWidth = 2

	paint.FillShape(gtx.Ops, COLOR_RED, clip.Rect{
		Min: image.Pt(-axisLength, -halfWidth),
		Max: image.Pt(axisLength, halfWidth),
	}.Op())

	paint.FillShape(gtx.Ops, COLOR_GREEN, clip.Rect{
		Min: image.Pt(-halfWidth, -axisLength),
		Max: image.Pt(halfWidth, axisLength),
	}.Op())
}
