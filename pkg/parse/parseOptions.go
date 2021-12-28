package parse

import "errors"

func parseOptions(p *Parser) (parseFn, error) {
	//option iter=1000 duration=1000 worker=50 debug=true er_calc=true damage_mode=true
	var err error

	//options debug=true iteration=5000 duration=90 workers=24;
	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemDebug:
			n, err = p.acceptSeqReturnLast(itemEqual, itemBool)
			if err == nil {
				p.opt.Debug = n.val == "true"
			}
		case itemIterations:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				p.opt.Iteration, err = itemNumberToInt(n)
			}
		case itemDuration:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				p.opt.Duration, err = itemNumberToInt(n)
			}
		case itemWorkers:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				p.opt.Workers, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			return parseRows, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing options")
}
