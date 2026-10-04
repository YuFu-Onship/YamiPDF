package trunk

import (
	rumia "app/src"
	"image"
)

// 渲染管理器 ---------------------------------

// 添加新的渲染实例
func (self Trunk) Add_new_document(doc_path string) (string, error) {
	return self.RenderManager.Add_new_document(doc_path)
}

// 删除渲染实例
func (self Trunk) Delete_document(id string) {
	self.RenderManager.Delete_document(id)
}

// 获取文件信息
//   - meta_data  map[string]string 元数据[title,author]
//   - page_num int 总页数
func (self Trunk) Get_info(id string) map[string]any {
	return self.RenderManager.Get_info(id)
}

// 渲染指定页面
func (self Trunk) Render_page(id string, page int) (*image.RGBA, error) {
	return self.RenderManager.Render_page(id, page)
}

// 得到pdf高度
func (self Trunk) Get_size(id string, page int) (image.Rectangle, error) {
	return self.RenderManager.Get_size(id, page)
}

func (self Trunk) Add_render_page(id string, page int, call_fn func(img *image.RGBA)) {
	self.RenderManager.Add_render_page(id, page, call_fn)
}

// func (self Trunk) Add_render_page_list(id string, dpi float64, pages []int, call_fn func(page int, dpi float64, img *image.RGBA)) {
func (self Trunk) Add_render_page_list(id string, tasks []*rumia.RenderFullTask) {
	self.RenderManager.Add_render_page_list(id, tasks)
}

// func (self Trunk) Add_render_tile(id string, dpi float64, page int, rect *types.Rectangle, call_fn func(page int, dpi float64, img *image.RGBA)) {
func (self Trunk) Add_render_tile_list(id string, tasks []*rumia.RenderTileTask) {
	self.RenderManager.Add_render_tile_list(id, tasks)
}
