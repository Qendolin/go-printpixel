package layout

type Unit int

const (
	Pixel Unit = iota
	Percent
)

type TrackDef struct {
	Value float32
	Unit  Unit
}

type Grid struct {
	Columns  []TrackDef
	Rows     []TrackDef
	Children [][]Layoutable
	SimpleBox
}

func NewGrid(cols []TrackDef, rows []TrackDef) Grid {
	children := make([][]Layoutable, len(cols))
	for x := range children {
		children[x] = make([]Layoutable, len(rows))
	}
	return Grid{
		Columns:  cols,
		Rows:     rows,
		Children: children,
	}
}

func (grid Grid) Layout() {
	colTrackPositions := make([]int, len(grid.Columns)+1)
	colTrackPositions[0] = 0
	rowTrackPositions := make([]int, len(grid.Rows)+1)
	rowTrackPositions[0] = 0

	var acc int
	for col, def := range grid.Columns {
		var delta int

		switch def.Unit {
		case Pixel:
			delta = int(def.Value)
		case Percent:
			delta = int(def.Value * float32(grid.Width()))
		}
		acc += delta
		colTrackPositions[col+1] = acc
	}
	acc = 0
	for col, def := range grid.Rows {
		var delta int

		switch def.Unit {
		case Pixel:
			delta = int(def.Value)
		case Percent:
			delta = int(def.Value * float32(grid.Height()))
		}
		acc += delta
		rowTrackPositions[col+1] = acc
	}

	for x, col := range grid.Children {
		if col == nil {
			continue
		}
		for y, child := range col {
			if child == nil {
				continue
			}
			posX := colTrackPositions[x]
			posY := rowTrackPositions[y]
			child.SetX(posX)
			child.SetY(posY)
			child.SetWidth(colTrackPositions[x+1] - posX)
			child.SetHeight(rowTrackPositions[y+1] - posY)
			if l, ok := child.(Layouter); ok {
				l.Layout()
			}
		}
	}
}
