package component

import "github.com/gizak/termui/v3/widgets"

type resourceDetailPanel struct {
	*widgets.Paragraph
}

func buildDetailParagraph() *resourceDetailPanel {
	p := widgets.NewParagraph()
	p.Title = "Detail"
	return &resourceDetailPanel{
		Paragraph: p,
	}
}

func (r *resourceDetailPanel) RefreshData(data string) {
	r.Text = data
}
