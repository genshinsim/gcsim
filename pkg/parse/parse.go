package parse

import (
	"errors"
)

func (p *Parser) Parse(text string) (*ActionList, error) {
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
	switch n := p.next(); n.typ {
	case itemCharacterKey:
		//parse character
		p.backup()
		return parseCharacter, nil
	//case options:
	//case active:
	//case target:
	//case energy:
	default: //default should be look for gcsl
		return parseProgram, nil
	}
}

func parseCharacter(p *Parser) (parseFn, error) {
	return nil, nil
}

func parseProgram(p *Parser) (parseFn, error) {

	//each line should start with one of the following:
	//	let
	//	a function call or action
	//	if
	//	while
	//loop through each line and add it to the block
	for {
		switch n := p.peek(); {
		case n.typ == keywordLet:
			node, err := p.parseLet()
			if err != nil {
				return nil, err
			}
			p.res.Program.append(node)
		case n.typ == itemIdentifier || n.typ == itemNumber:
			node, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			p.res.Program.append(node)
		case n.typ == itemCharacterKey:
			t1 := p.next() //should be char key
			if p.peek().typ <= itemActionKey {
				//not an ActionStmt
				p.backup2(t1)
				return parseCharacter, nil
			}
			//parse action item
			p.backup2(t1)
			node, err := p.parseAction()
			if err != nil {
				return nil, err
			}
			p.res.Program.append(node)
		//case options:
		//case active:
		//case target:
		//case energy:
		default:
			return nil, errors.New("unrecognized token")

		}
	}
}

//parseAction returns a node contain a character action, or a block of node containing
//a list of character actions
func (p *Parser) parseAction() (Stmt, error) {

	return nil, nil
}

// "let" has already been consumed.
func (p *Parser) parseLet() (Stmt, error) {
	//can be one of
	//ident = fn
	//ident = expr

	ident, err := p.consume(itemIdentifier)
	if err != nil {
		return nil, err //next token not and identifier
	}

	_, err = p.consume(itemAssign)
	if err != nil {
		return nil, err //next token not and identifier
	}

	switch n := p.peek(); {
	case n.typ == keywordFn:
	case n.typ == itemIdentifier || n.typ == itemNumber:
	default:
		return nil, errors.New("unrecognized token")
	}

	return nil, nil
}

func (p *Parser) parseExpr() (Expr, error) {
	return nil, nil
}

func (p *Parser) parseFn() (*FnStmt, error) {
	return nil, nil
}
