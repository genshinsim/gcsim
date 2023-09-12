package ast

import (
	"errors"
	"fmt"
)

func parseEnergy(p *Parser) (parseFn, error) {
	// energy once interval=300 amount=1 #once at frame 300
	// energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	n := p.next()
	switch n.Typ {
	case itemIdentifier:
		switch n.Val {
		case "once":
			return parseEnergyOnce, nil
		case "every":
			return parseEnergyEvery, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.line, n.Val)
		}
	case itemTerminateLine:
		return parseRows, nil
	default:
		return nil, fmt.Errorf("ln%v: unrecognized token parsing options: %v", n.line, n)
	}
}

func parseEnergyOnce(p *Parser) (parseFn, error) {
	// energy once interval=300 amount=1 #once at frame 300
	var err error
	p.res.EnergySettings.Active = true
	p.res.EnergySettings.Once = true

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "interval":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.EnergySettings.Start, err = itemNumberToInt(n)
				}
			case "amount":
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				count, err := itemNumberToInt(item)
				if err != nil {
					return nil, err
				}
				p.res.EnergySettings.Amount = count
			default:
				return nil, fmt.Errorf("ln%v: unrecognized energy event specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing energy event: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing energy event")
}

func parseEnergyEvery(p *Parser) (parseFn, error) {
	// energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	var err error
	p.res.EnergySettings.Active = true
	p.res.EnergySettings.Once = false

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "interval":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				p.res.EnergySettings.Start, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}

				n, err = p.acceptSeqReturnLast(itemComma, itemNumber)
				if err != nil {
					return nil, err
				}
				p.res.EnergySettings.End, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}
			case "amount":
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				count, err := itemNumberToInt(item)
				if err != nil {
					return nil, err
				}
				p.res.EnergySettings.Amount = count
			default:
				return nil, fmt.Errorf("ln%v: unrecognized energy event specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing energy event: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing energy event")
}
