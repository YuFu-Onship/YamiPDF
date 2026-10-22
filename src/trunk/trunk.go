package trunk

import (
	"app/src/render"
	"app/src/tool"
	"app/src/ui/page"
	"fmt"
	"time"
)

type Trunk struct {
	RootPath          string
	Page              *page.Page
	RenderManager     *render.RenderManager
	RenderInstManager *render.RenderInstManager

	// ShaderRender render.ShaderRender
}

func NewTrunk(path string) *Trunk {
	self := Trunk{
		RootPath: path,
	}

	// self.RenderManager = render.NewRenderManager()
	self.RenderInstManager = render.NewRenderInstManager()

	self.Page = page.NewPage()
	self.Page.API_set_trunk(self)

	return &self
}

func (self *Trunk) Run() {
	mem_monitor, _ := tool.NewMemoryMonitor()
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {

			wsMB, _, goAllocMB, goSysMB, _, _ := mem_monitor.GetMem()
			text := fmt.Sprintf("Sys: %.2f MB | GoAlloc: %.2f MB | GoSys: %.2f MB", wsMB, goAllocMB, goSysMB)
			self.Page.Set_deco_text(text)
		}
	}()

	// self.RenderManager.Run()
	self.RenderInstManager.Run()

	self.Page.Run()
}
