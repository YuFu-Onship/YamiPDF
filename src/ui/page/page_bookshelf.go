package page

import "app/src/ui"

type PageBookshelf struct {
	Style     *ui.Style
	bookcards []*ui.BookCard
	bookshelf ui.BookShelf
}

func NewPageShelf(style *ui.Style) *PageBookshelf {
	self := PageBookshelf{
		Style:     style,
		bookshelf: *ui.NewBookShelf(style),
		bookcards: []*ui.BookCard{},
	}

	self.bookshelf.
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"}).
		API_add_book(ui.BookInfo{Title: "miku miku oo ee"})

	return &self
}

func (self *PageBookshelf) Layout(gtx C) D {
	return self.bookshelf.Layout(gtx)
}

func (self *PageBookshelf) OnEnter() {}
func (self *PageBookshelf) OnExit()  {}
