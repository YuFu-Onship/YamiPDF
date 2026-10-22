package render

import (
	rumia "app/src"

	"github.com/gen2brain/go-fitz"
)

// -------------------------------------------------------------------------

// 渲染实例类别
type RenderInstType int

const (
	RenderInstType_PDF RenderInstType = iota
	RenderInstType_TXT
	RenderInstType_EPUB
)

// 渲染接口
type RenderAPI interface {
	ID() string
	Type() RenderInstType

	Get_total_count() int
	Destory()
}

// --------------------------------------------------------------------------

func NewRender(file_type RenderInstType) (RenderAPI, error) {
	switch file_type {
	case RenderInstType_PDF:
		return &RenderInstPDF{}, nil
	}
	return nil, nil
}

// --------------------------------------------------------------------------

// 渲染实例管理器
type RenderInstManager struct {
	doc   *fitz.Document
	insts map[string]RenderAPI
}

func NewRenderInstManager() *RenderInstManager {
	self := RenderInstManager{
		insts: make(map[string]RenderAPI),
	}
	return &self
}

func (self *RenderInstManager) Run() {

}

func (self *RenderInstManager) API_get_inst(id string) RenderAPI {
	return self.insts[id]
}

func (self *RenderInstManager) API_add_inst(id string, file_type rumia.FileType) {
	switch file_type {
	case rumia.FileType_PDF:
		if self.doc != nil {
			self.doc.Close()
		}
		self.doc, _ = fitz.New(id)

		self.insts[id] = NewRenderInstPDF().
			API_add_path(id).
			API_set_doc(self.doc)

	case rumia.FileType_EPUB:
	case rumia.FileType_TXT:
		self.insts[id] = NewRenderInstTXT()
	}
}
