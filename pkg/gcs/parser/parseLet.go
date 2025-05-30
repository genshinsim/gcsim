package parser

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

// expecting let ident = expr;
func (p *Parser) parseLet() (ast.Stmt, error) {
	n := p.next()

	ident, err := p.consume(ast.ItemIdentifier)
	if err != nil {
		// next token not an identifier
		return nil, fmt.Errorf("ln%v: expecting identifier after let, got %v", ident.Line, ident.Val)
	}

	stmt := &ast.LetStmt{
		Pos:   n.Pos,
		Ident: ident,
	}

	// optional typing info; if not present assume number
	if l := p.peek(); l.Typ != ast.ItemAssign {
		stmt.Type, err = p.parseTyping()
		if err != nil {
			return nil, err
		}
	}
	if stmt.Type == nil {
		stmt.Type = &ast.NumberType{Pos: n.Pos}
	}

	a, err := p.consume(ast.ItemAssign)
	if err != nil {
		// next token not and identifier
		return nil, fmt.Errorf("ln%v: expecting = after identifier in let statement, got %v", a.Line, a.Val)
	}

	stmt.Val, err = p.parseExpr(ast.Lowest)

	return stmt, err
}
