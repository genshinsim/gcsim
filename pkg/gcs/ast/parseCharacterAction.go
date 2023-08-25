package ast

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

// parseAction returns a node contain a character action, or a block of node containing
// a list of character actions
func (p *Parser) parseAction() (Stmt, error) {
	char, err := p.consume(itemCharacterKey)
	if err != nil {
		//this really shouldn't happen since we already checked
		return nil, fmt.Errorf("ln%v: expecting character key, got %v", char.line, char.Val)
	}
	charKey := shortcut.CharNameToKey[char.Val]

	//should be multiple action keys next
	var actions []*CallExpr
	if n := p.peek(); n.Typ != itemActionKey {
		return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.line, char.Val, n.Val)
	}

	//all actions needs to come before any + flags
Loop:
	for {
		switch n := p.next(); n.Typ {
		case itemTerminateLine:
			//stop here
			break Loop
		case itemActionKey:
			actionKey := action.StringToAction(n.Val)
			expr := &CallExpr{
				Pos: char.pos,
				Fun: &Ident{
					Pos:   n.pos,
					Value: "execute_action",
				},
				Args: make([]Expr, 0),
			}
			expr.Args = append(expr.Args,
				// char
				&NumberLit{
					Pos:      char.pos,
					IntVal:   int64(charKey),
					FloatVal: float64(charKey),
				},
				// action
				&NumberLit{
					Pos:      n.pos,
					IntVal:   int64(actionKey),
					FloatVal: float64(actionKey),
				},
			)
			//check for param -> then repeat
			param, err := p.acceptOptionalParamReturnMap()
			if err != nil {
				return nil, err
			}
			if param != nil {
				expr.Args = append(expr.Args, param)
			}

			//optional : and a number
			count, err := p.acceptOptionalRepeaterReturnCount()
			if err != nil {
				return nil, err
			}
			//add to array
			for i := 0; i < count; i++ {
				//TODO: all the repeated action will access the same map
				//ability implement should avoid modifying the maps
				actions = append(actions, expr)
			}

			n = p.next()
			if n.Typ != itemComma {
				p.backup()
				break Loop
			}
		default:
			//TODO: fix invalid key error
			return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.line, char.Val, n.Val)
		}
	}
	//check for optional flags

	//build stmt
	b := newBlockStmt(char.pos)
	for _, v := range actions {
		b.append(v)
	}
	return b, nil
}

func (p *Parser) acceptOptionalParamReturnMap() (Expr, error) {
	//check for params
	n := p.peek()
	if n.Typ != itemLeftSquareParen {
		return nil, nil
	}

	return p.parseMap()
}

func (p *Parser) acceptOptionalParamReturnOnlyIntMap() (map[string]int, error) {
	r := make(map[string]int)

	result, err := p.acceptOptionalParamReturnMap()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return r, nil
	}

	for k, v := range result.(*MapExpr).Fields {
		switch v.(type) {
		case *NumberLit:
			// skip
		default:
			return nil, fmt.Errorf("expected number in the map, got %v", v.String())
		}
		r[k] = int(v.(*NumberLit).IntVal)
	}
	return r, nil
}

func (p *Parser) acceptOptionalRepeaterReturnCount() (int, error) {
	count := 1
	n := p.next()
	if n.Typ != itemColon {
		p.backup()
		return count, nil
	}
	//should be a number next
	n = p.next()
	if n.Typ != itemNumber {
		return count, fmt.Errorf("ln%v: expected a number after : but got %v", n.line, n)
	}
	//parse number
	count, err := itemNumberToInt(n)
	return count, err
}
