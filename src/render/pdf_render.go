package render

import (
	rumia "app/src"
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
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

// 新渲染器 -----------------------------------------------------------------------------------
// 错误信息
var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidDocument  = errors.New("document instance is nil")
)

// 回调任务 --------------------
type renderFullTask struct {
	page    int
	dpi     float64
	call_fn func(img *image.RGBA)
}

type renderTileTask struct {
	page    int
	dpi     float64
	rect    *types.Rectangle
	call_fn func(img *image.RGBA)
}

// 渲染模式 --------------------
type RenderMode int

const (
	render_unknow RenderMode = iota
	render_fill              //渲染整页
	render_tile              //渲染部分瓦片
)

// 渲染实例 -------------------------------------------------------------
type RenderInstance struct {
	file_path string         //文件路径
	doc       *fitz.Document //渲染实例
	inst_mu   sync.Mutex     //防止同时操作实例或任务列表

	done         chan struct{}           //渲染完成通知,关闭实例
	pending_full []*rumia.RenderFullTask //整页任务列表
	pending_tile []*rumia.RenderTileTask //瓦片任务列表
	tasks        chan *renderFullTask    //渲染任务通道
	notify_full  chan struct{}           //通知,用于更新渲染列表
	notify_tile  chan struct{}           //通知,用于更新渲染列表

	mode        RenderMode      //渲染模式
	page_buffer map[int]*[]byte //pdf页缓冲区

	model_conf *model.Configuration //pdfcpu 配置
}

func NewRenderInstance(doc_path string) *RenderInstance {
	self := RenderInstance{
		done:        make(chan struct{}),
		tasks:       make(chan *renderFullTask, 100),
		notify_full: make(chan struct{}),
		notify_tile: make(chan struct{}),
		file_path:   doc_path,

		mode: render_fill,
	}
	self.doc, _ = fitz.New(doc_path) // 2ms

	self.model_conf = model.NewDefaultConfiguration()
	self.model_conf.Optimize = false
	self.model_conf.ValidationMode = model.ValidationRelaxed

	self.page_buffer = map[int]*[]byte{}
	return &self
}

// 创建doc
func (self *RenderInstance) New_doc_inst() {
	self.doc, _ = fitz.New(self.file_path)
}

// 增加页缓冲区
func (self *RenderInstance) Add_page_buffer(page int) *[]byte {
	f, err := os.Open(self.file_path)
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

	self.page_buffer[page] = &page_buf
	return &page_buf
}

// 得到页缓冲区
func (self *RenderInstance) Get_page_buffer(page int) *[]byte {
	if buf, exist := self.page_buffer[page]; exist {
		return buf
	}

	f, err := os.Open(self.file_path)
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

	self.page_buffer[page] = &page_buf
	return &page_buf
}

// 渲染瓦片
func (self *RenderInstance) Tile_render() {
	for {
		self.inst_mu.Lock()
		if len(self.pending_tile) == 0 {
			self.inst_mu.Unlock()
			break
		}

		task := self.pending_tile[0]
		self.pending_tile = self.pending_tile[1:]

		box := model.Box{Rect: task.Rect}
		page_buf := self.Add_page_buffer(task.Page)

		var tile_buf bytes.Buffer
		err := api.Crop(bytes.NewReader(*page_buf), &tile_buf, []string{"1"}, &box, self.model_conf)
		if err != nil {
			log.Fatal(err)
		}

		doc, _ := fitz.NewFromReader(&tile_buf)
		img, _ := doc.ImageDPI(0, task.Dpi)
		doc.Close()
		task.Call_fn(img)
		self.inst_mu.Unlock()
	}
}

// 渲染整页
func (self *renderFullTask) Full_render() {

}

// 页缓冲区
func (self *RenderInstance) update_page_buffer(tasks []*rumia.RenderTileTask) {
	keep_pages := map[int]bool{}
	for _, t := range tasks {
		count := t.Page
		keep_pages[count] = true
		if _, exist := self.page_buffer[count]; !exist {
			self.Add_page_buffer(count)
		}
	}

	for page := range self.page_buffer {
		if !keep_pages[page] {
			delete(self.page_buffer, page)
		}
	}
}

// 释放整页内存
func (self *RenderInstance) release_memory() {

}

// 渲染管理器 --------------------------------------------------------------
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
func (self *RenderManager) Run() {}

// 添加新的渲染实例
func (self *RenderManager) Add_new_document(doc_path string) (string, error) {
	id := uuid.New().String()

	inst := NewRenderInstance(doc_path)

	self.mu.Lock()
	self.insts[id] = inst
	self.mu.Unlock()

	go self.render_inst(inst)

	return id, nil
}

// 内部功能 --------------------

// 渲染局部信息
func (self *RenderManager) render_local() {
	doc_path := ""

	start := time.Now()

	f, err := os.Open(doc_path)
	if err != nil {
	}
	defer f.Close()

	conf := model.NewDefaultConfiguration()
	conf.Optimize = false
	conf.ValidationMode = model.ValidationRelaxed

	pageNum := 1
	var buf []byte
	pages := []string{fmt.Sprintf("%d", pageNum)}
	err = api.ExtractPages(f, pages, func(r io.Reader, i int) error {
		buf, err = io.ReadAll(r)
		return err
	}, conf)

	// log.Printf("size:%.2f mb", float64(len(buf))/(1<<20))
	rect := types.NewRectangle(50, 50, 300, 400)
	box := model.Box{
		Rect: rect,
	}

	var n_buf bytes.Buffer
	api.Crop(bytes.NewReader(buf), &n_buf, []string{"1"}, &box, nil)

	d, _ := fitz.NewFromMemory(n_buf.Bytes())
	d.Image(0)
	log.Printf("渲染耗时:%v", time.Since(start))
}

// 单个实例的渲染,运行时使用go func
func (self *RenderManager) render_inst(inst *RenderInstance) {

	for {
		select {
		case <-inst.notify_tile:
			inst.update_page_buffer(inst.pending_tile)
			// inst.Tile_render()
			for {
				inst.inst_mu.Lock()
				if len(inst.pending_tile) == 0 {
					inst.inst_mu.Unlock()
					break
				}

				task := inst.pending_tile[0]
				inst.pending_tile = inst.pending_tile[1:]
				inst.inst_mu.Unlock()

				// 开始渲染瓦片
				box := model.Box{Rect: task.Rect}
				page_buf := inst.Get_page_buffer(task.Page)

				var tile_buf bytes.Buffer
				err := api.Crop(bytes.NewReader(*page_buf), &tile_buf, []string{"1"}, &box, inst.model_conf)
				if err != nil {
					log.Fatal(err)
				}

				doc, _ := fitz.NewFromReader(&tile_buf)
				img, _ := doc.ImageDPI(0, task.Dpi)
				doc.Close()

				if task.Call_fn != nil {
					task.Call_fn(img)
				}
			}

		case <-inst.notify_full:
			for {
				inst.inst_mu.Lock()
				if len(inst.pending_full) == 0 {
					inst.inst_mu.Unlock()
					break
				}

				task := inst.pending_full[0]
				inst.pending_full = inst.pending_full[1:]
				inst.inst_mu.Unlock()

				// 开始渲染整页
				if inst.doc == nil {
					inst.New_doc_inst()
				}

				img, _ := inst.doc.ImageDPI(task.Page, float64(task.DPI))
				if task.Call_fn != nil {
					task.Call_fn(img)
				}

				// 1. 获取图像宽高
				// 				bounds := img.Bounds()
				// 				width := bounds.Dx()
				// 				height := bounds.Dy()
				//
				// 				// 2. 计算内存字节数 (RGBA 每像素占用 4 字节)
				// 				bytesCount := width * height * 4
				//
				// 				// 3. 转换为 MB
				// 				sizeInMB := float64(bytesCount) / (1024 * 1024)
				// 				fmt.Printf("图像分辨率: %dx%d, 内存占用: %.2f MB\n", width, height, sizeInMB)
				img = nil
			}

		case <-inst.done:
		}
	}
}

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

		r, err := GetCustomPageCount(inst.doc)
		if err != nil {
			log.Println(err)
		}
		log.Println(r)

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

	task := renderFullTask{
		page:    page,
		call_fn: call_fn,
	}

	select {
	case inst.tasks <- &task:
	default:
	}
}

// 增加需要渲染的页列表(回调)
func (self *RenderManager) Add_render_page_list(id string, tasks []*rumia.RenderFullTask) {

	self.mu.Lock()
	inst, exist := self.insts[id]
	self.mu.Unlock()

	if !exist {
		return
	}

	if inst.mode != render_fill {
		inst.mode = render_fill
		inst.inst_mu.Lock()
		inst.doc, _ = fitz.New(inst.file_path)
		inst.inst_mu.Unlock()
	}

	inst.inst_mu.Lock()
	inst.pending_full = tasks
	inst.inst_mu.Unlock()

	select {
	case inst.notify_full <- struct{}{}:
	default:
	}
}

// 需要渲染的瓦片区域
func (self *RenderManager) Add_render_tile_list(id string, tasks []*rumia.RenderTileTask) {
	self.mu.Lock()
	inst, exist := self.insts[id]
	self.mu.Unlock()

	if !exist {
		return
	}

	inst.pending_tile = tasks
	inst.mode = render_tile

	select {
	case inst.notify_tile <- struct{}{}:
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

// 释放内存
func (self *RenderManager) API_release_memory(id string) {
	if inst, exist := self.insts[id]; exist && inst.doc != nil {

	}
}
