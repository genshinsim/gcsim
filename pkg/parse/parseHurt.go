package parse

import (
	"errors"
	"fmt"
)

func parseHurtEvent(p *Parser) (parseFn, error) {
	//hurt+=once interval=300 amount=100,200 ele=pyro #once at frame 300 (or nearest)
	//hurt+=every interval=300,600 amount=100,200 ele=physical #randomly 100 to 200 dmg every 300 to 600 frames
	n := p.next()
	switch n.typ {
	case itemOnce:
		return parseHurtOnce, nil
	case itemEvery:
		return parseHurtEvery, nil
	default:
		return nil, fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
	}

}

func parseHurtOnce(p *Parser) (parseFn, error) {
	//interval=300 amount=100,200 ele=pyro #once at frame 300 (or nearest)
	var err error
	p.cfg.Hurt.Active = true
	p.cfg.Hurt.Once = true

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemInterval:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.cfg.Hurt.Start, err = itemNumberToInt(n)
			}
		case itemAmount:
			err = p.parseHurtAmount()
		case itemEle:
			err = p.parseHurtEle()
		case itemTerminateLine:
			return parseRows, nil
		default:
			err = fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing hurt event")
}

func parseHurtEvery(p *Parser) (parseFn, error) {
	//interval=300,600 amount=100,200 ele=physical #randomly 100 to 200 dmg every 300 to 600 frames
	var err error
	p.cfg.Hurt.Active = true
	p.cfg.Hurt.Once = false

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemInterval:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			p.cfg.Hurt.Start, err = itemNumberToInt(n)
			if err != nil {
				return nil, err
			}

			n, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			p.cfg.Hurt.End, err = itemNumberToInt(n)
			if err != nil {
				return nil, err
			}
		case itemAmount:
			err = p.parseHurtAmount()
		case itemEle:
			err = p.parseHurtEle()
		case itemTerminateLine:
			return parseRows, nil
		default:
			err = fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing hurt event")
}

func (p *Parser) parseHurtAmount() error {
	item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
	if err != nil {
		return err
	}
	min, err := itemNumberToFloat64(item)
	if err != nil {
		return err
	}

	item, err = p.acceptSeqReturnLast(itemComma, itemNumber)
	if err != nil {
		return err
	}
	max, err := itemNumberToFloat64(item)
	if err != nil {
		return err
	}

	p.cfg.Hurt.Min = min
	p.cfg.Hurt.Max = max

	return nil
}

func (p *Parser) parseHurtEle() error {
	_, err := p.consume(itemAssign)
	if err != nil {
		return err
	}
	n := p.next()
	if n.typ <= eleTypeKeyword {
		return fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
	}
	p.cfg.Hurt.Ele = eleKeys[n.val]
	return nil
}
