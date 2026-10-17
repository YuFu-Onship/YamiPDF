package page

import (
	"app/assets"
	"app/src/ui"
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
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions
type W = layout.Widget

type Page struct {
	Deco         *ui.CustomDeco
	Style        *ui.Style
	bridge_trunk ui.BridgeTrunk

	win_size      image.Point
	is_win_change bool
	appwindow     *app.Window

	page_manager *PageManager
}

func NewPage() *Page {
	theme := material.NewTheme()
	style := ui.NewStyle(theme, ui.Palette_teto)
	style.SetColorMode(true)

	self := Page{
		Deco:         ui.NewCustomDeco(&style),
		Style:        &style,
		page_manager: NewPageManager(&style),
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
	self.Style.Lang = ui.CN
	self.Style.Bridge_ui = self

	// 加载文字
	if fonts, err := self.load_fonts(); err == nil {
		shaper := text.NewShaper(text.WithCollection(fonts))
		self.Style.Theme.Shaper = shaper
	}

	// self.page_manager.API_add_page("main", NewPageMain(self.Style))
	self.page_manager.Add_page("main", true, func() PageAPI { return NewPageMain(self.Style) })
	self.page_manager.Add_page("pdf_viewer", false, func() PageAPI { return NewPagePDF(self.Style) })
	self.page_manager.API_set_cur_page("main")

	self.page_manager.API_set_bridge_ui("main", self)

	// pdf阅读器
	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Archives of Witch Hat Atelier (Kamome Shirahama).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Witch Hat Atelier, Vol.5 (Kamome Shirahama).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/单页.pdf"

	// doc_id, _ := self.bridge_trunk.Add_new_document(pdf_path)
	// pdfview := ui.NewPDFView(doc_id, self.Style, self.bridge_trunk, self)

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

			deco_back_state := Ternary(self.page_manager.API_get_cur_id() == "main", false, true)
			self.Deco.API_set_back(deco_back_state)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(self.Deco.Layout),
				layout.Flexed(1, self.page_manager.Layout),
			)

			self.Style.AniSys.Run(gtx)
			typ.Frame(gtx.Ops)
		}
	}
}

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

// 更改自定义标题栏文本
func (self Page) Set_deco_text(text string) {
	self.Deco.API_set_text(text)
	self.WindowRefresh()
}

// 更改原生标题栏文本
func (self *Page) Set_title_text(text string) {
	if self.appwindow != nil {
		self.appwindow.Option(app.Title(text))
	}
}

func (self Page) WindowRefresh() {
	self.appwindow.Invalidate()
}

// 窗口大小是否变动
//   - true 发生变化
//   - false 没有变化
func (self Page) IsWinSizeChange() bool {
	return self.is_win_change
}

// 设置trunk根接口
func (self *Page) API_set_trunk(trunk ui.BridgeTrunk) *Page {
	self.bridge_trunk = trunk
	return self
}

// 切换颜色模式
func (self Page) ToggleTheme() {
	self.Style.Toggle()
}

// 设置界面
func (self *Page) Set_cur_page(id string) {
	self.page_manager.API_set_cur_page(id)
}
