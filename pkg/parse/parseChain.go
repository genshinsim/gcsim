package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseChain(p *Parser) (parseFn, error) {
	//chain bt_seq,wait_for_bt,xq_seq,asdf /if=.xx.xx.xx>1 /swap=xiangling /limit=1 /try=1

	var err error

	block := core.ActionBlock{
		Type: core.ActionBlockTypeChain,
	}

	block.ChainSequences, err = p.acceptMacrosReturnList()
	if err != nil {
		return nil, err
	}

	//loop till end of line and check for flags
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemPlus:
			//expecting a + sign since it's all optional fields
			err = p.acceptActionBlockFlags(&block)
			if err != nil {
				return nil, err
			}
		case itemTerminateLine:
			//check block then add to action list
			ok := validateBlock(&block)
			if !ok {
				return nil, fmt.Errorf("ln%v: invalid action at %v", n.line, p.tokens)
			}
			p.cfg.Rotation = append(p.cfg.Rotation, block)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("(parse chain) unexpected token %v at pos %v line %v", n, p.pos, p.tokens)
		}
	}
	return nil, errors.New("unexpected end of line parsing character action")
}
