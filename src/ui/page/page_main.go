package page

import (
	"app/assets"
	"app/src/ui"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageMain struct {
	Style *ui.Style

	page_router *PageRouter

	side_bar *ui.SideBar

	bridge_ui ui.BridgeUI
}

func NewPageMain(style *ui.Style, page_router *PageRouter) *PageMain {
	self := PageMain{
		Style:       style,
		page_router: page_router,
	}

	self.page_router.API_add_manager("main", map[string]*PageEntry{
		"bookshelf": &PageEntry{keep_alive: true, builder: func() PageAPI { return NewPageShelf(style, self.page_router) }},
		"setting":   &PageEntry{keep_alive: true, builder: func() PageAPI { return NewPageSetting(style, self.page_router) }},
	})
	self.page_router.API_get_manager("main").API_set_cur_page("bookshelf")

	self.side_bar = ui.NewSideBar()
	self.side_bar.API_add_btn("book_shelf", true, true,
		ui.NewSideBarButton(self.Style, assets.Icon_books).
			API_set_text(func() string { return ui.GetWord("bookshelf", self.Style.Lang) }).
			API_set_callfn(func() { self.jump("bookshelf") }))

	self.side_bar.API_add_btn("experimental", true, true,
		ui.NewSideBarButton(self.Style, assets.Icon_experimental).
			API_set_text(func() string { return ui.GetWord("experimental", self.Style.Lang) }).
			API_set_callfn(func() { self.jump("experimental") }))

	self.side_bar.API_add_btn("color_mode", false, false,
		ui.NewSideBarButton(self.Style, assets.Icon_moon).
			API_set_text(func() string { return ui.GetWord("color_mode", self.Style.Lang) }).
			API_set_callfn(self.Style.Toggle))

	self.side_bar.API_add_btn("setting", false, true,
		ui.NewSideBarButton(self.Style, assets.Icon_setting).
			API_set_text(func() string { return ui.GetWord("setting", self.Style.Lang) }).
			API_set_callfn(func() { self.jump("setting") }))

	self.side_bar.API_set_active("book_shelf")

	return &self
}

func (self *PageMain) Layout(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(layout.Spacer{Width: unit.Dp(60)}.Layout),
				layout.Flexed(1, self.page_router.API_get_manager("main").Layout),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{Bottom: unit.Dp(12), Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, self.side_bar.Layout)
		}),
	)
}

func (self *PageMain) OnEnter() {

}
func (self *PageMain) OnExit()  {}
func (self *PageMain) Destory() {}

func (self *PageMain) jump(id string) {
	self.page_router.API_get_manager("main").API_set_cur_page(id)
}

func (self *PageMain) API_set_bridge_ui(bridge_ui ui.BridgeUI) {
	self.bridge_ui = bridge_ui
}

func (self *PageMain) API_set_parent_manager(pm interface{}) {
}
