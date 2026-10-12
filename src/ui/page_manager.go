package ui

import "gioui.org/layout"

type PageAPI interface {
	Layout(layout.Context) layout.Dimensions
	OnEnter()
	OnExit()
}

type PangeManager struct {
	Style *Style

	page_cur string
	pages    map[string]PageAPI
}

func NewPageManager(style *Style) *PangeManager {
	self := PangeManager{
		Style: style,

		page_cur: "",
		pages:    map[string]PageAPI{},
	}
	return &self
}

func (self *PangeManager) Layout(gtx C) D {
	return D{}
}

// 跳转到对应的页面
func (self *PangeManager) API_set_cur_page(id string) *PangeManager {
	_, exist := self.pages[id]
	if exist {
		self.page_cur = id
	}
	return self
}

// 添加新的页面
func (self *PangeManager) API_add_page(id string, page PageAPI) *PangeManager {
	self.pages[id] = page
	return self
}
