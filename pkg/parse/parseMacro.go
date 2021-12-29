package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func parseMacro(p *Parser) (parseFn, error) {

	var block core.ActionBlock
	var err error

	//check type of macro
	n := p.next()
	switch n.typ {
	case itemCharacterKey:
		//lex should have checked this already
		key, ok := keys.CharNameToKey[n.val]
		if !ok {
			return nil, fmt.Errorf("unexpected error, should be a recognized character key: %v", n)
		}
		if _, ok := p.chars[key]; !ok {
			p.newChar(key)
		}
		p.currentCharKey = key
		block, err = p.acceptCharAction()
	case itemWaitFor:
		block, err = p.acceptWait()
	case itemResetLimit:
		block, err = p.acceptResetLimit()
	default:
		//invalid
		return nil, fmt.Errorf("invalid token for macro %v, line %v", n, p.tokens)
	}

	if err != nil {
		return nil, err
	}

	//id for this macro should be first token
	n = p.tokens[0]

	p.macros[n.val] = block

	return parseRows, nil
}
