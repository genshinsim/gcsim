package parse

import "errors"

type Parser struct {
	lex *lexer
	res *ActionList

	//lookahead
	token []Token
	pos   int

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
		token:          make([]Token, 0, 20),
		pos:            -1,
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
	p.prefixParseFns[itemLeftParen] = p.parseParen
	p.prefixParseFns[keywordFn] = p.parseFn

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
	p.pos++
	if p.pos == len(p.token) {
		//grab more from the stream
		n := p.lex.nextItem()
		p.token = append(p.token, n)
	}
	return p.token[p.pos]
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.pos--
	//no op if at beginning
	if p.pos < -1 {
		p.pos = -1
	}
}

// peek returns but does not consume the next token.
func (p *Parser) peek() Token {
	n := p.next()
	p.backup()
	return n
}
