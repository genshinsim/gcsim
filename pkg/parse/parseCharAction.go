package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseCharActions(p *Parser) (parseFn, error) {
	//xiangling burst,skill;
	//raidenshongun attack:4,charge,attack:4,charge <conditions>

	block, err := p.acceptCharAction()
	if err != nil {
		return nil, err
	}

	p.cfg.Rotation = append(p.cfg.Rotation, block)

	return parseRows, nil
}

func (p *Parser) acceptCharAction() (core.ActionBlock, error) {
	var err error

	block := core.ActionBlock{
		Type:         core.ActionBlockTypeSequence,
		SequenceChar: p.currentCharKey,
	}

	//create an action block and add abilities to it
	block.Sequence, err = p.acceptAbilitiesReturnList()
	if err != nil {
		return block, err
	}

	//loop till end of line and check for flags
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemPlus:
			//expecting a + sign since it's all optional fields
			err = p.acceptActionBlockFlags(&block)
			if err != nil {
				return block, err
			}
		case itemTerminateLine:
			//check block then add to action list
			ok := validateBlock(&block)
			if !ok {
				return block, fmt.Errorf("ln%v: invalid action at %v", n.line, p.tokens)
			}
			return block, nil
		default:
			return block, fmt.Errorf("ln%v: (accept char action) unexpected token %v", n.line, n)
		}
	}
	return block, errors.New("unexpected end of line parsing character action")
}

func validateBlock(a *core.ActionBlock) bool {
	//TODO: add action block validation
	return true
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
		key, ok := core.CharNameToKey[x.val]
		if !ok {
			err = fmt.Errorf("ln%v: bad token at line %v - %v: %v; invalid char name; line: %v", n.line, x.line, x.pos, x, p.tokens)
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
		//next should be try or drop
		n = p.next()
		switch n.typ {
		case itemDrop:
			a.TryDropIfNotReady = true
		case itemWait:
			a.TryDropIfNotReady = false
		default:
			err = fmt.Errorf("ln%v: unexpected token after try=, expecting drop or wait, got %v", n.line, n)
		}
	}
	return err
}
