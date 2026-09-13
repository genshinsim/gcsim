package ast

import "fmt"

type Error struct {
	Pos Position
	Msg string
}

func (e Error) Error() string {
	if !e.Pos.IsValid() {
		return e.Msg
	}
	return fmt.Sprintf("ln%v: %v", e.Pos, e.Msg)
}

func NewError(pos Position, msg string) Error {
	return Error{
		Pos: pos,
		Msg: msg,
	}
}

func NewErrorf(pos Position, format string, a ...any) Error {
	return NewError(pos, fmt.Sprintf(format, a...))
}
