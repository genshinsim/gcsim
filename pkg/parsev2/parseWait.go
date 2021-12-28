package parse

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseWait(p *Parser) (parseFn, error) {
	// wait_for particles value=xingqiu max=100
	// wait_for mods value=.xiangling.bennettbuff==1 max=-1
	// wait_for time max=100
	w := core.ActionBlock{}
	w.Wait.Max = -2

	n, err := p.consume(itemIdentifier)
	if err != nil {
		return nil, fmt.Errorf("<parse wait> invalid tokens after wait_for (expecting identifier) at line %v", p.tokens)
	}
	switch n.val {
	case "particles":
		w.Wait.For = core.CmdWaitTypeParticle
	case "mods":
		w.Wait.For = core.CmdWaitTypeMods
	case "time":
		w.Wait.For = core.CmdWaitTypeTimed
	default:
		return nil, fmt.Errorf("<parse wait> invalid option %v after wait_for at line %v", n.val, p.tokens)
	}

	condOk := false

	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemValue:
			//value is either a string or fields
			//if mods expect fields, if particles expect string
			n, err = p.consume(itemEqual)
			if err != nil {
				return nil, fmt.Errorf("<parse wait> expected = after value, got %v. line %v", n, p.tokens)
			}

			//next should be either an identifier or a field
			n := p.next()

			switch n.typ {
			case itemIdentifier:
				//only valid if particles
				if w.Wait.For != core.CmdWaitTypeParticle {
					return nil, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
				}
				w.Wait.Source = n.val
			case itemField:
				//only valid if mods
				if w.Wait.For != core.CmdWaitTypeMods {
					return nil, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
				}
				//backup and parse condition
				p.backup()
				//scan for fields
				c, err := p.parseCondition()
				if err != nil {
					return nil, err
				}
				w.Wait.Conditions = c
				condOk = true
			default:
				return nil, fmt.Errorf("<parse wait> invalid value %v, line %v", n, p.tokens)
			}
		case itemMax:
			//this has to be after for
			if w.Wait.For == core.CmdWaitTypeInvalid {
				return nil, fmt.Errorf("<parse wait> missing for. val must be after for. line %v", p.tokens)
			}
			n, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
			if err == nil {
				w.Wait.Max, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			//make sure mandatory fields are present
			if w.Wait.Max == -2 {
				return nil, fmt.Errorf("<parse wait> missing max")
			}
			if !condOk {
				return nil, fmt.Errorf("<parse wait> missing condition")
			}
			return parseRows, nil
		default:
			err = fmt.Errorf("<parse wait> unrecognized token %v at line %v", n, p.tokens)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing wait commmand")

}

//bennett skill,attack,burst,attack +if=.xx.xx.xx>1 +swap=xiangling +if_onfield +limit=1 +try=1 +timeout=100 swap_to swap_lock
