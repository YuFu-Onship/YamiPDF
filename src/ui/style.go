package ui

import (
	"image/color"
	"math"

	"gioui.org/widget/material"
)

// 标题栏配色
type DecoColorPalette struct {
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

// 正常色板配色
type ColorPalette struct {
	Fg_1 color.NRGBA
	Bg_1 color.NRGBA

	Fg_2 color.NRGBA
	Bg_2 color.NRGBA

	Fg_3 color.NRGBA
	Bg_3 color.NRGBA
}

// 样式配色
type StylePalette struct {
	Dark  ColorPalette
	Light ColorPalette
}

// 临时float64颜色存储
type TempFloatNRGBA struct{ R, G, B, A float64 }
type TempColorPalette struct{ Fg_1, Fg_2, Fg_3, Bg_1, Bg_2, Bg_3 TempFloatNRGBA }

func from_uintColor_to_floatColor(c color.NRGBA) TempFloatNRGBA {
	return TempFloatNRGBA{R: float64(c.R), G: float64(c.G), B: float64(c.B), A: float64(c.A)}
}

func from_floatColor_to_uintColor(c TempFloatNRGBA) color.NRGBA {
	return color.NRGBA{
		R: uint8(math.Round(c.R)),
		G: uint8(math.Round(c.G)),
		B: uint8(math.Round(c.B)),
		A: uint8(math.Round(c.A)),
	}
}

func from_uintPalette_to_floatPalette(p ColorPalette) TempColorPalette {
	return TempColorPalette{
		Fg_1: from_uintColor_to_floatColor(p.Fg_1),
	}
}

// 样式 -------------------------------------------
type Style struct {
	Theme  *material.Theme
	isDark bool

	Palette     ColorPalette
	Palette_tar ColorPalette
	Palette_cur TempColorPalette

	Palette_deco DecoColorPalette // 标题栏配色
	PaletteID    StylePalette
	AniSys       *AniSystem
	Lang         Language
}

func NewStyle(theme *material.Theme, palette StylePalette) Style {
	isDark := false
	self := Style{
		Theme:     theme,
		isDark:    isDark,
		PaletteID: palette,
		AniSys:    NewAniSystem(),
		Lang:      CN,
	}

	self.apply()
	return self
}

func (self *Style) Run(gtx C) {
	self.AniSys.Run(gtx)

}

func (self *Style) Update(gtx C) {
	self.color_update()
	self.Theme.Bg = self.Palette.Bg_1
	self.Theme.Fg = self.Palette.Fg_1
	self.Theme.ContrastFg = self.Palette.Bg_3
	self.Theme.ContrastBg = self.Palette.Fg_3
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
		self.Palette_tar = self.PaletteID.Dark
		self.Palette_deco = Deco_palette_dark
	} else {
		self.Palette_tar = self.PaletteID.Light
		self.Palette_deco = Deco_palette_light
	}

	self.Palette.Fg_1 = self.Palette_tar.Fg_1
	self.Palette.Fg_2 = self.Palette_tar.Fg_2
	self.Palette.Fg_3 = self.Palette_tar.Fg_3

}

func (self *Style) color_update() {
	self.color_lerp(&self.Palette_cur.Bg_1, self.Palette_tar.Bg_1)
	self.color_lerp(&self.Palette_cur.Bg_2, self.Palette_tar.Bg_2)
	self.color_lerp(&self.Palette_cur.Bg_3, self.Palette_tar.Bg_3)

	self.Palette.Bg_1 = from_floatColor_to_uintColor(self.Palette_cur.Bg_1)
	self.Palette.Bg_2 = from_floatColor_to_uintColor(self.Palette_cur.Bg_2)
	self.Palette.Bg_3 = from_floatColor_to_uintColor(self.Palette_cur.Bg_3)
}

func (self *Style) color_lerp(color_cur *TempFloatNRGBA, color_tar color.NRGBA) {
	cr, cg, cb, ca := color_cur.R, color_cur.G, color_cur.B, color_cur.A
	tr, tg, tb, ta := float64(color_tar.R), float64(color_tar.G), float64(color_tar.B), float64(color_tar.A)

	result := self.AniSys.EaseGroup(color_cur, []float64{cr, cg, cb, ca}, []float64{tr, tg, tb, ta}, 0.1, 0.3)

	color_cur.R = result[0]
	color_cur.G = result[1]
	color_cur.B = result[2]
	color_cur.A = result[3]
}

// 设置当前色板
func (self *Style) API_set_palette(palette StylePalette) {
	self.PaletteID = palette
	self.apply()
}
