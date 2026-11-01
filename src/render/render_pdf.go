package render

import (
	rumia "app/src"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type NavigatorPDF interface {
	Render_full(id string, page int, dpi float64)
}

// pdf渲染实例
type RenderInstPDF struct {
	Path    string
	Id      string
	doc     *fitz.Document
	Manager RenderInstManagerAPI

	task_list []any

	model_conf *model.Configuration //pdfcpu 配置

	page_list []*image.RGBA
	tile_list []*image.RGBA
	mu        sync.Mutex
}

func NewRenderInstPDF() *RenderInstPDF {
	self := RenderInstPDF{
		task_list: []any{},
		page_list: []*image.RGBA{},
		tile_list: []*image.RGBA{},
	}
	self.model_conf = model.NewDefaultConfiguration()
	self.model_conf.Optimize = false
	return &self
}

func (self *RenderInstPDF) New() {}
func (self *RenderInstPDF) Destory() {
	// 确实需要手动解引用
	for i, ptr := range self.page_list {
		if ptr != nil {
			self.page_list[i] = nil
		}
	}

	for i, ptr := range self.tile_list {
		if ptr != nil {
			self.tile_list[i] = nil
		}
	}

	self.mu.Lock()
	if self.doc != nil {
		self.doc.Close()
		self.doc = nil
	}
	self.mu.Unlock()
}

func (self *RenderInstPDF) Get_total_count() int { return 0 }
func (self *RenderInstPDF) ID() string           { return self.Id }
func (self *RenderInstPDF) Type() RenderInstType { return RenderInstType_PDF }

func (self *RenderInstPDF) Run() {}

// 内部调用 -------------------------------------------------------------

func (self *RenderInstPDF) calc_page_data(page int) *[]byte {

	f, err := os.Open(self.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 当前页缓冲区
	var page_buf []byte
	pages := []string{fmt.Sprintf("%d", page+1)}

	if api.ExtractPages(f, pages, func(r io.Reader, i int) error {
		page_buf, err = io.ReadAll(r)
		return err
	}, self.model_conf) != nil {
		log.Fatal(err)
	}
	return &page_buf
}

// 外部调用 -------------------------------------------------------------

// 渲染整页
func (self *RenderInstPDF) Render_full_page(task *rumia.RenderTask_PDF_Full) {
	self.mu.Lock()
	if self.doc == nil {
		self.doc, _ = fitz.New(self.Path)
	}

	img := Custom_pdf_render_full(self.doc, task.Page, task.DPI)
	self.mu.Unlock()

	self.page_list = append(self.page_list, img)
	task.Call_fn(img)
}

// 渲染瓦片
func (self *RenderInstPDF) Render_tile(task *rumia.RenderTask_PDF_Tile) {
	self.mu.Lock()
	img := Custom_pdf_render_tile(self.doc, task.Page, task.Dpi, [4]float64{
		task.Rect.LL.X,
		task.Rect.LL.Y,
		task.Rect.UR.X,
		task.Rect.UR.Y,
	})
	self.mu.Unlock()

	self.page_list = append(self.page_list, img)
	task.Call_fn(img)
}

// 设置路径
func (self *RenderInstPDF) API_set_path(path string) *RenderInstPDF {
	self.Path = path
	return self
}

// 设置mupdf实例
func (self *RenderInstPDF) API_set_doc(doc *fitz.Document) *RenderInstPDF {
	self.doc = doc
	return self
}

// 得到doc实例
func (self *RenderInstPDF) API_get_doc() *fitz.Document {
	return self.doc
}

// 得到指定页尺寸
func (self *RenderInstPDF) API_get_page_size(page int) image.Rectangle {
	size, err := self.doc.Bound(page)
	if err != nil {
		log.Fatal(err)
	}
	return size
}

// 得到pdf信息
func (self *RenderInstPDF) API_get_info() map[string]string {
	return self.doc.Metadata()
}

// 得到pdf页数
func (self *RenderInstPDF) API_get_num() int {
	return self.doc.NumPage()
}

// 设置上一层级的渲染实例管理器接口
func (self *RenderInstPDF) API_set_manager_api(manager_api RenderInstManagerAPI) {
	self.Manager = manager_api
}

func (self *RenderInstPDF) API_initization() *fitz.Document {
	var err error
	self.doc, err = fitz.New(self.Path)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return self.doc
}
