package ast

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/shortcut"
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

//parseAction returns a node contain a character action, or a block of node containing
//a list of character actions
func (p *Parser) parseAction() (Stmt, error) {
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
		//this really shouldn't happen since we already checked
		return nil, fmt.Errorf("ln%v: expecting character key, got %v", char.line, char.Val)
	}
	charKey := shortcut.CharNameToKey[char.Val]

	//should be multiple action keys next
	var actions []*ActionStmt
	if n := p.peek(); n.Typ != itemActionKey {
		return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.line, char.Val, n.Val)
	}

	//all actions needs to come before any + flags
Loop:
	for {
		switch n := p.next(); n.Typ {
		case itemTerminateLine:
			//stop here
			break Loop
		case itemActionKey:
			a := &ActionStmt{
				Pos:    char.pos,
				Char:   charKey,
				Action: actionKeys[n.Val],
			}
			//check for param -> then repeat
			a.Param, err = p.acceptOptionalParamReturnMap()
			if err != nil {
				return nil, err
			}
			//optional : and a number
			count, err := p.acceptOptionalRepeaterReturnCount()
			if err != nil {
				return nil, err
			}
			//add to array
			for i := 0; i < count; i++ {
				//TODO: all the repeated action will access the same map
				//ability implement should avoid modifying the maps
				actions = append(actions, a)
			}

			n = p.next()
			if n.Typ != itemComma {
				p.backup()
				break Loop
			}
		default:
			//TODO: fix invalid key error
			return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.line, char.Val, n.Val)
		}
	}
	//check for optional flags

	//build stmt

	if len(actions) == 1 {
		return actions[0], nil
	} else {
		b := newBlockStmt(char.pos)
		for _, v := range actions {
			b.append(v)
		}
		return b, nil
	}
}

func (p *Parser) acceptOptionalParamReturnMap() (map[string]int, error) {
	r := make(map[string]int)

	//check for params
	n := p.next()
	if n.Typ != itemLeftSquareParen {
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
		switch n.Typ {
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
	if n.Typ != itemColon {
		p.backup()
		return count, nil
	}
	//should be a number next
	n = p.next()
	if n.Typ != itemNumber {
		return count, fmt.Errorf("ln%v: expected a number after : but got %v", n.line, n)
	}
	//parse number
	count, err := itemNumberToInt(n)
	return count, err
}
