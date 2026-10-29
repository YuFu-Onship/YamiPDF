package page

import (
	rumia "app/src"
	"app/src/ui"

	"gioui.org/layout"
	"gioui.org/unit"
)

type PageBookshelf struct {
	Style       *ui.Style
	page_router *PageRouter

	book_cards []*ui.BookCard
	bookshelf  ui.BookShelf

	list             ui.CustomListStruct
	width_spacer_tar float64
	width_spacer_cur float64
	is_apply_pos_now bool
	callback         func(book_info rumia.BookInfo)
}

func NewPageShelf(style *ui.Style, page_router *PageRouter) *PageBookshelf {
	self := PageBookshelf{
		Style:       style,
		page_router: page_router,
		bookshelf:   *ui.NewBookShelf(style),
		book_cards:  []*ui.BookCard{},

		list:     ui.CustomListStruct{},
		callback: func(book_info rumia.BookInfo) {},

		width_spacer_tar: 0.0,
		width_spacer_cur: 0.0,
		is_apply_pos_now: false,
	}
	self.list.Axis = layout.Vertical

	self.book_cards = append(self.book_cards,
		ui.NewBookCard(style).
			API_set_book_info(rumia.BookInfo{Title: "miku miku oo ee"}).
			API_set_callfn(self.callback_card))
	self.book_cards = append(self.book_cards,
		ui.NewBookCard(style).
			API_set_book_info(rumia.BookInfo{Title: "miku miku oo ee"}).
			API_set_callfn(self.callback_card))
	self.book_cards = append(self.book_cards,
		ui.NewBookCard(style).
			API_set_book_info(rumia.BookInfo{Title: "miku miku oo ee"}))
	self.book_cards = append(self.book_cards,
		ui.NewBookCard(style).
			API_set_book_info(rumia.BookInfo{Title: "miku miku oo ee"}))
	self.book_cards = append(self.book_cards,
		ui.NewBookCard(style).
			API_set_book_info(rumia.BookInfo{Title: "miku miku oo ee"}))

	self.bookshelf.API_set_callback(func(book_info rumia.BookInfo) {
		self.page_router.API_get_manager("page").API_set_cur_page("viewer_pdf")
	})

	return &self
}

func (self *PageBookshelf) Layout(gtx C) D {
	total_books := len(self.book_cards)
	if total_books == 0 {
		return layout.Dimensions{}
	}

	card_width := gtx.Dp(150 + 8)
	cols := max((gtx.Constraints.Max.X-gtx.Dp(24))/card_width, 1)
	row_count := (total_books + cols - 1) / cols

	row_children := make([][]layout.FlexChild, row_count)

	// 每行左侧的间隙,/2时居中, /4时靠左
	self.width_spacer_tar = float64(max(0, (gtx.Constraints.Max.X-cols*card_width)/2))
	self.width_spacer_tar = Ternary(row_count <= 1, 0, self.width_spacer_tar)
	if self.is_apply_pos_now {
		self.is_apply_pos_now = false
		self.width_spacer_cur = self.width_spacer_tar
	} else {
		self.Style.AniSys.Ease64(&self.width_spacer_cur, 1, &self.width_spacer_cur, self.width_spacer_tar, 0.3)
	}

	for row := range row_count {
		index_start := row * cols
		index_end := min(index_start+cols, total_books)

		for i := index_start; i < index_end; i += 1 {
			book_card := self.book_cards[i]

			// 左侧间距
			if i == index_start {
				row_children[row] = append(row_children[row], layout.Rigid(ui.Spacer(int(self.width_spacer_cur), 0)))
			}

			// 中间书籍卡片布局
			row_children[row] = append(row_children[row], layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, book_card.Layout)
			}))

			// 右侧间距
			if i == index_end-1 {
				row_children[row] = append(row_children[row], layout.Flexed(1, ui.FlexerX()))
			}
		}
	}
	list := ui.CustomList(self.Style.Theme, &self.list)
	return list.Layout(gtx, row_count, func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, row_children[i]...)
	})

	// return self.bookshelf.Layout(gtx)
}

func (self *PageBookshelf) callback_card(book_info rumia.BookInfo) {
	self.page_router.API_get_manager("page").API_set_cur_page("viewer_pdf")
}

func (self *PageBookshelf) OnEnter() {
	self.is_apply_pos_now = true
}
func (self *PageBookshelf) OnExit()  {}
func (self *PageBookshelf) Destory() {}
