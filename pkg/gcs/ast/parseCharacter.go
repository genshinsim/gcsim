package ast

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func parseCharacter(p *Parser) (parseFn, error) {
	//expecting one of:
	//	char lvl etc
	// 	add
	//should be any action here
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
	r := profile.CharacterProfile{}
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
	//xiangling c lvl=80/90 cons=4 talent=6,9,9;
	c := p.chars[p.currentCharKey]
	var err error
	var x Token
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case keywordLvl:
			c.Base.Level, c.Base.MaxLevel, err = p.acceptLevelReturnBaseMax()
			//err check below
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
		case ItemPlus: //optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				//expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param;", n.line)
				}
				p.backup()
				//overriding here if it already exists
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
	//after add we expect either weapon, set, or stats
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
	//xiangling add set="seal of insulation" count=4;
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
		case ItemPlus: //optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				//expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.line)
				}
				p.backup()
				//overriding here if it already exists
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
	//weapon="string name" lvl=??/?? refine=xx
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
		case ItemPlus: //optional flags
			n = p.next()
			switch n.Typ {
			case keywordParams:
				//expecting =[
				_, err = p.acceptSeqReturnLast(itemAssign, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.line)
				}
				p.backup()
				//overriding here if it already exists
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
	//xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
	c := p.chars[p.currentCharKey]

	//each line will be parsed separately into the map
	var line = make([]float64, attributes.EndStatType)
	var key string
	useRolls := false
	rollOpt := "avg"
	rarity := 5

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
		case itemIdentifier:
			switch n.Val {
			case "roll":
				x, err := p.acceptSeqReturnLast(itemAssign, itemIdentifier)
				if err != nil {
					return nil, err
				}
				//should be min, max, avg
				switch x.Val {
				case "avg", "min", "max":
					useRolls = true
					rollOpt = x.Val
				default:
					return nil, fmt.Errorf("ln%v: invalid roll option: %v", n.line, x.Val)
				}
			case "rarity":
				x, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToInt(x)
				if err != nil {
					return nil, err
				}
				if amt > 5 {
					amt = 5
				}
				if amt < 1 {
					amt = 1
				}
				rarity = amt
			default:
				return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats: %v", n.line, n)
			}
		case itemTerminateLine:
			//add stats into label
			m, ok := c.StatsByLabel[key]
			if !ok {
				m = make([]float64, attributes.EndStatType)
			}
			for i, v := range line {
				if useRolls {
					c.Stats[i] += v * rolls[rarity-1][rollOpt][i]
					m[i] += v * rolls[rarity-1][rollOpt][i]
				} else {
					c.Stats[i] += v
					m[i] += v
				}
			}
			c.StatsByLabel[key] = m
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats: %v", n.line, n)
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add stats")
}

var rolls = []map[string][]float64{
	{
		"min": {0, 0.0146, 1.85, 23.9, 0.0117, 1.56, 0.0117, 0.013, 4.66, 0.0078, 0.0155, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"avg": {0, 0.0164, 2.08, 26.89, 0.01315, 1.755, 0.01315, 0.0146, 5.245, 0.00875, 0.017450001, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"max": {0, 0.0182, 2.31, 29.88, 0.0146, 1.95, 0.0146, 0.0162, 5.83, 0.0097, 0.0194, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		"min": {0, 0.0204, 3.89, 50.19, 0.0163, 3.27, 0.0163, 0.0181, 6.53, 0.0109, 0.0218, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"avg": {0, 0.024766669, 4.7233334, 60.946667, 0.0198, 3.97, 0.0198, 0.022, 7.9300003, 0.0132, 0.026433334, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"max": {0, 0.0291, 5.56, 71.7, 0.0233, 4.67, 0.0233, 0.0259, 9.33, 0.0155, 0.0311, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		"min": {0, 0.0306, 7.78, 100.38, 0.0245, 6.54, 0.0245, 0.0272, 9.79, 0.0163, 0.0326, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"avg": {0, 0.03715, 9.445, 121.89, 0.02975, 7.9375, 0.02975, 0.03305, 11.889999, 0.0198, 0.039625, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"max": {0, 0.0437, 11.11, 143.4, 0.035, 9.34, 0.035, 0.0389, 13.99, 0.0233, 0.0466, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		"min": {0, 0.0408, 12.96, 167.3, 0.0326, 10.89, 0.0326, 0.0363, 13.06, 0.0218, 0.0435, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"avg": {0, 0.04955, 15.742499, 203.15, 0.039625, 13.225, 0.039625, 0.044025, 15.855, 0.026449999, 0.052849997, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"max": {0, 0.0583, 18.52, 239, 0.0466, 15.56, 0.0466, 0.0518, 18.65, 0.0311, 0.0622, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		"min": {0, 0.051, 16.2, 209.13, 0.0408, 13.62, 0.0408, 0.0453, 16.32, 0.0272, 0.0544, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"avg": {0, 0.06195, 19.675001, 253.94, 0.04955, 16.535, 0.04955, 0.05505, 19.815, 0.03305, 0.06605, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"max": {0, 0.0729, 23.15, 298.75, 0.0583, 19.45, 0.0583, 0.0648, 23.31, 0.0389, 0.0777, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}

func (p *Parser) acceptLevelReturnBaseMax() (base, max int, err error) {
	//expect =xx/yy
	var x Token
	x, err = p.consume(itemAssign)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token after lvl. expecting = got %v", x.line, x)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting a number for base lvl, got %v", x.line, x)
		return
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token for base lvl. got %v", x.line, x)
		return
	}
	x, err = p.consume(ItemForwardSlash)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting / separator for lvl, got %v", x.line, x)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("ln%v: expecting a number for max lvl, got %v", x.line, x)
		return
	}
	max, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("ln%v: unexpected token for lvl. got %v", x.line, x)
		return
	}
	if max < base {
		err = fmt.Errorf("ln%v: max level %v cannot be less than base level %v", x.line, max, base)
		return
	}
	return
}
