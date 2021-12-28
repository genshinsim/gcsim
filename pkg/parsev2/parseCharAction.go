package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func parseCharActions(p *Parser) (parseFn, error) {
	//xiangling burst,skill;
	//raidenshongun attack:4,charge,attack:4,charge <conditions>

	var err error

	block := core.ActionBlock{
		Type:         core.ActionBlockTypeSequence,
		SequenceChar: p.currentCharKey,
	}

	//create an action block and add abilities to it
	block.Sequence, err = p.acceptAbilitiesReturnList()
	if err != nil {
		return nil, err
	}

	//loop till end of line and check for conditions
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemPlus:
			//expecting a + sign since it's all optional fields
			err = p.consumeActionFlags(&block)
			if err != nil {
				return nil, err
			}
		case itemTerminateLine:
			//check block then add to action list
			ok := validateBlock(&block)
			if !ok {
				return nil, fmt.Errorf("invalid action at %v", p.tokens)
			}
			p.cfg.Rotation = append(p.cfg.Rotation, block)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("(parse char action) unexpected token %v at line %v", n, p.tokens)
		}
	}
	return nil, errors.New("unexpected end of line parsing character action")
}

func validateBlock(a *core.ActionBlock) bool {
	//TODO: add action block validation
	return true
}

func (p *Parser) consumeActionFlags(a *core.ActionBlock) error {
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
		x, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
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
