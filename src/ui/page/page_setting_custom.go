package page

import (
	"app/assets"
	"app/src/ui"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CustomConfig struct {
	DarkMode bool
	Language ui.Language
	Palette  ui.StylePalette
}

type PageSettingCustom struct {
	Style       *ui.Style
	page_router *PageRouter

	custom_config CustomConfig
	space         *ui.AniEaseList

	// 配置项
	dark_mode_switch widget.Bool

	language_radio *ui.CustomRadio
	palette_radio  *ui.CustomPaletteRadio
}

func NewPageSettingCustom(style *ui.Style, page_router *PageRouter) *PageSettingCustom {
	self := PageSettingCustom{
		Style:       style,
		page_router: page_router,
		space:       ui.NewAniEaseList(style, ui.AniEaseList{Num: 3, Space: 100, Tigger_percent: 0.3, Coe: 0.3}),
	}

	self.init_language()
	self.init_palette()

	return &self
}

func (self *PageSettingCustom) Layout(gtx C) D {
	self.space.Layout(gtx)

	gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(600))

	widgets := []layout.FlexChild{}
	widgets = append(widgets, layout.Rigid(func(gtx C) D { return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, self.layout_language) }))
	widgets = append(widgets, layout.Rigid(func(gtx C) D { return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, self.layout_palette) }))

	children := make([]layout.FlexChild, 0, len(widgets)*2)
	for i, widget := range widgets {
		if i != 0 {
			children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout))
		}
		children = append(children, widget)
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(self.space.API_get_value(i))}.Layout))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (self *PageSettingCustom) OnEnter() { self.space.Reset() }
func (self *PageSettingCustom) OnExit()  {}
func (self *PageSettingCustom) Destory() {}

// 个性化参数 ------------------------------------

// 暗色模式
func (self *PageSettingCustom) init_darkmode() {
	self.dark_mode_switch = widget.Bool{}
	self.dark_mode_switch.Value = self.Style.API_get_darkmode()
}
func (self *PageSettingCustom) layout_darkmode(gtx C) D {
	s := material.Switch(self.Style.Theme, &self.dark_mode_switch, "")
	self.Style.API_set_darkmode(self.dark_mode_switch.Value)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.H6(self.Style.Theme, assets.Icon_moon).Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(material.H6(self.Style.Theme, ui.GetWord("dark_mode", self.Style.Lang)).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(s.Layout),
	)
}

// 配色
func (self *PageSettingCustom) init_palette() {
	self.palette_radio = ui.NewCustomPaletteRadio(self.Style).
		API_add_choice("windows", func() color.NRGBA { return self.Style.API_get_theme_color(ui.Palette_windows) }).
		API_add_choice("teto", func() color.NRGBA { return self.Style.API_get_theme_color(ui.Palette_teto) }).
		API_add_choice("koishi", func() color.NRGBA { return self.Style.API_get_theme_color(ui.Palette_koishi) }).
		API_set_callback(func(key string) {
			switch key {
			case "windows":
				self.Style.API_set_palette(ui.Palette_windows)
			case "teto":
				self.Style.API_set_palette(ui.Palette_teto)
			case "koishi":
				self.Style.API_set_palette(ui.Palette_koishi)
			}
		}).
		API_set_value(func() string {
			switch self.Style.PaletteID {
			case ui.Palette_windows:
				return "windows"
			case ui.Palette_teto:
				return "teto"
			case ui.Palette_koishi:
				return "koishi"
			}
			return "windows"
		}())
}
func (self *PageSettingCustom) layout_palette(gtx C) D {
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.H6(self.Style.Theme, assets.Icon_palette).Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(material.H6(self.Style.Theme, ui.GetWord("palette", self.Style.Lang)).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(self.palette_radio.Layout),
	)

	return dims
}

// 语言
func (self *PageSettingCustom) init_language() {
	self.language_radio = ui.NewCustomRadio(self.Style).
		API_set_text("cn", func() string { return "中文" }).
		API_set_text("en", func() string { return "English" }).
		API_set_text("jp", func() string { return "日本語" }).
		API_set_callback(func(key string) {
			switch key {
			case "cn":
				self.Style.Lang = ui.CN
			case "en":
				self.Style.Lang = ui.EN
			case "jp":
				self.Style.Lang = ui.JP
			}
		}).
		API_set_value(func() string {
			switch self.Style.Lang {
			case ui.CN:
				return "cn"
			case ui.EN:
				return "en"
			case ui.JP:
				return "jp"
			}
			return "cn"
		}())
}

func (self *PageSettingCustom) layout_language(gtx C) D {
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.H6(self.Style.Theme, assets.Icon_language).Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(material.H6(self.Style.Theme, ui.GetWord("language", self.Style.Lang)).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(self.language_radio.Layout),
	)

	return dims
}
