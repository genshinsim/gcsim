package parser

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func parseHurt(p *Parser) (parseFn, error) {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300 (or nearest)
	// hurt every interval=480,720 amount=1,300 element=physical #randomly 1 to 300 dmg every 480 to 720 frames
	n := p.next()
	switch n.Typ {
	case ast.ItemIdentifier:
		switch n.Val {
		case "once":
			return parseHurtOnce, nil
		case "every":
			return parseHurtEvery, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.Line, n.Val)
		}
	case ast.ItemTerminateLine:
		return parseRows, nil
	default:
		return nil, fmt.Errorf("ln%v: unrecognized token parsing options: %v", n.Line, n)
	}
}

func parseHurtOnce(p *Parser) (parseFn, error) {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300
	var err error
	p.res.HurtSettings.Active = true
	p.res.HurtSettings.Once = true

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemIdentifier:
			switch n.Val {
			case ast.IntervalVal:
				n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err == nil {
					p.res.HurtSettings.Start, err = itemNumberToInt(n)
				}
			case ast.AmountVal:
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
				return nil, fmt.Errorf("ln%v: unrecognized hurt event specified: %v", n.Line, n.Val)
			}
		case ast.ItemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing hurt event: %v", n.Line, n)
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

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemIdentifier:
			switch n.Val {
			case ast.IntervalVal:
				n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				p.res.HurtSettings.Start, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}

				n, err = p.acceptSeqReturnLast(ast.ItemComma, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				p.res.HurtSettings.End, err = itemNumberToInt(n)
				if err != nil {
					return nil, err
				}
			case ast.AmountVal:
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
				return nil, fmt.Errorf("ln%v: unrecognized hurt event specified: %v", n.Line, n.Val)
			}
		case ast.ItemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing hurt event: %v", n.Line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing hurt event")
}

func parseHurtAmount(p *Parser) error {
	item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
	if err != nil {
		return err
	}
	minhurt, err := itemNumberToFloat64(item)
	if err != nil {
		return err
	}

	item, err = p.acceptSeqReturnLast(ast.ItemComma, ast.ItemNumber)
	if err != nil {
		return err
	}
	maxhurt, err := itemNumberToFloat64(item)
	if err != nil {
		return err
	}

	p.res.HurtSettings.Min = minhurt
	p.res.HurtSettings.Max = maxhurt

	return nil
}

func parseHurtElement(p *Parser) error {
	_, err := p.consume(ast.ItemAssign)
	if err != nil {
		return err
	}
	n := p.next()
	if n.Typ != ast.ItemElementKey {
		return fmt.Errorf("<hurt> bad token at line %v - %v: %v", n.Line, n.Pos, n)
	}
	p.res.HurtSettings.Element = ast.EleKeys[n.Val]
	return nil
}
