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
	Cols     []TrackDef
	Rows     []TrackDef
	Children [][]Layoutable
	SimpleBox
}

func NewGrid(cols []TrackDef, rows []TrackDef) Grid {
	g := Grid{
		Cols: cols,
		Rows: rows,
	}
	g.Init()

	return g
}

func (grid *Grid) Init() {
	children := make([][]Layoutable, len(grid.Cols))
	for x := range children {
		children[x] = make([]Layoutable, len(grid.Rows))
	}
	grid.Children = children
}

func (grid Grid) Layout() []Layoutable {
	colTrackPositions := make([]int, len(grid.Cols)+1)
	colTrackPositions[0] = 0
	rowTrackPositions := make([]int, len(grid.Rows)+1)
	rowTrackPositions[0] = 0

	var acc int
	for col, def := range grid.Cols {
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

	childs := make([]Layoutable, 0)
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
			childs = append(childs, child)
		}
	}
	return childs
}
