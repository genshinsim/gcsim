package parser

//go:generate go tool pigeon -optimize-parser -o gcsim_parser.go gcsim.peg

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Parser struct {
	file  *ast.File
	input string
	res   *info.ActionList
	prog  *ast.BlockStmt

	chars          map[keys.Char]*info.CharacterProfile
	charOrder      []keys.Char
	currentCharKey keys.Char

	constantFolding bool
}

func New(file *ast.File, input string) *Parser {
	p := &Parser{
		chars:          make(map[keys.Char]*info.CharacterProfile),
		constantFolding: true,
	}
	p.file = file
	file.SetInput(input)
	p.input = input
	p.res = &info.ActionList{
		Settings: info.SimulatorSettings{
			EnableHitlag:    true,
			DefHalt:         true,
			NumberOfWorkers: 20,
			Iterations:      1000,
			Delays: info.Delays{
				Swap: 1,
			},
		},
		InitialPlayerPos: info.Coord{
			R: 0.3,
		},
	}
	p.prog = ast.NewBlockStmt(0)
	return p
}

func (p *Parser) Parse() (*info.ActionList, ast.Node, error) {
	_, err := Parse("config", []byte(p.input),
		GlobalStore("parser", p),
	)
	if err != nil {
		if list, ok := err.(ErrorLister); ok {
			for _, e := range list.Errors() {
				switch e := e.(type) {
				case ast.Error:
					return nil, nil, e
				case *ast.Error:
					return nil, nil, *e
				case ParserError:
					_, _, off := e.Pos()
					return nil, nil, ast.NewErrorf(p.file.Position(ast.Pos(off)), "%v", e.Error())
				default:
					return nil, nil, fmt.Errorf("parse error: %v", e)
				}
			}
		}
		return nil, nil, err
	}

	if len(p.charOrder) > 4 {
		p.res.Errors = append(p.res.Errors, fmt.Errorf("config contains a total of %v characters; cannot exceed 4", len(p.charOrder)))
	}

	if p.res.InitialChar == keys.NoChar {
		p.res.Errors = append(p.res.Errors, errors.New("config does not contain active char"))
	}

	initialCharFound := false
	for _, v := range p.charOrder {
		p.res.Characters = append(p.res.Characters, *p.chars[v])
		if v == p.res.InitialChar {
			initialCharFound = true
		}
		count := 0
		for _, c := range p.chars[v].Sets {
			count += c
		}
		if count > 5 {
			p.res.Errors = append(p.res.Errors, fmt.Errorf("character %v has more than 5 total set items", v.String()))
		}
	}

	if !initialCharFound && p.res.InitialChar != 0 {
		p.res.Errors = append(p.res.Errors, fmt.Errorf("active char %v not found in team", p.res.InitialChar))
	}

	if len(p.res.Targets) == 0 {
		p.res.Errors = append(p.res.Errors, errors.New("config does not contain any targets"))
	}

	for i := range p.res.Targets {
		if p.res.Targets[i].Pos.R == 0 {
			p.res.Targets[i].Pos.R = 1
		}
	}

	if p.res.Settings.DamageMode {
		for i := range p.res.Targets {
			if p.res.Targets[i].HP == 0 {
				p.res.Errors = append(p.res.Errors, fmt.Errorf("damage mode is activated; target #%v does not have hp set", i+1))
			}
		}
	}

	p.res.ErrorMsgs = make([]string, 0, len(p.res.Errors))
	for _, v := range p.res.Errors {
		p.res.ErrorMsgs = append(p.res.ErrorMsgs, v.Error())
	}

	return p.res, p.prog, nil
}


