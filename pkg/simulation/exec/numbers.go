package exec

import "go/token"

type Number interface {
	~int | float64
}

func processNumber[V Number](x, y V, op token.Token) V {
	switch op {
	case token.ADD:
		return x + y
	case token.SUB:
		return x - y
	case token.MUL:
		return x * y
	case token.QUO:
		return x / y
	default:
		return 0
	}
}

func compareNumber[V Number](x, y V, op token.Token) bool {
	switch op {
	case token.EQL:
		return x == y
	case token.LSS:
		return x < y
	case token.GTR:
		return x > y
	case token.LEQ:
		return x <= y
	case token.GEQ:
		return x >= y
	default:
		return false
	}
}
