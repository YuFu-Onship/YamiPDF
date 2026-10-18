package page

import (
	"app/assets"
	"app/src/ui"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageSetting struct {
	Style       *ui.Style
	page_router *PageRouter

	menus      []string
	menu_space ui.AniEaseList
	menu_cur   string
	menu_btns  map[string]*ui.SideBarWideButton

	icon_index map[string]string // [0] 图标索引

	bg_idle, fg_idle     color.NRGBA
	bg_active, fg_active color.NRGBA
}

func NewPageSetting(style *ui.Style, page_router *PageRouter) *PageSetting {
	self := PageSetting{
		Style:       style,
		page_router: page_router,

		menus:      []string{"custom", "info"},
		menu_cur:   "custom",
		menu_space: *ui.NewAniEaseList(style, ui.AniEaseList{Coe: 0.3, Num: 2, Space: 100, Tigger_percent: 0.3}),
		menu_btns:  map[string]*ui.SideBarWideButton{},

		icon_index: map[string]string{"custom": assets.Icon_palette, "info": assets.Icon_info},
	}

	self.page_router.
		API_add_manager("setting", map[string]*PageEntry{
			"custom": &PageEntry{keep_alive: false, builder: func() PageAPI { return NewPageSettingCustom(style, self.page_router) }},
			"info":   &PageEntry{keep_alive: false, builder: func() PageAPI { return NewPageSettingInfo(style) }},
		}).
		API_get_manager("setting").API_set_cur_page("custom")

	// 初始化侧边栏按钮
	for _, text := range self.menus {
		self.menu_btns[text] = ui.NewSideBarWideButton(style, text).
			API_set_text(func() string { return ui.GetWord(text, self.Style.Lang) }).
			API_set_icon(self.icon_index[text]).
			API_set_callback(self.callback_sidebtn)
	}

	return &self
}

func (self *PageSetting) Layout(gtx C) D {
	self.bg_idle = self.Style.Palette.Bg_1
	self.fg_idle = self.Style.Palette.Fg_1
	self.bg_active = self.Style.Palette.Bg_3
	self.fg_active = self.Style.Palette.Fg_3

	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(self.layout_menu),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Flexed(1, self.page_router.API_get_manager("setting").Layout),
	)
	return dims
}

func (self *PageSetting) OnEnter() {
	if len(self.menus) > 0 {
		self.menu_cur = self.menus[0]
		self.jump(self.menus[0])

		for id, btn := range self.menu_btns {
			btn.API_set_active(id == self.menus[0])
		}
	}
	self.menu_space.Reset()
}

func (self *PageSetting) OnExit()  {}
func (self *PageSetting) Destory() {}

// 页面跳转
func (self *PageSetting) jump(id string) {
	self.page_router.API_get_manager("setting").API_set_cur_page(id)
}

// 侧边栏按钮回调
func (self *PageSetting) callback_sidebtn(id string) {
	self.menu_cur = id
	self.jump(id)
	for menu, btn := range self.menu_btns {
		btn.API_set_active(id == menu)
	}
}

// 左侧侧边栏布局
func (self *PageSetting) layout_menu(gtx C) D {
	gtx.Constraints.Max.X = gtx.Dp(150)
	children := []layout.FlexChild{}

	for index, id := range self.menus {
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(self.menu_space.API_get_value(index))}.Layout))
		children = append(children, layout.Rigid(self.menu_btns[id].Layout))
	}

	self.menu_space.Layout(gtx)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
