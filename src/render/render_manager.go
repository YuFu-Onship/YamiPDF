package render

import (
	rumia "app/src"
	"log"
	"runtime"
	"sync"

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
	doc   *fitz.Document
	insts map[string]RenderAPI

	task_chan   chan rumia.BaseTask
	task_list   []rumia.BaseTask
	task_notify chan struct{}

	mu sync.Mutex
}

func NewRenderInstManager() *RenderInstManager {
	self := RenderInstManager{
		insts:       make(map[string]RenderAPI),
		task_chan:   make(chan rumia.BaseTask, 100),
		task_list:   []rumia.BaseTask{},
		task_notify: make(chan struct{}, 1),
	}
	return &self
}

func (self *RenderInstManager) Run() {
	for {
		self.mu.Lock()
		if len(self.task_list) == 0 {
			self.mu.Unlock()

			<-self.task_notify
			continue
		}

		task := self.task_list[0]
		self.task_list = self.task_list[1:]
		self.mu.Unlock()

		switch v := task.Task.(type) {
		case func():
			v()
		case *rumia.RenderTask_PDF_Full:
			self.API_get_inst(task.Id).(*RenderInstPDF).Render_full_page(v)
		case *rumia.RenderTask_PDF_Tile:
			self.API_get_inst(task.Id).(*RenderInstPDF).Render_tile(v)
		default:
			log.Println("test")
		}

		// 		select {
		// 		case task := <-self.task_chan:
		// 			switch v := task.Task.(type) {
		// 			case func():
		// 				v()
		// 			case *rumia.RenderTask_PDF_Full:
		// 				// self.API_get_inst(task.Id).(*RenderInstPDF).Render_full()
		// 			}
		//
		// 		default:
		// }
	}
}

func (self *RenderInstManager) notify() {
	select {
	case self.task_notify <- struct{}{}:
	default:
	}
}

func (self *RenderInstManager) API_close_inst(id string) {
	if self.doc != nil {
		// go-fitz无法完全释放占用的内存, 所以会导致内存占用问题
		// 有时候多次重新打开/关闭文件能实现
		self.doc.Close()
		runtime.GC()
	}
}

// 追加添加任务
func (self *RenderInstManager) API_add_task(task rumia.BaseTask) {
	self.mu.Lock()
	self.task_list = append(self.task_list, task)
	self.mu.Unlock()
	self.notify()
}

// 覆盖替换任务列表
func (self *RenderInstManager) API_set_tasks(tasks []rumia.BaseTask) {
	self.mu.Lock()
	self.task_list = tasks
	self.mu.Unlock()
	self.notify()
}

// 清空任务列表
func (self *RenderInstManager) API_clear_tasks() {
	self.mu.Lock()
	self.task_list = make([]rumia.BaseTask, 0)
	self.mu.Unlock()
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
			API_set_path(id).
			API_set_doc(self.doc)

		// self.insts[id].API_set_manager_api(self)

	case rumia.FileType_EPUB:
	case rumia.FileType_TXT:
		self.insts[id] = NewRenderInstTXT()
	}
}
