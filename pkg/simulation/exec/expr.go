package exec

import "github.com/genshinsim/gcsim/pkg/parse"

func (e *Executor) evalExpr(ex parse.Expr) Obj {
	switch v := ex.(type) {
	case *parse.NumberLit:
		return e.evalNumberLit(v)
	case *parse.StringLit:
		return e.evalStringLit(v)
	case *parse.Ident:
		return e.evalIdent(v)
	case *parse.BinaryExpr:
		return e.evalBinaryExpr(v)
	default:
		return &null{}
	}
}

func (e *Executor) evalNumberLit(n *parse.NumberLit) Obj {
	return &number{
		isInt: n.IsInt,
		ival:  n.IntVal,
		fval:  n.FloatVal,
	}
}

func (e *Executor) evalStringLit(n *parse.StringLit) Obj {
	return &strval{
		str: n.Value,
	}
}

func (e *Executor) evalIdent(n *parse.Ident) Obj {
	//TODO: this should be a variable
	return &null{}
}

func (e *Executor) evalCallExpr(c *parse.CallExpr) Obj {
	//check if it's a system function
	//otherwise check the function map

	return &null{}
}

func (e *Executor) evalBinaryExpr(b *parse.BinaryExpr) Obj {
	return &null{}
}
