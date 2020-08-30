package font_test

import (
	"strings"
	"testing"

	"github.com/Qendolin/go-printpixel/core/font"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = font.BMF{
	Base:       29,
	Width:      512,
	Height:     512,
	LineHeight: 36,
	Face:       "Arial",
	Size:       36,
	Characters: map[rune]font.CharDef{
		'A': {
			Rune: 'A',
			X:    0, Y: 0,
			Width: 23, Height: 23,
			BearingX: -1, BearingY: 6,
			Advance: 21,
			Page:    0,
		},
		'L': {
			Rune: 'L',
			X:    47, Y: 0,
			Width: 16, Height: 23,
			BearingX: 2, BearingY: 6,
			Advance: 18,
			Page:    0,
		},
		'V': {
			Rune: 'V',
			X:    24, Y: 0,
			Width: 22, Height: 23,
			BearingX: 0, BearingY: 6,
			Advance: 21,
			Page:    0,
		},
	},
	Kernings: map[[2]rune]int{
		{'A', 'V'}: -2,
		{'V', 'A'}: -2,
		{'L', 'V'}: -2,
	},
}

var given = `info face="Arial" size=36 bold=0 italic=0 charset="" unicode=1 stretchH=100 smooth=1 aa=1 padding=0,0,0,0 spacing=1,1 outline=0
common lineHeight=36 base=29 scaleW=512 scaleH=512 pages=1 packed=0 alphaChnl=0 redChnl=4 greenChnl=4 blueChnl=4
page id=0 file="Arial2_0.png"
chars count=3
char id=65   x=0     y=0     width=23    height=23    xoffset=-1    yoffset=6     xadvance=21    page=0  chnl=15
char id=76   x=47    y=0     width=16    height=23    xoffset=2     yoffset=6     xadvance=18    page=0  chnl=15
char id=86   x=24    y=0     width=22    height=23    xoffset=0     yoffset=6     xadvance=21    page=0  chnl=15
kernings count=3
kerning first=86  second=65  amount=-2  
kerning first=76  second=86  amount=-2  
kerning first=65  second=86  amount=-2`

func TestParse(t *testing.T) {
	bmf, err := font.Parse(strings.NewReader(given))
	require.NoError(t, err)

	assert.Equal(t, expected.Base, bmf.Base)
	assert.Equal(t, expected.Width, bmf.Width)
	assert.Equal(t, expected.Height, bmf.Height)
	assert.Equal(t, expected.LineHeight, bmf.LineHeight)
	assert.Equal(t, expected.Size, bmf.Size)
	assert.Equal(t, expected.Face, bmf.Face)
	assert.EqualValues(t, []string{"Arial2_0.png"}, bmf.Pages)
	assert.EqualValues(t, expected.Characters, bmf.Characters)
	assert.Equal(t, expected.Kernings, bmf.Kernings)
}
