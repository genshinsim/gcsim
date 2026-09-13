package ast

import (
	"fmt"
	"sort"
)

type Position struct {
	Line   int
	Column int
}

func (p Position) IsValid() bool { return p.Line > 0 }

func (p Position) String() string {
	if !p.IsValid() {
		return "-"
	}
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}

type File struct {
	size  int
	lines []int
}

func (f *File) AddLine(offset int) {
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

func (f *File) Position(offset Pos) Position {
	line := sort.Search(len(f.lines), func(i int) bool { return f.lines[i] > int(offset) }) - 1
	if line < 0 {
		return Position{}
	}

	return Position{
		Line:   line + 1,
		Column: int(offset) - f.lines[line] + 1,
	}
}

func NewFile() *File {
	return &File{lines: []int{0}}
}
