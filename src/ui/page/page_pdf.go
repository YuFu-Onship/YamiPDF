package page

import "app/src/ui"

type PagePDF struct {
	Style  *ui.Style
	viewer *ui.PDFViewer
}

func NewPagePDF(style *ui.Style) *PagePDF {
	self := PagePDF{
		Style: style,
	}
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
func (self *PagePDF) OnExit()                                 {}
func (self *PagePDF) Destory()                                {}
func (self *PagePDF) API_set_bridge_ui(bridge_ui ui.BridgeUI) {}
func (self *PagePDF) API_set_parent_manager(pm interface{})   {}
