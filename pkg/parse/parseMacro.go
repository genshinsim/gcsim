package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseMacro(p *Parser) (parseFn, error) {

	block := core.ActionBlock{
		Type:         core.ActionBlockTypeSequence,
		SequenceChar: p.currentCharKey,
	}

	//check type of macro
	n := p.next()
	switch n.typ {
	case itemCharacterKey:
		next, err := p.acceptCharAction()
		if err != nil {
			return nil, err
		}
		block.ChainSequences = append(block.ChainSequences, next)
	case itemWaitFor:
		next, err := p.acceptWait()
		if err != nil {
			return nil, err
		}
		block.ChainSequences = append(block.ChainSequences, next)
	case itemResetLimit:
		next, err := p.acceptResetLimit()
		if err != nil {
			return nil, err
		}
		block.ChainSequences = append(block.ChainSequences, next)
	default:
		//invalid
		return nil, fmt.Errorf("invalid token for macro %v, line %v", n, p.tokens)
	}

	//id for this macro should be first token
	n = p.tokens[0]

	p.macros[n.val] = block

	return parseRows, nil
}
