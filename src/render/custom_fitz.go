package render

import (
	"fmt"
	"image"
	"sync"
	"unsafe"

	"github.com/gen2brain/go-fitz"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// #cgo CFLAGS: -I${SRCDIR}/../../
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

func Custom_render_full(doc *fitz.Document, page_num int, dpi float64, tile_w int, tile_h int) (*image.RGBA, error) {
	if doc == nil {
		return nil, nil
	}

	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc

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
