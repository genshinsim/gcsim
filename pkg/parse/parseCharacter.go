package parse

import "fmt"

func parseCharacter(p *Parser) (parseFn, error) {
	//expecting one of:
	//	char lvl etc
	// 	add
	//should be any action here
	switch n := p.next(); n.typ {
	case keywordChar:
	case keywordAdd:
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character>: %v", n.line, n)
	}
	return nil, nil
}

func parseCharacterAdd(p *Parser) (parseFn, error) {
	return nil, nil
}
