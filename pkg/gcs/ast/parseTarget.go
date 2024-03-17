package ast

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r info.EnemyProfile
	r.Resist = make(map[attributes.Element]float64)
	r.ParticleElement = attributes.NoElement
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			switch n.Val {
			case "pos": // pos will end up defaulting to 0,0 if not set
				// pos=1.00,2,00
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				x, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				item, err = p.acceptSeqReturnLast(itemComma, itemNumber)
				if err != nil {
					return nil, err
				}
				y, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.X = x
				r.Pos.Y = y
			case "radius":
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.R = amt
			case "type":
				item, err := p.acceptSeqReturnLast(itemAssign, itemIdentifier)
				if err != nil {
					return nil, err
				}
				params, err := p.acceptOptionalTargetParams()
				if err != nil {
					return nil, err
				}
				err = enemy.ConfigureTarget(&r, item.Val, params)
				if err != nil {
					return nil, err
				}
				p.res.Settings.DamageMode = true
			case "freeze_resist":
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				v, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.FreezeResist = v
				r.Modified = true
			default:
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
		case keywordLvl:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case itemStatKey:
			// should be hp
			if statKeys[n.Val] != attributes.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				if err != nil {
					return nil, err
				}
				p.res.Settings.DamageMode = true
				r.Modified = true
			}
		case keywordResist:
			// this sets all resistance
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			res := []attributes.Element{attributes.Electro, attributes.Cryo, attributes.Hydro, attributes.Physical, attributes.Pyro, attributes.Geo, attributes.Dendro, attributes.Anemo}
			for _, attr := range res {
				r.Resist[attr] += amt
			}
			r.Modified = true
		case keywordParticleThreshold:
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropThreshold = amt
			r.ParticleDrops = nil // separate particle system
			r.Modified = true
		case keywordParticleDropCount:
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropCount = amt
			r.Modified = true
		case keywordParticleElement:
			item, err := p.acceptSeqReturnLast(itemAssign, itemElementKey)
			if err != nil {
				return nil, err
			}
			if ele, ok := eleKeys[item.Val]; ok {
				r.ParticleElement = ele
			}
			r.Modified = true
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
			r.Modified = true
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

func (p *Parser) acceptOptionalTargetParams() (enemy.TargetParams, error) {
	result := enemy.TargetParams{
		HpMultiplier: 0.0,
		Particles:    true,
	}

	// check for params
	n := p.next()
	if n.Typ != itemLeftSquareParen {
		p.backup()
		return result, nil
	}

	// loop until we hit square paren
	for {
		// we're expecting ident = int
		i, err := p.consume(itemIdentifier)
		if err != nil {
			return result, err
		}

		item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
		if err != nil {
			return result, err
		}

		switch i.Val {
		case "hp_mult":
			result.HpMultiplier, err = itemNumberToFloat64(item)
			if err != nil {
				return result, err
			}
		case "particles":
			val, err := itemNumberToInt(item)
			if err != nil {
				return result, err
			}
			result.Particles = val != 0
		}

		// if we hit ], return; if we hit , keep going, other wise error
		n := p.next()
		switch n.Typ {
		case itemRightSquareParen:
			return result, nil
		case itemComma:
			// do nothing, keep going
		default:
			return result, fmt.Errorf("ln%v: <action param> bad token %v", n.line, n)
		}
	}
}
