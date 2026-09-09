package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	fitz "github.com/gen2brain/go-fitz"
)

const pdfPath = "E:/Users/YUFU/Documents/Books/comic/Witch Hat Atelier, Vol.5 (Kamome Shirahama) (z-library.sk, 1lib.sk, z-lib.sk).pdf"
const renderDPI = 32

// viewer 保存当前显示页,并在后台渲染页面。
type viewer struct {
	doc *fitz.Document

	mtx        sync.Mutex
	page       int           // 当前显示页 (0-based)
	img        *widget.Image // 当前显示图像
	busy       bool          // 是否有渲染任务进行中
	pending    int           // 渲染期间记录的最新请求页, -1 表示无
	invalidate func()        // 请求窗口重绘, 由渲染 goroutine 调用
	closed     atomic.Bool   // 窗口是否已关闭

	wg sync.WaitGroup // 等待在途渲染结束, 便于安全关闭 doc
}

// newPageImage 构造按窗口大小等比缩放的图像组件。
func newPageImage(img image.Image) *widget.Image {
	return &widget.Image{
		Src:      paint.NewImageOp(img),
		Fit:      widget.Contain, // 保持纵横比缩放到可用区域
		Position: layout.Center,
	}
}

// request 请求显示第 p 页 (0-based)。
func (v *viewer) request(p int) {
	v.mtx.Lock()
	defer v.mtx.Unlock()
	if p < 0 || p >= v.doc.NumPage() {
		return
	}
	if v.busy {
		v.pending = p // 渲染期间按键只记录最新目标
		return
	}
	v.busy = true
	v.wg.Add(1)
	go v.render(p)
}

// render 在后台渲染第 p 页。
func (v *viewer) render(p int) {
	defer v.wg.Done()

	img, err := v.doc.ImageDPI(p, renderDPI)
	if err != nil {
		log.Printf("渲染第 %d 页失败: %v", p+1, err)
	} else {
		debugf("render page %d done %dx%d", p, img.Bounds().Dx(), img.Bounds().Dy())
		v.mtx.Lock()
		v.page = p
		v.img = newPageImage(img)
		v.mtx.Unlock()
		if !v.closed.Load() && v.invalidate != nil {
			v.invalidate() // 新图像就绪, 请求重绘
		}
	}

	// 渲染期间如有新请求则继续
	v.mtx.Lock()
	next := v.pending
	v.pending = -1
	v.busy = false
	v.mtx.Unlock()
	if next >= 0 && next != p {
		v.request(next)
	}
}

// flipPage 相对当前页翻页。
func (v *viewer) flipPage(delta int) {
	v.mtx.Lock()
	cur := v.page
	v.mtx.Unlock()
	v.request(cur + delta)
}

func run(w *app.Window, v *viewer) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			v.closed.Store(true)
			debugf("DestroyEvent err=%v", e.Err)
			return e.Err
		case app.FrameEvent:
			debugf("FrameEvent")
			gtx := app.NewContext(&ops, e)

			// 上下键翻页
			for {
				ev, ok := gtx.Event(
					key.Filter{Name: key.NameUpArrow},
					key.Filter{Name: key.NameDownArrow},
				)
				if !ok {
					break
				}
				if ke, ok := ev.(key.Event); ok && ke.State == key.Press {
					switch ke.Name {
					case key.NameUpArrow:
						v.flipPage(-1)
					case key.NameDownArrow:
						v.flipPage(+1)
					}
				}
			}

			v.mtx.Lock()
			img := v.img
			page := v.page
			total := v.doc.NumPage()
			v.mtx.Unlock()

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if img == nil {
						return layout.Dimensions{Size: gtx.Constraints.Min}
					}
					return img.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Label(theme, unit.Sp(16), fmt.Sprintf("第 %d / %d 页   (↑/↓ 翻页)", page+1, total))
					label.Alignment = text.Middle
					return label.Layout(gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

// debugf 输出调试日志。
func debugf(format string, args ...any) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\n", fmt.Sprintf(format, args...))
}

func main() {
	// 1. 打开 PDF
	doc, err := fitz.New(pdfPath)
	if err != nil {
		log.Fatalf("打开PDF失败: %v", err)
	}
	defer doc.Close()

	v := &viewer{doc: doc, pending: -1}

	// 2. 文档信息
	pageCount := doc.NumPage()
	log.Printf("总页数: %d", pageCount)
	if pageCount == 0 {
		log.Fatal("文档没有页面")
	}

	// 3. 打开窗口 (后台渲染第一页, 不阻塞启动)
	w := new(app.Window)
	v.invalidate = w.Invalidate
	go func() {
		w.Option(
			app.Title("PDF Reader"),
			app.Size(unit.Dp(900), unit.Dp(1200)),
		)
		if err := run(w, v); err != nil {
			log.Fatal(err)
		}
	}()
	v.request(0)
	app.Main()
	debugf("app.Main returned")

	// 等待在途渲染结束再关闭文档
	v.wg.Wait()
	debugf("wg.Wait done")
}
