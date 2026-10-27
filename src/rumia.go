// 是的,公共包 common_package 名字太长了

// 公共结构体
package rumia

import (
	"image"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// FileType -------------------------------------------------------------
type FileType int

const (
	FileType_PDF FileType = iota
	FileType_EPUB
	FileType_TXT
)

var FileType_table map[FileType]string = map[FileType]string{
	FileType_PDF:  "PDF",
	FileType_EPUB: "EPUB",
	FileType_TXT:  "TXT",
}

// ----------------------------------------------------------------------

// basetask
type BaseTask struct {
	Id   string
	Task any
}

// 整页渲染任务
type RenderTask_PDF_Full struct {
	Page    int
	DPI     float64
	Call_fn func(img *image.RGBA)
}

// 瓦片渲染任务
type RenderTask_PDF_Tile struct {
	Page       int
	Dpi        float64
	Rect       *types.Rectangle //pdf页面的坐标系,原点在左下角
	Rect_ratio *types.Rectangle //图形学坐标系,原点在左上角
	Call_fn    func(img *image.RGBA)
}

// 书本信息
type BookInfo struct {
	Title    string
	Author   string
	Path     string
	FileType FileType
}
