package page

import (
	"app/assets"
	"app/src/ui"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageSetting struct {
	Style        *ui.Style
	page_manager *PageManager

	menus      []string
	menu_space ui.AniEaseList
	menu_cur   string

	menu_btns map[string]*ui.SideBarWideButton

	text_index map[string][]string

	bg_idle, fg_idle     color.NRGBA
	bg_active, fg_active color.NRGBA
}

func NewPageSetting(style *ui.Style) *PageSetting {
	self := PageSetting{
		Style:        style,
		page_manager: NewPageManager(style),

		menus:      []string{},
		menu_space: *ui.NewAniEaseList(style, ui.AniEaseList{Coe: 0.3, Num: 2, Space: 100, Tigger_percent: 0.3}),
		menu_cur:   "",
		menu_btns:  map[string]*ui.SideBarWideButton{},
		// menu_clickables: map[string]*widget.Clickable{},
		text_index: map[string][]string{},
	}

	self.page_manager.API_set_id("setting")

	self.menus = []string{"custom", "info"}
	self.menu_cur = "custom"

	self.text_index["custom"] = []string{assets.Icon_palette}
	self.text_index["info"] = []string{assets.Icon_info}

	// self.page_manager.Add_page("custom", false, func() PageAPI { return NewPageSettingCustom(style) })
	// self.page_manager.Add_page("info", false, func() PageAPI { return NewPageSettingInfo(style) })
	self.page_manager.API_set_cur_page(self.menu_cur)

	// 初始化侧边栏按钮
	for _, text := range self.menus {
		self.menu_btns[text] = ui.NewSideBarWideButton(style, text).
			API_set_text(func() string { return ui.GetWord(text, self.Style.Lang) }).
			API_set_icon(self.text_index[text][0])

		self.menu_btns[text].API_set_callback(func(id string) {
			self.menu_cur = id
			self.page_manager.API_set_cur_page(id)
			for _, k := range self.menus {
				self.menu_btns[k].API_set_active(k == self.menu_cur)
			}
		})

	}
	for _, k := range self.menus {
		self.menu_btns[k].API_set_active(k == self.menu_cur)
	}

	return &self
}

func (self *PageSetting) Layout(gtx C) D {
	self.bg_idle = self.Style.Palette.Bg_1
	self.fg_idle = self.Style.Palette.Fg_1
	self.bg_active = self.Style.Palette.Bg_3
	self.fg_active = self.Style.Palette.Fg_3

	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(self.menu),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Flexed(1, self.page_manager.Layout),
	)
	return dims
}

func (self *PageSetting) OnEnter() {
	if len(self.menus) > 0 {
		self.menu_cur = self.menus[0]
		self.page_manager.API_set_cur_page(self.menus[0])
	}
	self.menu_space.Reset()
	self.page_manager.OnEnter()

}
func (self *PageSetting) OnExit() {
	self.page_manager.OnExit()
}
func (self *PageSetting) Destory()                                {}
func (self *PageSetting) API_set_bridge_ui(bridge_ui ui.BridgeUI) {}
func (self *PageSetting) API_set_parent_manager(interface{})      {}

// 左侧侧边栏
func (self *PageSetting) menu(gtx C) D {
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
