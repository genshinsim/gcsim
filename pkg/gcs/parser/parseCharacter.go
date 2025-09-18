package parser

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func parseCharacter(p *Parser) (parseFn, error) {
	// expecting one of:
	//	char lvl etc
	// 	add
	// should be any action here
	switch n := p.next(); n.Typ {
	case ast.KeywordChar:
		return parseCharDetails, nil
	case ast.KeywordAdd:
		return parseCharacterAdd, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character>: %v", n.Line, n)
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
	var x ast.Token
	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.KeywordLvl:
			c.Base.Level, c.Base.MaxLevel, err = p.acceptLevelReturnBaseMax()
			// err check below
		case ast.KeywordCons:
			x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				c.Base.Cons, err = itemNumberToInt(x)
			}
		case ast.KeywordTalent:
			x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Attack, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}

			x, err = p.acceptSeqReturnLast(ast.ItemComma, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Skill, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}

			x, err = p.acceptSeqReturnLast(ast.ItemComma, ast.ItemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Burst, err = itemNumberToInt(x)
			if err != nil {
				return nil, err
			}
		case ast.ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case ast.KeywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param;", n.Line)
				}
				p.backup()
				// overriding here if it already exists
				c.Params, err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.Line, n)
			}
		case ast.ItemTerminateLine:
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
	case ast.KeywordWeapon:
		return parseCharAddWeapon, nil
	case ast.KeywordSet:
		return parseCharAddSet, nil
	case ast.KeywordStats:
		return parseCharAddStats, nil
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character> add: %v", n.Line, n)
	}
}

func parseCharAddSet(p *Parser) (parseFn, error) {
	// xiangling add set="seal of insulation" count=4;
	c := p.chars[p.currentCharKey]
	var err error
	var x ast.Token
	x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemString)
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

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.KeywordCount:
			x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				count, err = itemNumberToInt(x)
			}
		case ast.ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case ast.KeywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.Line)
				}
				p.backup()
				// overriding here if it already exists
				c.SetParams[label], err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.Line, n)
			}
		case ast.ItemTerminateLine:
			c.Sets[label] = count
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unexpected token after in parsing sets: %v", n.Line, n)
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
	var x ast.Token
	x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemString)
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

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.KeywordLvl:
			c.Weapon.Level, c.Weapon.MaxLevel, err = p.acceptLevelReturnBaseMax()
			lvlOk = true
		case ast.KeywordRefine:
			x, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
			if err == nil {
				c.Weapon.Refine, err = itemNumberToInt(x)
				refineOk = true
			}
		case ast.ItemPlus: // optional flags
			n = p.next()
			switch n.Typ {
			case ast.KeywordParams:
				// expecting =[
				_, err = p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param", n.Line)
				}
				p.backup()
				// overriding here if it already exists
				c.Weapon.Params, err = p.acceptOptionalParamReturnOnlyIntMap()
			default:
				err = fmt.Errorf("ln%v: unexpected token after +: %v", n.Line, n)
			}
		case ast.ItemTerminateLine:
			if !lvlOk {
				return nil, fmt.Errorf("ln%v: weapon %v missing lvl", n.Line, s)
			}
			if !refineOk {
				return nil, fmt.Errorf("ln%v: weapon %v missing refine", n.Line, s)
			}
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add weapon: %v", n.Line, n)
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
	line := make([]float64, attributes.EndStatType)
	var key string

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemStatKey:
			if _, err := p.consume(ast.ItemAssign); err != nil {
				return nil, err
			}
			amt, err := p.parseFloat64Const()
			if err != nil {
				return nil, err
			}
			// TODO: use attributes.StrToStatType?
			pos := ast.StatKeys[n.Val]
			line[pos] += amt
		case ast.KeywordLabel:
			x, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemIdentifier)
			if err != nil {
				return nil, err
			}
			key = x.Val
		case ast.ItemTerminateLine:
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
		case ast.ItemIdentifier:
			if n.Val == "random" {
				return parseCharAddRandomStats(p)
			}
			fallthrough
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats: %v", n.Line, n)
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

	for n := p.next(); n.Typ != ast.ItemEOF; n = p.next() {
		switch n.Typ {
		case ast.ItemTerminateLine:
			// check to make sure all values are valid
			err := rs.Validate()
			if err != nil {
				return nil, fmt.Errorf("ln%v: %w", n.Line, err)
			}
			c := p.chars[p.currentCharKey]
			c.RandomSubstats = rs
			return parseRows, nil
		case ast.ItemIdentifier:
			switch n.Val {
			case "rarity":
				x, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemNumber)
				if err != nil {
					return nil, err
				}
				rs.Rarity, err = itemNumberToInt(x)
				if err != nil {
					return nil, err
				}
			case "sand":
				x, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Sand = ast.StatKeys[x.Val]
			case "goblet":
				x, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Goblet = ast.StatKeys[x.Val]
			case "circlet":
				x, err := p.acceptSeqReturnLast(ast.ItemAssign, ast.ItemStatKey)
				if err != nil {
					return nil, err
				}
				rs.Circlet = ast.StatKeys[x.Val]
			default:
				return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats random: %v", n.Line, n)
			}
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing add stats random: %v", n.Line, n)
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add stats (with random subs)")
}

func (p *Parser) acceptLevelReturnBaseMax() (int, int, error) {
	base := 0
	maxlvl := 0
	var err error
	// expect =xx/yy
	var x ast.Token
	x, err = p.consume(ast.ItemAssign)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: unexpected token after lvl. expecting = got %v", x.Line, x)
	}
	x, err = p.consume(ast.ItemNumber)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: expecting a number for base lvl, got %v", x.Line, x)
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: unexpected token for base lvl. got %v", x.Line, x)
	}
	x, err = p.consume(ast.ItemForwardSlash)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: expecting / separator for lvl, got %v", x.Line, x)
	}
	x, err = p.consume(ast.ItemNumber)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: expecting a number for max lvl, got %v", x.Line, x)
	}
	maxlvl, err = itemNumberToInt(x)
	if err != nil {
		return base, maxlvl, fmt.Errorf("ln%v: unexpected token for lvl. got %v", x.Line, x)
	}
	if maxlvl < base {
		return base, maxlvl, fmt.Errorf("ln%v: max level %v cannot be less than base level %v", x.Line, maxlvl, base)
	}
	return base, maxlvl, nil
}
