package parser

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/validation"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

// parseAction returns a node contain a character action, or a block of node containing
// a list of character actions
func (p *Parser) parseAction() (ast.Stmt, error) {
	char, err := p.consume(ast.ItemCharacterKey)
	if err != nil {
		// this really shouldn't happen since we already checked
		return nil, fmt.Errorf("ln%v: expecting character key, got %v", char.Line, char.Val)
	}
	charKey := shortcut.CharNameToKey[char.Val]

	// should be multiple action keys next
	var actions []*ast.CallExpr
	if n := p.peek(); n.Typ != ast.ItemActionKey {
		return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.Line, char.Val, n.Val)
	}

	// all actions needs to come before any + flags
Loop:
	for {
		switch n := p.next(); n.Typ {
		case ast.ItemTerminateLine:
			// stop here
			break Loop
		case ast.ItemActionKey:
			actionKey := action.StringToAction(n.Val)
			expr := &ast.CallExpr{
				Pos: char.Pos,
				Fun: &ast.Ident{
					Pos:   n.Pos,
					Value: "execute_action",
				},
				Args: make([]ast.Expr, 0),
			}
			expr.Args = append(expr.Args,
				// char
				&ast.NumberLit{
					Pos:      char.Pos,
					IntVal:   int64(charKey),
					FloatVal: float64(charKey),
				},
				// action
				&ast.NumberLit{
					Pos:      n.Pos,
					IntVal:   int64(actionKey),
					FloatVal: float64(actionKey),
				},
			)
			// check for param -> then repeat
			param, err := p.acceptOptionalParamReturnMap()
			if err != nil {
				return nil, err
			}
			if param == nil {
				param = &ast.MapExpr{Pos: n.Pos}
			}
			// validate params
			// TODO: this is inefficient but we don't have a "compile" step yet
			m := param.(*ast.MapExpr).Fields
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			err = validation.ValidateCharParamKeys(charKey, actionKey, keys)
			if err != nil {
				return nil, fmt.Errorf("ln%v: character %v: %w", n.Line, charKey, err)
			}
			expr.Args = append(expr.Args, param)

			// optional : and a number
			count, err := p.acceptOptionalRepeaterReturnCount()
			if err != nil {
				return nil, err
			}
			// add to array
			for range count {
				// TODO: all the repeated action will access the same map
				// ability implement should avoid modifying the maps
				actions = append(actions, expr)
			}

			n = p.next()
			if n.Typ != ast.ItemComma {
				p.backup()
				break Loop
			}
		default:
			// TODO: fix invalid key error
			return nil, fmt.Errorf("ln%v: expecting actions for character %v, got %v", n.Line, char.Val, n.Val)
		}
	}
	// check for optional flags

	// build stmt
	b := ast.NewBlockStmt(char.Pos)
	for _, v := range actions {
		b.Append(v)
	}
	return b, nil
}

func (p *Parser) acceptOptionalParamReturnMap() (ast.Expr, error) {
	// check for params
	n := p.peek()
	if n.Typ != ast.ItemLeftSquareParen {
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

	for k, v := range result.(*ast.MapExpr).Fields {
		switch v.(type) {
		case *ast.NumberLit:
			// skip
		default:
			return nil, fmt.Errorf("expected number in the map, got %v", v.String())
		}
		r[k] = int(v.(*ast.NumberLit).IntVal)
	}
	return r, nil
}

func (p *Parser) acceptOptionalRepeaterReturnCount() (int, error) {
	count := 1
	n := p.next()
	if n.Typ != ast.ItemColon {
		p.backup()
		return count, nil
	}
	// should be a number next
	n = p.next()
	if n.Typ != ast.ItemNumber {
		return count, fmt.Errorf("ln%v: expected a number after : but got %v", n.Line, n)
	}
	// parse number
	count, err := itemNumberToInt(n)
	return count, err
}
