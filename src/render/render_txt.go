package render

import (
	"image"

	"github.com/gen2brain/go-fitz"
)

type NavigatorTXT interface {
	Render_full(id string, page int, dpi float64)
}

// pdf渲染实例
type RenderInstTXT struct {
	Path string
	Id   string
	doc  *fitz.Document
}

func NewRenderInstTXT() *RenderInstPDF {
	self := RenderInstPDF{}
	return &self
}

func (self *RenderInstTXT) New()     {}
func (self *RenderInstTXT) Destory() {}

func (self *RenderInstTXT) Get_total_count() int { return 0 }
func (self *RenderInstTXT) ID() string           { return "" }
func (self *RenderInstTXT) Type() RenderInstType { return RenderInstType_PDF }

func (self *RenderInstTXT) Render_full(id string, page int, dpi float64) *image.RGBA {
	return nil
}
