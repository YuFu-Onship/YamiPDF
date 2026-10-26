package render

import (
	rumia "app/src"
	"log"

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
	API_set_manager_api(manager_api RenderInstManagerAPI)
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

// 渲染实例管理器接口
// 用于向子实例提供管理器的相关功能
type RenderInstManagerAPI interface {
	API_add_task(task func())
}

// 渲染实例管理器
type RenderInstManager struct {
	doc       *fitz.Document
	insts     map[string]RenderAPI
	task_chan chan any
	task_list []rumia.BaseTask
}

func NewRenderInstManager() *RenderInstManager {
	self := RenderInstManager{
		insts:     make(map[string]RenderAPI),
		task_chan: make(chan any, 100),
		task_list: make([]rumia.BaseTask, 100),
	}
	return &self
}

func (self *RenderInstManager) Run() {
	// for task := range self.task_chan {
	// 	task()
	// }

	for {
		select {
		case task := <-self.task_chan:
			switch v := task.(type) {
			case func():
				v()
			case *rumia.RenderTask_PDF_Full:

			}

		default:
		}
	}
}

func (self *RenderInstManager) Delete_inst(id string) {
	if self.doc != nil {
		self.doc.Close()
	}
	delete(self.insts, id)
}

func (self *RenderInstManager) API_add_task(id string, task func()) {

	self.task_list = append(self.task_list, rumia.BaseTask{})
	// self.task_chan <- task
}

// 替换列表
func (self *RenderInstManager) API_set_tasks(tasks []rumia.BaseTask) {
	self.task_list = tasks
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
		var err error
		self.doc, err = fitz.New(id)
		if err != nil {
			log.Fatal(err)
		}

		self.insts[id] = NewRenderInstPDF().
			API_add_path(id).
			API_set_doc(self.doc)

		// self.insts[id].API_set_manager_api(self)

	case rumia.FileType_EPUB:
	case rumia.FileType_TXT:
		self.insts[id] = NewRenderInstTXT()
	}
}
