package page

import (
	"app/src/ui"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type PageSettingInfo struct {
	Style      *ui.Style
	info_space *ui.AniEaseList
	text_list  [][]string
}

func NewPageSettingInfo(style *ui.Style) *PageSettingInfo {
	self := PageSettingInfo{
		Style:      style,
		info_space: ui.NewAniEaseList(style, ui.AniEaseList{Coe: 0.3, Num: 5, Space: 50, Tigger_percent: 0.4}),
	}

	self.text_list = [][]string{
		[]string{"version", "v0.1"},
		[]string{"developer", "yufu"},
	}

	return &self
}

func (self *PageSettingInfo) Layout(gtx C) D {
	self.info_space.Layout(gtx)

	widgets := []layout.FlexChild{}
	widgets = append(widgets, layout.Rigid(material.H4(self.Style.Theme, "Yami PDF").Layout))

	for _, t := range self.text_list {
		title := material.H6(self.Style.Theme, ui.GetWord(t[0], self.Style.Lang))
		desc := material.Body1(self.Style.Theme, ui.GetWord(t[1], self.Style.Lang))

		title.Color = self.Style.Palette.Fg_1
		desc.Color = self.Style.Palette.Fg_1

		widgets = append(widgets, layout.Rigid(title.Layout))
		widgets = append(widgets, layout.Rigid(desc.Layout))
	}

	children := make([]layout.FlexChild, 0, len(widgets)*2)
	for i, w := range widgets {
		children = append(children, w)
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(self.info_space.API_get_value(i))}.Layout))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (self *PageSettingInfo) OnEnter() {
	self.info_space.Reset()
}
func (self *PageSettingInfo) OnExit()  {}
func (self *PageSettingInfo) Destory() {}

func (self *PageSettingInfo) API_set_bridge_ui(bridge_ui ui.BridgeUI) {}
