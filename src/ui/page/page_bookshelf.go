package page

import (
	rumia "app/src"
	"app/src/ui"
)

type PageBookshelf struct {
	Style       *ui.Style
	page_router *PageRouter

	bookcards []*ui.BookCard
	bookshelf ui.BookShelf
}

func NewPageShelf(style *ui.Style, page_router *PageRouter) *PageBookshelf {
	self := PageBookshelf{
		Style:       style,
		page_router: page_router,
		bookshelf:   *ui.NewBookShelf(style),
		bookcards:   []*ui.BookCard{},
	}

	self.bookshelf.
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(rumia.BookInfo{Title: "miku miku oo ee"})

	self.bookshelf.API_set_callback(func(book_info rumia.BookInfo) {
		self.page_router.API_get_manager("page").API_set_cur_page("viewer_pdf")
	})

	return &self

}

func (self *PageBookshelf) Layout(gtx C) D {
	return self.bookshelf.Layout(gtx)
}

func (self *PageBookshelf) OnEnter() {}
func (self *PageBookshelf) OnExit()  {}
func (self *PageBookshelf) Destory() {}
