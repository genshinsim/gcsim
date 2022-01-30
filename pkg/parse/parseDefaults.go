package parse

import "errors"

func parseDefaults(p *Parser) (parseFn, error) {
	var err error

	// Should handle something like the below:
	// defaults swapdelay=8 framedefaults="human";
	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemSwapDelay:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				p.cfg.Settings.SwapDelay, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			return parseRows, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing defaults")
}
