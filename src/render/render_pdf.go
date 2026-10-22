package render

import (
	"image"
	"log"

	"github.com/gen2brain/go-fitz"
)

type NavigatorPDF interface {
	Render_full(id string, page int, dpi float64)
}

// pdf渲染实例
type RenderInstPDF struct {
	Path string
	Id   string
	doc  *fitz.Document
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

func (self *RenderInstPDF) Render_full(page int, dpi float64) *image.RGBA {
	img, err := self.doc.ImageDPI(page, dpi)
	if err != nil {
		log.Fatal(err)
	}
	return img
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
