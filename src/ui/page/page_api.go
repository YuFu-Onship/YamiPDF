package page

import (
	"app/src/ui"
	"image"
	"log"

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
	API_set_bridge_ui(bridge_ui ui.BridgeUI)
	API_set_parent_manager(interface{})
}

// 管理不同页面管理器的管理器 ----------------------------------------------------------
type PageRouter struct {
	managers map[string]*PageManager
}

func NewPageRouter() *PageRouter {
	self := PageRouter{
		managers: make(map[string]*PageManager),
	}
	return &self
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
		entry.instance.API_set_parent_manager(self) // 传递自身,实现下级页面管理器访问上级
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

// 设置当前页面管理器的id
func (self *PageManager) API_set_id(id string) *PageManager {
	self.id = id
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
		if page.instance != nil {
			page.instance.API_set_bridge_ui(bridge_ui)
		}
	}
}

// 得到上一级的页面管理器
func (self *PageManager) API_get_last_manager() *PageManager {
	log.Println(self.id)
	return self.last_manager
}
