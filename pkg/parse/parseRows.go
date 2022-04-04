package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseRows(p *Parser) (parseFn, error) {
	p.tokens = make([]item, 0, 20)
	p.pos = -1

	//consume the entire line
	for n := p.l.nextItem(); n.typ != itemEOF; n = p.l.nextItem() {
		if n.typ == itemError {
			return nil, errors.New(n.val)
		}
		p.tokens = append(p.tokens, n)
		if n.typ == itemTerminateLine {
			break
		}
	}

	if len(p.tokens) == 0 {
		return nil, nil //at end of the file
	}

	/**
	beginning of the line we're expecting one of the following:
		- <character name>
		- options
		- chain
		- wait
		- reset_limit
		- <identifier> + :
		-
	**/

	n := p.next()

	switch n.typ {
	case itemTarget:
		return parseTarget, nil
	case itemEnergy:
		return parseEnergyEvent, nil
	case itemHurt:
		return parseHurtEvent, nil
	case itemCharacterKey:
		//after character name it shoudl be one of the following:
		//- char
		//- add
		//abilities

		//lex should have checked this already
		key, ok := core.CharNameToKey[n.val]
		if !ok {
			return nil, fmt.Errorf("ln%v: unexpected error, should be a recognized character key: %v", n.line, n)
		}
		if _, ok := p.chars[key]; !ok {
			p.newChar(key)
		}
		p.currentCharKey = key

		return parseChar, nil
	case itemOptions:
		return parseOptions, nil
	case itemChain:
		return parseChain, nil
	case itemWaitFor:
		return parseWait, nil
	case itemResetLimit:
		return parseResetLimit, nil
	case itemActive:
		return parseActiveChar, nil
	case itemRestart:
		return parseRestart, nil
	case itemIdentifier:
		//this is for macros. next has to be a colon
		x, err := p.consume(itemColon)
		if err != nil {
			return nil, fmt.Errorf("<parse row expecting : after an identifier but got %v; line %v", x, p.tokens)
		}
		return parseMacro, nil
	case itemWait: //THIS IS FOR CALC MODE ONLY
		return parseCalcModeWait, nil
	}

	return nil, fmt.Errorf("<parse row> invalid token at start of line: %v", p.tokens)
}
