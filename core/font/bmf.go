package font

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type FormatError struct {
	LineNumber int
	Line       string
	Err        error
}

func (e FormatError) Error() string {
	return fmt.Sprintf("Error in line %v: '%v'", e.LineNumber, e.Line)
}

func (e FormatError) Unwrap() error {
	return e.Err
}

type BMF struct {
	LineHeight    int
	Base          int
	Width, Height int
	Pages         []string
	Characters    map[rune]CharDef
	Kernings      map[[2]rune]int
}

type CharDef struct {
	Rune               rune
	X, Y               int
	Width, Height      int
	BearingX, BearingY int
	Advance            int
	Page               int
}

func Parse(file io.Reader) (bmf *BMF, err error) {
	var lineNr int
	var line string
	defer func() {
		if err != nil {
			err = FormatError{
				Line:       line,
				LineNumber: lineNr,
				Err:        err,
			}
		}
	}()
	bmf = &BMF{
		Characters: map[rune]CharDef{},
		Kernings:   map[[2]rune]int{},
	}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lineNr++
		line = sc.Text()
		tag := strings.SplitN(line, " ", 2)
		switch tag[0] {
		case "char":
			var char CharDef
			for _, attr := range strings.Split(tag[1], " ") {
				kv := strings.Split(attr, "=")
				if len(kv) != 2 {
					continue
				}
				num, err := strconv.Atoi(kv[1])
				if err != nil {
					return nil, err
				}
				switch kv[0] {
				case "id":
					char.Rune = rune(num)
				case "x":
					char.X = num
				case "y":
					char.Y = num
				case "width":
					char.Width = num
				case "height":
					char.Height = num
				case "yoffset":
					char.BearingY = num
				case "xoffset":
					char.BearingX = num
				case "xadvance":
					char.Advance = num
				case "page":
					char.Page = num
				}
			}
			bmf.Characters[char.Rune] = char
		case "common":
			for _, attr := range strings.Split(tag[1], " ") {
				kv := strings.Split(attr, "=")
				if len(kv) != 2 {
					continue
				}
				num, err := strconv.Atoi(kv[1])
				if err != nil {
					return nil, err
				}
				switch kv[0] {
				case "lineHeight":
					bmf.LineHeight = num
				case "base":
					bmf.Base = num
				case "scaleW":
					bmf.Width = num
				case "scaleH":
					bmf.Height = num
				case "pages":
					bmf.Pages = make([]string, num)
				}
			}
		case "page":
			rd := bufio.NewReader(strings.NewReader(tag[1]))
			var id int
			var file string
			var eol error
			for {
				if b, err := rd.Peek(1); err == nil && b[0] == ' ' {
					rd.ReadByte()
					continue
				}

				var s string
				s, eol = rd.ReadString('=')
				if eol != nil {
					break
				}
				switch s {
				case "id=":
					s, eol = rd.ReadString(' ')
					s = strings.TrimRight(s, " ")
					id, err = strconv.Atoi(s)
					if err != nil {
						return nil, err
					}
					if eol != nil {
						break
					}
				case "file=":
					if b, err := rd.Peek(1); err == nil && b[0] == '"' {
						rd.ReadByte()
						file, eol = rd.ReadString('"')
						if eol != nil {
							return nil, eol
						}
						file = file[:len(file)-1]
					} else {
						file, eol = rd.ReadString(' ')
						if eol != nil {
							break
						}
						file = file[:len(file)-1]
					}
				}
			}
			bmf.Pages[id] = file
		case "kerning":
			var first, second, amount int
			for _, attr := range strings.Split(tag[1], " ") {
				kv := strings.Split(attr, "=")
				if len(kv) != 2 {
					continue
				}
				num, err := strconv.Atoi(kv[1])
				if err != nil {
					return nil, err
				}
				switch kv[0] {
				case "first":
					first = num
				case "second":
					second = num
				case "amount":
					amount = num
				}
			}
			bmf.Kernings[[2]rune{rune(first), rune(second)}] = amount
		}
	}
	return bmf, sc.Err()
}
