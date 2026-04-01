package scene

type Unit int

const (
	Pixel Unit = iota
	Percent
)

type Layouter interface {
	Layout() (children []Layoutable)
}

type Layoutable interface {
	X() int
	SetX(int)
	Y() int
	SetY(int)
	Width() int
	SetWidth(int)
	Height() int
	SetHeight(int)
}

type Box struct {
	width  int
	height int
	x      int
	y      int
}

func (box *Box) Width() int {
	return box.width
}

func (box *Box) Height() int {
	return box.height
}

func (box *Box) SetWidth(width int) {
	box.width = width
}

func (box *Box) SetHeight(height int) {
	box.height = height
}

func (box *Box) SetX(x int) {
	box.x = x
}

func (box *Box) SetY(y int) {
	box.y = y
}

func (box *Box) X() int {
	return box.x
}

func (box *Box) Y() int {
	return box.y
}

type Tree struct {
	Root  Node
	Nodes []*Node
}

type Node struct {
	Value      interface{}
	Layouter   Layouter
	Layoutable Layoutable
	Children   []Node
	Depth      int
}

func Layout(root Layouter) Tree {
	tree := Tree{
		Root:  Node{Value: root, Layouter: root, Depth: 0},
		Nodes: make([]*Node, 1),
	}
	tree.Nodes[0] = &tree.Root

	if l, ok := root.(Layoutable); ok {
		tree.Root.Layoutable = l
	}

	var walk func(root *Node, depth int)
	walk = func(parent *Node, depth int) {
		children := parent.Layouter.Layout()
		if children == nil {
			return
		}
		parent.Children = make([]Node, len(children))

		q := make([]*Node, len(children))
		j := 0
		for i, c := range children {
			if l, ok := c.(Layouter); ok {
				parent.Children[i] = Node{Value: c, Layouter: l, Layoutable: c, Depth: depth}
				n := &parent.Children[i]
				tree.Nodes = append(tree.Nodes, n)
				q[j] = n
				j++
			} else {
				parent.Children[i] = Node{Value: c, Layoutable: c, Depth: depth}
				tree.Nodes = append(tree.Nodes, &parent.Children[i])
			}
		}

		for _, n := range q[:j] {
			walk(n, depth+1)
		}
	}
	walk(&tree.Root, 1)

	return tree
}
