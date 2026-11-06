package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// 工具面板
type ToolPanel struct {
	Style *Style
	List  *CustomListStruct

	children []layout.FlexChild

	Width     float64
	width_tar float64
	width_cur float64
	is_expand bool
}

func NewToolPanel(style *Style) *ToolPanel {
	self := ToolPanel{
		Style: style,
		List:  &CustomListStruct{},

		children: []layout.FlexChild{},

		Width:     0.0,
		width_tar: 0.0,
		width_cur: 0.0,
		is_expand: false,
	}
	self.List.Axis = layout.Vertical
	return &self
}

func (self *ToolPanel) Layout(gtx C) D {
	gtx.Constraints.Max.X = gtx.Dp(unit.Dp(self.width_cur))
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	self.width_tar = Ternary(self.is_expand, self.Width, 0)
	dims := layout.Flex{}.Layout(gtx)

	self.Style.AniSys.Ease64(&self.width_cur, 0.1, &self.width_cur, self.width_tar, 0.4)

	return dims
}

func (self *ToolPanel) API_set_width(width float64) *ToolPanel {
	self.Width = width
	return self
}

func (self *ToolPanel) API_set_state(is_expand bool) *ToolPanel {
	self.is_expand = is_expand
	return self
}

func (self *ToolPanel) API_get_state() bool {
	return self.is_expand
}
