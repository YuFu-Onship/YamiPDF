package page

import (
	"app/assets"
	"app/src/ui"

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
	language_enum    widget.Enum
	palette_enum     widget.Enum
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

	widgets := []layout.FlexChild{}
	widgets = append(widgets, layout.Rigid(func(gtx C) D { return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, self.layout_language) }))
	widgets = append(widgets, layout.Rigid(func(gtx C) D { return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, self.layout_palette) }))

	children := make([]layout.FlexChild, 0, len(widgets)*2)
	for i, widget := range widgets {
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
	self.palette_enum = widget.Enum{}

	switch self.Style.PaletteID {
	case ui.Palette_windows:
		self.palette_enum.Value = "windows"
	case ui.Palette_teto:
		self.palette_enum.Value = "teto"
	case ui.Palette_koishi:
		self.palette_enum.Value = "koishi"

	}
}
func (self *PageSettingCustom) layout_palette(gtx C) D {

	windows := material.RadioButton(self.Style.Theme, &self.palette_enum, "windows", "")
	windows.IconColor = self.Style.API_get_theme_color(ui.Palette_windows)

	teto := material.RadioButton(self.Style.Theme, &self.palette_enum, "teto", "")
	teto.IconColor = self.Style.API_get_theme_color(ui.Palette_teto)

	koishi := material.RadioButton(self.Style.Theme, &self.palette_enum, "koishi", "")
	koishi.IconColor = self.Style.API_get_theme_color(ui.Palette_koishi)

	self.palette_enum.Hovered()

	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.H6(self.Style.Theme, assets.Icon_palette).Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(material.H6(self.Style.Theme, ui.GetWord("palette", self.Style.Lang)).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(windows.Layout),
		layout.Rigid(teto.Layout),
		layout.Rigid(koishi.Layout),
	)

	switch self.palette_enum.Value {
	case "windows":
		self.Style.API_set_palette(ui.Palette_windows)
	case "teto":
		self.Style.API_set_palette(ui.Palette_teto)
	case "koishi":
		self.Style.API_set_palette(ui.Palette_koishi)
	}

	return dims
}

// 语言
func (self *PageSettingCustom) init_language() {
	self.language_enum = widget.Enum{}
	switch self.Style.Lang {
	case ui.CN:
		self.language_enum.Value = "cn"
	case ui.EN:
		self.language_enum.Value = "en"
	case ui.JP:
		self.language_enum.Value = "jp"

	}
}
func (self *PageSettingCustom) layout_language(gtx C) D {
	cn := material.RadioButton(self.Style.Theme, &self.language_enum, "cn", ui.GetWord("中文", 0))
	cn.Color = self.Style.Palette.Fg_1
	cn.IconColor = self.Style.Palette.Fg_1

	en := material.RadioButton(self.Style.Theme, &self.language_enum, "en", ui.GetWord("English", 0))
	en.Color = self.Style.Palette.Fg_1
	en.IconColor = self.Style.Palette.Fg_1

	jp := material.RadioButton(self.Style.Theme, &self.language_enum, "jp", ui.GetWord("日本語", 0))
	jp.Color = self.Style.Palette.Fg_1
	jp.IconColor = self.Style.Palette.Fg_1

	self.language_enum.Hovered()
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.H6(self.Style.Theme, assets.Icon_language).Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(material.H6(self.Style.Theme, ui.GetWord("language", self.Style.Lang)).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(cn.Layout),
		layout.Rigid(en.Layout),
		layout.Rigid(jp.Layout),
	)

	switch self.language_enum.Value {
	case "cn":
		self.Style.Lang = ui.CN
	case "en":
		self.Style.Lang = ui.EN
	case "jp":
		self.Style.Lang = ui.JP
	}

	return dims
}
