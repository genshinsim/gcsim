package ast

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func parseCharacter(p *Parser) (parseFn, error) {
	// expecting one of:
	//	char lvl etc
	// 	add
	// should be any action here
	switch n := p.next(); n.Typ {
	case keywordChar:
		return parseCharDetails, nil
	case keywordAdd:
		return parseCharacterAdd, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character>: %v", n.line, n)
	}
}

func (p *Parser) newChar(key keys.Char) {
	r := info.CharacterProfile{}
	r.Base.Key = key
	r.Stats = make([]float64, attributes.EndStatType)
	r.StatsByLabel = make(map[string][]float64)
	r.Params = make(map[string]int)
	r.Sets = make(map[keys.Set]int)
	r.SetParams = make(map[keys.Set]map[string]int)
	r.Weapon.Params = make(map[string]int)
	r.Base.Element = keys.CharKeyToEle[key]
	p.chars[key] = &r
	p.charOrder = append(p.charOrder, key)
}

func parseCharDetails(p *Parser) (parseFn, error) {
	// xiangling c lvl=80/90 cons=4 talent=6,9,9;
	c := p.chars[p.currentCharKey]
	var err error
	var x Token
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case keywordLvl:
			c.Base.Level, c.Base.MaxLevel, err = p.acceptLevelReturnBaseMax()
			// err check below
		case keywordCons:
			x, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.Cons, err = itemNumberToInt(x)
			}
		case keywordTalent:
			x, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Attack, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}

			x, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Skill, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}

			x, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Burst, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}
		case ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param;", n.line)
				}
				p.backup()
				// overriding here if it already exists
				c.Params, err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.line, n)
			}
		case itemTerminateLine:
			return parseRows, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing character")
}

func parseCharacterAdd(p *Parser) (parseFn, error) {
	// after add we expect either weapon, set, or stats
	n := p.next()
	switch n.Typ {
	case keywordWeapon:
		return parseCharAddWeapon, nil
	case keywordSet:
		return parseCharAddSet, nil
	case keywordStats:
		return parseCharAddStats, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character> add: %v", n.line, n)
	}
}

func parseCharAddSet(p *Parser) (parseFn, error) {
	// xiangling add set="seal of insulation" count=4;
	c := p.chars[p.currentCharKey]
	var err error
	var x Token
	x, err = p.acceptSeqReturnLast(itemAssign, itemString)
	if err != nil {
		return nil, err
	}
	s := x.Val
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	label, ok := shortcut.SetNameToKey[s]
	if !ok {
		return nil, fmt.Errorf("invalid set %v", s)
	}
	count := 0

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case keywordCount:
			x, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				count, err = itemNumberToInt(x)
			}
		case ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.line)
				}
				p.backup()
				// overriding here if it already exists
				c.SetParams[label], err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.line, n)
			}
		case itemTerminateLine:
			c.Sets[label] = count
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unexpected token after in parsing sets: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add set")
}

func parseCharAddWeapon(p *Parser) (parseFn, error) {
	// weapon="string name" lvl=??/?? refine=xx
	c := p.chars[p.currentCharKey]
	var err error
	var x Token
	x, err = p.acceptSeqReturnLast(itemAssign, itemString)
	if err != nil {
		return nil, err
	}
	s := x.Val
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	c.Weapon.Key = shortcut.WeaponNameToKey[s]
	c.Weapon.Name = c.Weapon.Key.String()

	lvlOk := false
	refineOk := false

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case keywordLvl:
			c.Weapon.Level, c.Weapon.MaxLevel, err = p.acceptLevelReturnBaseMax()
			lvlOk = true
		case keywordRefine:
			x, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Weapon.Refine, err = itemNumberToInt(x)
				refineOk = true
			}
		case ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.line)
				}
				p.backup()
				// overriding here if it already exists
				c.Weapon.Params, err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.line, n)
			}
		case itemTerminateLine:
			if !lvlOk {
				return nil, fmt.Errorf("ln%v: weapon %v missing lvl", n.line, s)
			}
			if !refineOk {
				return nil, fmt.Errorf("ln%v: weapon %v missing refine", n.line, s)
			}
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add weapon: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add weapon")
}

func parseCharAddStats(p *Parser) (parseFn, error) {
	// xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
	c := p.chars[p.currentCharKey]

	// each line will be parsed separately into the map
	var line = make([]float64, attributes.EndStatType)
	var key string

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemStatKey:
			x, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(x)
			if err != nil {
				return nil, err
			}
			pos := statKeys[n.Val]
			line[pos] += amt
		case keywordLabel:
			x, err := p.acceptSeqReturnLast(itemAssign, itemIdentifier)
			if err != nil {
				return nil, err
			}
			key = x.Val
		case itemTerminateLine:
			// add stats into label
			m, ok := c.StatsByLabel[key]
			if !ok {
				m = make([]float64, attributes.EndStatType)
			}
			for i, v := range line {
				c.Stats[i] += v
				m[i] += v
			}
			c.StatsByLabel[key] = m
			return parseRows, nil
		case itemIdentifier:
			if n.Val == "random" {
				return parseCharAddRandomStats(p)
			}
			fallthrough
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats: %v", n.line, n)
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add stats")
}

func parseCharAddRandomStats(p *Parser) (parseFn, error) {
	// xiangling add stats random rarity=5 sand=hp% goblet=pyro% circlet=cr

	// note that plume/flower not specified and will be ignored
	rs := &info.RandomSubstats{
		Rarity: 5, // default to 5 star
	}

	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemTerminateLine:
			// check to make sure all values are valid
			err := rs.Validate()
			if err != nil {
				return nil, fmt.Errorf("ln%v: %w", n.line, err)
			}
			c := p.chars[p.currentCharKey]
			c.RandomSubstats = rs
			return parseRows, nil
		case itemIdentifier:
			switch n.Val {
			case "rarity":
				x, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				rs.Rarity, err = itemNumberToInt(x)
				if err != nil {
					return nil, err
				}
			case "sand":
				x, err := p.acceptSeqReturnLast(itemAssign, itemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Sand = statKeys[x.Val]
			case "goblet":
				x, err := p.acceptSeqReturnLast(itemAssign, itemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Goblet = statKeys[x.Val]
			case "circlet":
				x, err := p.acceptSeqReturnLast(itemAssign, itemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Circlet = statKeys[x.Val]
			default:
				return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats random: %v", n.line, n)
			}
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats random: %v", n.line, n)
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add stats (with random subs)")
}

func (p *Parser) acceptLevelReturnBaseMax() (int, int, error) {
	base := 0
	max := 0
	var err error
	// expect =xx/yy
	var x Token
	x, err = p.consume(itemAssign)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: unexpected token after lvl. expecting = got %v", x.line, x)
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: expecting a number for base lvl, got %v", x.line, x)
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: unexpected token for base lvl. got %v", x.line, x)
	}
	x, err = p.consume(ItemForwardSlash)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: expecting / separator for lvl, got %v", x.line, x)
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: expecting a number for max lvl, got %v", x.line, x)
	}
	max, err = itemNumberToInt(x)
	if err != nil {
		return base, max, fmt.Errorf("ln%v: unexpected token for lvl. got %v", x.line, x)
	}
	if max < base {
		return base, max, fmt.Errorf("ln%v: max level %v cannot be less than base level %v", x.line, max, base)
	}
	return base, max, nil
}
