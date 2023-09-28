package ast

import "fmt"

// expecting let ident = expr;
func (p *Parser) parseLet() (Stmt, error) {
	n := p.next()

	ident, err := p.consume(itemIdentifier)
	if err != nil {
		// next token not an identifier
		return nil, fmt.Errorf("ln%v: expecting identifier after let, got %v", ident.line, ident.Val)
	}

	stmt := &LetStmt{
		Pos:   n.pos,
		Ident: ident,
	}

	// optional typing info; if not present assume number
	if l := p.peek(); l.Typ != itemAssign {
		stmt.Type, err = p.parseTyping()
		if err != nil {
			return nil, err
		}
	}
	if stmt.Type == nil {
		stmt.Type = &NumberType{Pos: n.pos}
	}

	a, err := p.consume(itemAssign)
	if err != nil {
		// next token not and identifier
		return nil, fmt.Errorf("ln%v: expecting = after identifier in let statement, got %v", a.line, a.Val)
	}

	stmt.Val, err = p.parseExpr(Lowest)

	return stmt, err
}
