package page

import "app/src/ui"

type PageBookshelf struct {
	Style *ui.Style
}

func NewPageShelf(style *ui.Style) *PageBookshelf {
	self := PageBookshelf{
		Style: style,
	}
	return &self
}

func (self *PageBookshelf) Layout(gtx C) D {
	return D{}
}

func (self *PageBookshelf) OnEnter() {}
func (self *PageBookshelf) OnExit()  {}
