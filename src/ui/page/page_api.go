package page

import (
	"app/src/ui"

	"gioui.org/layout"
)

type PageAPI interface {
	Layout(layout.Context) layout.Dimensions
	OnEnter()
	OnExit()
}

type PageManager struct {
	Style *ui.Style

	page_cur string
	pages    map[string]PageAPI
}

func NewPageManager(style *ui.Style) *PageManager {
	self := PageManager{
		Style: style,

		page_cur: "",
		pages:    map[string]PageAPI{},
	}
	return &self
}

func (self *PageManager) Layout(gtx C) D {
	if page, exist := self.pages[self.page_cur]; exist {
		return page.Layout(gtx)
		// return self.pages[self.page_cur].Layout(gtx)
	}
	return D{}
}

// 跳转到对应的页面
func (self *PageManager) API_set_cur_page(id string) *PageManager {
	_, exist := self.pages[id]
	if exist {
		self.page_cur = id
	}
	return self
}

// 添加新的页面
func (self *PageManager) API_add_page(id string, page PageAPI) *PageManager {
	self.pages[id] = page
	return self
}
