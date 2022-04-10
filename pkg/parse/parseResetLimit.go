package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseResetLimit(p *Parser) (parseFn, error) {
	block, err := p.acceptResetLimit()
	if err != nil {
		return nil, err
	}
	p.cfg.Rotation = append(p.cfg.Rotation, block)
	return parseRows, nil
}

func (p *Parser) acceptResetLimit() (core.ActionBlock, error) {
	block := core.ActionBlock{
		Type: core.ActionBlockTypeResetLimit,
	}
	//next token should be ;
	n := p.next()
	if n.typ != itemTerminateLine {
		return block, fmt.Errorf("ln%v: reset_limit expecting ; got %v", n.line, n)
	}

	return block, nil
}
