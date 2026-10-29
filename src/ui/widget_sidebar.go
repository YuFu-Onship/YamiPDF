package ui

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SideBar struct {
	btns_top    []*SideBarButton
	btns_bottom []*SideBarButton
	btns        map[string]*SideBarButton
	pad         float64
}

// 侧边栏
func NewSideBar() *SideBar {
	self := SideBar{
		btns_top:    []*SideBarButton{},
		btns_bottom: []*SideBarButton{},
		btns:        map[string]*SideBarButton{},
		pad:         4,
	}
	return &self
}

func (self *SideBar) Layout(gtx C) D {
	self.update_active(gtx)

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D { return self.arrange_btn_list(gtx, self.btns_top) }),
		layout.Flexed(1, func(gtx C) D { return layout.Spacer{}.Layout(gtx) }),
		layout.Rigid(func(gtx C) D { return self.arrange_btn_list(gtx, self.btns_bottom) }),
	)
	return dims
}

// 排列按钮竖直方向
func (self *SideBar) arrange_btn_list(gtx C, btns []*SideBarButton) D {
	children := make([]layout.FlexChild, len(btns))

	length := len(btns)
	for i, btn := range btns {
		bottom_pad := Ternary(i != length-1, self.pad, 0)
		children[i] = layout.Rigid(func(gtx C) D { return layout.Inset{Bottom: unit.Dp(bottom_pad)}.Layout(gtx, btn.Layout) })
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

// 检查激活
func (self *SideBar) update_active(gtx C) {
	allBtns := append(self.btns_top, self.btns_bottom...)
	for _, btn := range allBtns {
		if btn.clickable_icon.Clicked(gtx) {
			if btn.is_allow_active {
				for _, b := range allBtns {
					b.API_set_active(b == btn)
				}
			}
			btn.fn()
		}
	}
}

// 增添新按钮
//   - is_top 是否在上半部分
//   - is_allow_active 是否出现激活状态
func (self *SideBar) API_add_btn(id string, is_top bool, is_allow_active bool, btn *SideBarButton) *SideBar {
	btn.API_set_allow_active(is_allow_active)
	self.btns[id] = btn
	if is_top {
		self.btns_top = append(self.btns_top, btn)
	} else {
		self.btns_bottom = append(self.btns_bottom, btn)
	}
	return self
}

// 设置按钮激活
func (self *SideBar) API_set_active(id string) {
	if btn, exist := self.btns[id]; exist {
		if btn.is_allow_active {
			for _, b := range self.btns {
				b.API_set_active(b == btn)
			}
		}
	}
}

// ------------------------------------------------------------------

type modeSideBarButton int

const (
	mode_sidebar_btn_idle modeSideBarButton = iota
	mode_sidebar_btn_hover
	mode_sidebar_btn_active
)

// 侧边栏说明按钮
type SideBarButton struct {
	Style *Style

	icon string

	text      func() string
	text_last string
	state     modeSideBarButton

	is_active       bool
	is_allow_active bool

	clickable_icon *widget.Clickable

	desc_width_tar float64
	desc_width_cur float64
	bar_height_tar float64
	bar_height_cur float64

	is_redraw_desc bool
	desc_size      image.Point

	fn func() //回调函数
}

// 创建新的侧边栏按钮
func NewSideBarButton(style *Style, icon_1 string) *SideBarButton {
	self := SideBarButton{
		Style:           style,
		icon:            icon_1,
		state:           mode_sidebar_btn_idle,
		is_active:       false,
		is_allow_active: true,
		fn:              func() {},
		is_redraw_desc:  false,
	}

	self.clickable_icon = &widget.Clickable{}

	return &self
}

func (self *SideBarButton) Layout(gtx C) D {
	radius := gtx.Dp(4)

	// 更新文本
	self.is_redraw_desc = Ternary(self.text_last == self.text(), false, true)
	self.text_last = self.text()

	if self.is_redraw_desc {
		self.is_redraw_desc = false
		gtxMeasure := gtx
		gtxMeasure.Constraints.Min = image.Point{}
		l := material.Label(self.Style.Theme, unit.Sp(20), self.text())
		l.MaxLines = 1
		macro := op.Record(gtxMeasure.Ops)
		t := l.Layout(gtxMeasure)
		macro.Stop()
		self.desc_size = t.Size
	}

	// 更新状态机
	self.state = Ternary(self.clickable_icon.Hovered(), mode_sidebar_btn_hover, mode_sidebar_btn_idle)
	self.state = Ternary(self.is_active, mode_sidebar_btn_active, self.state)
	self.bar_height_tar = Ternary(self.is_active, float64(gtx.Dp(24)), 0.0)

	// element
	ele_bg := Ternary(self.state == mode_sidebar_btn_idle, self.Style.Palette.Bg_2, self.Style.Palette.Bg_3)
	ele_fg := Ternary(self.state == mode_sidebar_btn_idle, self.Style.Palette.Fg_2, self.Style.Palette.Fg_3)

	style_btn := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(40), gtx.Dp(40))
		var btn material.ButtonStyle
		btn = material.Button(self.Style.Theme, self.clickable_icon, self.icon)
		btn.Font.Style = Ternary(self.state == mode_sidebar_btn_idle, font.Regular, font.Italic)
		btn.Background = ele_bg
		btn.Color = ele_fg
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(24)
		btn.CornerRadius = unit.Dp(radius)

		return btn.Layout(gtx)
	}

	// 左侧指示条
	style_bar := func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Dp(6)
		if self.state != mode_sidebar_btn_active {
			return D{Size: gtx.Constraints.Min}
		}
		// draw_rectangle_2(gtx, 0, int(unit.Dp(20-self.bar_height_cur*0.5)), gtx.Dp(6), int(unit.Dp(self.bar_height_cur)), ele_bg, gtx.Dp(2))
		draw_rectangle_2(gtx, 0, gtx.Dp(unit.Dp(20))-int(self.bar_height_cur/2), gtx.Dp(6), int(unit.Dp(self.bar_height_cur)), ele_bg, gtx.Dp(2))
		return D{Size: gtx.Constraints.Min}
	}

	// 右侧描述条
	style_desc := func(gtx C) D {
		w, h := self.desc_size.X+gtx.Dp(60), gtx.Dp(40)
		self.desc_width_tar = Ternary(self.clickable_icon.Hovered(), float64(w+gtx.Dp(1)), 0)

		if self.desc_width_cur == 0 {
			return D{}
		}

		desc_clip := clip.RRect{Rect: image.Rect(0, 0, int(self.desc_width_cur), h)}.Push(gtx.Ops)
		dims := layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				draw_rectangle(gtx, image.Pt(w, h), ele_bg, radius)
				return D{}
			}),
			layout.Stacked(func(gtx C) D {
				offset := op.Offset(image.Pt((w-self.desc_size.X)/2, (h-self.desc_size.Y)/2)).Push(gtx.Ops)
				l := material.Label(self.Style.Theme, unit.Sp(20), self.text())
				l.Color = ele_fg
				l.MaxLines = 1
				l.Layout(gtx)

				offset.Pop()
				return D{}
			}),
		)
		desc_clip.Pop()

		return dims
	}

	dims := layout.Flex{}.Layout(gtx,
		layout.Rigid(style_bar),
		layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
		layout.Rigid(style_btn),
		layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
		layout.Rigid(style_desc),
	)
	self.Style.AniSys.Ease64(&self.desc_width_cur, 1, &self.desc_width_cur, self.desc_width_tar, 0.4)
	self.Style.AniSys.Ease64(&self.bar_height_cur, 0.1, &self.bar_height_cur, self.bar_height_tar, 0.4)

	return dims
}

// 设置回调函数
func (self *SideBarButton) API_set_callfn(fn func()) *SideBarButton {
	self.fn = fn
	return self
}

// 设置文本
func (self *SideBarButton) API_set_text(text func() string) *SideBarButton {
	self.text = text
	self.is_redraw_desc = true
	return self
}

// 设置激活
func (self *SideBarButton) API_set_active(is_active bool) *SideBarButton {
	if !self.is_active && is_active {
		self.bar_height_cur = 0.0
	}
	self.is_active = is_active
	return self
}

// 设置是否允许激活
func (self *SideBarButton) API_set_allow_active(is_allow_active bool) *SideBarButton {
	self.is_allow_active = is_allow_active
	return self
}

// 设置状态
func (self *SideBarButton) API_set_state(state modeSideBarButton) *SideBarButton {
	self.state = state
	return self
}

// 得到状态
func (self *SideBarButton) API_get_state() modeSideBarButton {
	return self.state
}

// 尺寸更新
func (self *SideBarButton) API_resize() {
	self.API_set_text(self.text)
}

// -----------------------------------------------------------------------

// 宽侧边栏按钮
type SideBarWideButton struct {
	Style *Style

	id   string
	icon string
	text func() string

	bar_height_tar float64
	bar_height_cur float64

	clickable *widget.Clickable

	is_active bool

	callback func(string)
}

func NewSideBarWideButton(style *Style, id string) *SideBarWideButton {
	self := SideBarWideButton{
		Style:     style,
		is_active: true,
		id:        id,

		icon: "",
		text: func() string { return "" },

		clickable: &widget.Clickable{},
		callback:  func(s string) {},
	}
	return &self
}

func (self *SideBarWideButton) Layout(gtx C) D {
	if self.clickable.Clicked(gtx) {
		self.callback(self.id)
	}

	is_hoverd := self.clickable.Hovered()

	self.bar_height_tar = Ternary(self.is_active, 30.0, 0.0)
	// self.bar_height_cur = Ternary(self.is_active, self.bar_height_tar, 0.0)
	self.Style.AniSys.Ease64(&self.bar_height_cur, 1, &self.bar_height_cur, self.bar_height_tar, 0.4)

	fg := self.Style.Palette.Fg_1
	bg := Ternary(self.is_active, self.Style.Palette.Bg_2, self.Style.Palette.Bg_1)
	bg_bar := self.Style.Palette.Bg_3
	// fg_icon := Ternary(self.is_active, self.Style.Palette.Bg_3, self.Style.Palette.Fg_1)

	// 左侧指示条
	fn_bar := func(gtx C) D {
		width := gtx.Dp(4)
		gtx.Constraints.Min.X = width
		if !self.is_active {
			return D{Size: gtx.Constraints.Min}
		}
		draw_rectangle_2(gtx, 0, -gtx.Dp(unit.Dp(self.bar_height_cur/2)), width, gtx.Dp(unit.Dp(self.bar_height_cur)), bg_bar, gtx.Dp(2))
		return D{Size: gtx.Constraints.Min}
	}

	// 右侧按钮(图标_文本)
	label_icon := material.H6(self.Style.Theme, self.icon)
	label_text := material.H6(self.Style.Theme, self.text())
	label_icon.Color = fg
	label_text.Color = fg

	// label_icon.Font.Style = Ternary(self.is_active, font.Italic, font.Regular)
	label_icon.Font.Style = Ternary(is_hoverd, font.Italic, label_icon.Font.Style)

	btn := material.Button(self.Style.Theme, self.clickable, "")
	btn.Background = bg

	fn_btn := func(gtx C) D {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(btn.Layout),
			layout.Stacked(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(label_icon.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(label_text.Layout),
					layout.Flexed(1, layout.Spacer{}.Layout),
				)
			}),
		)

	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(fn_bar),
		layout.Rigid(layout.Spacer{Width: unit.Dp(2)}.Layout),
		layout.Flexed(1, fn_btn),
	)
}

func (self *SideBarWideButton) API_set_icon(icon string) *SideBarWideButton {
	self.icon = icon
	return self
}

func (self *SideBarWideButton) API_set_text(text func() string) *SideBarWideButton {
	self.text = text
	return self
}

func (self *SideBarWideButton) API_set_active(state bool) *SideBarWideButton {
	self.is_active = state
	return self
}

func (self *SideBarWideButton) API_set_callback(fn func(string)) *SideBarWideButton {
	self.callback = fn
	return self
}
