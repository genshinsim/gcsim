package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func parseChar(p *Parser) (parseFn, error) {

	n := p.next()
	switch n.typ {
	case itemChar:
		return parseCharDetails, nil
	case itemAdd:
		return parseCharAdd, nil
	case itemActionKey:
		return parseCharActions, nil
	default:
		return nil, fmt.Errorf("unexpected token after <character>: %v at line %v", n, p.tokens)
	}
}

func (p *Parser) newChar(key keys.Char) {
	r := core.CharacterProfile{}
	r.Base.Key = key
	r.Stats = make([]float64, len(core.StatTypeString))
	r.Sets = make(map[string]int)
	r.Base.StartHP = -1
	p.chars[key] = &r
}

func parseCharDetails(p *Parser) (parseFn, error) {
	//xiangling c lvl=80/90 cons=4 talent=6,9,9;
	c := p.chars[p.currentCharKey]
	var err error
	var x item
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemLvl:
			c.Base.Level, c.Base.MaxLevel, err = p.acceptLevelReturnBaseMax()
			//err check below
		case itemCons:
			x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				c.Base.Cons, err = itemNumberToInt(x)
			}
		case itemTalent:
			x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
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
		case itemTerminateLine:
			return parseRows, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("unexpected end of line while parsing character")
}

func parseCharAdd(p *Parser) (parseFn, error) {
	//after add we expect either weapon, set, or stats
	n := p.next()
	switch n.typ {
	case itemWeapon:
		return parseCharAddWeapon, nil
	case itemSet:
		return parseCharAddSet, nil
	case itemStats:
		return parseCharAddStats, nil
	default:
		return nil, fmt.Errorf("unexpected token after <character> add: %v at line %v", n, p.tokens)
	}
}

func (p *Parser) acceptLevelReturnBaseMax() (base, max int, err error) {
	//expect =xx/yy
	var x item
	x, err = p.consume(itemEqual)
	if err != nil {
		err = fmt.Errorf("unexpected token after lvl. expecting = got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("expecting a number for base lvl, got %v at line %v", x, p.tokens)
		return
	}
	base, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("unexpected token for base lvl. got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemForwardSlash)
	if err != nil {
		err = fmt.Errorf("expecting / separator for lvl, got %v at line %v", x, p.tokens)
		return
	}
	x, err = p.consume(itemNumber)
	if err != nil {
		err = fmt.Errorf("expecting a number for max lvl, got %v at line %v", x, p.tokens)
		return
	}
	max, err = itemNumberToInt(x)
	if err != nil {
		err = fmt.Errorf("unexpected token for lvl. got %v at line %v", x, p.tokens)
		return
	}
	if max < base {
		err = fmt.Errorf("max level %v cannot be less than base level %v at line %v", max, base, p.tokens)
		return
	}
	return
}

func parseCharAddSet(p *Parser) (parseFn, error) {
	//xiangling add set="seal of insulation" count=4;
	c := p.chars[p.currentCharKey]
	var err error
	var x item
	x, err = p.acceptSeqReturnLast(itemEqual, itemString)
	if err != nil {
		return nil, err
	}
	s := x.val
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	label := s
	count := 0

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemCount:
			x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				count, err = itemNumberToInt(x)
			}
		case itemTerminateLine:
			c.Sets[label] = count
			return parseRows, nil
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
	var x item
	x, err = p.acceptSeqReturnLast(itemEqual, itemString)
	if err != nil {
		return nil, err
	}
	s := x.val
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	c.Weapon.Name = s

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemLvl:
			c.Weapon.Level, c.Weapon.MaxLevel, err = p.acceptLevelReturnBaseMax()
		case itemRefine:
			x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				c.Weapon.Refine, err = itemNumberToInt(x)
			}
		case itemWeapon:

		case itemTerminateLine:
			return parseRows, nil
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

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemStatKey:
			x, err := p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(x)
			if err != nil {
				return nil, err
			}
			pos := statKeys[n.val]
			c.Stats[pos] += amt
		case itemTerminateLine:
			return parseRows, nil
		}
	}
	return nil, errors.New("unexpected end of line while parsing character add set")
}
