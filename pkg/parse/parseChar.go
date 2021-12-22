package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (p *Parser) newChar(key keys.Char) {
	r := core.CharacterProfile{}
	r.Base.Key = key
	r.Stats = make([]float64, len(core.StatTypeString))
	r.Sets = make(map[string]int)
	r.Base.StartHP = -1
	p.chars[key] = &r
}

//return a char name
func (p *Parser) acceptChar() (*core.CharacterProfile, error) {
	n, err := p.consume(itemIdentifier)
	// log.Println(n)
	if err != nil {
		return nil, err
	}
	key, ok := keys.CharNameToKey[n.val]
	if !ok {
		return nil, fmt.Errorf("bad token at line %v - %v: %v; invalid char name", n.line, n.pos, n)
	}
	if _, ok := p.chars[key]; !ok {
		p.newChar(key)
	}
	return p.chars[key], nil
}

func parseChar(p *Parser) (parseFn, error) {
	var err error
	//char+=bennett ele=pyro lvl=70 hp=8352 atk=165 def=619 cr=0.05 cd=0.50 er=.2 cons=6 talent=1,8,8;
	c, err := p.acceptChar()
	if err != nil {
		return nil, err
	}

	var item item

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemStartHP:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.StartHP, err = itemNumberToFloat64(item)
			}
		case itemEle:
			_, err = p.consume(itemAssign)
			if err != nil {
				return nil, err
			}
			n := p.next()
			if n.typ <= eleTypeKeyword {
				return nil, fmt.Errorf("<char> expecting element, got bad token at line %v - %v: %v", n.line, n.pos, n)
			}
			c.Base.Element = eleKeys[n.val]
		case itemLvl:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.Level, err = itemNumberToInt(item)
			}
		case itemCons:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.Cons, err = itemNumberToInt(item)
			}
		case statHP:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.HP, err = itemNumberToFloat64(item)
			}
		case statATK:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.Atk, err = itemNumberToFloat64(item)
			}
		case statDEF:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Base.Def, err = itemNumberToFloat64(item)
			}
		case itemTalent:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Attack, err = itemNumberToInt(item)
			if err != nil {
				return nil, err
			}

			item, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Skill, err = itemNumberToInt(item)
			if err != nil {
				return nil, err
			}

			item, err = p.acceptSeqReturnLast(itemComma, itemNumber)
			if err != nil {
				return nil, err
			}
			c.Talents.Burst, err = itemNumberToInt(item)
			if err != nil {
				return nil, err
			}

		case itemTerminateLine:
			return parseRows, nil
		default:
			if n.typ > statKeyword && n.typ < eleTypeKeyword {
				s := n.val
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}

				pos := core.StrToStatType(s)
				c.Stats[pos] += amt
			} else {
				return nil, fmt.Errorf("<char> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing char")
}

func parseStats(p *Parser) (parseFn, error) {

	//stats+=bennett label=flower hp=4780 def=44 er=.065 cr=.097 cd=.124;
	c, err := p.acceptChar()
	if err != nil {
		return nil, err
	}

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch {
		case n.typ == itemLabel:
			_, err = p.acceptSeqReturnLast(itemAssign, itemIdentifier)
			if err != nil {
				return nil, err
			}
		case n.typ == itemTerminateLine:
			return parseRows, nil
		case n.typ > statKeyword && n.typ < eleTypeKeyword:
			s := n.val
			item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			amt, err := itemNumberToFloat64(item)
			if err != nil {
				return nil, err
			}

			pos := core.StrToStatType(s)
			c.Stats[pos] += amt
		default:
			return nil, fmt.Errorf("<stats> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
	}
	return nil, errors.New("unexpected end of line while parsing stats")
}

func parseWeapon(p *Parser) (parseFn, error) {
	var err error
	//weapon+=bennett label="festering desire" atk=401 er=0.559 refine=3;
	c, err := p.acceptChar()
	if err != nil {
		return nil, err
	}
	var item item

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemLabel:
			n, err = p.acceptSeqReturnLast(itemAssign, itemString)
			if err != nil {
				return nil, err
			}
			s := n.val
			if len(s) > 0 && s[0] == '"' {
				s = s[1:]
			}
			if len(s) > 0 && s[len(s)-1] == '"' {
				s = s[:len(s)-1]
			}
			c.Weapon.Name = s
		case itemRefine:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Weapon.Refine, err = itemNumberToInt(item)
			}
		case statATK:
			item, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				c.Weapon.Atk, err = itemNumberToFloat64(item)
			}
		case itemParam:
			param, err := p.parseWeaponParams()
			if err != nil {
				return nil, err
			}
			c.Weapon.Param = param
		case itemTerminateLine:
			return parseRows, nil
		default:
			if n.typ > statKeyword && n.typ < eleTypeKeyword {
				s := n.val
				item, err := p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err != nil {
					return nil, err
				}
				amt, err := itemNumberToFloat64(item)
				if err != nil {
					return nil, err
				}

				pos := core.StrToStatType(s)
				c.Stats[pos] += amt
			} else {
				return nil, fmt.Errorf("<weapon> bad token at line %v - %v: %v", n.line, n.pos, n)
			}
		}

		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing weapon")
}

func (p *Parser) parseWeaponParams() (map[string]int, error) {
	param := make(map[string]int)
	_, err := p.consume(itemAssign)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(itemLeftSquareParen)
	if err != nil {
		return nil, err
	}

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemIdentifier:
			s := n.val
			//should be ident=123
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err != nil {
				return nil, err
			}
			val, err := itemNumberToInt(n)
			if err != nil {
				return nil, err
			}
			param[s] = val
		case itemComma:
			//ignore and keep going
		case itemRightSquareParen:
			return param, nil
		default:
			return nil, fmt.Errorf("<weapon param> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
	}

	return nil, errors.New("unexpected end of line while parsing weapon param")
}

func parseArtifacts(p *Parser) (parseFn, error) {
	var err error
	//art+=xiangling label="gladiator's finale" count=2;
	c, err := p.acceptChar()
	if err != nil {
		return nil, err
	}
	var label string
	var count int

	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemLabel:
			n, err = p.acceptSeqReturnLast(itemAssign, itemString)
			if err != nil {
				return nil, err
			}
			s := n.val
			if len(s) > 0 && s[0] == '"' {
				s = s[1:]
			}
			if len(s) > 0 && s[len(s)-1] == '"' {
				s = s[:len(s)-1]
			}
			label = s
		case itemCount:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				count, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			c.Sets[label] = count
			return parseRows, nil
		default:
			return nil, fmt.Errorf("<art> bad token at line %v - %v: %v", n.line, n.pos, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing artifacts")
}
