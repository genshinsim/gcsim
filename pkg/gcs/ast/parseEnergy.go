package ast

import (
	"errors"
	"fmt"
)

func parseEnergy(p *Parser) (parseFn, error) {
	//energy every=?? amount=??
	var err error
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "every":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Energy.Every, err = itemNumberToFloat64(n)
				}
			case "amount":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Energy.Amount, err = itemNumberToInt(n)
				}
			default:
				return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing options: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing options")
}
