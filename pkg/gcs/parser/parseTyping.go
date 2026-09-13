package parser

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (p *Parser) parseOptionalType() (ast.ExprType, error) {
	// the only time typing information should be present is either in a fn signature or after a let stmt
	// should be safe to assume that if the next token is either an identifier or fn, it should be typing info
	// if it's not either then we return nil
	n := p.peek()
	switch n.Typ {
	case ast.ItemIdentifier:
	case ast.KeywordFn:
	default:
		return nil, nil
	}
	return p.parseTyping()
}

func (p *Parser) parseTyping() (ast.ExprType, error) {
	// next should be an ident with one of the following value, otherwise error
	// - number
	// - string
	// - fn(...) : ...
	n := p.peek()
	switch n.Typ {
	case ast.ItemIdentifier:
		return p.parseBasicType()
	case ast.KeywordFn:
		return p.parseFnType()
	default:
		return nil, fmt.Errorf("ln%v: error parsing type info, unexpected value after :, got %v", n.Line, n.Val)
	}
}

func (p *Parser) parseBasicType() (ast.ExprType, error) {
	n := p.next()
	if n.Typ != ast.ItemIdentifier {
		return nil, fmt.Errorf("ln%v: error parsing basic type, expecting identifier, got %v", n.Line, n.Val)
	}
	switch n.Val {
	case "string":
		return &ast.StringType{Pos: n.Pos}, nil
	case "number":
		return &ast.NumberType{Pos: n.Pos}, nil
	case "map":
		return &ast.MapType{Pos: n.Pos}, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected basic type parsing type info; got %v", n.Line, n.Val)
	}
}

func (p *Parser) parseFnType() (ast.ExprType, error) {
	// expecting something like: fn(number) : number
	// this is also valid:
	//   fn(fn(number):number, number) : fn(number)
	// going to end up with lots of recursive calls...
	var err error
	n := p.next()
	if n.Typ != ast.KeywordFn {
		return nil, fmt.Errorf("ln%v: error parsing fn type, expecting fn, got %v", n.Line, n.Val)
	}
	res := &ast.FuncType{
		Pos: n.Pos,
	}
	// we're expecting, in order:
	// - (
	// - comma separate types
	// - )
	// - optional return value
	n = p.next()
	if n.Typ != ast.ItemLeftParen {
		return nil, fmt.Errorf("ln%v: expecting ( after fn parsing typing, got %v", n.Line, n.Val)
	}
	done := false
	// if next is rightparen then this fn has no arguments
	if l := p.peek(); l.Typ == ast.ItemRightParen {
		// consume the token
		p.next()
		done = true
	}
	for !done {
		// we expect the first token or group of token to be typing info
		typ, err := p.parseTyping()
		if err != nil {
			return nil, err
		}
		res.ArgsType = append(res.ArgsType, typ)

		// the next token is either a ) signifiying we're done, or a comma meaning we should
		// continue parsing
		n = p.next()
		switch n.Typ {
		case ast.ItemRightParen:
			// if next is ) then we're done
			done = true
		case ast.ItemComma:
			// comma means we keep going
		default:
			// unexpected token
			return nil, fmt.Errorf("ln%v: unexpected token parsing fn type: %v", n.Line, n.Val)
		}
	}
	// check for optional return type
	res.ResultType, err = p.parseOptionalType()
	if err != nil {
		return nil, err
	}
	// TODO: this is only here for compatability reasons; to be removed??
	if res.ResultType == nil {
		res.ResultType = &ast.NumberType{Pos: n.Pos}
	}
	return res, nil
}
