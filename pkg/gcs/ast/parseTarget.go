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
			case "resist_frozen":
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				v, err := itemNumberToInt(item)
				if err != nil {
					return nil, err
				}
				r.ResistFrozen = v != 0
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

			// TODO: make this more elegant...
			r.Resist[attributes.Electro] += amt
			r.Resist[attributes.Cryo] += amt
			r.Resist[attributes.Hydro] += amt
			r.Resist[attributes.Physical] += amt
			r.Resist[attributes.Pyro] += amt
			r.Resist[attributes.Geo] += amt
			r.Resist[attributes.Dendro] += amt
			r.Resist[attributes.Anemo] += amt
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
			r.ParticleDrops = nil // separate particle system
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
		case "mult", "hp_multiplier":
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
