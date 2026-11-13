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
	API_get_doc() *fitz.Document
	API_initization() *fitz.Document
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
}

// 渲染实例管理器
type RenderInstManager struct {
	doc   *fitz.Document
	insts map[string]RenderAPI

	task_chan chan rumia.BaseTask

	task_list   []rumia.BaseTask
	task_list_1 []rumia.BaseTask // 覆写任务列表
	task_list_2 []rumia.BaseTask // 追加任务列表
	chan_notify chan struct{}

	mu sync.Mutex
}

func NewRenderInstManager() *RenderInstManager {
	self := RenderInstManager{
		insts:       make(map[string]RenderAPI),
		task_chan:   make(chan rumia.BaseTask, 100),
		task_list:   []rumia.BaseTask{},
		task_list_1: []rumia.BaseTask{},
		task_list_2: []rumia.BaseTask{},
		chan_notify: make(chan struct{}, 1),
	}
	return &self
}

func (self *RenderInstManager) Run() {
	for {

		self.mu.Lock()
		if len(self.task_list) == 0 && len(self.task_list_1) == 0 && len(self.task_list_2) == 0 {
			self.mu.Unlock()
			runtime.GC()
			<-self.chan_notify // 没有通知时堵塞
			continue
		}

		if len(self.task_list_1) != 0 {
			self.task_list = append(self.task_list, self.task_list_1[0])
			self.task_list_1 = self.task_list_1[1:]
		}

		if len(self.task_list_2) != 0 {
			self.task_list = append(self.task_list, self.task_list_2[0])
			self.task_list_2 = self.task_list_2[1:]
		}

		// task_1 := self.task_list[0]
		// self.task_list = self.task_list[1:]
		// self.task_list = append(self.task_list, task_1)

		task := self.task_list[0]
		var zero rumia.BaseTask
		self.task_list[0] = zero
		self.task_list = self.task_list[1:]

		self.mu.Unlock()

		switch v := task.Task.(type) {
		case func():
			v()

		case *rumia.RenderTask_PDF_Full:
			if inst := self.API_get_inst(task.Id); inst != nil {
				inst.(*RenderInstPDF).Render_full_page(v)
			}
		case *rumia.RenderTask_PDF_Tile:
			if inst := self.API_get_inst(task.Id); inst != nil {
				inst.(*RenderInstPDF).Render_tile(v)
			}
		case *rumia.RenderTask_PDF_ReleasePix:
			if v.Img != nil {
				v.Img.Pix = nil
				v.Img = nil
			}
		default:
			log.Println("test")
		}

		// 		task_a := self.task_list_2[0]
		// 		self.task_list_2 = self.task_list_2[1:]
		//
		// 		switch v := task_a.Task.(type) {
		// 		case *rumia.RenderTask_PDF_ClearPix:
		// 			v.Img = nil
		// 		}

	}
}

func (self *RenderInstManager) notify() {
	select {
	case self.chan_notify <- struct{}{}:
	default:
	}
}

// 覆盖替换任务列表
func (self *RenderInstManager) API_set_tasks(tasks []rumia.BaseTask) {
	self.mu.Lock()
	self.task_list_1 = tasks
	self.mu.Unlock()
	self.notify()
}

// 追加添加任务列表
func (self *RenderInstManager) API_add_tasks(tasks []rumia.BaseTask) {
	self.mu.Lock()
	self.task_list_2 = append(self.task_list, tasks...)
	self.mu.Unlock()
	self.notify()
}

func (self *RenderInstManager) API_close_inst(id string) {
	if self.doc != nil {
		// go-fitz无法完全释放占用的内存, 所以会导致内存占用问题
		// 有时候多次重新打开/关闭文件能实现

		// doc实例的关闭是由渲染实例自己处理的
		if inst := self.API_get_inst(id); inst != nil {
			inst.Destory()
		}

		// self.doc.Close()
		// self.doc = nil

		runtime.GC()
	}
	if _, exist := self.insts[id]; exist {
		delete(self.insts, id)
	}
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

		inst := NewRenderInstPDF().API_set_path(id)
		inst.API_set_manager_api(self)

		self.doc = inst.API_initization()
		self.insts[id] = inst

	case rumia.FileType_EPUB:
	case rumia.FileType_TXT:
		self.insts[id] = NewRenderInstTXT()
	}
}
