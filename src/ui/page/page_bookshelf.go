package page

import (
	rumia "app/src"
	"app/src/ui"
)

type PageBookshelf struct {
	Style     *ui.Style
	bookcards []*ui.BookCard
	bookshelf ui.BookShelf

	bridge_ui ui.BridgeUI
}

func NewPageShelf(style *ui.Style) *PageBookshelf {
	self := PageBookshelf{
		Style:     style,
		bookshelf: *ui.NewBookShelf(style),
		bookcards: []*ui.BookCard{},
	}

	self.bookshelf.
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"})

	self.bookshelf.API_set_callback(func(book_info rumia.BookInfo) {
		if self.bridge_ui != nil {
			self.bridge_ui.Set_cur_page("pdf_viewer")
		}
	})

	return &self

}

func (self *PageBookshelf) Layout(gtx C) D {
	return self.bookshelf.Layout(gtx)
}

func (self *PageBookshelf) OnEnter() {}
func (self *PageBookshelf) OnExit()  {}
func (self *PageBookshelf) Destory() {}

func (self *PageBookshelf) API_set_bridge_ui(bridge_ui ui.BridgeUI) {
	self.bridge_ui = bridge_ui
}
