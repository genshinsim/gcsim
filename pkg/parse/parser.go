package parse

import "errors"

type Parser struct {
	lex *lexer
	res *ActionList

	//lookahead
	token     [3]Token
	peekCount int

	//parseFn
	prefixParseFns map[TokenType]func() Expr
	infixParseFns  map[TokenType]func(Expr) Expr
}

type ActionList struct {
	FnMap   map[string]Node
	Program *BlockStmt
}

type parseFn func(*Parser) (parseFn, error)

func New(input string) *Parser {
	p := &Parser{
		prefixParseFns: make(map[TokenType]func() Expr),
		infixParseFns:  make(map[TokenType]func(Expr) Expr),
	}
	p.lex = lex(input)
	p.res = &ActionList{
		FnMap:   make(map[string]Node),
		Program: newBlockStmt(0),
	}
	//expr functions
	p.prefixParseFns[itemIdentifier] = p.parseIdent
	p.prefixParseFns[itemNumber] = p.parseNumber
	p.prefixParseFns[LogicNot] = p.parseUnaryExpr
	p.prefixParseFns[itemMinus] = p.parseUnaryExpr

	p.infixParseFns[itemPlus] = p.parseBinaryExpr
	p.infixParseFns[itemMinus] = p.parseBinaryExpr
	p.infixParseFns[itemSlash] = p.parseBinaryExpr
	p.infixParseFns[itemAsterisk] = p.parseBinaryExpr
	p.infixParseFns[OpEqual] = p.parseBinaryExpr
	p.infixParseFns[OpNotEqual] = p.parseBinaryExpr
	p.infixParseFns[OpLessThan] = p.parseBinaryExpr
	p.infixParseFns[OpLessThanOrEqual] = p.parseBinaryExpr
	p.infixParseFns[OpGreaterThan] = p.parseBinaryExpr
	p.infixParseFns[OpGreaterThanOrEqual] = p.parseBinaryExpr

	return p
}

// consume returns err if next token does not match expected
// otherwise return next token and nil error
func (p *Parser) consume(i TokenType) (Token, error) {
	n := p.next()
	if n.typ != i {
		return n, errors.New("unexpected token")
	}
	return n, nil
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
