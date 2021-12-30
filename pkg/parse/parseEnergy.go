package parse

import (
	"errors"
	"fmt"
)

func parseEnergyEvent(p *Parser) (parseFn, error) {
	//energy once interval=300 amount=1 #once at frame 300
	//energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	n := p.next()
	switch n.typ {
	case itemOnce:
		return parseEnergyOnce, nil
	case itemEvery:
		return parseEnergyEvery, nil
	default:
		return nil, fmt.Errorf("<energy> bad token at line %v - %v: %v", n.line, n.pos, n)
	}

}

func parseEnergyOnce(p *Parser) (parseFn, error) {
	//energy once interval=300 amount=1 #once at frame 300
	var err error
	p.cfg.Energy.Active = true
	p.cfg.Energy.Once = true

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemInterval:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				p.cfg.Energy.Start, err = itemNumberToInt(n)
			}
		case itemAmount:
			item, err := p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err != nil {
				return nil, err
			}
			count, err := itemNumberToInt(item)
			if err != nil {
				return nil, err
			}
			p.cfg.Energy.Particles = count
		case itemTerminateLine:
			return parseRows, nil
		default:
			err = fmt.Errorf("<energy> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing energy event")
}

func parseEnergyEvery(p *Parser) (parseFn, error) {
	//energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	var err error
	p.cfg.Energy.Active = true
	p.cfg.Energy.Once = false

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemInterval:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err != nil {
				return nil, err
			}
			p.cfg.Energy.Start, err = itemNumberToInt(n)
			if err != nil {
				return nil, err
			}

			n, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			p.cfg.Energy.End, err = itemNumberToInt(n)
			if err != nil {
				return nil, err
			}
		case itemAmount:
			item, err := p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err != nil {
				return nil, err
			}
			count, err := itemNumberToInt(item)
			if err != nil {
				return nil, err
			}
			p.cfg.Energy.Particles = count
		case itemTerminateLine:
			return parseRows, nil
		default:
			err = fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing energy event")
}
