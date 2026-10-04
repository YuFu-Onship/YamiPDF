package render

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/gen2brain/go-fitz"
)

// #cgo CFLAGS: -I${SRCDIR}/../../
// #include "custom_fitz.h"
import "C"

// 严格按照 Document 结构体顺序和类型对齐内存, 越权强读fitz内部的指针
type fitzDocInternal struct {
	ctx    unsafe.Pointer // *C.struct_fz_context
	data   []byte         // []byte
	doc    unsafe.Pointer // *C.struct_fz_document
	mtx    sync.Mutex     // sync.Mutex
	stream unsafe.Pointer // *C.fz_stream
}

func GetCustomPageCount(doc *fitz.Document) (int, error) {
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

func Drop_down_tiles(doc *fitz.Document) error {
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}
	internal := (*fitzDocInternal)(unsafe.Pointer(doc))
	ctx_ptr, doc_ptr := internal.ctx, internal.doc
	if ctx_ptr == nil || doc_ptr == nil {
		return fmt.Errorf("ctx or doc is nil")
	}
	C.custom_drop_drawn_tiles(ctx_ptr, doc_ptr)
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
