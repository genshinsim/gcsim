package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (p *Parser) acceptLevelReturnBaseMax() (base, max int, err error) {
	//expect =xx/yy
	var x item
	x, err = p.consume(itemEqual)
	if err != nil {
		err = fmt.Errorf("unexpected token after lvl. expecting = got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("expecting a number for base lvl, got %v at line %v", x, p.tokens)
		return
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("unexpected token for base lvl. got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemForwardSlash)
	if err != nil {
		err = fmt.Errorf("expecting / separator for lvl, got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("expecting a number for max lvl, got %v at line %v", x, p.tokens)
		return
	}
	max, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("unexpected token for lvl. got %v at line %v", x, p.tokens)
		return
	}
	if max < base {
		err = fmt.Errorf("max level %v cannot be less than base level %v at line %v", max, base, p.tokens)
		return
	}
	return
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

		item, err := p.acceptSeqReturnLast(itemEqual, itemNumber)
		if err != nil {
			return r, err
		}

		r[i.val], err = itemNumberToInt(item)
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
			return r, fmt.Errorf("<action param> bad token at line %v - %v: %v", n.line, n.pos, n)
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
		return count, fmt.Errorf("expected a number after : but got %v on line %v", n, p.tokens)
	}
	//parse number
	count, err := itemNumberToInt(n)
	return count, err
}

func (p *Parser) acceptAbilitiesReturnList() ([]core.ActionItem, error) {
	//raidenshongun attack:4,charge,attack:4,charge
	var res []core.ActionItem
	var err error
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemActionKey:
			act := core.ActionItem{
				Typ:    actionKeys[n.val],
				Target: p.currentCharKey,
			}
			//optional params
			act.Param, err = p.acceptOptionalParamReturnMap()
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
				res = append(res, act)
			}
			//stop if not ,
			n = p.next()
			if n.typ != itemComma {
				p.backup()
				return res, nil
			}
		default:
			return nil, fmt.Errorf("unexpected token %v in line %v", n, p.tokens)
		}
	}
	return nil, errors.New("unexpected end of line")
}

func (p *Parser) acceptMacrosReturnList() ([]core.ActionBlock, error) {
	//id:4,id,id,id
	var res []core.ActionBlock
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemIdentifier:
			//see if macro exists
			block, ok := p.macros[n.val]
			if !ok {
				return nil, fmt.Errorf("macro %v not defined at line %v", n.val, p.tokens)
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
				res = append(res, block)
			}
			//stop if not ,
			n = p.next()
			if n.typ != itemComma {
				p.backup()
				return res, nil
			}
		default:
			return nil, fmt.Errorf("unexpected token %v in line %v", n, p.tokens)
		}
	}
	return nil, errors.New("unexpected end of line")
}

func (p *Parser) acceptActionBlockFlags(a *core.ActionBlock) error {
	/**
	- `+if`: condition that must be fulfiled before this line can be executed
	- `+swap_to`: forcefully swap to another char after this line has finished
	- `+swap_lock`: force the sim to stay on this character for x frames
	- `+is_onfield`: this line can only be used if this character is already on field
	- `+label`: a name for this line; can't have duplicated labels
	- `+needs`: this line can only be executed if the previous action is the referred to label
	- `+limit`: number of times this line can be executed (replaces once)
	- `+timeout`: this line cannot be executed again for x number of frames
	- `+try`: if this flag is set, then this line will execute as long as first ability in the list is executable.
			  If flag is set to 1, then if the sim will keep trying to execute the next ability in the list even if
			  it's not ready (even if it means waiting forever). If flag is set to 0, then if the next ability is not
			  ready immediately after the previous ability, the next ability along with the rest of the sequence will be dropped
	**/
	var x item
	var err error
	n := p.next()
	switch n.typ {
	case itemIf:
		a.Conditions, err = p.parseIf()
	case itemSwap:
		x, err = p.acceptSeqReturnLast(itemEqual, itemCharacterKey)
		key, ok := keys.CharNameToKey[x.val]
		if !ok {
			err = fmt.Errorf("bad token at line %v - %v: %v; invalid char name", n.line, n.pos, n)
		}
		a.SwapTo = key
	case itemSwapLock:
		x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
		if err == nil {
			a.SwapLock, err = itemNumberToInt(x)
		}
	case itemOnField:
		a.OnField = true
	case itemLabel:
		x, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
		a.Label = x.val
	case itemLimit:
		x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
		if err == nil {
			a.Limit, err = itemNumberToInt(x)
		}
	case itemTimeout:
		x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
		if err == nil {
			a.Timeout, err = itemNumberToInt(x)
		}
	case itemNeeds:
		x, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
		a.Needs = x.val
	case itemTry:
		a.Try = true
		//equal is optional
		n = p.next()
		if n.typ != itemEqual {
			p.backup()
			break
		}
		x, err = p.consume(itemNumber)
		if err == nil {
			var val int
			val, err = itemNumberToInt(x)
			if val == 0 {
				a.TryDropIfNotReady = true
			}
		}
	}
	return err
}
