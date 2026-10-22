package page

import (
	rumia "app/src"
	"app/src/ui"
	"image"

	"gioui.org/op"
)

type ParentManager interface {
	Set_parent(PageManager)
}

type PageBuilder func() PageAPI

type PageAPI interface {
	Layout(C) D
	OnEnter()
	OnExit()
	Destory()
}

// 管理不同页面管理器的管理器 ----------------------------------------------------------
type PageRouter struct {
	Style    *ui.Style
	managers map[string]*PageManager

	Tools       ui.BridgeUI
	BridgeTrunk rumia.BridgeTrunk
}

func NewPageRouter(style *ui.Style) *PageRouter {
	self := PageRouter{
		Style:    style,
		managers: make(map[string]*PageManager),
	}
	return &self
}

func (self *PageRouter) API_add_manager(id string, pages map[string]*PageEntry) *PageRouter {
	if manager, exist := self.managers[id]; exist {
		manager.Add_pages(pages)
	} else {
		self.managers[id] = NewPageManager(self.Style).API_set_id(id).Add_pages(pages)
	}

	return self
}

func (self *PageRouter) API_get_manager(id string) *PageManager {
	if manager, exist := self.managers[id]; exist {
		return manager
	}
	return nil
}

// 得到指定页面管理器的当前页面id
func (self *PageRouter) API_get_manager_cur_page_id(manager_id string) string {
	if manager, exist := self.managers[manager_id]; exist {
		return manager.API_get_cur_id()
	}
	return ""
}

// 页面管理器 --------------------------------------------------------------------------
type PageEntry struct {
	keep_alive bool
	builder    PageBuilder
	instance   PageAPI
}

type PageManager struct {
	Style *ui.Style
	id    string // 当前页面管理器的id,可选可不选

	page_cur string
	pages    map[string]*PageEntry

	is_change    bool
	offset_y_tar float64
	offset_y_cur float64

	last_manager *PageManager
}

func NewPageManager(style *ui.Style) *PageManager {
	self := PageManager{
		Style: style,
		id:    "",

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

// func (self *PageManager) OnEnter() {
// 	for _, page := range self.pages {
// 		page.instance.OnEnter()
// 	}
// }
// func (self *PageManager) OnExit() {
// 	for _, page := range self.pages {
// 		page.instance.OnExit()
// 	}
// }

// 设置当前页面管理器的id
func (self *PageManager) API_set_id(id string) *PageManager {
	self.id = id
	return self
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
		if old_entry.instance != nil {
			if old_entry.keep_alive {
				old_entry.instance.OnExit()
			} else {
				old_entry.instance.Destory()
				old_entry.instance = nil
			}
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

// 添加数个页面
func (self *PageManager) Add_pages(pages map[string]*PageEntry) *PageManager {
	if self.pages == nil {
		self.pages = make(map[string]*PageEntry)
	}

	for k, v := range pages {
		self.pages[k] = v
	}
	return self
}
