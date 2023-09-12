package ast

import "fmt"

func (p *Parser) parseOptionalType(identPos Pos) (ExprType, error) {
	// the only time typing information should be present is either in a fn signature or after a let stmt
	// should be safe to assume that if the next token is either an identifier or fn, it should be typing info
	// if it's not either then we return nil
	n := p.peek()
	switch n.Typ {
	case itemIdentifier:
	case keywordFn:
	default:
		return nil, nil
	}
	return p.parseTyping()
}

func (p *Parser) parseTyping() (ExprType, error) {
	// next should be an ident with one of the following value, otherwise error
	// - number
	// - string
	// - fn(...) : ...
	n := p.peek()
	switch n.Typ {
	case itemIdentifier:
		return p.parseBasicType()
	case keywordFn:
		return p.parseFnType()
	default:
		return nil, fmt.Errorf("ln%v: error parsing type info, unexpected value after :, got %v", n.line, n.Val)
	}
}

func (p *Parser) parseBasicType() (ExprType, error) {
	n := p.next()
	if n.Typ != itemIdentifier {
		return nil, fmt.Errorf("ln%v: error parsing basic type, expecting identifier, got %v", n.line, n.Val)
	}
	switch n.Val {
	case "string":
		return &StringType{Pos: n.pos}, nil
	case "number":
		return &NumberType{Pos: n.pos}, nil
	case "map":
		return &MapType{Pos: n.pos}, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected basic type parsing type info; got %v", n.line, n.Val)
	}
}

func (p *Parser) parseFnType() (ExprType, error) {
	// expecting something like: fn(number) : number
	// this is also valid:
	//   fn(fn(number):number, number) : fn(number)
	// going to end up with lots of recursive calls...
	var err error
	n := p.next()
	if n.Typ != keywordFn {
		return nil, fmt.Errorf("ln%v: error parsing fn type, expecting fn, got %v", n.line, n.Val)
	}
	res := &FuncType{
		Pos: n.pos,
	}
	// we're expecting, in order:
	// - (
	// - comma separate types
	// - )
	// - optional return value
	n = p.next()
	if n.Typ != itemLeftParen {
		return nil, fmt.Errorf("ln%v: expecting ( after fn parsing typing, got %v", n.line, n.Val)
	}
	done := false
	// if next is rightparen then this fn has no arguments
	if l := p.peek(); l.Typ == itemRightParen {
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
		case itemRightParen:
			// if next is ) then we're done
			done = true
		case itemComma:
			// comma means we keep going
		default:
			// unexpected token
			return nil, fmt.Errorf("ln%v: unexpected token parsing fn type: %v", n.line, n.Val)
		}
	}
	// check for optional return type
	res.ResultType, err = p.parseOptionalType(n.pos)
	if err != nil {
		return nil, err
	}
	//TODO: this is only here for compatability reasons; to be removed??
	if res.ResultType == nil {
		res.ResultType = &NumberType{Pos: n.pos}
	}
	return res, nil
}
