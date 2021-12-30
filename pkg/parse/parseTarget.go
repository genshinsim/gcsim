package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r core.EnemyProfile
	r.Resist = make(map[core.EleType]float64)
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemLvl:
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case itemStatKey:
			//should be hp
			if statKeys[n.val] != core.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				p.cfg.DamageMode = true
			}
		case itemElementKey:
			s := n.val
			item, err := p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			r.Resist[eleKeys[s]] += amt
		case itemTerminateLine:
			p.cfg.Targets = append(p.cfg.Targets, r)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing target")
}
