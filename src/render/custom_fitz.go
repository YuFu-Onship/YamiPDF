package render

import (
	"fmt"
	"image"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// #cgo CFLAGS: -I${SRCDIR}/../../
// #include "mupdf/fitz.h"
// #include "custom_fitz.h"
import "C"

// go-fitz内部的类在这里的对照
// 严格按照 Document 结构体顺序和类型对齐内存, 越权强读fitz内部的指针
type fitzDocInternal struct {
	ctx    unsafe.Pointer // *C.struct_fz_context
	data   []byte         // []byte
	doc    unsafe.Pointer // *C.struct_fz_document
	mtx    sync.Mutex     // sync.Mutex
	stream unsafe.Pointer // *C.fz_stream
}

func Custom_get_page_count(doc *fitz.Document) (int, error) {
	if doc == nil {
		return 0, fmt.Errorf("document is nil")
	}

	// 转换为本地内存镜像
	internal := (*fitzDocInternal)(unsafe.Pointer(doc))

	if internal.ctx == nil || internal.doc == nil {
		return 0, fmt.Errorf("invalid document or context (ctx or doc is nil)")
	}

	// 调用 C 扩展函数
	count := C.custom_get_page_count(internal.ctx, internal.doc)
	if count < 0 {
		return 0, fmt.Errorf("failed to count pages via mupdf C API")
	}

	return int(count), nil
}

func Custom_drop_path(doc *fitz.Document, pix_ptr unsafe.Pointer) {
	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	if internal.doc == nil || internal.ctx == nil {
		return
	}

	C.custom_drop_pixmap(internal.ctx, pix_ptr)
}

func Drop_down_tiles(doc *fitz.Document) error {
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}
	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc
	if ctx_ptr == nil || doc_ptr == nil {
		return fmt.Errorf("ctx or doc is nil")
	}
	C.custom_drop_down_tiles(ctx_ptr, doc_ptr)
	return nil
}

func Drop_down_pages(doc *fitz.Document) error {
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}
	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc
	if ctx_ptr == nil || doc_ptr == nil {
		return fmt.Errorf("ctx or doc is nil")
	}
	// C.custom_drop_drawn_tiles(ctx_ptr, doc_ptr)
	// C.custom_drop_drawn_pages()
	return nil
}

func RenderTile(doc *fitz.Document, page_num int, dpi float64, rect types.Rectangle) (*image.RGBA, error) {
	if doc == nil {
		return nil, fmt.Errorf("doc is nil")
	}

	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	internal.mtx.Lock()
	defer internal.mtx.Unlock()

	// ctx := (*C.struct_fitz_context)(internal.ctx)
	// docPtr := (*C.struct_fz_document)(internal.doc)

	return nil, nil
}

// 自定义: 渲染pdf整页
func Custom_pdf_render_full(doc_inst *fitz.Document, page_num int, dpi float64) *image.RGBA {
	// 结构体映射
	internal := (*fitzDocInternal)(unsafe.Pointer(doc_inst))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc

	// 转化为C指针
	ctx := (*C.fz_context)(ctx_ptr)
	doc := (*C.fz_document)(doc_ptr)

	// 加载页
	page := C.fz_load_page(ctx, doc, C.int(page_num))
	defer C.fz_drop_page(ctx, page)

	// 构建渲染区域矩阵
	ctm := C.fz_scale(C.float(dpi)/72, C.float(dpi)/72)
	bounds := C.fz_bound_page(ctx, page)
	bounds = C.fz_transform_rect(bounds, ctm)
	bbox := C.fz_round_rect(bounds)

	// 创建pix
	pix := C.fz_new_pixmap_with_bbox(ctx, C.fz_device_rgb(ctx), bbox, nil, 1)
	defer C.fz_drop_pixmap(ctx, pix)
	C.fz_clear_pixmap_with_value(ctx, pix, 0xff)

	// 绘制功能,内容,区域,目标区域
	device := C.fz_new_draw_device(ctx, ctm, pix)
	defer C.fz_drop_device(ctx, device)
	C.fz_enable_device_hints(ctx, device, C.FZ_NO_CACHE)

	// 进行渲染
	drawMatrix := C.fz_identity
	C.fz_run_page_contents(ctx, page, device, drawMatrix, nil)

	// 得到图像行跨距
	stride := int(C.fz_pixmap_stride(ctx, pix))

	// 拿到图像像素数组的起始地址
	pixels := C.fz_pixmap_samples(ctx, pix)

	// 创建go图像
	w, h := int(pix.w), int(pix.h)
	img := image.NewRGBA(image.Rect(
		int(bbox.x0),
		int(bbox.y0),
		int(bbox.x1),
		int(bbox.y1),
	))

	// 复制内存并操作
	// stride 跨距, 等于img每行的的字节长
	// C中为了性能考量,字节长是32的整数倍,每一行可能会有多余的无意义填充
	// 而go的img是紧密排列的,所以需要按各自的跨距拷贝内存
	if pixels != nil {
		for y := 0; y < h; y += 1 {
			row_dest := img.Pix[y*img.Stride : y*img.Stride+w*4]

			// 指针转数字+跨距, 之后转化为指针
			row_src := unsafe.Pointer(uintptr(unsafe.Pointer(pixels)) + uintptr(y*stride))

			// 裸内存转化为go切片,给裸指针包含元数据(类型,长度)
			srcSlice := unsafe.Slice((*byte)(row_src), w*4)
			copy(row_dest, srcSlice)
		}
	}

	C.fz_close_device(ctx, device)
	return img
}

// 自定义:渲染pdf瓦片
func Custom_pdf_render_tile(doc_inst *fitz.Document, page_num int, dpi float64, rect [4]float64) *image.RGBA {
	internal := (*fitzDocInternal)(unsafe.Pointer(doc_inst))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc

	ctx := (*C.fz_context)(ctx_ptr)
	doc := (*C.fz_document)(doc_ptr)

	page := C.fz_load_page(ctx, doc, C.int(page_num))
	defer C.fz_drop_page(ctx, page)
	ctm := C.fz_scale(C.float(dpi)/72, C.float(dpi)/72)

	fz_rect := C.fz_rect{
		x0: C.float(rect[0]),
		y0: C.float(rect[1]),
		x1: C.float(rect[2]),
		y1: C.float(rect[3]),
	}

	trans_rect := C.fz_transform_rect(fz_rect, ctm)
	bbox := C.fz_round_rect(trans_rect)

	pix := C.fz_new_pixmap_with_bbox(ctx, C.fz_device_rgb(ctx), bbox, nil, 1)
	defer C.fz_drop_pixmap(ctx, pix)
	C.fz_clear_pixmap_with_value(ctx, pix, 0xff)

	device := C.fz_new_draw_device(ctx, ctm, pix)
	defer C.fz_drop_device(ctx, device)
	C.fz_enable_device_hints(ctx, device, C.FZ_NO_CACHE)

	drawMatrix := C.fz_identity
	C.fz_run_page_contents(ctx, page, device, drawMatrix, nil)

	stride := int(C.fz_pixmap_stride(ctx, pix))
	pixels := C.fz_pixmap_samples(ctx, pix)

	w, h := int(pix.w), int(pix.h)

	img := image.NewRGBA(image.Rect(
		int(bbox.x0),
		int(bbox.y0),
		int(bbox.x1),
		int(bbox.y1),
	))

	if pixels != nil {
		for y := 0; y < h; y += 1 {
			row_dest := img.Pix[y*img.Stride : y*img.Stride+w*4]
			row_src := unsafe.Pointer(uintptr(unsafe.Pointer(pixels)) + uintptr(y*stride))
			srcSlice := unsafe.Slice((*byte)(row_src), w*4)
			copy(row_dest, srcSlice)
		}
	}

	C.fz_close_device(ctx, device)
	return img
}

// 废弃:渲染pdf整页
func Custom_render_full(doc *fitz.Document, page_num int, dpi float64, tile_w int, tile_h int) (*image.RGBA, error) {
	if doc == nil {
		return nil, nil
	}

	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc

	start := time.Now()
	C.custom_pdf_load_page(ctx_ptr, doc_ptr, C.int(page_num))
	end := time.Since(start)
	log.Println(end)

	img := image.NewRGBA(image.Rect(0, 0, tile_w, tile_h))

	C.custom_render_full(
		ctx_ptr, doc_ptr,
		C.int(page_num),
		C.float(dpi),
		C.int(tile_w),
		C.int(tile_h),
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
	)

	return img, nil
}
