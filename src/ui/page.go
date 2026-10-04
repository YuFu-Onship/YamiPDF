package ui

import (
	"app/assets"
	_ "embed"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Page struct {
	Deco   *CustomDeco
	Style  *Style
	bridge TrunkBridge

	win_size      image.Point
	is_win_change bool
	appwindow     *app.Window
}

func NewPage(b TrunkBridge) *Page {
	theme := material.NewTheme()
	style := NewStyle(theme, Palette_koishi)
	style.SetColorMode(true)

	self := Page{
		bridge: b,
		Deco:   NewCustomDeco(&style),
		Style:  &style,
	}

	return &self
}

func (self *Page) Run() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Simple PDF"),
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.MinSize(unit.Dp(400), unit.Dp(300)),
			app.Decorated(false),
		)

		if err := self.draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func (self *Page) draw(w *app.Window) error {
	var ops op.Ops
	self.appwindow = w
	self.Style.SetColorMode(true)
	self.Style.Lang = EN

	// 加载文字
	if fonts, err := self.load_fonts(); err == nil {
		shaper := text.NewShaper(text.WithCollection(fonts))
		self.Style.Theme.Shaper = shaper
	}

	// pdf阅读器
	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Archives of Witch Hat Atelier (Kamome Shirahama).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Witch Hat Atelier, Vol.5 (Kamome Shirahama).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/单页.pdf"

	// doc_id, _ := self.bridge.Add_new_document(pdf_path)
	// pdfview := NewPDFView(doc_id, self.bridge, self)

	// 平滑列表
	// list := CustomListStruct{}
	// list.Axis = layout.Vertical
	// l := CustomList(self.Style.Theme, &list)

	// 开关控件
	sw := widget.Bool{}

	head_clickable := widget.Clickable{}

	// 侧边栏
	side_bar := NewSideBar().
		API_add_btn("home", true, true, NewSideBarButton(self.Style, Icon_home).API_set_text(func() string { return GetWord("home", self.Style.Lang) }).API_set_callfn(func() {})).
		API_add_btn("books", true, true, NewSideBarButton(self.Style, Icon_setting).API_set_text(func() string { return GetWord("bookshelf", self.Style.Lang) }).API_set_callfn(func() {})).
		API_add_btn("color_mode", false, false, NewSideBarButton(self.Style, Icon_moon).API_set_text(func() string { return GetWord("color_mode", self.Style.Lang) }).API_set_callfn(func() { self.Style.SetColorMode(!self.Style.isDark) })).
		API_add_btn("setting", false, true, NewSideBarButton(self.Style, Icon_setting).API_set_text(func() string { return GetWord("setting", self.Style.Lang) }).API_set_callfn(func() {}))
	side_bar.API_set_active("home")

	for {
		switch typ := w.Event().(type) {
		case app.DestroyEvent:
			return typ.Err

		case app.ViewEvent:

		case app.ConfigEvent:
			self.Deco.Deco.Maximized = typ.Config.Mode == app.Maximized || typ.Config.Mode == app.Fullscreen

		case app.FrameEvent:
			ops.Reset()

			self.is_win_change = typ.Size != self.win_size
			self.win_size = typ.Size

			if self.is_win_change {
				// side_btn.API_resize()
			}

			gtx := app.NewContext(&ops, typ)
			paint.Fill(gtx.Ops, self.Style.Palette.Bg_1)

			self.Style.Update(gtx)

			self.Deco.Actions(gtx, w)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(self.Deco.Layout),
				layout.Flexed(1, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: unit.Dp(12), Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, side_bar.Layout)
						}),

						layout.Rigid(func(gtx C) D {
							return Border{
								Text: "draw call",
								Bg:   self.Style.Palette.Bg_1,
								Fg:   self.Style.Palette.Fg_1,
							}.Layout(gtx, self.Style, func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Rigid(material.Label(self.Style.Theme, unit.Sp(20), "E:/Users/YUFU/Documents/Books").Layout),
									layout.Rigid(material.Switch(self.Style.Theme, &sw, "").Layout),
								)
							})
						}),

						layout.Rigid(func(gtx C) D {
							return DropDown{}.Layout(gtx, self.Style, &head_clickable,
								func(gtx C) D {
									l := material.Body1(self.Style.Theme, "texddddddddddddddddddddddddt")
									return l.Layout(gtx)
								},
								func(gtx C) D {
									l := material.Body2(self.Style.Theme, "texdsadssssssssssssssst")
									return l.Layout(gtx)
								},
							)
						}),

						// layout.Flexed(1, func(gtx C) D { return pdfview.Layout(gtx, self.Style) }),
						// 						layout.Flexed(1, func(gtx C) D {
						// 							l.Layout(gtx, 1000, func(gtx layout.Context, index int) layout.Dimensions {
						// 								// 渲染单个 Item 的 UI 逻辑
						// 								text := fmt.Sprintf("这是第 %d 个平滑滚动列表项", index+1)
						//
						// 								return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// 									return material.Body1(self.Style.Theme, text).Layout(gtx)
						// 								})
						// 							})
						// 							return D{}
						// 						}),
					)
				}),

				// layout.Flexed(1, func(gtx C) D { return pdfview.Layout(gtx, self.Style) }),
			)

			self.Style.AniSys.Run(gtx)
			typ.Frame(gtx.Ops)
		}
	}
}

// 接口部分 ------------------------------------------

// 切换颜色模式
func (self Page) ToggleTheme() {
	if self.Style.isDark {
		self.Style.SetColorMode(false)
	} else {
		self.Style.SetColorMode(true)
	}
}

// 窗口大小是否变动
//   - true 发生变化
//   - false 没有变化
func (self Page) IsWinSizeChange() bool {
	return self.is_win_change
}

func (self Page) WindowRefresh() {
	self.appwindow.Invalidate()
}

// 显示信息
func (self Page) ShowInfo(id any, text string) {}

// 更改标题栏文本
func (self Page) Set_deco_text(text string) {
	self.Deco.API_set_text(text)
	self.WindowRefresh()
}

// 更改软件标题文本
func (self *Page) Set_title_text(text string) {
	if self.appwindow != nil {
		self.appwindow.Option(app.Title(text))
	}
}

// 加载文字 --------------------------------------------

// 从字节码加载文字
func (self *Page) load_fonts() ([]font.FontFace, error) {
	pares_icon_outline, err := opentype.Parse(assets.FontIconOutline)
	pares_icon_fill, err := opentype.Parse(assets.FontIconFill)
	pares_emoji, err := opentype.Parse(assets.FontEmoji)
	pares_sans, err := opentype.Parse(assets.FontSans)

	if err != nil {
		return nil, err
	}
	fontFace := []font.FontFace{
		font.FontFace{
			Face: pares_sans,
			Font: font.Font{Typeface: "sans"},
		},
		font.FontFace{
			Face: pares_emoji,
			Font: font.Font{Typeface: "noto"},
		},
		font.FontFace{
			Face: pares_icon_outline,
			Font: font.Font{Typeface: "icons", Style: font.Regular},
		},
		font.FontFace{
			Face: pares_icon_fill,
			Font: font.Font{Typeface: "icons", Style: font.Italic},
		},
	}

	return fontFace, nil
}
