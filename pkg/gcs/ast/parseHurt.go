package ast

import (
	"errors"
	"fmt"
)

func parseHurt(p *Parser) (parseFn, error) {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300 (or nearest)
	// hurt every interval=480,720 amount=1,300 element=physical #randomly 1 to 300 dmg every 480 to 720 frames
	n := p.next()
	switch n.Typ {
	case itemIdentifier:
		switch n.Val {
		case "once":
			return parseHurtOnce, nil
		case "every":
			return parseHurtEvery, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.line, n.Val)
		}
	case itemTerminateLine:
		return parseRows, nil
	default:
		return nil, fmt.Errorf("ln%v: unrecognized token parsing options: %v", n.line, n)
	}
}

func parseHurtOnce(p *Parser) (parseFn, error) {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300
	var err error
	p.res.HurtSettings.Active = true
	p.res.HurtSettings.Once = true

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "interval":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.HurtSettings.Start, err = itemNumberToInt(n)
				}
			case "amount":
				err := parseHurtAmount(p)
				if err != nil {
					return nil, err
				}
			case "element":
				err := parseHurtElement(p)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("ln%v: unrecognized hurt event specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing hurt event: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing hurt event")
}

func parseHurtEvery(p *Parser) (parseFn, error) {
	// hurt every interval=480,720 amount=1,300 element=physical #randomly 1 to 300 dmg every 480 to 720 frames
	var err error
	p.res.HurtSettings.Active = true
	p.res.HurtSettings.Once = false

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "interval":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				p.res.HurtSettings.Start, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}

				n, err = p.acceptSeqReturnLast(itemComma, itemNumber)
				if err != nil {
					return nil, err
				}
				p.res.HurtSettings.End, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}
			case "amount":
				err := parseHurtAmount(p)
				if err != nil {
					return nil, err
				}
			case "element":
				err := parseHurtElement(p)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("ln%v: unrecognized hurt event specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing hurt event: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing hurt event")
}

func parseHurtAmount(p *Parser) error {
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

	p.res.HurtSettings.Min = min
	p.res.HurtSettings.Max = max

	return nil
}

func parseHurtElement(p *Parser) error {
	_, err := p.consume(itemAssign)
	if err != nil {
		return err
	}
	n := p.next()
	if n.Typ != itemElementKey {
		return fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.line, n.pos, n)
	}
	p.res.HurtSettings.Element = eleKeys[n.Val]
	return nil
}
