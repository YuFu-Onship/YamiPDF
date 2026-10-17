package ui

import (
	"app/assets"
	rumia "app/src"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 书架排版控件
type BookShelf struct {
	Style *Style
	list  CustomListStruct

	book_cards []*BookCard

	callback_fn func(rumia.BookInfo)

	width_spacer_tar float64
	width_spacer_cur float64
}

func NewBookShelf(style *Style) *BookShelf {
	self := BookShelf{
		Style: style,
		list:  CustomListStruct{},

		book_cards: []*BookCard{},

		callback_fn: func(rumia.BookInfo) {},

		width_spacer_tar: 0,
		width_spacer_cur: 0,
	}

	self.list.Axis = layout.Vertical

	return &self
}

func (self *BookShelf) Layout(gtx layout.Context) layout.Dimensions {

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
	self.Style.AniSys.Ease64(&self.width_spacer_cur, 1, &self.width_spacer_cur, self.width_spacer_tar, 0.15)

	for row := range row_count {
		index_start := row * cols
		index_end := min(index_start+cols, total_books)

		for i := index_start; i < index_end; i += 1 {
			book_card := self.book_cards[i]

			// 左侧间距
			if i == index_start {
				// row_children[row] = append(row_children[row], layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout))
				row_children[row] = append(row_children[row], layout.Rigid(layout.Spacer{Width: unit.Dp(self.width_spacer_cur)}.Layout))
			}

			// 中间书籍卡片布局
			row_children[row] = append(row_children[row], layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, book_card.Layout)
			}))
			// row_children[row] = append(row_children[row], layout.Rigid(layout.Spacer{Width: unit.Dp(self.width_spacer_cur)}.Layout))

			// 右侧间距
			if i == index_end-1 {
				row_children[row] = append(row_children[row], layout.Flexed(1, FlexerX()))
			}
		}
	}
	list := CustomList(self.Style.Theme, &self.list)
	return list.Layout(gtx, row_count, func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, row_children[i]...)
	})
}

// 添加书籍
func (self *BookShelf) API_add_book(book_info rumia.BookInfo) *BookShelf {
	book_card := NewBookCard(self.Style).API_set_book_info(book_info)
	self.book_cards = append(self.book_cards, book_card)
	return self
}

// 设置回调函数
func (self *BookShelf) API_set_callback(fn func(rumia.BookInfo)) *BookShelf {
	self.callback_fn = fn
	for _, b := range self.book_cards {
		b.API_set_callfn(func(book_info rumia.BookInfo) {
			self.callback_fn(book_info)
		})
	}
	return self
}

// ------------------------------------------------------------------------------------------------

// 书籍卡片
type BookCard struct {
	Style *Style

	img        *image.RGBA
	img_widget *widget.Image
	clickable  *widget.Clickable

	book_info rumia.BookInfo

	extend_title *ExtendTagBar
	extend_tag   *ExtendTagBar

	file_type rumia.FileType

	callback_fn func(rumia.BookInfo) // 回调函数

	fg, bg color.NRGBA
}

// 书本卡片
func NewBookCard(style *Style) *BookCard {
	self := BookCard{
		clickable: &widget.Clickable{},
		Style:     style,
		file_type: rumia.FileType_PDF,

		callback_fn: func(rumia.BookInfo) {},
	}

	self.extend_tag = NewExtendTagBar(style)
	self.extend_title = NewExtendTagBar(style).API_set_max_lines(2)

	return &self
}

func (self *BookCard) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max = image.Pt(gtx.Dp(150), gtx.Dp(225))
	gtx.Constraints.Min = gtx.Constraints.Max

	// 回调函数
	if self.clickable.Clicked(gtx) {
		self.callback_fn(self.book_info)
	}

	is_hover := self.clickable.Hovered()

	self.extend_title.is_extend = Ternary(is_hover, true, false)
	self.extend_tag.is_extend = Ternary(is_hover, true, false)
	self.fg = Ternary(is_hover, self.Style.Palette.Fg_3, self.Style.Palette.Fg_2)
	self.bg = Ternary(is_hover, self.Style.Palette.Bg_3, self.Style.Palette.Bg_2)

	self.extend_title.API_set_color(self.fg, self.bg)
	self.extend_tag.API_set_color(self.fg, self.bg)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			draw_rectangle(gtx, gtx.Constraints.Max, self.Style.Palette.Bg_2, 4)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(self.Style.Theme, self.clickable, "")
			btn.Background.A = 0
			return btn.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(140)
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, FlexerY()),
				layout.Rigid(self.extend_title.Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(self.extend_tag.Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			)
		}),
	)
}

// 设置书本信息
func (self *BookCard) API_set_book_info(book_info rumia.BookInfo) *BookCard {
	self.book_info = book_info

	self.extend_title.API_set_text([]string{self.book_info.Title})
	self.extend_tag.API_set_text([]string{rumia.FileType_table[rumia.FileType_PDF], "12%"})

	self.extend_title.is_refresh = true
	self.extend_tag.is_refresh = true
	return self
}

// 设置显示的文件类型
func (self *BookCard) API_set_file_type(file_type rumia.FileType) *BookCard {
	self.book_info.FileType = file_type
	return self
}

// 设置回调函数
func (self *BookCard) API_set_callfn(callback_fn func(book_info rumia.BookInfo)) *BookCard {
	self.callback_fn = callback_fn
	return self
}

// ------------------------------------------------------------------------------------------------

// 含有一组标签的伸缩条
type ExtendTagBar struct {
	Style *Style
	tags  []string
	sizes []image.Point

	Height    float64
	Width     float64
	width_tar float64
	width_cur float64
	spacer    unit.Dp

	is_extend  bool
	is_refresh bool

	max_lines int
	font_size unit.Sp
	fg, bg    color.NRGBA
}

func NewExtendTagBar(style *Style) *ExtendTagBar {
	self := ExtendTagBar{
		Style:      style,
		is_refresh: true,
		is_extend:  false,
		spacer:     unit.Dp(4),
		max_lines:  1,
		font_size:  unit.Sp(15),
	}
	self.fg = self.Style.Palette.Fg_2
	self.fg = self.Style.Palette.Bg_2
	return &self
}

func (self *ExtendTagBar) Layout(gtx layout.Context) layout.Dimensions {
	if self.is_refresh {
		self.is_refresh = false
		self.sizes = []image.Point{}

		self.Width = 0

		for i, tag := range self.tags {
			label_text := material.Label(self.Style.Theme, self.font_size, tag)
			label_text.MaxLines = self.max_lines
			label_icon := material.Label(self.Style.Theme, unit.Sp(14), assets.Icon_vertical_line)

			macro := op.Record(gtx.Ops)
			lt := layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, label_text.Layout)
			li := label_icon.Layout(gtx)
			macro.Stop()

			self.sizes = append(self.sizes, lt.Size)

			self.Width += float64(lt.Size.X)                    // 文本宽度
			self.Width += Ternary(i > 0, float64(li.Size.X), 0) // 间隙宽度
			self.Height = float64(lt.Size.Y)                    // 文本高度
		}
		self.Width += float64(int(self.spacer) * (len(self.tags) - 1))
	}

	self.width_tar = Ternary(self.is_extend, self.Width, 0)
	self.width_cur = Ternary(self.is_extend, self.width_cur, 0)

	if self.width_cur == 0 && self.width_tar == 0 {
		self.Style.AniSys.Ease64(&self.width_cur, 1, &self.width_cur, self.width_tar, 0.4)
		return D{Size: image.Pt(0, gtx.Dp(unit.Dp(self.Height)))}
	}

	trans_clip := clip.Rect(image.Rect(0, 0, int(self.width_cur), int(self.Height))).Push(gtx.Ops)
	paint.FillShape(gtx.Ops, self.bg, clip.RRect{Rect: image.Rectangle{Max: image.Pt(int(self.Width), int(self.Height))}, SW: 0, NW: 0, SE: gtx.Dp(4), NE: gtx.Dp(4)}.Op(gtx.Ops))
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, self.build_tags(gtx)...)
	trans_clip.Pop()

	self.Style.AniSys.Ease64(&self.width_cur, 1, &self.width_cur, self.width_tar, 0.4)
	return layout.Dimensions{Size: dims.Size}
}

// 创建标签元素
func (self *ExtendTagBar) build_tags(gtx layout.Context) []layout.FlexChild {
	children := []layout.FlexChild{}
	for i, tag := range self.tags {
		if i > 0 {
			li := material.Label(self.Style.Theme, unit.Sp(14), assets.Icon_vertical_line)
			li.Color = self.fg
			li.Color.A = 128
			children = append(children, layout.Rigid(li.Layout))
		}

		label_text := material.Label(self.Style.Theme, self.font_size, tag)
		label_text.Color = self.fg
		label_text.MaxLines = self.max_lines
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, label_text.Layout)
		}))
	}
	return children
}

// 覆盖替换标签
func (self *ExtendTagBar) API_set_text(tags []string) *ExtendTagBar {
	self.tags = tags
	self.is_refresh = true
	return self
}

func (self *ExtendTagBar) API_set_max_lines(max_lines int) *ExtendTagBar {
	self.max_lines = max_lines
	self.is_refresh = true
	return self
}

// 设置字体大小
func (self *ExtendTagBar) API_set_font_size(font_size unit.Sp) *ExtendTagBar {
	self.font_size = font_size
	self.is_refresh = true
	return self
}

// 设置颜色
func (self *ExtendTagBar) API_set_color(fg, bg color.NRGBA) *ExtendTagBar {
	self.fg = fg
	self.bg = bg
	return self
}
