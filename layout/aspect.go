package layout

import "math"

type AspectMode int

const (
	Contain = AspectMode(iota)
	Cover
	FitWidth
	FitHieght
)

type Aspect struct {
	SimpleBox
	Child Layoutable
	Ratio float64
	Mode  AspectMode
}

func (a *Aspect) Layout() {
	if a.Ratio == 0 {
		a.Ratio = 1
	}

	a.Child.SetX(a.x)
	a.Child.SetY(a.y)
	switch a.Mode {
	case Contain:
		r := float64(a.width) / float64(a.height)
		if math.Abs(r-a.Ratio) < 1e-6 {
			a.Child.SetWidth(a.width)
			a.Child.SetHeight(a.height)
		} else if r > a.Ratio {
			a.Child.SetHeight(a.height)
			a.Child.SetWidth(int(float64(a.height) * a.Ratio))
		} else {
			a.Child.SetWidth(a.width)
			a.Child.SetHeight(int(float64(a.width) / a.Ratio))
		}
	case Cover:
		r := float64(a.width) / float64(a.height)
		if math.Abs(r-a.Ratio) < 1e-6 {
			a.Child.SetWidth(a.width)
			a.Child.SetHeight(a.height)
		} else if r > a.Ratio {
			a.Child.SetWidth(a.width)
			a.Child.SetHeight(int(float64(a.width) / a.Ratio))
		} else {
			a.Child.SetHeight(a.height)
			a.Child.SetWidth(int(float64(a.height) * a.Ratio))
		}
	case FitWidth:
		a.Child.SetWidth(a.width)
		a.Child.SetHeight(int(float64(a.width) / a.Ratio))
	case FitHieght:
		a.Child.SetHeight(a.height)
		a.Child.SetWidth(int(float64(a.height) * a.Ratio))
	}

	if l, ok := a.Child.(Layouter); ok {
		l.Layout()
	}
}
