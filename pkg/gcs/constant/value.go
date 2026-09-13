package constant

import (
	"fmt"
	"strconv"
)

type ValueTyp int

type Value interface {
	Inspect() string
	Typ() ValueTyp
}

const (
	Number ValueTyp = iota
	String
)

var typStrings = map[ValueTyp]string{
	Number: "number",
	String: "string",
}

func (v ValueTyp) String() string {
	if name, ok := typStrings[v]; ok {
		return name
	}
	return "unknown"
}

// various Obj types
type (
	number struct {
		ival    int64
		fval    float64
		isFloat bool
	}

	strval struct {
		str string
	}
)

// number.
func (n *number) Inspect() string {
	if n.isFloat {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	}
	return strconv.FormatInt(n.ival, 10)
}
func (n *number) Typ() ValueTyp { return Number }

// strval.
func (s *strval) Inspect() string { return fmt.Sprintf("\"%v\"", s.str) }
func (s *strval) Typ() ValueTyp   { return String }

func Val(x Value) any {
	switch x := x.(type) {
	case *number:
		if x.isFloat {
			return ntof(x)
		}
		return ntoi(x)
	case *strval:
		return x.str
	default:
		return nil
	}
}

func Make(x any) Value {
	switch x := x.(type) {
	case bool:
		return bton(x)
	case int:
		return &number{
			ival:    int64(x),
			isFloat: false,
		}
	case int64:
		return &number{
			ival:    x,
			isFloat: false,
		}
	case float64:
		return &number{
			fval:    x,
			isFloat: true,
		}
	case string:
		return &strval{str: x}
	case Value:
		return x
	default:
		return nil
	}
}
