package ui

import (
	rumia "app/src"
	"image"
)

// 内部接口
type UIBridge interface {
	// 切换颜色模式
	ToggleTheme()

	// 窗口大小是否变动
	//   - true 发生变化
	//   - false 没有变化
	IsWinSizeChange() bool

	// 请求窗口级重绘
	WindowRefresh()

	// 增加显示信息
	ShowInfo(id any, text string)

	// 更改标题栏文本
	Set_deco_text(text string)

	// 更改软件窗口文本
	Set_title_text(text string)
}

// 外部项目接口
type TrunkBridge interface {
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
	Add_render_page_list(id string, tasks []*rumia.RenderFullTask)

	// 添加新的渲染瓦块
	// Add_render_tile(id string, dpi float64, page int, rect *types.Rectangle, call_fn func(page int, dpi float64, img *image.RGBA))
	Add_render_tile_list(id string, tasks []*rumia.RenderTileTask)
}
