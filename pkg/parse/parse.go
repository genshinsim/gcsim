package parse

import "errors"

type Parser struct {
	lex  *lexer
	tree *Tree

	//lookahead
	token     [3]Token
	peekCount int
}

type Tree struct {
	FnMap map[string]Node
	Node  Node
}

type parseFn func(*Parser) (parseFn, error)

func New(name, input string) *Parser {
	p := &Parser{}
	p.lex = lex(name, input)
	p.tree = &Tree{
		FnMap: make(map[string]Node),
	}
	return p
}

func (p *Parser) Parse(text string) (*Tree, error) {
	var err error
	for state := parseProgram; state != nil; {
		state, err = state(p)
		if err != nil {
			return nil, err
		}
	}
	return p.tree, nil
}

func parseProgram(p *Parser) (parseFn, error) {
	//each line should start with one of the following:
	//	let
	//	a function call or action
	//	if
	//	while

	switch n := p.next(); {
	case n.typ == keywordLet:
		return parseLet, nil
	default:
		//unknown
		return nil, errors.New("unrecognized token")
	}
}

// next returns the next token.
func (p *Parser) next() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	return p.token[p.peekCount]
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.peekCount++
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (p *Parser) backup2(t1 Token) {
	p.token[1] = t1
	p.peekCount = 2
}

// backup3 backs the input stream up three tokens
// The zeroth token is already there.
func (p *Parser) backup3(t2, t1 Token) { // Reverse order: we're pushing back.
	p.token[1] = t1
	p.token[2] = t2
	p.peekCount = 3
}

// peek returns but does not consume the next token.
func (p *Parser) peek() Token {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.nextItem()
	return p.token[0]
}
