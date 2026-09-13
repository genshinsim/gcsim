package parser

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (p *Parser) parseFnStmt() (ast.Stmt, error) {
	// fn ident(...deint){ block }
	n := p.next()
	if n.Typ != ast.KeywordFn {
		return nil, fmt.Errorf("ln %v: expecting fn, got %v", n.Line, n.Val)
	}
	n = p.next()
	if n.Typ != ast.ItemIdentifier {
		return nil, fmt.Errorf("ln %v: expecting identifier after fn, got %v", n.Line, n.Val)
	}
	// expecting function body
	lit, err := p.parseFn()
	if err != nil {
		return nil, err
	}
	return &ast.FnStmt{
		Pos:   n.Pos,
		Ident: n,
		Func:  lit,
	}, nil
}

func (p *Parser) parseFnExpr() (ast.Expr, error) {
	// fn (...ident) { block }
	// consume the fn
	n := p.next()
	if n.Typ != ast.KeywordFn {
		return nil, fmt.Errorf("ln %v: expecting fn, got %v", n.Line, n.Val)
	}
	// expecting function body
	lit, err := p.parseFn()
	if err != nil {
		return nil, err
	}
	return &ast.FuncExpr{
		Pos:  n.Pos,
		Func: lit,
	}, nil
}

func (p *Parser) parseFn() (*ast.FuncLit, error) {
	// (...ident){ block }
	var err error

	// expect n to be left parent
	n := p.peek()
	if n.Typ != ast.ItemLeftParen {
		return nil, fmt.Errorf("ln%v: expecting ( after identifier, got %v", n.Line, n.Val)
	}

	lit := &ast.FuncLit{
		Pos: n.Pos,
		Signature: &ast.FuncType{
			Pos: n.Pos,
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
			return nil, fmt.Errorf("ln%v: fn contains duplicated param name %v", n.Line, v.Value)
		}
		chk[v.Value] = true
	}

	// if next is not left brace then we're expecting typing info
	if l := p.peek(); l.Typ != ast.ItemLeftBrace {
		lit.Signature.ResultType, err = p.parseTyping()
		if err != nil {
			return nil, err
		}
	}
	// TODO: if nil we are assuming number for compatbility reasons
	if lit.Signature.ResultType == nil {
		// TODO: the position here is wrong... really shouldn't be the position of the open bracket
		// TODO: should fix this by adding a current position the parser
		lit.Signature.ResultType = &ast.NumberType{Pos: n.Pos}
	}

	lit.Body, err = p.parseBlock()
	if err != nil {
		return nil, err
	}

	return lit, nil
}

func (p *Parser) parseFnArgs() ([]*ast.Ident, []ast.ExprType, error) {
	// consume (
	var args []*ast.Ident
	var argsType []ast.ExprType
	p.next()
	for n := p.next(); n.Typ != ast.ItemRightParen; n = p.next() {
		a := &ast.Ident{}
		// expecting ident, comma
		if n.Typ != ast.ItemIdentifier {
			return nil, nil, fmt.Errorf("ln%v: expecting identifier in param list, got %v", n.Line, n.Val)
		}
		a.Pos = n.Pos
		a.Value = n.Val

		args = append(args, a)

		// check for optional typing
		// if not present assume number
		typ, err := p.parseOptionalType()
		if err != nil {
			return nil, nil, err
		}
		// TODO: if nil we are assuming number for compatbility reasons
		if typ == nil {
			typ = &ast.NumberType{Pos: n.Pos}
		}

		argsType = append(argsType, typ)

		// if next token is a comma, then there should be another ident after that
		// otherwise we have a problem
		if l := p.peek(); l.Typ == ast.ItemComma {
			p.next() // consume the comma
			if l = p.peek(); l.Typ != ast.ItemIdentifier {
				return nil, nil, fmt.Errorf("ln%v: expecting another identifier after comma in param list, got %v", n.Line, n.Val)
			}
		}
	}
	return args, argsType, nil
}
