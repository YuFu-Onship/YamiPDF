package page

import (
	"app/src/ui"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
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

	is_change    bool
	offset_y_tar float64
	offset_y_cur float64
}

func NewPageManager(style *ui.Style) *PageManager {
	self := PageManager{
		Style: style,

		page_cur: "",
		pages:    map[string]PageAPI{},

		is_change:    false,
		offset_y_tar: 0,
		offset_y_cur: 0,
	}
	return &self
}

func (self *PageManager) Layout(gtx C) D {

	if self.is_change {
		self.is_change = false
		self.offset_y_cur = float64(gtx.Constraints.Max.Y)
	}

	self.Style.AniSys.Ease64(&self.offset_y_cur, 1, &self.offset_y_cur, self.offset_y_tar, 0.4)

	if page, exist := self.pages[self.page_cur]; exist {
		trans_offset := op.Offset(image.Pt(0, int(self.offset_y_cur))).Push(gtx.Ops)
		dims := page.Layout(gtx)
		trans_offset.Pop()

		return dims
	}
	return D{}
}

// 跳转到对应的页面
func (self *PageManager) API_set_cur_page(id string) *PageManager {
	if self.page_cur == id {
		return self
	}

	if _, exist := self.pages[id]; !exist {
		return self
	}

	// self.pages[self.page_cur].OnExit()
	self.page_cur = id
	self.pages[self.page_cur].OnEnter()
	self.is_change = true

	return self
}

// 添加新的页面
func (self *PageManager) API_add_page(id string, page PageAPI) *PageManager {
	self.pages[id] = page
	return self
}
