// 是的,公共包 common_package 名字太长了

// 公共结构体
package rumia

import (
	"image"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// 整页渲染任务
type RenderFullTask struct {
	Page    int
	DPI     float64
	Call_fn func(img *image.RGBA)
}

// 瓦片渲染任务
type RenderTileTask struct {
	Page       int
	Dpi        float64
	Rect       *types.Rectangle //pdf页面的坐标系,原点在左下角
	Rect_ratio *types.Rectangle //图形学坐标系,原点在左上角
	Call_fn    func(img *image.RGBA)
}
