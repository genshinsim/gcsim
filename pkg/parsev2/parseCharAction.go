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
				return block, fmt.Errorf("invalid action at %v", p.tokens)
			}
			return block, nil
		default:
			return block, fmt.Errorf("(accept char action) unexpected token %v at line %v", n, p.tokens)
		}
	}
	return block, errors.New("unexpected end of line parsing character action")
}

func validateBlock(a *core.ActionBlock) bool {
	//TODO: add action block validation
	return true
}
