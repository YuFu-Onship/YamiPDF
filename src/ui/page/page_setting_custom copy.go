package page

import "app/src/ui"

type PageSettingInfo struct {
	Style *ui.Style
}

func NewPageSettingInfo(style *ui.Style) *PageSettingInfo {
	self := PageSettingInfo{
		Style: style,
	}
	return &self
}

func (self *PageSettingInfo) Layout(gtx C) D {
	return D{}
}

func (self *PageSettingInfo) OnEnter() {}
func (self *PageSettingInfo) OnExit()  {}
