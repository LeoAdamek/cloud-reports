package printer

import "errors"

// Errors

var (
// ErrTableTooWide means the table is too wide to print
ErrTableTooWide = errors.New("table is too wide to print")
)

const rowHeight = 25

// Table represents a data table for printing
type Table struct {
	Columns []TableColumn
	Cells [][]string
	// Enables Borders
	Boarders bool
}

// TableColumn represents a single column of a table
type TableColumn struct {
	Title string
	Width int
}

func (t Table) TotalWidth() int {
	w := 0

	for _, c := range t.Columns {
		w += c.Width
	}

	return w
}


func (p *Printer) DrawTable(t Table) error {
	if t.TotalWidth() > 480 {
		return ErrTableTooWide
	}

	p.SetOrigin()

	if t.Boarders {
		p.drawTableGrid(t)
	}

	tw := 0
	for _, c := range t.Columns {
		p.MoveTo(tw+10, 0)
		p.GraphicalWrite(c.Title, 1, TextLeftToRight)
		tw += c.Width
	}


	for rowi, row := range t.Cells {
		tw = 0
		for celli, cell := range row {
			p.MoveTo(tw+10, -rowHeight*(rowi+1))

			cellW := t.Columns[celli].Width

			cellLength := cellW / 12

			text := cell
			size := 1

			if len(cell) > cellLength {
				size = 0

				if len(cell) > 2*cellLength {
					text = cell[:2*cellLength]
				}
			}

			p.GraphicalWrite(text, size, TextLeftToRight)
			tw += cellW
		}
	}

	return nil
}


func (p *Printer) drawTableGrid(t Table) {

	totalWidth := t.TotalWidth()
	for i := 0; i < len(t.Cells)+1; i++ {
		p.MoveTo(0, -rowHeight*i)
		p.DrawTo(totalWidth, -rowHeight*i)
	}

	tw := 0

	for _, c := range t.Columns {
		p.MoveTo(tw, 0)
		p.DrawTo(tw, -rowHeight*len(t.Cells)+1)

		tw += c.Width
	}

	p.MoveTo(totalWidth, 0)
	p.DrawTo(totalWidth, -rowHeight*len(t.Cells)+1)
}

