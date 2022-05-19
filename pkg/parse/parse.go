package parse

import (
	"fmt"
	"strconv"
)

type precedence int

const (
	_ precedence = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[TokenType]precedence{
	OpEqual:              EQUALS,
	OpNotEqual:           EQUALS,
	OpLessThan:           LESSGREATER,
	OpGreaterThan:        LESSGREATER,
	OpLessThanOrEqual:    LESSGREATER,
	OpGreaterThanOrEqual: LESSGREATER,
	itemPlus:             SUM,
	itemMinus:            SUM,
	itemSlash:            PRODUCT,
	itemAsterisk:         PRODUCT,
}

func (t Token) precedence() precedence {
	if p, ok := precedences[t.typ]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Parse() (*ActionList, error) {
	var err error
	for state := parseText; state != nil; {
		state, err = state(p)
		if err != nil {
			return nil, err
		}
	}
	return p.res, nil
}

func parseText(p *Parser) (parseFn, error) {
	switch n := p.peek(); n.typ {
	case itemCharacterKey:
		p.next()
		//check if this is character stats etc or an action
		if p.peek().typ <= itemActionKey {
			//not an ActionStmt
			p.backup()
			return parseCharacter, nil
		}
		p.backup()
		//parse action item
		fallthrough
	case keywordFunc:
		fallthrough
	case keywordLet:
		return parseProgram, nil
	case itemEOF:
		return nil, nil
	default: //default should be look for gcsl
		return parseProgram, nil
	}
}

func parseCharacter(p *Parser) (parseFn, error) {
	return nil, nil
}

func parseProgram(p *Parser) (parseFn, error) {
	node := p.parseStatement()
	n, err := p.consume(itemTerminateLine)
	if err != nil {
		return nil, fmt.Errorf("expecting end of line token parsing stmt, got %v", n)
	}
	p.res.Program.append(node)
	return parseText, nil
}

func (p *Parser) parseStatement() Node {
	n := p.peek()
	switch n.typ {
	case keywordLet:
		return p.parseLet()
	case itemCharacterKey:
		return p.parseAction()
	case keywordFunc:
		return p.parseFn()
	default:
		return p.parseExpr(LOWEST)
	}
}

//parseAction returns a node contain a character action, or a block of node containing
//a list of character actions
func (p *Parser) parseAction() Stmt {

	return nil
}

// "let" has already been consumed.
func (p *Parser) parseLet() Stmt {
	//ident = expr;

	ident, err := p.consume(itemIdentifier)
	if err != nil {
		return nil //next token not and identifier
	}

	fmt.Print(ident)
	_, err = p.consume(itemAssign)
	if err != nil {
		return nil //next token not and identifier
	}

	expr := p.parseExpr(LOWEST)

	stmt := &LetStmt{
		Pos:   ident.pos,
		Ident: ident,
		Val:   expr,
	}

	return stmt
}

func (p *Parser) parseExpr(pre precedence) Expr {
	t := p.next()
	prefix := p.prefixParseFns[t.typ]
	if prefix == nil {
		return nil
	}
	p.backup()
	leftExp := prefix()

	for n := p.peek(); n.typ != itemTerminateLine && pre < n.precedence(); n = p.peek() {
		infix := p.infixParseFns[n.typ]
		if infix == nil {
			return leftExp
		}

		leftExp = infix(leftExp)
	}

	return leftExp
}

//next is an identifier
func (p *Parser) parseIdent() Expr {
	n := p.next()
	return &Ident{Pos: n.pos, Value: n.Val}

}

func (p *Parser) parseNumber() Expr {
	//string, int, float, or bool
	n := p.next()
	num := &NumberLit{Pos: n.pos}
	//try parse int, if not ok then try parse float
	iv, err := strconv.ParseInt(n.Val, 10, 64)
	if err == nil {
		num.IntVal = iv
		num.IsInt = true
		num.FloatVal = float64(iv)
	} else {
		fv, err := strconv.ParseFloat(n.Val, 64)
		if err != nil {
			panic("invalid number")
		}
		num.FloatVal = fv
	}
	return num
}

func (p *Parser) parseUnaryExpr() Expr {
	n := p.next()
	switch n.typ {
	case LogicNot:
	case itemMinus:
	default:
		panic("unrecognized unary operator")
	}
	expr := &UnaryExpr{
		Pos: n.pos,
		Op:  n,
	}
	expr.Right = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parseBinaryExpr(left Expr) Expr {
	n := p.next()
	expr := &BinaryExpr{
		Pos:  n.pos,
		Op:   n,
		Left: left,
	}
	pr := n.precedence()
	expr.Right = p.parseExpr(pr)
	return expr
}

func (p *Parser) parseParen() Expr {
	//skip the paren
	p.next()

	exp := p.parseExpr(LOWEST)

	if n := p.peek(); n.typ != itemRightParen {
		return nil
	}
	p.next() // consume the right paren

	return exp
}

func (p *Parser) parseFn() *FnStmt {
	return nil
}
