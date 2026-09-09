package ui

import (
	"image/color"

	"gioui.org/widget/material"
)

// PICO8 色板
var (
	COLOR_BLACK       = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // 0: 黑色
	COLOR_DARK_BLUE   = color.NRGBA{R: 29, G: 43, B: 83, A: 255}    // 1: 深蓝
	COLOR_DARK_PURPLE = color.NRGBA{R: 126, G: 37, B: 83, A: 255}   // 2: 深紫
	COLOR_DARK_GREEN  = color.NRGBA{R: 0, G: 135, B: 81, A: 255}    // 3: 深绿
	COLOR_BROWN       = color.NRGBA{R: 171, G: 82, B: 54, A: 255}   // 4: 棕色
	COLOR_DARK_GRAY   = color.NRGBA{R: 95, G: 87, B: 79, A: 255}    // 5: 深灰
	COLOR_LIGHT_GRAY  = color.NRGBA{R: 194, G: 195, B: 199, A: 255} // 6: 浅灰
	COLOR_WHITE       = color.NRGBA{R: 255, G: 241, B: 232, A: 255} // 7: 白色
	COLOR_RED         = color.NRGBA{R: 255, G: 0, B: 77, A: 255}    // 8: 红色
	COLOR_ORANGE      = color.NRGBA{R: 255, G: 163, B: 0, A: 255}   // 9: 橙色
	COLOR_YELLOW      = color.NRGBA{R: 255, G: 236, B: 39, A: 255}  // 10: 黄色
	COLOR_GREEN       = color.NRGBA{R: 0, G: 228, B: 54, A: 255}    // 11: 绿色
	COLOR_BLUE        = color.NRGBA{R: 41, G: 173, B: 255, A: 255}  // 12: 蓝色
	COLOR_INDIGO      = color.NRGBA{R: 131, G: 118, B: 156, A: 255} // 13: 靛蓝/灰紫
	COLOR_PINK        = color.NRGBA{R: 255, G: 119, B: 168, A: 255} // 14: 粉红
	COLOR_PEACH       = color.NRGBA{R: 255, G: 204, B: 170, A: 255} // 15: 桃红/肤色
)

type ColorPalette struct {
	Fg         color.NRGBA
	BG         color.NRGBA
	ContrastFg color.NRGBA
	ContrastBg color.NRGBA

	Deco_fg                 color.NRGBA
	Deco_bg                 color.NRGBA
	Decoicon_fg_idle        color.NRGBA
	Decoicon_fg_hover       color.NRGBA
	Decoicon_fg_press       color.NRGBA
	Decoicon_bg_idle        color.NRGBA
	Decoicon_bg_hover       color.NRGBA
	Decoicon_bg_press       color.NRGBA
	Decoicon_fg_idle_close  color.NRGBA
	Decoicon_fg_hover_close color.NRGBA
	Decoicon_fg_press_close color.NRGBA
	Decoicon_bg_idle_close  color.NRGBA
	Decoicon_bg_hover_close color.NRGBA
	Decoicon_bg_press_close color.NRGBA
}

type StylePalette struct {
	Dark  ColorPalette
	Light ColorPalette
}

func colorNRGBA(r, g, b, a uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

var WindowsBlue StylePalette = StylePalette{
	Light: ColorPalette{
		Fg:         color.NRGBA{R: 26, G: 26, B: 26, A: 255},
		BG:         color.NRGBA{R: 243, G: 243, B: 243, A: 255},
		ContrastFg: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		ContrastBg: color.NRGBA{R: 0, G: 95, B: 184, A: 255},

		Deco_fg:                 colorNRGBA(243, 243, 243, 255),
		Deco_bg:                 colorNRGBA(0, 95, 184, 255),
		Decoicon_fg_idle:        color.NRGBA{R: 209, G: 209, B: 209, A: 255},
		Decoicon_fg_hover:       color.NRGBA{R: 209, G: 209, B: 209, A: 255},
		Decoicon_fg_press:       color.NRGBA{R: 209, G: 209, B: 209, A: 200},
		Decoicon_bg_idle:        color.NRGBA{R: 0, G: 0, B: 0, A: 0},
		Decoicon_bg_hover:       color.NRGBA{R: 0, G: 0, B: 0, A: 60},
		Decoicon_bg_press:       color.NRGBA{R: 0, G: 0, B: 0, A: 80},
		Decoicon_fg_idle_close:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
		Decoicon_fg_hover_close: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Decoicon_fg_press_close: color.NRGBA{R: 255, G: 255, B: 255, A: 180},
		Decoicon_bg_idle_close:  color.NRGBA{R: 0, G: 0, B: 0, A: 0},
		Decoicon_bg_hover_close: color.NRGBA{R: 232, G: 17, B: 35, A: 255},
		Decoicon_bg_press_close: color.NRGBA{R: 196, G: 43, B: 28, A: 255},
	},
	Dark: ColorPalette{
		Fg:         color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BG:         color.NRGBA{R: 32, G: 32, B: 32, A: 255},
		ContrastFg: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		ContrastBg: color.NRGBA{R: 76, G: 194, B: 255, A: 255},

		Deco_fg:                 color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Deco_bg:                 color.NRGBA{R: 32, G: 32, B: 32, A: 255},
		Decoicon_fg_idle:        color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Decoicon_fg_hover:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Decoicon_fg_press:       color.NRGBA{R: 255, G: 255, B: 255, A: 200},
		Decoicon_bg_idle:        color.NRGBA{R: 0, G: 0, B: 0, A: 0},
		Decoicon_bg_hover:       color.NRGBA{R: 255, G: 255, B: 255, A: 26},
		Decoicon_bg_press:       color.NRGBA{R: 255, G: 255, B: 255, A: 65},
		Decoicon_fg_idle_close:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Decoicon_fg_hover_close: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Decoicon_fg_press_close: color.NRGBA{R: 255, G: 255, B: 255, A: 180},
		Decoicon_bg_idle_close:  color.NRGBA{R: 0, G: 0, B: 0, A: 0},
		Decoicon_bg_hover_close: color.NRGBA{R: 232, G: 17, B: 35, A: 255},
		Decoicon_bg_press_close: color.NRGBA{R: 241, G: 112, B: 122, A: 255},
	},
}

type Style struct {
	Theme     *material.Theme
	isDark    bool
	Palette   ColorPalette
	PaletteID StylePalette
	AniSys    *AniSystem

	palette_tar ColorPalette
	palette_cur ColorPalette
}

func NewStyle(theme *material.Theme, palette StylePalette) Style {
	isDark := false
	self := Style{
		Theme:     theme,
		isDark:    isDark,
		PaletteID: palette,
		AniSys:    NewAniSystem(),
	}
	self.apply()
	return self
}

func (self *Style) Run(gtx C) {
	self.AniSys.Run(gtx)
}

func (self *Style) Toggle() {
	self.isDark = !self.isDark
	self.apply()
}

// 设置颜色模式
//   - true: 暗色模式
//   - false: 亮色模式
func (self *Style) SetColorMode(isDark bool) {
	self.isDark = isDark
	self.apply()
}

func (self *Style) SetPalette(palette StylePalette) {
	self.PaletteID = palette
	self.apply()
}

func (self *Style) apply() {
	if self.isDark {
		self.Palette = self.PaletteID.Dark
	} else {
		self.Palette = self.PaletteID.Light
	}

	self.Theme.Bg = self.Palette.BG
	self.Theme.Fg = self.Palette.Fg
	self.Theme.ContrastFg = self.Palette.ContrastFg
	self.Theme.ContrastBg = self.Palette.ContrastBg
}
