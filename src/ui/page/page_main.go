package page

import (
	"app/assets"
	"app/src/ui"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageMain struct {
	Style        *ui.Style
	page_manager PageManager
	side_bar     *ui.SideBar
}

func NewPageMain(style *ui.Style) *PageMain {
	self := PageMain{
		Style: style,
	}

	self.side_bar = ui.NewSideBar().
		// API_add_btn("home", true, true, NewSideBarButton(self.Style, Icon_home).API_set_text(func() string { return GetWord("home", self.Style.Lang) }).API_set_callfn(func() {})).
		API_add_btn("book_shelf", true, true, ui.NewSideBarButton(self.Style, assets.Icon_books).API_set_text(func() string { return ui.GetWord("bookshelf", self.Style.Lang) }).API_set_callfn(func() {})).
		API_add_btn("color_mode", false, false, ui.NewSideBarButton(self.Style, assets.Icon_moon).API_set_text(func() string { return ui.GetWord("color_mode", self.Style.Lang) }).API_set_callfn(self.Style.Toggle)).
		API_add_btn("setting", false, true, ui.NewSideBarButton(self.Style, assets.Icon_setting).API_set_text(func() string { return ui.GetWord("setting", self.Style.Lang) }).API_set_callfn(func() {}))
	self.side_bar.API_set_active("book_shelf")

	self.page_manager = *NewPageManager(style)
	self.page_manager.API_add_page("book_shelf", NewPageShelf(style))

	return &self
}

func (self *PageMain) Layout(gtx C) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Bottom: unit.Dp(12), Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, self.side_bar.Layout)
		}),
	)

}

func (self *PageMain) OnEnter() {}
func (self *PageMain) OnExit()  {}
