package page

import (
	"app/assets"
	"app/src/ui"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PageSetting struct {
	Style        *ui.Style
	page_manager *PageManager

	menus           []string
	menu_space      ui.AniEaseList
	menu_cur        string
	menu_clickables map[string]*widget.Clickable
	menu_callback   func(string)

	text_index map[string][]string

	bg_idle, fg_idle     color.NRGBA
	bg_active, fg_active color.NRGBA
}

func NewPageSetting(style *ui.Style) *PageSetting {
	self := PageSetting{
		Style:        style,
		page_manager: NewPageManager(style),

		menus:           []string{},
		menu_space:      *ui.NewAniEaseList(style, ui.AniEaseList{Coe: 0.3, Num: 2, Space: 100, Tigger_percent: 0.3}),
		menu_cur:        "",
		menu_clickables: map[string]*widget.Clickable{},
		text_index:      map[string][]string{},
	}

	self.menus = append(self.menus, "custom", "info")
	self.text_index["custom"] = []string{assets.Icon_palette}
	self.text_index["info"] = []string{assets.Icon_info}

	self.page_manager.Add_page("custom", false, func() PageAPI { return NewPageSettingCustom(style) })
	self.page_manager.Add_page("info", false, func() PageAPI { return NewPageSettingInfo(style) })

	self.page_manager.API_set_cur_page("custom")

	self.menu_callback = func(id string) {
		self.menu_cur = id
		self.page_manager.API_set_cur_page(id)
	}

	for _, text := range self.menus {
		self.menu_clickables[text] = &widget.Clickable{}
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

// 左侧侧边栏
func (self *PageSetting) menu(gtx C) D {
	gtx.Constraints.Max.X = gtx.Dp(150)
	children := []layout.FlexChild{}

	for index, id := range self.menus {
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
		children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(self.menu_space.API_get_value(index))}.Layout))
		children = append(children, layout.Rigid(func(gtx C) D {
			return self.build_muen_btn(gtx, self.menu_clickables[id], id)
		}))
	}

	self.menu_space.Layout(gtx)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

// 构建 侧边栏按钮
func (self *PageSetting) build_muen_btn(gtx C, clickable *widget.Clickable, id string) D {
	var icon string
	if _, exist := self.text_index[id]; exist {
		icon = self.text_index[id][0]
	} else {
		icon = assets.Icon_loading
	}

	// icon := Ternary(exist, self.text_index[id][0], assets.Icon_loading)

	text := ui.GetWord(id, self.Style.Lang)

	if clickable.Clicked(gtx) {
		self.menu_callback(id)
	}

	is_hover := clickable.Hovered()
	is_active := id == self.menu_cur

	btn := material.Button(self.Style.Theme, clickable, "")
	btn.Background = Ternary(is_active, self.bg_active, self.bg_idle)

	label_icon := material.H6(self.Style.Theme, icon)
	label_text := material.H6(self.Style.Theme, text)

	label_text.Color = Ternary(is_active, self.fg_active, self.fg_idle)
	label_icon.Color = Ternary(is_active, self.fg_active, self.fg_idle)

	label_text.MaxLines = 1

	label_icon.Font.Style = Ternary(is_hover, font.Italic, font.Regular)
	label_icon.Font.Style = Ternary(self.menu_cur == id, font.Italic, label_icon.Font.Style)

	gtx.Constraints.Min.Y = gtx.Dp(40)
	gtx.Constraints.Max.Y = gtx.Dp(40)
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(btn.Layout),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
				layout.Rigid(label_icon.Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
				layout.Rigid(label_text.Layout),
				layout.Flexed(1, layout.Spacer{}.Layout),
			)
		}),
	)
}
