package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (p *Parser) acceptLevelReturnBaseMax() (base, max int, err error) {
	//expect =xx/yy
	var x item
	x, err = p.consume(itemEqual)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token after lvl. expecting = got %v", x.line, x)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting a number for base lvl, got %v", x.line, x)
		return
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token for base lvl. got %v", x.line, x)
		return
	}
	x, err = p.consume(itemForwardSlash)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting / separator for lvl, got %v", x.line, x)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting a number for max lvl, got %v", x.line, x)
		return
	}
	max, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token for lvl. got %v", x.line, x)
		return
	}
	if max < base {
		err = fmt.Errorf("ln%v: max level %v cannot be less than base level %v", x.line, max, base)
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
			return nil, fmt.Errorf("ln%v: unexpected token %v", n.line, n)
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
				return nil, fmt.Errorf("ln%v: macro %v not defined", n.line, n.val)
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
			return nil, fmt.Errorf("ln%v: unexpected token %v", n.line, n)
		}
	}
	return nil, errors.New("unexpected end of line")
}
