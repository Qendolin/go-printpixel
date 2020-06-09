package layout

type Box interface {
	Width() int
	SetWidth(int)
	Height() int
	SetHeight(int)
}

type Layouter interface {
	Layout()
}

type Layoutable interface {
	SetX(int)
	SetY(int)
	X() int
	Y() int
	Box
}

type SimpleBox struct {
	width  int
	height int
	x      int
	y      int
}

func (box *SimpleBox) Width() int {
	return box.width
}

func (box *SimpleBox) Height() int {
	return box.height
}

func (box *SimpleBox) SetWidth(width int) {
	box.width = width
}

func (box *SimpleBox) SetHeight(height int) {
	box.height = height
}

func (box *SimpleBox) SetX(x int) {
	box.x = x
}

func (box *SimpleBox) SetY(y int) {
	box.y = y
}

func (box *SimpleBox) X() int {
	return box.x
}

func (box *SimpleBox) Y() int {
	return box.y
}
