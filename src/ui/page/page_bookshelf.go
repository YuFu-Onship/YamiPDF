package page

import (
	rumia "app/src"
	"app/src/ui"
	"log"
)

type PageBookshelf struct {
	Style     *ui.Style
	bookcards []*ui.BookCard
	bookshelf ui.BookShelf

	bridge_ui ui.BridgeUI

	root_manager *PageManager
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
		if self.root_manager != nil {
			self.root_manager.API_get_last_manager().API_set_cur_page("pdf_viewer")
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

func (self *PageBookshelf) API_set_parent_manager(pm interface{}) {
	log.Println("inst books")
	self.root_manager = pm.(*PageManager)
}
