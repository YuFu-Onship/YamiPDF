package page

import (
	"app/src/ui"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageSettingCustom struct {
	Style *ui.Style

	dropdown_language *ui.DropDownChoice
	dropdown_color    *ui.DropDownChoice
}

func NewPageSettingCustom(style *ui.Style) *PageSettingCustom {
	self := PageSettingCustom{
		Style: style,
	}

	self.dropdown_language = ui.NewDropDownChoice(style).
		API_add_choice("Chinese").
		API_add_choice("English").
		API_add_choice("Japanese").
		API_set_callfn(func(choice string) {
			switch choice {
			case "Chinese":
				self.Style.Lang = ui.CN
			case "English":
				self.Style.Lang = ui.EN
			case "Japanese":
				self.Style.Lang = ui.JP
			}
		})

	self.dropdown_color = ui.NewDropDownChoice(style).
		API_add_choice("Windows").
		API_add_choice("Teto").
		API_add_choice("Koishi").
		API_set_callfn(func(choice string) {
			switch choice {
			case "Windows":
				self.Style.API_set_palette(ui.Palette_windows)
			case "Teto":
				self.Style.API_set_palette(ui.Palette_teto)
			case "Koishi":
				self.Style.API_set_palette(ui.Palette_koishi)

			}
		})

	return &self
}

func (self *PageSettingCustom) Layout(gtx C) D {
	self.dropdown_language.API_set_border_text(ui.GetWord("language", self.Style.Lang))
	self.dropdown_color.API_set_border_text(ui.GetWord("palette", self.Style.Lang))

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(self.dropdown_language.Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
		layout.Rigid(self.dropdown_color.Layout),
	)
}

func (self *PageSettingCustom) OnEnter() {}
func (self *PageSettingCustom) OnExit()  {}
