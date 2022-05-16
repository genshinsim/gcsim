package parse

type Parser struct {
	l      *lexer
	tokens []item
	pos    int //current position
}

type parseFn func(*Parser) (parseFn, error)

func New(name, input string) *Parser {
	p := &Parser{}
	p.l = lex(name, input)
	p.pos = -1
	return p
}
