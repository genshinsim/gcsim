package parse

import (
	"fmt"
)

type actionItem struct {
	typ   Token
	param map[string]int
}

type actionAPLOpt struct {
	onField           bool
	limit             int
	timeout           int
	swapTo            Token //character to swap to
	swapLock          int
	try               bool
	tryDropIfNotReady bool
}

func (p *Parser) parseCharacterAction() Stmt {
	//actions can be
	//apl options:
	//	+if
	//	+swap_to
	//	+swap_lock
	//	+is_onfield
	//	+label
	//	+needs
	//	+limit
	//	+timeout
	//	+try
	char, err := p.consume(itemCharacterKey)
	if err != nil {
		panic("parse char action expects character key, got " + char.String())
	}

	//should be multiple action keys next
	var actions []actionItem
	if p.peek().typ != itemActionKey {
		//TODO: fix error logging
		return nil
	}

	//all actions needs to come before any + flags
Loop:
	for {
		switch n := p.next(); n.typ {
		case itemTerminateLine:
			//stop here
			break Loop
		case itemActionKey:
			a := actionItem{
				typ: n,
			}
			//check for param -> then repeat
			a.param, err = p.acceptOptionalParamReturnMap()
			if err != nil {
				//TODO: fix error logging
				return nil
			}
			//optional : and a number
			count, err := p.acceptOptionalRepeaterReturnCount()
			if err != nil {
				//TODO: fix error logging
				return nil
			}
			//add to array
			for i := 0; i < count; i++ {
				//TODO: all the repeated action will access the same map
				//ability implement should avoid modifying the maps
				actions = append(actions, a)
			}

			n = p.next()
			if n.typ != itemComma {
				p.backup()
				break Loop
			}
		default:
			//TODO: fix invalid key error
			return nil
		}
	}
	//check for optional flags

	//build stmt

	return nil
}

func (p *Parser) acceptOptionalParamReturnMap() (map[string]int, error) {
	r := make(map[string]int)

	//check for params
	n := p.next()
	if n.typ != itemLeftSquareParen {
		p.backup()
		return r, nil
	}

	//loop until we hit square paren
	for {
		//we're expecting ident = int
		i, err := p.consume(itemIdentifier)
		if err != nil {
			return r, err
		}

		item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
		if err != nil {
			return r, err
		}

		r[i.Val], err = itemNumberToInt(item)
		if err != nil {
			return r, err
		}

		//if we hit ], return; if we hit , keep going, other wise error
		n := p.next()
		switch n.typ {
		case itemRightSquareParen:
			return r, nil
		case itemComma:
			//do nothing, keep going
		default:
			return r, fmt.Errorf("ln%v: <action param> bad token %v", n.line, n)
		}
	}
}

func (p *Parser) acceptOptionalRepeaterReturnCount() (int, error) {
	count := 1
	n := p.next()
	if n.typ != itemColon {
		p.backup()
		return count, nil
	}
	//should be a number next
	n = p.next()
	if n.typ != itemNumber {
		return count, fmt.Errorf("ln%v: expected a number after : but got %v", n.line, n)
	}
	//parse number
	count, err := itemNumberToInt(n)
	return count, err
}
