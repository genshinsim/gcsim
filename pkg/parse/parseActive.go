package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseActiveChar(p *Parser) (parseFn, error) {
	n := p.next()
	if n.typ != itemCharacterKey {
		return nil, fmt.Errorf("ln%v: active char expecting char key got %v", n.line, n)
	}
	p.cfg.Characters.Initial = core.CharNameToKey[n.val]
	n = p.next()
	if n.typ != itemTerminateLine {
		return nil, fmt.Errorf("ln%v: active char expecting ; got %v", n.line, n)
	}

	return parseRows, nil
}
