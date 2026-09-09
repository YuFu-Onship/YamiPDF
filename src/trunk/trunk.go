package trunk

import (
	"app/src/render"
	"app/src/ui"
)

type Trunk struct {
	RootPath      string
	Page          *ui.Page
	RenderManager *render.RenderManager
	// ShaderRender render.ShaderRender
}

func NewTrunk(path string) *Trunk {
	self := Trunk{
		RootPath: path,
	}
	self.RenderManager = render.NewRenderManager()
	self.Page = ui.NewPage(self)
	return &self
}

func (self *Trunk) Run() {
	self.RenderManager.Run()
	self.Page.Run()
}
