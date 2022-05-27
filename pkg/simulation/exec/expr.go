package exec

import "github.com/genshinsim/gcsim/pkg/parse"

func (e *Executor) evalExpr(ex parse.Expr) {
	switch v := ex.(type) {
	case *parse.BinaryExpr:
		e.evalBinaryExpr(v)
	}
}

func (e *Executor) evalBinaryExpr(b *parse.BinaryExpr) {

}
