package ui

import (
	"image"
	"image/color"
	"slices"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DropDown struct {
	Header W
	Panel  W

	Text      string // 经典边框文本,需要显示
	Is_fill   bool   // 是否填充绘制背景
	Is_border bool   // 是否绘制边框
	Is_expand bool   // 是否展开

	Bg color.NRGBA // 背景颜色
	Fg color.NRGBA // 边框颜色

	height_tar   float64
	height_cur   float64
	height_panel float64

	clickable_header *widget.Clickable
}

func NewDropDown() *DropDown {
	self := DropDown{
		Is_expand: false,
		Is_fill:   true,
		Is_border: true,

		height_tar:   0,
		height_cur:   0,
		height_panel: 0,

		clickable_header: &widget.Clickable{},
	}

	return &self
}

// 下拉栏布局, 布局中只能包含绘制逻辑
//   - header 顶部布局
//   - panel 下拉内容布局
func (self *DropDown) Layout(gtx C, style *Style, header W, panel W) D {

	if self.clickable_header.Clicked(gtx) {
		self.Is_expand = !self.Is_expand
	}

	self.height_tar = Ternary(self.Is_expand, self.height_panel+float64(unit.Dp(4)), -float64(gtx.Dp(3)))

	btn_header := material.Button(style.Theme, self.clickable_header, "")
	btn_header.Background.A = 0

	icon_label := material.Label(style.Theme, unit.Sp(20), Ternary(self.Is_expand, Icon_up, Icon_down))
	icon_label.Font.Style = Ternary(self.clickable_header.Hovered(), font.Italic, font.Regular)

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					if self.Is_fill {
						draw_rectangle(gtx, gtx.Constraints.Min, self.Bg, gtx.Dp(4))
					}
					return D{Size: gtx.Constraints.Min}
				}),

				layout.Expanded(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, FlexerX()),
							layout.Rigid(icon_label.Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						)
					})
				}),

				layout.Stacked(func(gtx C) D {
					border := Border{
						Width: unit.Dp(1),
						Bg:    self.Bg,
						Fg:    self.Fg,
					}

					if self.Is_border {
						border.Text = self.Text
					} else {
						border.Fg.A = 0
					}

					return border.Layout(gtx, style, header)
				}),

				layout.Expanded(btn_header.Layout),
			)
		}),

		layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),

		layout.Rigid(func(gtx C) D {
			macro := op.Record(gtx.Ops)
			p := panel(gtx)
			record := macro.Stop()

			self.height_panel = float64(p.Size.Y)

			clip := clip.RRect{Rect: image.Rect(-gtx.Dp(1), gtx.Dp(-1), p.Size.X+gtx.Dp(2), int(self.height_cur)+gtx.Dp(2))}.Push(gtx.Ops)
			self.draw_backgroud(gtx, p.Size, gtx.Dp(8))
			record.Add(gtx.Ops)
			clip.Pop()

			return D{Size: image.Pt(p.Size.X, int(self.height_cur))}
		}),
	)

	style.AniSys.Ease64(&self.height_cur, 1, &self.height_cur, self.height_tar, 0.5)

	return dims
}

// 绘制背景和边框
func (self *DropDown) draw_backgroud(gtx C, size image.Point, radius int) {
	if self.Is_fill {
		draw_rectangle(gtx, size, self.Bg, radius)
	}
	if self.Is_border {
		draw_rectangle_line(gtx, size, self.Fg, radius, float32(gtx.Dp(1)))
	}
}

// -------------------------------------------------------------------------------------------------------------

// 下拉单选
type DropDownChoice struct {
	choices      []string
	choice_cur   string            // 当前的值
	choice_index int               // 当前的值的索引
	choices_fn   map[string]func() // 选项对应的回调函数

	clickables []widget.Clickable

	list CustomListStruct

	call_fn func(string)

	drop_down *DropDown
	Style     *Style
}

func NewDropDownChoice() *DropDownChoice {
	self := DropDownChoice{
		drop_down:  NewDropDown(),
		choices:    []string{},
		clickables: []widget.Clickable{},
		list:       CustomListStruct{},
		call_fn:    func(string) {},
	}

	self.list.Axis = layout.Vertical
	self.choices_fn = map[string]func(){}

	return &self
}

// 单选布局
func (self *DropDownChoice) Layout(gtx C, style *Style) D {
	self.Style = style
	self.drop_down.Fg = self.Style.Palette.Fg_1
	self.drop_down.Bg = self.Style.Palette.Bg_1

	// self.drop_down.Fg = COLOR_BLUE
	// self.drop_down.Bg = COLOR_DARK_BLUE

	gtx.Constraints.Max.X = gtx.Dp(200)
	dims := self.drop_down.Layout(gtx, self.Style,
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			l := material.Label(self.Style.Theme, unit.Sp(20), GetWord(self.choice_cur, self.Style.Lang))
			return l.Layout(gtx)
		},

		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(2), Right: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
				return self.list.Layout(gtx, len(self.choices), func(gtx C, index int) D {
					return layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1)}.Layout(gtx, func(gtx C) D {
						if self.clickables[index].Clicked(gtx) {
							self.choice_index = index
							self.choice_cur = self.choices[index]

							self.call_fn(self.choice_cur)

							if fn, exist := self.choices_fn[self.choice_cur]; exist {
								fn()
							}
						}
						return self.build_choice_btn(gtx, &self.clickables[index], self.choices[index], true)
					})
				})
			})
		},
	)
	return dims
}

// 构建没有边框时的按钮

// 构建单选按钮
func (self *DropDownChoice) build_choice_btn(gtx C, clickable *widget.Clickable, text string, is_check bool) D {
	gtx.Constraints.Min.Y = gtx.Dp(40)
	gtx.Constraints.Max.Y = gtx.Dp(40)

	btn := material.Button(self.Style.Theme, clickable, "")
	btn.Background.A = 0
	btn.CornerRadius = unit.Dp(4)

	label_icon := material.Label(self.Style.Theme, unit.Sp(20), Icon_check)
	label_icon.Color.A = Ternary(self.choices[self.choice_index] == text, label_icon.Color.A, 0)

	label_text := material.Label(self.Style.Theme, unit.Sp(20), GetWord(text, self.Style.Lang))
	label_text.MaxLines = 1

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(btn.Layout),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
				layout.Rigid(label_icon.Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
				layout.Rigid(label_text.Layout),
				layout.Flexed(1, FlexerX()),
			)
		}),
	)
}

// 添加选项,内部会从语言索引表上获取
func (self *DropDownChoice) API_add_choice(choice string) *DropDownChoice {
	self.choices = append(self.choices, choice)
	self.clickables = append(self.clickables, widget.Clickable{})
	return self
}

// 添加选项以及每个选项附带的回调函数
func (self *DropDownChoice) API_add_choice_fn(choice string, fn func()) *DropDownChoice {
	self.choices = append(self.choices, choice)
	self.clickables = append(self.clickables, widget.Clickable{})
	self.choices_fn[choice] = fn
	return self
}

// 添加/设置回调函数
func (self *DropDownChoice) API_set_callfn(fn func(choice string)) *DropDownChoice {
	self.call_fn = fn
	return self
}

// 设置默认值
func (self *DropDownChoice) API_set_value(choice string) *DropDownChoice {
	self.choice_cur = Ternary(slices.Contains(self.choices, choice), choice, "")
	return self
}

// 设置边框文本
func (self *DropDownChoice) API_set_border_text(text string) *DropDownChoice {
	self.drop_down.Text = text
	return self
}

// 得到当前值
func (self *DropDownChoice) API_get_cur_choice() string {
	return self.choice_cur
}
