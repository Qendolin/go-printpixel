package font

import (
	"io"
	"io/ioutil"
	"github.com/Qendolin/go-bmf"
)

type Font struct {
	LineHeight    int
	Base          int
	Width, Height int
	Pages         []string
	Characters    map[rune]CharDef
	Kernings      map[[2]rune]int
	Size          int
	Face          string
}

type CharDef struct {
	Rune               rune
	X, Y               int
	Width, Height      int
	BearingX, BearingY int
	Advance            int
	Page               int
}

func Parse(file io.Reader) (*Font, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	bmfont, err := bmf.Parse(data)
	if err != nil {
		return nil, err
	}

	pages := make([]string, len(bmfont.Pages))
	for i, p := range bmfont.Pages {
		pages[i] = p.File
	}

	chars := make(map[rune]CharDef, len(bmfont.Chars))
	for _, chr := range bmfont.Chars {
		chars[chr.Id] = CharDef{
			Rune: chr.Id,
			X:    chr.X, Y: chr.Y,
			Width: chr.Width, Height: chr.Height,
			BearingX: chr.XOffset, BearingY: chr.YOffset,
			Advance: chr.XAdvance,
			Page:    chr.Page,
		}
	}

	kerns := make(map[[2]rune]int, len(bmfont.Kernings))
	for _, k := range bmfont.Kernings {
		kerns[[2]rune{k.First, k.Second}] = k.Amount
	}

	return &Font{
		LineHeight: bmfont.Common.LineHeight,
		Base:       bmfont.Common.Base,
		Width:      bmfont.Common.ScaleW,
		Height:     bmfont.Common.ScaleH,
		Pages:      pages,
		Characters: chars,
		Kernings:   kerns,
		Size:       bmfont.Info.Size,
		Face:       bmfont.Info.Face,
	}, nil
}