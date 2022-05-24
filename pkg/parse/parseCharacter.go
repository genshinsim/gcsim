package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func parseCharacter(p *Parser) (parseFn, error) {
	//expecting one of:
	//	char lvl etc
	// 	add
	//should be any action here
	switch n := p.next(); n.typ {
	case keywordChar:
	case keywordAdd:
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after <character>: %v", n.line, n)
	}
	return nil, nil
}

func (p *Parser) newChar(key keys.Char) {
	r := character.CharacterProfile{}
	r.Base.Key = key
	r.Stats = make([]float64, attributes.EndStatType)
	r.StatsByLabel = make(map[string][]float64)
	r.Params = make(map[string]int)
	r.Sets = make(map[keys.Set]int)
	r.SetParams = make(map[keys.Set]map[string]int)
	r.Weapon.Params = make(map[string]int)
	r.Base.StartHP = -1
	r.Base.Element = keys.CharKeyToEle[key]
	p.chars[key] = &r
	p.charOrder = append(p.charOrder, key)
}

func parseCharDetails(p *Parser) (parseFn, error) {
	//xiangling c lvl=80/90 cons=4 talent=6,9,9;
	c := p.chars[p.currentCharKey]
	var err error
	var x Token
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
		case itemStartHP:
			x, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				c.Base.StartHP, err = itemNumberToFloat64(x)
			}
		case itemPlus: //optional flags
			n = p.next()
			switch n.typ {
			case itemParams:
				//expecting =[
				_, err = p.acceptSeqReturnLast(itemEqual, itemLeftSquareParen)
				if err != nil {
					return nil, fmt.Errorf("ln%v: invalid token after param; line %v", n.line, p.tokens)
				}
				p.backup()
				//overriding here if it already exists
				c.Params, err = p.acceptOptionalParamReturnMap()
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
	return nil, nil
}
