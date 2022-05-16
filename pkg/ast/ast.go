package ast

type AST struct {
	Variables map[string]interface{} //map of variables
	Tree      Node
}

type Node interface {
	Copy() Node
}

type IdentKind int

const (
	Bad IdentKind = iota
	Var
	Fun
	Lbl
)

type Program struct {
	Stmts []Node
}

func (p *Program) Copy() *Program {
	next := &Program{}
	next.Stmts = make([]Node, len(p.Stmts))
	for i := range p.Stmts {
		next.Stmts[i] = p.Stmts[i].Copy()
	}
	return next
}
