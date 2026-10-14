package page

import "app/src/ui"

type PageSettingCustom struct {
	Style *ui.Style
}

func NewPageSettingCustom(style *ui.Style) *PageSettingCustom {
	self := PageSettingCustom{
		Style: style,
	}
	return &self
}

func (self *PageSettingCustom) Layout(gtx C) D {
	return D{}
}

func (self *PageSettingCustom) OnEnter() {}
func (self *PageSettingCustom) OnExit()  {}
