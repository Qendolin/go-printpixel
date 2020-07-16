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

func (stack Stack) Layout() []Layoutable {
	if stack.Children == nil {
		return nil
	}

	for _, c := range stack.Children {
		c.SetWidth(stack.width)
		c.SetHeight(stack.height)
		c.SetX(stack.x)
		c.SetY(stack.y)
	}

	return stack.Children
}

type Center struct {
	SimpleBox
	Child Layoutable
}

func (c Center) Layout() []Layoutable {
	if c.Child == nil {
		return nil
	}

	c.Child.SetWidth(c.width)
	c.Child.SetHeight(c.height)
	c.Child.SetX(c.x - c.width/2)
	c.Child.SetY(c.y - c.height/2)

	return []Layoutable{c.Child}
}
