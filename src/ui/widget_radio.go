package ui

import (
	"app/assets"
	"image"
	"image/color"
	"log"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CustomRadio struct {
	Style *Style

	keys          []string
	value         string // 当前值
	texts         map[string]func() string
	detecion_area *widget.Clickable // 光标检测区域,用于展开检测
	clickables    map[string]*widget.Clickable

	spacer   float32
	callback func(key string)
}

func NewCustomRadio(style *Style) *CustomRadio {
	self := CustomRadio{
		Style: style,

		value: "",
		keys:  []string{},
		texts: map[string]func() string{},

		detecion_area: &widget.Clickable{},
		clickables:    make(map[string]*widget.Clickable),

		spacer:   4,
		callback: func(key string) {},
	}
	return &self
}

func (self *CustomRadio) Layout(gtx C) D {
	children := []layout.FlexChild{}

	for i, key := range self.keys {
		if self.clickables[key].Clicked(gtx) {
			self.value = key
			self.callback(key)
			log.Println(gtx.Constraints.Max.Y)
		}

		if i != 0 {
			children = append(children, layout.Rigid(layout.Spacer{Width: unit.Dp(self.spacer)}.Layout))
		}
		children = append(children, layout.Rigid(func(gtx C) D { return self.build_btn(gtx, key) }))
	}

	return self.detecion_area.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
	})
}

func (self *CustomRadio) build_btn(gtx C, key string) D {
	is_active := self.value == key

	btn := material.Button(self.Style.Theme, self.clickables[key], "")
	btn.Background = Ternary(is_active, self.Style.Palette.Bg_3, self.Style.Palette.Bg_2)
	btn.Inset = layout.Inset{}

	label_text := material.Body1(self.Style.Theme, self.texts[key]())
	label_text.Color = Ternary(is_active, self.Style.Palette.Fg_3, self.Style.Palette.Fg_2)

	fn_text := func(gtx C) D {
		gtx.Constraints.Min.Y = gtx.Dp(30)
		gtx.Constraints.Max.Y = gtx.Dp(30)

		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(label_text.Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			)
		})
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(btn.Layout),
		layout.Stacked(fn_text),
	)
}

func (self *CustomRadio) API_set_text(key string, text func() string) *CustomRadio {
	self.keys = append(self.keys, key)
	self.texts[key] = text
	self.clickables[key] = &widget.Clickable{}
	return self
}
func (self *CustomRadio) API_set_callback(fn func(key string)) *CustomRadio {
	self.callback = fn
	return self
}

func (self *CustomRadio) API_set_value(key string) *CustomRadio {
	self.value = key
	return self
}

// ----------------------------------------------------------------------------------------------

type CustomPaletteRadio struct {
	Style *Style

	value      string
	keys       []string
	texts      map[string]func() string
	colors     map[string]func() color.NRGBA
	clickables map[string]*widget.Clickable

	spacer   float32
	callback func(key string)
}

func NewCustomPaletteRadio(style *Style) *CustomPaletteRadio {
	self := CustomPaletteRadio{
		Style: style,

		value:      "",
		keys:       []string{},
		texts:      make(map[string]func() string),
		colors:     make(map[string]func() color.NRGBA),
		clickables: make(map[string]*widget.Clickable),

		spacer:   4,
		callback: func(key string) {},
	}
	return &self
}

func (self *CustomPaletteRadio) Layout(gtx C) D {
	children := []layout.FlexChild{}
	for i, key := range self.keys {
		if self.clickables[key].Clicked(gtx) {
			self.value = key
			self.callback(key)
			log.Println(gtx.Constraints.Max.Y)
		}

		if i != 0 {
			children = append(children, layout.Rigid(layout.Spacer{Width: unit.Dp(self.spacer)}.Layout))
		}
		children = append(children, layout.Rigid(func(gtx C) D { return self.build_btn(gtx, key) }))
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
}

func (self *CustomPaletteRadio) build_btn(gtx C, key string) D {
	is_active := key == self.value

	btn := material.Button(self.Style.Theme, self.clickables[key], assets.Icon_check)
	btn.Color = self.Style.Palette.Fg_3
	btn.Color.A = Ternary(is_active, btn.Color.A, 0)
	btn.Background = self.colors[key]()
	btn.Inset = layout.Inset{}
	btn.TextSize = unit.Sp(20)
	fn_btn := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(30), gtx.Dp(30))
		gtx.Constraints.Max = image.Pt(gtx.Dp(30), gtx.Dp(30))
		return btn.Layout(gtx)
	}

	return fn_btn(gtx)
}

func (self *CustomPaletteRadio) API_add_choice(key string, c func() color.NRGBA) *CustomPaletteRadio {
	self.keys = append(self.keys, key)
	self.colors[key] = c
	self.clickables[key] = &widget.Clickable{}
	return self
}

func (self *CustomPaletteRadio) API_set_callback(fn func(key string)) *CustomPaletteRadio {
	self.callback = fn
	return self
}

func (self *CustomPaletteRadio) API_set_value(key string) *CustomPaletteRadio {
	self.value = key
	return self
}
