package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	//next should be a string
	_, err = p.consume(itemString)
	if err != nil {
		return nil, err
	}
	var r core.EnemyProfile
	r.Resist = make(map[core.EleType]float64)
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch {
		case n.typ == itemLvl:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case n.typ == statHP:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				p.cfg.DamageMode = true
			}
		case n.typ > eleTypeKeyword:
			s := n.val
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			r.Resist[eleKeys[s]] += amt
		case n.typ == itemTerminateLine:
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
