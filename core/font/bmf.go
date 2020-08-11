package font

import (
	"bufio"
	"errors"
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
		tag, attribs, strs, err := parseTag(line)
		if err != nil {
			return nil, err
		}
		switch tag {
		case "info":
			for k, v := range attribs {
				switch k {
				case "size":
					bmf.Size = v
				case "face":
					bmf.Face = strs[v]
				}
			}
		case "char":
			var char CharDef
			for k, v := range attribs {
				switch k {
				case "id":
					char.Rune = rune(v)
				case "x":
					char.X = v
				case "y":
					char.Y = v
				case "width":
					char.Width = v
				case "height":
					char.Height = v
				case "yoffset":
					char.BearingY = v
				case "xoffset":
					char.BearingX = v
				case "xadvance":
					char.Advance = v
				case "page":
					char.Page = v
				}
			}
			bmf.Characters[char.Rune] = char
		case "common":
			for k, v := range attribs {
				switch k {
				case "lineHeight":
					bmf.LineHeight = v
				case "base":
					bmf.Base = v
				case "scaleW":
					bmf.Width = v
				case "scaleH":
					bmf.Height = v
				case "pages":
					bmf.Pages = make([]string, v)
				}
			}
		case "page":
			var id int
			var file string
			for k, v := range attribs {
				switch k {
				case "id":
					id = v
				case "file":
					file = strs[v]
				}
			}
			bmf.Pages[id] = file
		case "kerning":
			var first, second, amount int
			for k, v := range attribs {
				switch k {
				case "first":
					first = v
				case "second":
					second = v
				case "amount":
					amount = v
				}
			}
			bmf.Kernings[[2]rune{rune(first), rune(second)}] = amount
		}
	}
	return bmf, sc.Err()
}

func parseTag(line string) (name string, values map[string]int, strs []string, err error) {
	values = map[string]int{}
	strs = []string{}

	var stripped string
	rd := bufio.NewReader(strings.NewReader(line))
	for {
		start, err := rd.ReadString('"')
		stripped += start
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return "", nil, nil, err
		}
		str, err := rd.ReadString('"')
		if errors.Is(err, io.EOF) {
			return "", nil, nil, fmt.Errorf("expected \"")
		}
		strs = append(strs, str[:len(str)-1])
	}

	fields := strings.Fields(stripped)
	if len(fields) == 0 {
		return "", nil, nil, fmt.Errorf("empty tag")
	}

	strIdx := 0
	for i, f := range fields {
		if i == 0 {
			name = f
			continue
		}

		kv := strings.Split(f, "=")
		if len(kv) != 2 {
			return "", nil, nil, fmt.Errorf("expected key-value pair")
		}
		key, value := kv[0], kv[1]
		if value == "\"" {
			values[key] = strIdx
			strIdx++
		} else if num, err := strconv.Atoi(value); err == nil {
			values[key] = num
		}
	}

	return
}
