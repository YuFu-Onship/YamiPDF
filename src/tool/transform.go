package tool

import "math"

// A点绕B点旋转一定角度
func Rotate_point(x, y, ox, oy, arc float64) (float64, float64) {
	sin_a, cos_a := math.Sin(arc), math.Cos(arc)
	dx := x - ox
	dy := y - oy

	rx := dx*cos_a - dy*sin_a
	ry := dx*sin_a + dy*cos_a
	return rx + ox, ry + oy
}
