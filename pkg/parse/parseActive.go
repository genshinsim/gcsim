package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func parseActiveChar(p *Parser) (parseFn, error) {
	n := p.next()
	if n.typ != itemCharacterKey {
		return nil, fmt.Errorf("active char expecting char key got %v at line %v", n, p.tokens)
	}
	p.cfg.Characters.Initial = keys.CharNameToKey[n.val]
	n = p.next()
	if n.typ != itemTerminateLine {
		return nil, fmt.Errorf("active char expecting ; got %v at line %v", n, p.tokens)
	}

	return parseRows, nil
}
