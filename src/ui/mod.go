package ui

import (
	"image/color"
	"log"
	"time"

	"gioui.org/layout"
)

type C = layout.Context
type D = layout.Dimensions
type W = layout.Widget

var PalettePico8 = color.Palette{
	color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // 0: Black
	color.RGBA{R: 0x1D, G: 0x2B, B: 0x53, A: 0xFF}, // 1: Dark Blue
	color.RGBA{R: 0x7E, G: 0x25, B: 0x53, A: 0xFF}, // 2: Dark Purple
	color.RGBA{R: 0x00, G: 0x87, B: 0x51, A: 0xFF}, // 3: Dark Green
	color.RGBA{R: 0xAB, G: 0x52, B: 0x36, A: 0xFF}, // 4: Brown
	color.RGBA{R: 0x5F, G: 0x57, B: 0x4F, A: 0xFF}, // 5: Dark Gray
	color.RGBA{R: 0xC2, G: 0xC3, B: 0xC7, A: 0xFF}, // 6: Light Gray
	color.RGBA{R: 0xFF, G: 0xF1, B: 0xE8, A: 0xFF}, // 7: White
	color.RGBA{R: 0xFF, G: 0x00, B: 0x4D, A: 0xFF}, // 8: Red
	color.RGBA{R: 0xFF, G: 0xA3, B: 0x00, A: 0xFF}, // 9: Orange
	color.RGBA{R: 0xFF, G: 0xEC, B: 0x27, A: 0xFF}, // 10: Yellow
	color.RGBA{R: 0x00, G: 0xE4, B: 0x36, A: 0xFF}, // 11: Green
	color.RGBA{R: 0x29, G: 0xAD, B: 0xFF, A: 0xFF}, // 12: Blue
	color.RGBA{R: 0x83, G: 0x76, B: 0x9C, A: 0xFF}, // 13: Lavender
	color.RGBA{R: 0xFF, G: 0x77, B: 0xA8, A: 0xFF}, // 14: Pink
	color.RGBA{R: 0xFF, G: 0xCC, B: 0xAA, A: 0xFF}, // 15: Peach
}

func ComputeTime(fn func()) {
	start := time.Now()
	fn()
	log.Printf("耗时:%v", time.Since(start))
}
