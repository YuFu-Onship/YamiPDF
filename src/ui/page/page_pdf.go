package page

import (
	rumia "app/src"
	"app/src/ui"
	"runtime"
)

type PagePDF struct {
	Style       *ui.Style
	page_router *PageRouter

	viewer *ui.PDFViewer
	id     string
}

func NewPagePDF(style *ui.Style, page_router *PageRouter) *PagePDF {
	self := PagePDF{
		Style:       style,
		page_router: page_router,
	}

	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Archives of Witch Hat Atelier (Kamome Shirahama).pdf"
	pdf_path := "E:/Users/YUFU/Documents/Books/comic/Witch Hat Atelier, Vol.1 (Kamome Shirahama) (z-library.sk, 1lib.sk, z-lib.sk).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/单页.pdf"

	// doc_id, _ := self.page_router.BridgeTrunk.Add_new_document(pdf_path)
	self.page_router.BridgeTrunk.Add_render_inst(pdf_path, rumia.FileType_PDF)

	// self.id = doc_id

	self.viewer = ui.NewPDFView(pdf_path, self.Style, self.page_router.BridgeTrunk, self.page_router.Tools)

	// self.viewer = ui.NewPDFView("id", style,self.page_router.Tools,)
	return &self
}

func (self *PagePDF) Layout(gtx C) D {
	if self.viewer != nil {
		return self.viewer.Layout(gtx)
	}
	return D{}
}

func (self *PagePDF) OnEnter() {

}
func (self *PagePDF) OnExit() {
}
func (self *PagePDF) Destory() {
	self.page_router.BridgeTrunk.Delete_document(self.id)
	runtime.GC()
}
func (self *PagePDF) API_set_bridge_ui(bridge_ui ui.BridgeUI) {}
func (self *PagePDF) API_set_parent_manager(pm interface{})   {}
