package font

import (
	"sort"
	"unicode"
)

type Breaker interface {
	Break(curr rune, next rune) BreakConstraint
}

type Style struct {
	Kerning    bool
	Size       float32
	LineHeight float32
	TabSize    int
}

// TODO: Rename "LayoutSpecs"
type LayoutSpecs struct {
	Breaker  Breaker
	Width    float32
	Height   float32
	Ellipsis bool
}

type BreakConstraint int

const (
	PreventBreak BreakConstraint = -1
	CanBreak     BreakConstraint = 0
	ShouldBreak  BreakConstraint = 1
	ForceBreak   BreakConstraint = 2
)

func containsRune(slice []rune, r rune) bool {
	i := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= r
	})
	return i < len(slice) && slice[i] == r
}

type ListBreaker struct {
	ForceAfter    []rune
	ShouldBefore  []rune
	ShouldAfter   []rune
	PreventBefore []rune
	PreventAfter  []rune
}

func NewListBreaker(forceAfter, shouldBefore, shouldAfter, preventBefore, preventAfter []rune) ListBreaker {
	br := ListBreaker{
		ForceAfter:    forceAfter,
		ShouldBefore:  shouldBefore,
		ShouldAfter:   shouldAfter,
		PreventBefore: preventBefore,
		PreventAfter:  preventAfter,
	}
	// sort slices to make search faster
	sort.Slice(br.ForceAfter, func(i, j int) bool {
		return br.ForceAfter[i] < br.ForceAfter[j]
	})
	sort.Slice(br.ShouldAfter, func(i, j int) bool {
		return br.ShouldAfter[i] < br.ShouldAfter[j]
	})
	sort.Slice(br.PreventAfter, func(i, j int) bool {
		return br.PreventAfter[i] < br.PreventAfter[j]
	})
	sort.Slice(br.ShouldBefore, func(i, j int) bool {
		return br.ShouldBefore[i] < br.ShouldBefore[j]
	})
	sort.Slice(br.PreventBefore, func(i, j int) bool {
		return br.PreventBefore[i] < br.PreventBefore[j]
	})
	return br
}

func (br ListBreaker) Break(curr, next rune) BreakConstraint {
	if containsRune(br.ForceAfter, curr) {
		if curr == '\r' && next == '\n' {
			return PreventBreak
		}
		return ForceBreak
	}
	if containsRune(br.PreventBefore, next) {
		return PreventBreak
	}
	if containsRune(br.PreventAfter, curr) {
		return PreventBreak
	}
	if containsRune(br.ShouldAfter, curr) {
		if curr == '\u002D' && unicode.IsDigit(next) {
			return PreventBreak
		}
		return ShouldBreak
	}
	if containsRune(br.ShouldBefore, next) {
		return ShouldBreak
	}
	return CanBreak
}

var DefaultBreaker = NewListBreaker(
	[]rune{'\n', '\u000C', '\u000B', '\r', '\u2028', '\u2029'},
	[]rune{'\u00B4', '\u1FFD', '\u02DF', '\u02C8', '\u02CC', '\u002F'},
	[]rune{'\u2014', '\u1680', '\u2000', '\u2001', '\u2002', ' ', '\u2003', '\u2004', '\u2005', '\u2006', '\u2008', '\u2009', '\u200A', '\u205F', '\u3000', '\u0009', '\u00AD', '\u058A', '\u2010', '\u2012', '\u2013', '\u05BE', '\u0F0B', '\u1316', '\u17D8', '\u17DA', '\u2027', '\u007C', '\u002D'},
	[]rune{'\r', '\n', '\u00A0', '\u202F', '\u180E', '\u034F', '\u2007', '\u2011', '\u0F08', '\u0F0C', '\u0F12', '\u035C', '\u035D', '\u035E', '\u035F', '\u0360', '\u0361', '\u0362', '\u200D', '\u2060', '\uFEFF'},
	[]rune{'\u00A0', '\u202F', '\u180E', '\u034F', '\u2007', '\u2011', '\u0F08', '\u0F0C', '\u0F12', '\u035C', '\u035D', '\u035E', '\u035F', '\u0360', '\u0361', '\u0362', '\u200D', '\u2060', '\uFEFF'},
)

type preventBreaker struct{}

func (br preventBreaker) Break(a, b rune) BreakConstraint {
	return PreventBreak
}

var PreventBreaker = preventBreaker{}

func calcTabWidth(tabSize int, bmf *Font) int {
	if chr, ok := bmf.Characters[' ']; ok {
		return tabSize * chr.Advance
	}
	chr := bmf.Characters[0]
	return tabSize * chr.Advance
}

func Layout(text []rune, bmf *Font, idk LayoutSpecs, style Style) []rune {
	if idk.Breaker == nil {
		idk.Breaker = DefaultBreaker
	}
	if style.TabSize == 0 {
		style.TabSize = 1
	}

	scale := style.Size / float32(bmf.Size)
	lineWidth := int(idk.Width / scale)
	lines := int(idk.Height / float32(bmf.LineHeight) / scale / style.LineHeight)
	var ellipsisWidth int
	var ellipsis []rune
	if idk.Ellipsis {
		// Pick ellipsis character if it exists
		if chr, ok := bmf.Characters['\u2026']; ok {
			ellipsisWidth = chr.Advance
			ellipsis = []rune{'\u2026'}
		} else if chr, ok = bmf.Characters['.']; ok {
			w := chr.Advance
			if k, ok := bmf.Kernings[[2]rune{'.', '.'}]; ok && style.Kerning {
				w += k
			}
			ellipsisWidth = 3 * w
			ellipsis = []rune{'.', '.', '.'}
		} else {
			idk.Ellipsis = false
		}
	}
	tabWidth := calcTabWidth(style.TabSize, bmf)

	var res []rune
	var word []rune
	var wordWidth int
	var px, line int
	for i, chr := range text {
		var next rune
		if i+1 < len(text) {
			next = text[i+1]
		}
		def := bmf.Characters[chr]
		if chr == '\t' {
			def.Advance = tabWidth - (px+wordWidth)%tabWidth
		}
		br := idk.Breaker.Break(chr, next)
		width := def.Advance
		if k, ok := bmf.Kernings[[2]rune{chr, next}]; ok && style.Kerning {
			width += k
		}
		futureLineWidth := px + wordWidth + width

		if line+1 == lines {
			if br == ForceBreak || futureLineWidth >= lineWidth-ellipsisWidth {
				line++
				break
			}
			word = append(word, chr)
			wordWidth += width
			continue
		}

		switch br {
		case ForceBreak:
			word = append(word, '\n')
			res = append(res, word...)
			line++
			px = 0
			word = []rune{}
			wordWidth = 0
		case PreventBreak:
			word = append(word, chr)
			wordWidth += width
		case ShouldBreak:
			if futureLineWidth > lineWidth {
				word = append(word, '\n', chr)
				res = append(res, word...)
				line++
				px = width
			} else {
				word = append(word, chr)
				res = append(res, word...)
				px += wordWidth + width
			}
			word = []rune{}
			wordWidth = 0
		case CanBreak:
			if wordWidth+width > lineWidth {
				// word is too long
				word = append(word, '\n')
				res = append(res, word...)
				line++
				px = 0
				word = []rune{chr}
				wordWidth = width
			} else if futureLineWidth > lineWidth {
				// line is too long
				word = append(word, chr)
				res = append(res, '\n')
				line++
				px = 0
				wordWidth += width
			} else {
				word = append(word, chr)
				wordWidth += width
			}
		}
	}
	if line == lines {
		word = append(word, ellipsis...)
	}
	res = append(res, word...)
	return res
}

type BrokenText struct {
	Text       string
	Kerning    bool
	Size       float32
	LineHeight float32
}

func Mesh(text []rune, bmf *Font, style Style) (v []float32, t []float32) {
	if style.TabSize == 0 {
		style.TabSize = 1
	}

	tabWidth := calcTabWidth(style.TabSize, bmf)
	sx, sy := float32(bmf.Width), float32(bmf.Height)
	px, py := 0, 0
	verts := make([]float32, 0, len(text)*12)
	texs := make([]float32, 0, len(text)*12)
	for i, chr := range text {
		if chr == '\n' {
			py -= int(float32(bmf.LineHeight) * style.LineHeight)
			px = 0
			continue
		}
		def, ok := bmf.Characters[chr]
		if !ok {
			if chr == '\t' {
				px += tabWidth - px%tabWidth
			}
			if unicode.IsControl(chr) {
				continue
			}
			if def, ok = bmf.Characters[0]; !ok {
				continue
			}
		}

		if def.Width != 0 && def.Height != 0 {
			rx := px + def.BearingX
			ry := py - def.BearingY
			verts = append(verts,
				float32(rx), float32(ry),
				float32(rx+def.Width), float32(ry),
				float32(rx+def.Width), float32(ry-def.Height),
				float32(rx+def.Width), float32(ry-def.Height),
				float32(rx), float32(ry-def.Height),
				float32(rx), float32(ry))
			pg := float32(def.Page)
			texs = append(texs,
				float32(def.X)/sx+pg, float32(def.Y)/sy,
				float32(def.X+def.Width)/sx+pg, float32(def.Y)/sy,
				float32(def.X+def.Width)/sx+pg, float32(def.Y+def.Height)/sy,
				float32(def.X+def.Width)/sx+pg, float32(def.Y+def.Height)/sy,
				float32(def.X)/sx+pg, float32(def.Y+def.Height)/sy,
				float32(def.X)/sx+pg, float32(def.Y)/sy)
		}

		px += def.Advance
		if i+1 < len(text) && style.Kerning {
			if amount, ok := bmf.Kernings[[2]rune{chr, text[i+1]}]; ok {
				px += amount
			}
		}
	}
	return verts, texs
}
