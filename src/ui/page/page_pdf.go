package page

import (
	rumia "app/src"
	"app/src/ui"
	"runtime"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PagePDF struct {
	Style       *ui.Style
	page_router *PageRouter

	viewer *ui.PDFViewer
	id     string

	panel_tool       *ToolPanel_PagePDF
	clickable_expand *widget.Clickable
}

func NewPagePDF(style *ui.Style, page_router *PageRouter) *PagePDF {
	self := PagePDF{
		Style:       style,
		page_router: page_router,

		clickable_expand: &widget.Clickable{},
	}
	self.panel_tool = NewToolPanel_PagePDF(style)

	pdf_path := "E:/Users/YUFU/Documents/Books/comic/Archives of Witch Hat Atelier (Kamome Shirahama).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/comic/Witch Hat Atelier, Vol.1 (Kamome Shirahama) (z-library.sk, 1lib.sk, z-lib.sk).pdf"
	// pdf_path := "E:/Users/YUFU/Documents/Books/单页.pdf"

	// doc_id, _ := self.page_router.BridgeTrunk.Add_new_document(pdf_path)
	self.page_router.BridgeTrunk.Add_render_inst(pdf_path, rumia.FileType_PDF)

	self.viewer = ui.NewPDFView(pdf_path, self.Style, self.page_router.BridgeTrunk, self.page_router.Tools)
	self.id = pdf_path
	return &self
}

func (self *PagePDF) Layout(gtx C) D {
	// if self.viewer != nil {
	// 	return self.viewer.Layout(gtx)
	// }

	if self.clickable_expand.Clicked(gtx) {
		self.panel_tool.API_set_state(!self.panel_tool.API_get_state())
	}
	if self.clickable_expand.Hovered() {
	}

	btn_panel := material.Button(self.Style.Theme, self.clickable_expand, "Expand")

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(btn_panel.Layout),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, self.viewer.Layout),
				layout.Rigid(self.panel_tool.Layout),
			)
		}),
	)

	return dims
}

func (self *PagePDF) OnEnter() {

}
func (self *PagePDF) OnExit() {
}
func (self *PagePDF) Destory() {
	// self.page_router.BridgeTrunk.Delete_document(self.id)
	self.page_router.BridgeTrunk.API_close_inst(self.id)
	runtime.GC()
}
func (self *PagePDF) API_set_bridge_ui(bridge_ui ui.BridgeUI) {}

func (self *PagePDF) API_set_parent_manager(pm interface{}) {}

// 目录面板

// 工具面板
type ToolPanel_PagePDF struct {
	Style *ui.Style

	Panel *ui.ToolPanel

	Clickable_next  *widget.Clickable
	Clickable_last  *widget.Clickable
	Dropdown_shader *ui.DropDownChoice
	Dropdown_layout *ui.DropDownChoice
	Dropdown_order  *ui.DropDownChoice
}

func NewToolPanel_PagePDF(style *ui.Style) *ToolPanel_PagePDF {
	self := ToolPanel_PagePDF{
		Style: style,
	}
	self.Panel = ui.NewToolPanel(style).API_set_width(300)
	return &self
}

func (self *ToolPanel_PagePDF) Layout(gtx C) D {
	return self.Panel.Layout(gtx)
}

func (self *ToolPanel_PagePDF) API_set_state(is_expand bool) *ToolPanel_PagePDF {
	self.Panel.API_set_state(is_expand)
	return self
}
func (self *ToolPanel_PagePDF) API_get_state() bool {
	return self.Panel.API_get_state()
}
