package ast

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r enemy.EnemyProfile
	r.Resist = make(map[attributes.Element]float64)
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case keywordLvl:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case itemStatKey:
			//should be hp
			if statKeys[n.Val] != attributes.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				p.res.Settings.DamageMode = true
			}
		case keywordResist:
			//this sets all resistance
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			//TODO: make this more elegant...
			r.Resist[attributes.Electro] += amt
			r.Resist[attributes.Cryo] += amt
			r.Resist[attributes.Hydro] += amt
			r.Resist[attributes.Physical] += amt
			r.Resist[attributes.Pyro] += amt
			r.Resist[attributes.Geo] += amt
			r.Resist[attributes.Dendro] += amt
			r.Resist[attributes.Anemo] += amt

		case itemElementKey:
			s := n.Val
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			r.Resist[eleKeys[s]] += amt
		case itemTerminateLine:
			p.res.Targets = append(p.res.Targets, r)
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
