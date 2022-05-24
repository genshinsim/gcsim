package parse

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type precedence int

const (
	_ precedence = iota
	Lowest
	Equals
	LessOrGreater
	Sum
	Product
	Prefix
	Call
)

var precedences = map[TokenType]precedence{
	OpEqual:              Equals,
	OpNotEqual:           Equals,
	OpLessThan:           LessOrGreater,
	OpGreaterThan:        LessOrGreater,
	OpLessThanOrEqual:    LessOrGreater,
	OpGreaterThanOrEqual: LessOrGreater,
	itemPlus:             Sum,
	itemMinus:            Sum,
	itemSlash:            Product,
	itemAsterisk:         Product,
	itemLeftParen:        Call,
}

func (t Token) precedence() precedence {
	if p, ok := precedences[t.typ]; ok {
		return p
	}
	return Lowest
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
		if p.peek().typ != itemActionKey {
			//not an ActionStmt
			//set up char and set key
			key, ok := keys.CharNameToKey[n.Val]
			if !ok {
				//TODO: better err handling
				panic("invalid char key " + n.Val)
			}
			if _, ok := p.chars[key]; !ok {
				p.newChar(key)
			}
			p.currentCharKey = key
			p.backup()
			return parseCharacter, nil
		}
		p.backup()
		//parse action item
		// return parseProgram, nil
		node := p.parseStatement()
		p.res.Program.append(node)
		return parseText, nil
	case itemEOF:
		return nil, nil
	default: //default should be look for gcsl
		node := p.parseStatement()
		p.res.Program.append(node)
		return parseText, nil
	}
}

func (p *Parser) parseStatement() Node {
	//some statements end in semi, other don't
	hasSemi := true
	var node Node
	switch n := p.peek(); n.typ {
	case keywordBreak:
		fallthrough
	case keywordFallthrough:
		fallthrough
	case keywordContinue:
		node = p.parseCtrl()
	case keywordLet:
		node = p.parseLet()
	case itemCharacterKey:
		node = p.parseAction()
	case keywordReturn:
		node = p.parseReturn()
	case keywordIf:
		node = p.parseIf()
		hasSemi = false
	case keywordSwitch:
		node = p.parseSwitch()
		hasSemi = false
	case keywordFn:
		node = p.parseFn()
		hasSemi = false
	case keywordWhile:
		node = p.parseWhile()
		hasSemi = false
	case itemIdentifier:
		p.next()
		//check if = after
		if x := p.peek(); x.typ == itemAssign {
			p.backup()
			node = p.parseAssign()
			break
		}
		//it's an expr if no assign
		p.backup()
		fallthrough
	default:
		node = p.parseExpr(Lowest)
	}
	if hasSemi {
		n, err := p.consume(itemTerminateLine)
		if err != nil {
			panic("expecting ; got " + n.String())
		}
	}
	return node
}

func (p *Parser) parseLet() Stmt {
	//var ident = expr;
	n := p.next()

	ident, err := p.consume(itemIdentifier)
	if err != nil {
		//next token not and identifier
		panic("expecting ident after nil, got " + ident.String())
	}

	a, err := p.consume(itemAssign)
	if err != nil {
		//next token not and identifier
		panic("expecting assign after nil, got " + a.String())
	}

	expr := p.parseExpr(Lowest)

	stmt := &LetStmt{
		Pos:   n.pos,
		Ident: ident,
		Val:   expr,
	}

	return stmt
}

// expecting ident = expr
func (p *Parser) parseAssign() Stmt {

	ident, err := p.consume(itemIdentifier)
	if err != nil {
		//next token not and identifier
		panic("expecting ident after nil, got " + ident.String())
	}

	a, err := p.consume(itemAssign)
	if err != nil {
		//next token not and identifier
		panic("expecting assign after nil, got " + a.String())
	}

	expr := p.parseExpr(Lowest)

	stmt := &AssignStmt{
		Pos:   ident.pos,
		Ident: ident,
		Val:   expr,
	}

	return stmt

}

func (p *Parser) parseIf() Stmt {
	n := p.next()

	stmt := &IfStmt{
		Pos: n.pos,
	}

	stmt.Condition = p.parseExpr(Lowest)

	//expecting a { next
	if n := p.peek(); n.typ != itemLeftBrace {
		return nil
	}

	stmt.IfBlock = p.parseBlock() //parse block here

	//stop if no else
	if n := p.peek(); n.typ != keywordElse {
		return stmt
	}

	//skip the else keyword
	p.next()

	//expecting another block
	stmt.ElseBlock = p.parseBlock()

	return stmt
}

func (p *Parser) parseSwitch() Stmt {

	//switch expr { }
	n, err := p.consume(keywordSwitch)
	if err != nil {
		panic("unreachable")
	}

	stmt := &SwitchStmt{
		Pos: n.pos,
	}

	stmt.Condition = p.parseExpr(Lowest)

	if p.next().typ != itemLeftBrace {
		//TODO: handle switch error
		return nil
	}

	//look for cases while not }
	for n := p.next(); n.typ != itemRightBrace; n = p.next() {
		//expecting case expr: block
		switch n.typ {
		case keywordCase:
			cs := &CaseStmt{
				Pos: n.pos,
			}
			cs.Condition = p.parseExpr(Lowest)
			//colon, then read until we hit next case
			if p.peek().typ != itemColon {
				panic("expecting : got " + p.peek().String())
			}
			cs.Body = p.parseCaseBody()
			stmt.Cases = append(stmt.Cases, cs)
		case keywordDefault:
			//colon, then read until we hit next case
			if p.peek().typ != itemColon {
				panic("expecting : got " + p.peek().String())
			}
			stmt.Default = p.parseCaseBody()
		default:
			panic("expecting case or default token, got " + n.String())
		}

	}

	return stmt
}

func (p *Parser) parseCaseBody() *BlockStmt {
	n := p.next() //start with :
	block := newBlockStmt(n.pos)
	var node Node
	//parse line by line until we hit }
	for {
		//make sure we don't get any illegal lines
		switch n := p.peek(); n.typ {
		case itemCharacterKey:
			if !p.peekValidCharAction() {
				panic("unexpected non action statement with char in block")
			}
		case keywordDefault:
			fallthrough
		case keywordCase:
			fallthrough
		case itemRightBrace:
			return block
		case itemEOF:
			panic("reached end of file without }")
		}
		//parse statement here
		node = p.parseStatement()
		block.append(node)
	}
}

// while { }
func (p *Parser) parseWhile() Stmt {
	n := p.next()

	stmt := &WhileStmt{
		Pos: n.pos,
	}

	stmt.Condition = p.parseExpr(Lowest)

	//expecting a { next
	if n := p.peek(); n.typ != itemLeftBrace {
		return nil
	}

	stmt.WhileBlock = p.parseBlock() //parse block here

	return stmt
}

func (p *Parser) parseFn() Stmt {
	//fn ident(...ident){ block }
	//consume fn
	n := p.next()
	stmt := &FnStmt{
		Pos: n.pos,
	}

	//ident next
	n, err := p.consume(itemIdentifier)
	if err != nil {
		panic("expecting identifier after fn, got " + n.String())
	}
	stmt.FunVal = n

	if l := p.peek(); l.typ != itemLeftParen {
		//TODO: error handling here?
		return nil
	}

	stmt.Args = p.parseFnArgs()
	stmt.Body = p.parseBlock()

	return stmt
}

func (p *Parser) parseFnArgs() []*Ident {
	//consume (
	var args []*Ident
	p.next()
	for n := p.next(); n.typ != itemRightParen; n = p.next() {
		a := &Ident{}
		//expecting ident, comma
		if n.typ != itemIdentifier {
			panic("expecting ident in param list, got " + n.String())
		}
		a.Pos = n.pos
		a.Value = n.Val

		args = append(args, a)

		//if next token is a comma, then there should be another ident after that
		//otherwise we have a problem
		if l := p.peek(); l.typ == itemComma {
			p.next() //consume the comma
			if l = p.peek(); l.typ != itemIdentifier {
				panic("expecting another identifier after comma, got " + l.String())
			}
		}
	}
	return args
}

func (p *Parser) parseReturn() Stmt {
	n := p.next() //return
	stmt := &ReturnStmt{
		Pos: n.pos,
	}
	stmt.Val = p.parseExpr(Lowest)
	return stmt
}

func (p *Parser) parseCtrl() Stmt {
	n := p.next()
	stmt := &CtrlStmt{
		Pos: n.pos,
	}
	switch n.typ {
	case keywordBreak:
		stmt.Typ = CtrlBreak
	case keywordContinue:
		stmt.Typ = CtrlContinue
	case keywordFallthrough:
		stmt.Typ = CtrlFallthrough
	default:
		panic("invalid token, expecting a ctrl token, got " + n.String())
	}
	return stmt
}

func (p *Parser) parseCall(fun Expr) Expr {
	// ident has aready been consumed
	// switch fun.(type) {
	// case *Ident:
	// case *FnExpr:
	// default:
	// 	panic("invalid fun expression to function call")
	// }

	//expecting (params)
	n, err := p.consume(itemLeftParen)
	if err != nil {
		panic("expecting call to start with (")
	}
	expr := &CallExpr{
		Pos: n.pos,
		Fun: fun,
	}
	expr.Args = p.parseCallArgs()

	return expr

}

func (p *Parser) parseCallArgs() []Expr {
	var args []Expr

	if p.peek().typ == itemRightParen {
		return args
	}

	//next should be an expression
	args = append(args, p.parseExpr(Lowest))

	for p.peek().typ == itemComma {
		p.next() //skip the comma
		args = append(args, p.parseExpr(Lowest))
	}

	if p.next().typ != itemRightParen {
		p.backup()
		//TODO: handle error here
		return nil
	}

	return args
}

//check if it's a valid character action, assuming current token is "character"
func (p *Parser) peekValidCharAction() bool {
	p.next()
	//check if this is character stats etc or an action
	if p.peek().typ != itemActionKey {
		p.backup()
		//not an ActionStmt
		return false
	}
	p.backup()
	return true
}

//parseBlock return a node contain and BlockStmt
func (p *Parser) parseBlock() *BlockStmt {
	//should be surronded by {}
	n, err := p.consume(itemLeftBrace)
	if err != nil {
		//TODO: better parser error handling
		panic("expecting block to start with {")
	}
	block := newBlockStmt(n.pos)
	var node Node
	//parse line by line until we hit }
	for {
		//make sure we don't get any illegal lines
		switch n := p.peek(); n.typ {
		case itemCharacterKey:
			if !p.peekValidCharAction() {
				panic("unexpected non action statement with char in block")
			}
		case itemRightBrace:
			p.next() //consume the braces
			return block
		case itemEOF:
			panic("reached end of file without }")
		}
		//parse statement here
		node = p.parseStatement()
		block.append(node)
	}

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
	expr.Right = p.parseExpr(Prefix)
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

	exp := p.parseExpr(Lowest)

	if n := p.peek(); n.typ != itemRightParen {
		return nil
	}
	p.next() // consume the right paren

	return exp
}
