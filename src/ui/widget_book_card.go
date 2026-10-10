package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 书架界面中书的控件

// 书籍卡片
type BookCard struct {
	Style *Style

	img        *image.RGBA
	img_widget *widget.Image
	clickable  *widget.Clickable

	book_percent float64     // 阅读百分比
	book_title   string      // 书本名称
	book_path    string      // 书本路径
	book_type    string      // 书本格式
	book_cover   *image.RGBA // 封面

	extend_bar         *ExtendBar
	extend_bar_percent *ExtendBar
}

// 书本卡片
func NewBookCard(style *Style) *BookCard {
	self := BookCard{
		clickable: &widget.Clickable{},
		Style:     style,
	}
	self.extend_bar = NewExtendBar(style)
	self.extend_bar_percent = NewExtendBar(style)

	self.extend_bar.API_set_text("新编地图学教程新编地图学教程新编地图学教程")
	self.extend_bar_percent.API_set_text("PDF  |  12%")

	return &self
}

func (self *BookCard) Layout(gtx C) D {
	gtx.Constraints.Max = image.Pt(gtx.Dp(150), gtx.Dp(225))
	gtx.Constraints.Min = gtx.Constraints.Max

	if self.clickable.Clicked(gtx) {
	}

	is_hover := self.clickable.Hovered()

	self.extend_bar.is_extend = Ternary(is_hover, true, false)
	self.extend_bar_percent.is_extend = Ternary(is_hover, true, false)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			draw_rectangle(gtx, gtx.Constraints.Max, COLOR_DARK_GRAY, 4)
			return D{Size: gtx.Constraints.Max}
		}),
		layout.Expanded(func(gtx C) D {
			btn := material.Button(self.Style.Theme, self.clickable, "")
			btn.Background.A = 0
			return btn.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			// gtx.Constraints.Min = gtx.Constraints.Max
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Flexed(1, FlexerY()),
				layout.Rigid(self.extend_bar.Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(self.extend_bar_percent.Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			)
		}),
	)
}

// 内含文本的伸缩条
type ExtendBar struct {
	Style  *Style
	Text   string
	Height float64
	Width  float64

	width_cur float64
	width_tar float64

	is_extend  bool
	is_refresh bool
}

func NewExtendBar(style *Style) *ExtendBar {
	self := ExtendBar{
		Style: style,

		is_extend:  false,
		is_refresh: true,

		Height:    30,
		Width:     0.0,
		width_cur: 0.0,
		width_tar: 0.0,
	}
	return &self
}

func (self *ExtendBar) Layout(gtx C) D {
	text_width := gtx.Dp(130)
	if self.is_refresh {
		self.is_refresh = false

		func(gtx C) {
			gtx.Constraints.Max.X = text_width

			label := material.Label(self.Style.Theme, unit.Sp(16), self.Text)
			label.MaxLines = 2

			macro := op.Record(gtx.Ops)
			l := label.Layout(gtx)
			macro.Stop()

			self.Width = float64(l.Size.X + gtx.Dp(16))
			self.Height = float64(l.Size.Y)
		}(gtx)

		// self.Width = math.Min(float64(l.Size.X+gtx.Dp(16)), float64(gtx.Dp(140)))
	}

	self.width_tar = Ternary(self.is_extend, self.Width, 0)

	rect := image.Rect(0, 0, int(self.width_cur), gtx.Dp(unit.Dp(self.Height)))
	trans_clip := clip.Rect(rect).Push(gtx.Ops)

	// draw_rectangle(gtx, image.Pt(int(self.Width), int(self.Height)), self.Style.Palette.Bg_3, 4)

	paint.FillShape(gtx.Ops, self.Style.Palette.Bg_3, clip.RRect{Rect: image.Rectangle{Max: image.Pt(int(self.Width), int(self.Height))},
		SW: 0,
		NW: 0,
		SE: gtx.Dp(4),
		NE: gtx.Dp(4),
	}.Op(gtx.Ops))

	label_text := material.Label(self.Style.Theme, unit.Sp(16), self.Text)
	label_text.MaxLines = 2
	label_text.Color = self.Style.Palette.Fg_3

	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
		layout.Rigid(func(gtx C) D { gtx.Constraints.Max.X = text_width; return label_text.Layout(gtx) }),
		layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
	)

	trans_clip.Pop()

	self.Style.AniSys.Ease64(&self.width_cur, 1, &self.width_cur, self.width_tar, 0.4)

	return dims
}

func (self *ExtendBar) API_set_text(text string) *ExtendBar {
	self.Text = text
	self.is_refresh = true
	return self
}

// 含有一组标签的伸缩条
type ExtendTagBar struct {
	Style *Style
	tags  []string

	Height    float64
	Width     float64
	width_tar float64
	width_cur float64
	spacer    float64

	is_extend bool
	is_recalc bool
}

func NewExtendTagBar(style *Style) *ExtendTagBar {
	self := ExtendTagBar{
		Style: style,
	}
	return &self
}

func (self *ExtendTagBar) Layout(gtx C) D {
	return layout.Flex{}.Layout(gtx)
}

// 覆盖替换标签
func (self *ExtendTagBar) API_set_text(tags []string) *ExtendTagBar {
	self.tags = tags
	return self
}

// 创建标签元素
func (self *ExtendTagBar) build_tags(gtx C) D {
	radius_west := gtx.Dp(4)
	radius_east := gtx.Dp(4)

	children := []layout.FlexChild{}
	for i, tag := range self.tags {
		radius_west = Ternary(i == 0, 0, radius_west)
		dims := func(gtx C) D {
			paint.FillShape(gtx.Ops, self.Style.Palette.Bg_3, clip.RRect{Rect: image.Rectangle{Max: image.Pt(int(self.Width), int(self.Height))},
				SW: radius_west,
				NW: radius_west,
				SE: radius_east,
				NE: radius_east,
			}.Op(gtx.Ops))

			return material.Label(self.Style.Theme, unit.Sp(16), tag).Layout(gtx)
		}
		children = append(children, layout.Rigid(dims))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}
