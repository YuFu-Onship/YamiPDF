package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/shirou/gopsutil/v4/process"
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
	style := NewStyle(theme, WindowsBlue)
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

	pdf_path := path.Join(self.bridge.GetRootPath(), "/assets/Witch Hat Atelier, Vol.5 (Kamome Shirahama).pdf")
	// pdf_path := path.Join(self.bridge.GetRootPath(), "E:/Users/YUFU/Documents/Books/comic/Archives of Witch Hat Atelier (Kamome Shirahama).pdf")

	doc_id, _ := self.bridge.Add_new_document(pdf_path)
	pdfview := NewPDFView(doc_id, self.bridge, self)

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

			gtx := app.NewContext(&ops, typ)
			paint.Fill(gtx.Ops, self.Style.Palette.BG)

			self.Deco.Actions(gtx, w)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D { return self.Deco.Layout(gtx) }),
				layout.Flexed(1, func(gtx C) D { return pdfview.Layout(gtx, self.Style) }),
			)

			self.Style.AniSys.Run(gtx)
			typ.Frame(gtx.Ops)
		}
	}
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func printOSMemUsage() {
	// 获取当前进程 PID
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		fmt.Println("获取进程失败:", err)
		return
	}

	// 获取内存占用信息
	memInfo, err := p.MemoryInfo()
	if err != nil {
		fmt.Println("获取内存信息失败:", err)
		return
	}

	// RSS (Resident Set Size) 即操作系统为该进程分配的物理内存大小（字节）
	rssMB := float64(memInfo.RSS) / 1024 / 1024
	// VMS (Virtual Memory Size) 虚拟内存大小
	vmsMB := float64(memInfo.VMS) / 1024 / 1024

	fmt.Printf("[PID: %d] 实际占用物理内存(RSS): %.2f MB | 虚拟内存(VMS): %.2f MB\n", pid, rssMB, vmsMB)
}

func printMem(tag string) {
	runtime.GC()         // 强制 GC
	debug.FreeOSMemory() // 尝试把释放的内存还给操作系统
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("[%s] Go 堆内存: %d MB\n", tag, m.Alloc/1024/1024)
}

// 图像灰度处理
func ToGrayscaleOp(src image.Image) paint.ImageOp {
	bounds := src.Bounds()

	dst := image.NewRGBA(bounds)
	for x := 0; x <= bounds.Dx(); x++ {
		for y := 0; y <= bounds.Dy(); y++ {
			r, g, b, a := src.At(x, y).RGBA()
			if a == 0 {
				dst.SetRGBA(x, y, color.RGBA{})
				continue
			}
			c := uint8(0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
			dst.SetRGBA(x, y, color.RGBA(color.RGBA{R: c, G: c, B: c, A: uint8(a >> 8)}))
		}
	}
	op := paint.NewImageOp(dst)
	op.Filter = paint.FilterNearest
	return op
}

// 接口部分 ------------------------------------------

// 切换颜色模式
func (self Page) ToggleTheme() {
	if self.Style.isDark {
		self.Style.SetColorMode(false)
	} else {
		self.Style.SetColorMode(true)
	}
	self.Style.apply()
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
