package rumia

import (
	"image"
)

// 外部项目接口
type BridgeTrunk interface {
	Test()

	// 获得exe或main所在的文件路径
	GetRootPath() string

	// 渲染管理器 ---------------------------------

	// 添加新的渲染实例
	Add_new_document(doc_path string) (string, error)

	// 删除渲染实例
	Delete_document(id string)

	// 获取文件信息
	//   - meta_data  map[string]string 元数据[title,author]
	//   - page_num int 总页数
	Get_info(id string) map[string]any

	// 渲染指定页面
	Render_page(id string, page int) (*image.RGBA, error)

	// 得到指定页尺寸
	Get_size(id string, page int) (image.Rectangle, error)

	// 添加新的渲染页
	Add_render_page(id string, page int, call_fn func(img *image.RGBA))

	// 添加新的渲染页列表
	// Add_render_page_list(id string, dpi float64, pages []int, call_fn func(page int, dpi float64, img *image.RGBA))
	Add_render_page_list(id string, tasks []*RenderTask_PDF_Full)

	// 添加新的渲染瓦块
	// Add_render_tile(id string, dpi float64, page int, rect *types.Rectangle, call_fn func(page int, dpi float64, img *image.RGBA))
	Add_render_tile_list(id string, tasks []*RenderTask_PDF_Tile)

	// 新的渲染实例管理器 ---------------------------------------------------------

	// 增加渲染实例
	Add_render_inst(id string, file_type FileType)

	// 渲染整页
	Render_pdf_full_page(id string, task RenderTask_PDF_Full)

	// 渲染整页列表
	Rebder_pdf_full_pages(id string, tasks []*RenderTask_PDF_Full)

	// 得到pdf元信息
	API_get_pdf_info(id string) map[string]any

	// 得到pdf指定页的尺寸
	API_get_pdf_page_size(id string, page int) image.Rectangle

	// 覆盖渲染任务列表
	API_set_tasks(tasks []BaseTask)

	// 追加渲染任务
	API_add_task(task BaseTask)

	// 关闭渲染实例
	API_close_inst(id string)

	// ui相关 ---------------------------------------------------------------------

	// 刷新窗口
	Refresh_window()
}
