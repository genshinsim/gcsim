package parse

import (
	"errors"
	"fmt"
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

	n := p.next()

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

	switch n.typ {
	case itemCharacterKey:
	case itemOptions:
		return parseOptions, nil
	case itemChain:
	case itemWait:
		return parseWait, nil
	case itemResetLimit:
	case itemIdentifier:
		//this is for macros. next has to be a colon
		x, err := p.consume(itemColon)
		if err != nil {
			return nil, fmt.Errorf("<parse row expecting : after an identifier but got %v; line %v", x, p.tokens)
		}
	}

	return nil, fmt.Errorf("<parse row> invalid token at start of line: %v", p.tokens)
}
