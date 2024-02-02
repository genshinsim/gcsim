package ast

import "fmt"

func (p *Parser) parseFnStmt() (Stmt, error) {
	// fn ident(...deint){ block }
	n := p.next()
	if n.Typ != keywordFn {
		return nil, fmt.Errorf("ln %v: expecting fn, got %v", n.line, n.Val)
	}
	n = p.next()
	if n.Typ != itemIdentifier {
		return nil, fmt.Errorf("ln %v: expecting identifier after fn, got %v", n.line, n.Val)
	}
	// expecting function body
	lit, err := p.parseFn()
	if err != nil {
		return nil, err
	}
	return &FnStmt{
		Pos:   n.pos,
		Ident: n,
		Func:  lit,
	}, nil
}

func (p *Parser) parseFnExpr() (Expr, error) {
	// fn (...ident) { block }
	// consume the fn
	n := p.next()
	if n.Typ != keywordFn {
		return nil, fmt.Errorf("ln %v: expecting fn, got %v", n.line, n.Val)
	}
	// expecting function body
	lit, err := p.parseFn()
	if err != nil {
		return nil, err
	}
	return &FuncExpr{
		Pos:  n.pos,
		Func: lit,
	}, nil
}

func (p *Parser) parseFn() (*FuncLit, error) {
	// (...ident){ block }
	var err error

	// expect n to be left parent
	n := p.peek()
	if n.Typ != itemLeftParen {
		return nil, fmt.Errorf("ln%v: expecting ( after identifier, got %v", n.line, n.Val)
	}

	lit := &FuncLit{
		Pos: n.pos,
		Signature: &FuncType{
			Pos: n.pos,
		},
	}

	// parse the arguments
	lit.Args, lit.Signature.ArgsType, err = p.parseFnArgs()
	if err != nil {
		return nil, err
	}

	// check that args are not duplicates
	chk := make(map[string]bool)
	for _, v := range lit.Args {
		if _, ok := chk[v.Value]; ok {
			return nil, fmt.Errorf("ln%v: fn contains duplicated param name %v", n.line, v.Value)
		}
		chk[v.Value] = true
	}

	// if next is not left brace then we're expecting typing info
	if l := p.peek(); l.Typ != itemLeftBrace {
		lit.Signature.ResultType, err = p.parseTyping()
		if err != nil {
			return nil, err
		}
	}
	//TODO: if nil we are assuming number for compatbility reasons
	if lit.Signature.ResultType == nil {
		//TODO: the position here is wrong... really shouldn't be the position of the open bracket
		//TODO: should fix this by adding a current position the parser
		lit.Signature.ResultType = &NumberType{Pos: n.pos}
	}

	lit.Body, err = p.parseBlock()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseFnArgs() ([]*Ident, []ExprType, error) {
	// consume (
	var args []*Ident
	var argsType []ExprType
	p.next()
	for n := p.next(); n.Typ != itemRightParen; n = p.next() {
		a := &Ident{}
		// expecting ident, comma
		if n.Typ != itemIdentifier {
			return nil, nil, fmt.Errorf("ln%v: expecting identifier in param list, got %v", n.line, n.Val)
		}
		a.Pos = n.pos
		a.Value = n.Val

		args = append(args, a)

		// check for optional typing
		// if not present assume number
		typ, err := p.parseOptionalType()
		if err != nil {
			return nil, nil, err
		}
		//TODO: if nil we are assuming number for compatbility reasons
		if typ == nil {
			typ = &NumberType{Pos: n.pos}
		}

		argsType = append(argsType, typ)

		// if next token is a comma, then there should be another ident after that
		// otherwise we have a problem
		if l := p.peek(); l.Typ == itemComma {
			p.next() // consume the comma
			if l = p.peek(); l.Typ != itemIdentifier {
				return nil, nil, fmt.Errorf("ln%v: expecting another identifier after comma in param list, got %v", n.line, n.Val)
			}
		}
	}
	return args, argsType, nil
}
