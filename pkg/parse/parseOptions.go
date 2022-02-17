package parse

import (
	"errors"
	"fmt"
)

func parseOptions(p *Parser) (parseFn, error) {
	//option iter=1000 duration=1000 worker=50 debug=true er_calc=true damage_mode=true
	var err error

	//options debug=true iteration=5000 duration=90 workers=24;
	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemIdentifier:
			//expecting identifier = some value
			switch n.val {
			case "debug":
				n, err = p.acceptSeqReturnLast(itemEqual, itemBool)
				// every run is going to have a debug from now on so we basically ignore what this flag says
			case "iteration":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Iterations, err = itemNumberToInt(n)
				}
			case "duration":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Duration, err = itemNumberToInt(n)
				}
			case "workers":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.NumberOfWorkers, err = itemNumberToInt(n)
				}
			case "mode":
				n, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
				if err == nil {
					//should be either apl or sl
					m, ok := queueModeKeys[n.val]
					if !ok {
						return nil, fmt.Errorf("invalid queue mode, got %v", n.val)
					}
					p.cfg.Settings.QueueMode = m
				}
			case "swap_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.SwapDelay, err = itemNumberToInt(n)
				}
			case "er_calc":
				//does nothing thus far...
			default:
				return nil, fmt.Errorf("unrecognized option specified: %v", n.val)
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
