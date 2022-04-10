package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseCalcModeWait(p *Parser) (parseFn, error) {
	//wait 10
	//wait until 600
	block := core.ActionBlock{
		Type: core.ActionBlockTypeCalcWait,
	}
	n := p.next()
	switch n.typ {
	case itemNumber:
		val, err := itemNumberToInt(n)
		if err != nil {
			return nil, err
		}
		block.CalcWait.Val = val
	case itemUntil:
		//should be number next
		n = p.next()
		if n.typ != itemNumber {
			return nil, fmt.Errorf("ln%v: unexpected token after wait, expecting a number got %v line %v", n.line, n, p.tokens)
		}
		val, err := itemNumberToInt(n)
		if err != nil {
			return nil, err
		}
		block.CalcWait.Frames = true
		block.CalcWait.Val = val
	default:
		return nil, fmt.Errorf("ln%v: unexpected token after wait, got %v line %v", n.line, n, p.tokens)
	}
	//expect end of line
	n = p.next()
	if n.typ != itemTerminateLine {
		return nil, fmt.Errorf("ln%v: wait expecting ; got %v", n.line, n)
	}
	p.cfg.Rotation = append(p.cfg.Rotation, block)
	return parseRows, nil
}

func parseWait(p *Parser) (parseFn, error) {
	block, err := p.acceptWait()
	if err != nil {
		return nil, err
	}
	p.cfg.Rotation = append(p.cfg.Rotation, block)

	return parseRows, nil
}

func (p *Parser) acceptWait() (core.ActionBlock, error) {
	// wait_for particles value=xingqiu max=100
	// wait_for mods value=.xiangling.bennettbuff==1 max=-1
	// wait_for time max=100
	w := core.ActionBlock{
		Type: core.ActionBlockTypeWait,
	}
	w.Wait.Max = -2

	condOk := false

	n, err := p.consume(itemIdentifier)
	if err != nil {
		return w, fmt.Errorf("<parse wait> invalid tokens after wait_for (expecting identifier) at line %v", p.tokens)
	}
	switch n.val {
	case "particles":
		w.Wait.For = core.CmdWaitTypeParticle
	case "mods":
		w.Wait.For = core.CmdWaitTypeMods
	case "time":
		w.Wait.For = core.CmdWaitTypeTimed
		condOk = true
	default:
		return w, fmt.Errorf("ln%v: <parse wait> invalid option %v after wait_for", n.line, n.val)
	}

	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemValue:
			//value is either a string or fields
			//if mods expect fields, if particles expect string
			n, err = p.consume(itemEqual)
			if err != nil {
				return w, fmt.Errorf("<parse wait> expected = after value, got %v. line %v", n, p.tokens)
			}

			//next should be either an identifier or a field
			n := p.next()

			switch n.typ {
			case itemIdentifier, itemCharacterKey:
				//only valid if particles
				if w.Wait.For != core.CmdWaitTypeParticle {
					return w, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
				}
				w.Wait.Source = n.val
				condOk = true
			case itemField:
				//only valid if mods
				if w.Wait.For != core.CmdWaitTypeMods {
					return w, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
				}
				//backup and parse condition
				p.backup()
				//scan for fields
				c, err := p.parseCondition()
				if err != nil {
					return w, err
				}
				w.Wait.Conditions = c
				condOk = true
			default:
				return w, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
			}
		case itemMax:
			//this has to be after for
			if w.Wait.For == core.CmdWaitTypeInvalid {
				return w, fmt.Errorf("<parse wait> missing for. val must be after for. line %v", p.tokens)
			}
			n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
			if err == nil {
				w.Wait.Max, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			//make sure mandatory fields are present
			if w.Wait.Max == -2 {
				return w, fmt.Errorf("<parse wait> missing max, line %v", p.tokens)
			}
			if !condOk {
				return w, fmt.Errorf("<parse wait> missing condition, line %v", p.tokens)
			}
			return w, nil
		case itemPlus:
			//+filler=skill[param=1]
			n = p.next()
			switch n.typ {
			case itemFiller:
				//consume an equal sign
				n, err = p.acceptSeqReturnLast(itemEqual, itemActionKey)
				if err != nil {
					return w, fmt.Errorf("<parse wait> unrecognized token after filler at line %v", p.tokens)
				}
				//expecting a
				act := core.ActionItem{
					Typ:    actionKeys[n.val],
					Target: core.NoChar, //since it's active char only
				}
				//optional params
				act.Param, err = p.acceptOptionalParamReturnMap()
				if err != nil {
					return w, err
				}
				w.Wait.FillAction = act
			default:
				return w, fmt.Errorf("ln%v: <parse wait> unrecognized token %v after +", n.line, n)
			}
		default:
			err = fmt.Errorf("ln%v: <parse wait> unrecognized token %v", n.line, n)
		}
		if err != nil {
			return w, err
		}
	}

	return w, errors.New("unexpected end of line while parsing wait commmand")

}

//bennett skill,attack,burst,attack +if=.xx.xx.xx>1 +swap=xiangling +if_onfield +limit=1 +try=1 +timeout=100 swap_to swap_lock
