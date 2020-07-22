package scene

type Absolute struct {
	SimpleBox
	Child Layoutable
	DX    float32
	DY    float32
	W     float32
	H     float32
	Unit  Unit
}

func (abs *Absolute) Layout() []Layoutable {
	if abs.Child == nil {
		return nil
	}
	switch abs.Unit {
	case Pixel:
		abs.Child.SetWidth(int(abs.W))
		abs.Child.SetHeight(int(abs.H))
		abs.Child.SetX(abs.x + int(abs.DX))
		abs.Child.SetY(abs.y + int(abs.DY))
	case Percent:
		abs.Child.SetWidth(int(float32(abs.width) * abs.W))
		abs.Child.SetHeight(int(float32(abs.height) * abs.H))
		abs.Child.SetX(abs.x + int(float32(abs.width)*abs.DX))
		abs.Child.SetY(abs.y + int(float32(abs.height)*abs.DY))
	}
	return []Layoutable{abs.Child}
}

type Stack struct {
	SimpleBox
	Children []Layoutable
}

func Stacked(c ...Layoutable) *Stack {
	return &Stack{
		Children: c,
	}
}

func (stack *Stack) Layout() []Layoutable {
	if stack.Children == nil {
		return nil
	}

	for _, c := range stack.Children {
		c.SetWidth(stack.width)
		c.SetHeight(stack.height)
		c.SetX(stack.x)
		c.SetY(stack.y)
	}

	//reverse children so that the first is ontop
	c := make([]Layoutable, len(stack.Children))
	copy(c, stack.Children)
	for i := len(c)/2 - 1; i >= 0; i-- {
		opp := len(c) - 1 - i
		c[i], c[opp] = c[opp], c[i]
	}

	return c
}

type Center struct {
	SimpleBox
	Child Layoutable
}

func Centered(c Layoutable) *Center {
	return &Center{
		Child: c,
	}
}

func (c *Center) Layout() []Layoutable {
	if c.Child == nil {
		return nil
	}

	c.Child.SetWidth(c.width)
	c.Child.SetHeight(c.height)
	c.Child.SetX(c.x - c.width/2)
	c.Child.SetY(c.y - c.height/2)

	return []Layoutable{c.Child}
}

type Layer struct {
	SimpleBox
	Child Layoutable
}

func (l *Layer) Layout() []Layoutable {
	if l.Child == nil {
		return nil
	}

	l.Child.SetWidth(l.width)
	l.Child.SetHeight(l.height)
	l.Child.SetX(l.x)
	l.Child.SetY(l.y)

	return []Layoutable{l.Child}
}
