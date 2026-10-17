package page

import (
	"app/src/ui"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

type PageBuilder func() PageAPI

type PageAPI interface {
	Layout(layout.Context) layout.Dimensions
	OnEnter()
	OnExit()
	Destory()
	API_set_bridge_ui(bridge_ui ui.BridgeUI)
}

type PageEntry struct {
	keep_alive bool
	builder    PageBuilder
	instance   PageAPI
}

type PageManager struct {
	Style *ui.Style

	page_cur string
	pages    map[string]*PageEntry

	is_change    bool
	offset_y_tar float64
	offset_y_cur float64
}

func NewPageManager(style *ui.Style) *PageManager {
	self := PageManager{
		Style: style,

		page_cur: "",
		pages:    make(map[string]*PageEntry),

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

	if page := self.get_or_create_page(self.page_cur); page != nil {
		trans_offset := op.Offset(image.Pt(0, int(self.offset_y_cur))).Push(gtx.Ops)
		dims := page.Layout(gtx)
		trans_offset.Pop()
		return dims
	}

	return D{}
}

func (self *PageManager) get_or_create_page(id string) PageAPI {
	entry, exist := self.pages[id]
	if !exist {
		return nil
	}

	if entry.instance == nil && entry.builder != nil {
		entry.instance = entry.builder()
	}
	return entry.instance
}

func (self *PageManager) OnEnter() {
	// for _, page := range self.pages {
	// 	page.instance.OnEnter()
	// }
}
func (self *PageManager) OnExit() {
	// for _, page := range self.pages {
	// 	page.instance.OnExit()
	// }
}

// 跳转到对应的页面
func (self *PageManager) API_set_cur_page(id string) *PageManager {
	if self.page_cur == id {
		return self
	}

	if _, exist := self.pages[id]; !exist {
		return self
	}

	// 清理上一个界面
	if old_entry, exist := self.pages[self.page_cur]; exist {
		if old_entry.keep_alive {
			old_entry.instance.OnExit()
		} else {
			old_entry.instance = nil
		}
	}

	self.page_cur = id

	// 打开或创建下一个界面
	if page := self.get_or_create_page(self.page_cur); page != nil {
		page.OnEnter()
	}

	self.is_change = true
	return self
}

// 得到当前页面id
func (self *PageManager) API_get_cur_id() string {
	return self.page_cur
}

// 注册页面
func (self *PageManager) Add_page(id string, keep_alive bool, builder PageBuilder) *PageManager {
	self.pages[id] = &PageEntry{
		builder:    builder,
		keep_alive: keep_alive,
	}
	return self
}

// 给某个元素设置page主界面接口
func (self *PageManager) API_set_bridge_ui(id string, bridge_ui ui.BridgeUI) {
	if page, exist := self.pages[id]; exist {
		page.instance.API_set_bridge_ui(bridge_ui)
	}
}
