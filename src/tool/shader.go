package tool

import (
	"image"
	"image/color"
)

// 色板 ----------------------------------------------------------------

// GameBoy 经典四色绿调色板
var Palette_gameboy = [4]color.RGBA{
	{R: 155, G: 188, B: 15, A: 255}, // 0: 最亮 (浅黄绿)
	{R: 139, G: 172, B: 15, A: 255}, // 1: 次亮 (黄绿)
	{R: 48, G: 98, B: 48, A: 255},   // 2: 次暗 (深绿)
	{R: 15, G: 56, B: 15, A: 255},   // 3: 最暗 (极深墨绿)
}

// 图像shader ----------------------------------------------------------

// gameboy风格
func Shader_gameboy(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			gray := 0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8)
			var level int
			switch {
			case gray >= 192:
				level = 0
			case gray >= 128:
				level = 1
			case gray >= 64:
				level = 2
			default:
				level = 3
			}

			c := Palette_gameboy[level]
			c.A = uint8(a >> 8)
			src.SetRGBA(x, y, c)
		}
	}
	return src
}

// 灰度风格
func Shader_gray(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			c := uint8(0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
			src.SetRGBA(x, y, color.RGBA{R: c, G: c, B: c, A: uint8(a >> 8)})
		}
	}
	return src
}

// 二值化风格
func Shader_binary(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			gray := uint8(0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
			var c uint8 = 0
			if gray > 128 {
				c = 255
			}

			src.SetRGBA(x, y, color.RGBA{R: c, G: c, B: c, A: uint8(a >> 8)})
		}
	}
	return src
}

// 暖色风格
func Shader_warm(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			rNew := r8 * 1.2  // 增加红色比例
			gNew := g8 * 1.05 // 略微增加绿色比例以叠加偏黄效果
			bNew := b8 * 0.85 // 降低蓝色比例

			if rNew > 255 {
				rNew = 255
			}
			if gNew > 255 {
				gNew = 255
			}

			src.SetRGBA(x, y, color.RGBA{
				R: uint8(rNew),
				G: uint8(gNew),
				B: uint8(bNew),
				A: uint8(a >> 8),
			})
		}
	}
	return src
}
