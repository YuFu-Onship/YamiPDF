package render

import (
	rumia "app/src"
	"image"
	"log"

	"github.com/gen2brain/go-fitz"
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
}

func NewRenderInstPDF() *RenderInstPDF {
	self := RenderInstPDF{}
	return &self
}

func (self *RenderInstPDF) New()     {}
func (self *RenderInstPDF) Destory() {}

func (self *RenderInstPDF) Get_total_count() int { return 0 }
func (self *RenderInstPDF) ID() string           { return "" }
func (self *RenderInstPDF) Type() RenderInstType { return RenderInstType_PDF }

func (self *RenderInstPDF) Render_full(task rumia.RenderTask_PDF_Full) {
	if self.Manager == nil {
		return
	}

	// 将渲染传入渲染通道
	self.Manager.API_add_task(func() {
		img, err := self.doc.ImageDPI(task.Page, task.DPI)
		if err != nil {
			log.Fatal(err)
		}
		task.Call_fn(img)
	})
}

func (self *RenderInstPDF) Render_full_pages(tasks []*rumia.RenderTask_PDF_Full) {

}

func (self *RenderInstPDF) API_add_path(path string) *RenderInstPDF {
	self.Path = path
	return self
}
func (self *RenderInstPDF) API_set_doc(doc *fitz.Document) *RenderInstPDF {
	self.doc = doc
	return self
}

func (self *RenderInstPDF) API_get_page_size(page int) image.Rectangle {
	size, err := self.doc.Bound(page)
	if err != nil {
		log.Fatal(err)
	}
	return size
}

func (self *RenderInstPDF) API_get_info() map[string]string {
	return self.doc.Metadata()
}

func (self *RenderInstPDF) API_get_num() int {
	return self.doc.NumPage()
}

func (self *RenderInstPDF) API_set_manager_api(manager_api RenderInstManagerAPI) {
	self.Manager = manager_api
}
