package render

import (
	"app/src/tool"
	"bytes"
	"errors"
	"image"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type PDFInfo struct {
	NumPage int
}

type PDFRender struct {
	doc *fitz.Document
}

// 渲染pdf图片并输出
//   - filepath 文件路径
//
// - page 页数(从1开始)
func (self *PDFRender) Render(filepath string, page int) (*image.RGBA, error) {
	totalStart := time.Now()

	// 1. 打开文件耗时
	stepStart := time.Now()
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.Printf("[DEBUG] 打开文件耗时: %v", time.Since(stepStart))

	// 2. Crop 裁剪耗时
	stepStart = time.Now()
	// 单位常量使用 types.POINTS
	box, _ := model.ParseBox("dim: 100% 100%, pos: c", types.POINTS)
	// box, _ := model.ParseBox("pos:tl, off:100 100, dim:1000 1000 abs", types.POINTS)
	var buf bytes.Buffer
	if err := api.Crop(f, &buf, []string{"1"}, box, nil); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] 页面裁剪(Crop)耗时: %v", time.Since(stepStart))

	// 3. 初始化 fitz 耗时
	stepStart = time.Now()
	doc, err := fitz.NewFromReader(&buf)
	if err != nil {
		return nil, err
	}
	defer doc.Close()
	log.Printf("[DEBUG] 初始化fitz(NewFromReader)耗时: %v", time.Since(stepStart))

	// 4. 图片渲染耗时
	stepStart = time.Now()
	// 直接返回 img 即可，无需类型断言
	// img, err := doc.Image(page)
	img, err := doc.ImageDPI(page, 36)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] 图像渲染(ImageDPI)耗时: %v", time.Since(stepStart))

	// 总耗时
	log.Printf("[DEBUG] PDF渲染完成，总耗时: %v", time.Since(totalStart))

	return img, nil
}

// 新渲染器 -----------------------------------------------------------------------------------

// 错误信息
var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidDocument  = errors.New("document instance is nil")
)

// 回调任务 --------------------
type renderTask struct {
	page    int
	call_fn func(img *image.RGBA)
}

// 渲染实例 --------------------
type RenderInstance struct {
	doc  *fitz.Document
	lock *tool.ChanMutex
	done chan struct{}

	tasks chan *renderTask

	task_mu sync.Mutex
	pending []*renderTask
	notify  chan struct{}
}

// 渲染管理器 ------------------
type RenderManager struct {
	mu    sync.RWMutex
	insts map[string]*RenderInstance
}

// 创建新的渲染管理器
func NewRenderManager() *RenderManager {
	self := &RenderManager{
		insts: map[string]*RenderInstance{},
	}
	return self
}

// 运行
func (self *RenderManager) Run() {
	// go func() {
	// self.render()
	// }()
}

// 添加新的渲染实例
func (self *RenderManager) Add_new_document(doc_path string) (string, error) {
	doc, err := fitz.New(doc_path)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()

	inst := &RenderInstance{
		doc:    doc,
		lock:   tool.NewChanMutex(),
		done:   make(chan struct{}),
		tasks:  make(chan *renderTask, 100),
		notify: make(chan struct{}),
	}

	self.mu.Lock()
	self.insts[id] = inst
	self.mu.Unlock()

	go self.render_inst(inst)

	return id, nil
}

// 内部功能 --------------------

// 单个实例的渲染
func (self *RenderManager) render_inst(inst *RenderInstance) {
	// 单页加载
	// for {
	// 	select {
	// 	case task, ok := <-inst.tasks:
	// 		if !ok {
	// 			return
	// 		}
	// 		img, err := inst.doc.ImageDPI(task.page, 32)
	// 		// img, err := inst.doc.Image(task.page)
	// 		if err != nil {
	// 			return
	// 		}
	// 		if task.call_fn != nil {
	// 			task.call_fn(img)
	// 		}

	// 	case <-inst.done:
	// 	}
	// }

	// 页列表
	for {
		select {
		case <-inst.notify:
			for {
				inst.task_mu.Lock()
				if len(inst.pending) == 0 {
					inst.task_mu.Unlock()
					break // 当前队列已处理完毕，退出内部循环，继续等待下一次 notify
				}

				task := inst.pending[0]
				inst.pending = inst.pending[1:]
				inst.task_mu.Unlock()

				// 在锁外执行耗时的渲染操作
				start := time.Now()
				img, _ := inst.doc.ImageDPI(task.page, 2)
				// img, _ := inst.doc.Image(task.page)
				log.Printf("渲染PDF%v", time.Since(start))

				if task.call_fn != nil {
					task.call_fn(img)
				}
			}
		case <-inst.done:
		}
	}
}

// 渲染页面
// func (self *RenderManager) render() {
// 	for _, inst := range self.insts {
// 		if inst.doc != nil && len(inst.load_pages) != 0 {
// 			if inst.lock.TryLock() {
// 				p := inst.load_pages[len(inst.load_pages)-1]
// 				inst.load_pages = inst.load_pages[:len(inst.load_pages)-1]
// 				go func(p int) {
// 					_, _ = inst.doc.Image(p)
// 					inst.lock.Unlock()
// 				}(p)
// 			}
// 		}
// 	}
// }

// 外部调用 -----------------------------------------------------

// 删除渲染实例
func (self *RenderManager) Delete_document(id string) {
	self.mu.Lock()
	if inst, exist := self.insts[id]; exist && inst.doc != nil {
		inst.doc.Close()
		delete(self.insts, id)
	}
	self.mu.Unlock()
}

// 获取文件信息
//   - meta_data  map[string]string 元数据[title,author]
//   - page_num int 总页数
func (self *RenderManager) Get_info(id string) map[string]any {
	info := map[string]any{
		"meta_data": map[string]string{},
		"page_num":  0,
	}

	self.mu.RLock()
	inst, exist := self.insts[id]
	self.mu.RUnlock()

	if exist && inst.doc != nil {
		info["meta_data"] = inst.doc.Metadata()
		info["page_num"] = inst.doc.NumPage()
	}
	return info
}

// 增加需要渲染的页(回调)
//   - call_func 渲染完成后需要需要执行的操作
func (self *RenderManager) Add_render_page(id string, page int, call_fn func(img *image.RGBA)) {
	self.mu.RLock()
	inst, exist := self.insts[id]
	self.mu.RUnlock()

	if !exist {
		return
	}

	task := renderTask{
		page:    page,
		call_fn: call_fn,
	}

	select {
	case inst.tasks <- &task:
	default:
	}
}

// 增加需要渲染的页列表(回调)
func (self *RenderManager) Add_render_page_list(id string, pages []int, call_fn func(page int, img *image.RGBA)) {
	self.mu.Lock()
	inst, exist := self.insts[id]
	self.mu.Unlock()

	if !exist {
		return
	}

	var newTasks []*renderTask
	for _, p := range pages {
		page := p
		newTasks = append(newTasks, &renderTask{
			page: page,
			call_fn: func(img *image.RGBA) {
				call_fn(page, img)
			},
		})
	}
	// log.Println(pages)

	inst.task_mu.Lock()
	inst.pending = newTasks
	inst.task_mu.Unlock()

	select {
	case inst.notify <- struct{}{}:
	default:
	}

}

// 渲染指定页面
func (self *RenderManager) Render_page(id string, page int) (*image.RGBA, error) {
	start := time.Now()
	if inst, exist := self.insts[id]; exist && inst.doc != nil {
		img, err := inst.doc.ImageDPI(page, 4)
		if err != nil {
			return nil, err
		}
		log.Printf("渲染第 %d 页, 耗时: %v", page, time.Since(start))
		return img, err
	}
	return nil, ErrDocumentNotFound
}

// 得到指定页的尺寸
func (self *RenderManager) Get_size(id string, page int) (image.Rectangle, error) {
	if inst, exist := self.insts[id]; exist && inst.doc != nil {
		size, err := inst.doc.Bound(page)
		if err != nil {
			return image.Rect(0, 0, 0, 0), err
		}
		return size, nil
	}
	return image.Rect(0, 0, 0, 0), ErrDocumentNotFound
}
