package component

import (
	"image"
	"log"

	. "github.com/gizak/termui/v3"
)

type resourcePanel struct {
	*Block

	Header []string
	Rows   [][]string

	ColWidths []int
	ColGap    int
	PadLeft   int

	ShowCursor  bool
	CursorColor Color

	UniqueCol    int    // the column used to uniquely identify each table row
	SelectedItem string // used to keep the cursor on the correct item if the data changes
	SelectedRow  int
	TopRow       int // used to indicate where in the table we are scrolled at

	ColResizer func()
}

// NewDebugTable returns a new resourcePanel instance

func BuildResourcePanel() *resourcePanel {
	rPanel := &resourcePanel{
		Block:       NewBlock(),
		ColGap:      3,
		SelectedRow: 0,
		TopRow:      0,
		UniqueCol:   0,
		ShowCursor:  true,
		CursorColor: ColorYellow,
		ColResizer:  func() {},
	}

	rPanel.Header = []string{""}
	rPanel.Rows = [][]string{{""}}
	rPanel.ColWidths = []int{0}
	return rPanel
}

func (rp *resourcePanel) Draw(buf *Buffer) {
	rp.Block.Draw(buf)

	rp.ColResizer()

	// finds exact column starting position
	colXPos := []int{}
	cur := 1 + rp.PadLeft
	for _, w := range rp.ColWidths {
		colXPos = append(colXPos, cur)
		cur += w
		cur += rp.ColGap
	}

	// prints header
	for i, h := range rp.Header {
		width := rp.ColWidths[i]
		if width == 0 {
			continue
		}
		// don't render column if it doesn't fit in widget
		if width > (rp.Inner.Dx()-colXPos[i])+1 {
			continue
		}
		buf.SetString(
			h,
			NewStyle(Theme.Default.Fg, ColorClear, ModifierBold),
			image.Pt(rp.Inner.Min.X+colXPos[i]-1, rp.Inner.Min.Y),
		)
	}

	if rp.TopRow < 0 {
		log.Printf("table widget TopRow value less than 0. TopRow: %v", rp.TopRow)
		return
	}

	// prints each row
	for rowNum := rp.TopRow; rowNum < rp.TopRow+rp.Inner.Dy()-1 && rowNum < len(rp.Rows); rowNum++ {
		row := rp.Rows[rowNum]
		y := (rowNum + 2) - rp.TopRow

		// prints cursor
		style := NewStyle(Theme.Default.Fg)
		if rp.ShowCursor {
			if rowNum == rp.SelectedRow {
				style.Fg = rp.CursorColor
				for i, width := range rp.ColWidths {
					if width == 0 {
						continue
					}
					buf.SetString(
						TrimString(row[i], width),
						style,
						image.Pt(rp.Inner.Min.X+colXPos[i]-1, rp.Inner.Min.Y+y-1),
					)
				}
				rp.SelectedItem = row[rp.UniqueCol]
				rp.SelectedRow = rowNum
			}
		}

		// prints each col of the row
		for i, width := range rp.ColWidths {
			if width == 0 {
				continue
			}
			// don't render column if width is greater than distance to end of widget
			if width > (rp.Inner.Dx()-colXPos[i])+1 {
				continue
			}
			r := TrimString(row[i], width)
			buf.SetString(
				r,
				style,
				image.Pt(rp.Inner.Min.X+colXPos[i]-1, rp.Inner.Min.Y+y-1),
			)
		}
	}
}

// Scrolling ///////////////////////////////////////////////////////////////////

// calcPos is used to calculate the cursor position and the current view into the table.
func (rp *resourcePanel) calcPos() {
	rp.SelectedItem = ""

	if rp.SelectedRow < 0 {
		rp.SelectedRow = 0
	}
	if rp.SelectedRow < rp.TopRow {
		rp.TopRow = rp.SelectedRow
	}

	if rp.SelectedRow > len(rp.Rows)-1 {
		rp.SelectedRow = len(rp.Rows) - 1
	}
	if rp.SelectedRow > rp.TopRow+(rp.Inner.Dy()-2) {
		rp.TopRow = rp.SelectedRow - (rp.Inner.Dy() - 2)
	}
}

func (rp *resourcePanel) ScrollUp() {
	rp.SelectedRow--
	rp.calcPos()
}

func (rp *resourcePanel) ScrollDown() {
	rp.SelectedRow++
	rp.calcPos()
}

func (rp *resourcePanel) ScrollTop() {
	rp.SelectedRow = 0
	rp.calcPos()
}

func (rp *resourcePanel) ScrollBottom() {
	rp.SelectedRow = len(rp.Rows) - 1
	rp.calcPos()
}

func (rp *resourcePanel) ScrollHalfPageUp() {
	rp.SelectedRow = rp.SelectedRow - (rp.Inner.Dy()-2)/2
	rp.calcPos()
}

func (rp *resourcePanel) ScrollHalfPageDown() {
	rp.SelectedRow = rp.SelectedRow + (rp.Inner.Dy()-2)/2
	rp.calcPos()
}

func (rp *resourcePanel) ScrollPageUp() {
	rp.SelectedRow -= (rp.Inner.Dy() - 2)
	rp.calcPos()
}

func (rp *resourcePanel) ScrollPageDown() {
	rp.SelectedRow += (rp.Inner.Dy() - 2)
	rp.calcPos()
}

func (rp *resourcePanel) HandleClick(x, y int) {
	x = x - rp.Min.X
	y = y - rp.Min.Y
	if (x > 0 && x <= rp.Inner.Dx()) && (y > 0 && y <= rp.Inner.Dy()) {
		rp.SelectedRow = (rp.TopRow + y) - 2
		rp.calcPos()
	}
}

func (rp *resourcePanel) RefreshPanelData(newHeader []string, newData [][]string, newColWidths []int) {
	rp.ColWidths = newColWidths
	rp.Header = newHeader
	rp.Rows = newData
}

func (rp *resourcePanel) selectedToggle() {
	if rp.BorderStyle == NewStyle(selectedPanelColor) {
		rp.BorderStyle = NewStyle(ColorClear)
	} else {
		rp.BorderStyle = NewStyle(selectedPanelColor)
	}
}
