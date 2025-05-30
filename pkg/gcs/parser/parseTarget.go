package parser

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func parseTarget(p *Parser) (parseFn, error) {
	var err error
	var r info.EnemyProfile
	r.Resist = make(map[attributes.Element]float64)
	r.ParticleElement = attributes.NoElement
	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemIdentifier:
			switch n.Val {
			case "pos": // pos will end up defaulting to 0,0 if not set
				// pos=1.00,2,00
				if _, err := p.consume(ast.ItemAssign); err != nil {
					return nil, err
				}
				x, err := p.parseFloat64Const()
				if err != nil {
					return nil, err
				}
				if _, err := p.consume(ast.ItemComma); err != nil {
					return nil, err
				}
				y, err := p.parseFloat64Const()
				if err != nil {
					return nil, err
				}
				r.Pos.X = x
				r.Pos.Y = y
			case "radius":
				item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}
				r.Pos.R = amt
			case "type":
				item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemIdentifier)
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
				if _, err := p.consume(ast.ItemAssign); err != nil {
					return nil, err
				}
				v, err := p.parseFloat64Const()
				if err != nil {
					return nil, err
				}
				r.FreezeResist = v
				r.Modified = true
			default:
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
			}
		case ast.KeywordLvl:
			n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				r.Level, err = itemNumberToInt(n)
			}
		case ast.ItemStatKey:
			// should be hp
			if ast.StatKeys[n.Val] != attributes.HP {
				return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
			}
			n, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				r.HP, err = itemNumberToFloat64(n)
				if err != nil {
					return nil, err
				}
				p.res.Settings.DamageMode = true
				r.Modified = true
			}
		case ast.KeywordResist:
			// this sets all resistance
			if _, err := p.consume(ast.ItemAssign); err != nil {
				return nil, err
			}
			amt, err := p.parseFloat64Const()
			if err != nil {
				return nil, err
			}

			res := []attributes.Element{attributes.Electro, attributes.Cryo, attributes.Hydro, attributes.Physical, attributes.Pyro, attributes.Geo, attributes.Dendro, attributes.Anemo}
			for _, attr := range res {
				r.Resist[attr] += amt
			}
			r.Modified = true
		case ast.KeywordParticleThreshold:
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropThreshold = amt
			r.ParticleDrops = nil // separate particle system
			r.ParticleElement = attributes.NoElement
			r.Modified = true
		case ast.KeywordParticleDropCount:
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}
			r.ParticleDropCount = amt
			r.Modified = true
		case ast.KeywordParticleElement:
			item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemElementKey)
			if err != nil {
				return nil, err
			}
			if ele, ok := ast.EleKeys[item.Val]; ok {
				r.ParticleElement = ele
			}
			r.Modified = true
		case ast.ItemElementKey:
			s := n.Val
			if _, err := p.consume(ast.ItemAssign); err != nil {
				return nil, err
			}
			amt, err := p.parseFloat64Const()
			if err != nil {
				return nil, err
			}

			r.Resist[ast.EleKeys[s]] += amt
			r.Modified = true
		case ast.ItemTerminateLine:
			p.res.Targets = append(p.res.Targets, r)
			return parseRows, nil
		default:
			return nil, fmt.Errorf("<target> bad token at line %v - %v: %v", n.Line, n.Pos, n)
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
	if n.Typ != ast.ItemLeftSquareParen {
		p.backup()
		return result, nil
	}

	// loop until we hit square paren
	for {
		// we're expecting ident = int
		i, err := p.consume(ast.ItemIdentifier)
		if err != nil {
			return result, err
		}

		item, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
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
		case ast.ItemRightSquareParen:
			return result, nil
		case ast.ItemComma:
			// do nothing, keep going
		default:
			return result, fmt.Errorf("ln%v: <action param> bad token %v", n.Line, n)
		}
	}
}
