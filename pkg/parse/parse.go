package parse

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Parser struct {
	l              *lexer
	tokens         []item
	pos            int          //current position
	currentCharKey core.CharKey //current character being parsed

	//results
	cfg       *core.SimulationConfig
	chars     map[core.CharKey]*core.CharacterProfile
	charOrder []core.CharKey
	macros    map[string]core.ActionBlock
}

type parseFn func(*Parser) (parseFn, error)

func New(name, input string) *Parser {
	p := &Parser{}
	p.l = lex(name, input)
	p.pos = -1
	return p
}

func (p *Parser) Parse() (core.SimulationConfig, error) {
	//initialize
	var err error

	p.cfg = &core.SimulationConfig{}
	p.chars = make(map[core.CharKey]*core.CharacterProfile)
	p.charOrder = make([]core.CharKey, 0, 4)
	p.macros = make(map[string]core.ActionBlock)

	//default run options
	p.cfg.Settings.Duration = 90
	p.cfg.Settings.Iterations = 1000
	p.cfg.Settings.NumberOfWorkers = 20

	// Invalid value as default - should enforce setting mode to avoid confusion
	p.cfg.Settings.QueueMode = -1

	state := parseRows
	for state != nil && err == nil {
		state, err = state(p)
		if err != nil {
			return *p.cfg, err
		}
	}

	if err != nil {
		return *p.cfg, err
	}

	if len(p.charOrder) > 4 {
		return *p.cfg, fmt.Errorf("config contains a total of %v characters; cannot exceed 4", len(p.charOrder))
	}

	for _, v := range p.charOrder {
		p.cfg.Characters.Profile = append(p.cfg.Characters.Profile, *p.chars[v])
		//check number of set
		count := 0
		for _, c := range p.chars[v].Sets {
			count += c
		}
		if count > 5 {
			return *p.cfg, fmt.Errorf("character %v have more than 5 total set items", v.String())
		}
	}

	//check target hp
	for i, v := range p.cfg.Targets {
		if p.cfg.DamageMode && v.HP <= 0 {
			return *p.cfg, errors.New("if any one target has hp > 0, then all target must have hp > 0")
		} else if !p.cfg.DamageMode {
			//we should never actually get here
			p.cfg.Targets[i].HP = 0 //make sure its 0 if not running hp mode
		}
	}

	//check number of targets (should have at least 1)

	if len(p.cfg.Targets) == 0 {
		return *p.cfg, errors.New("no target set; you should have at least 1 target added")
	}

	return *p.cfg, nil
}

func (p *Parser) acceptSeqReturnLast(items ...ItemType) (item, error) {
	var n item
	for _, v := range items {
		n = p.next()
		if n.typ != v {
			_, file, no, _ := runtime.Caller(1)
			return n, fmt.Errorf("(%s#%d) expecting %v, got token %v at line: %v", file, no, v, n, p.tokens)
		}
	}
	return n, nil
}

func (p *Parser) consume(i ItemType) (item, error) {
	n := p.next()
	// log.Println(n)
	if n.typ != i {
		_, file, no, _ := runtime.Caller(1)
		return n, fmt.Errorf("(%s#%d) expecting %v, got token %v at line: %v", file, no, i, n, p.tokens)
	}
	return n, nil
}

func (p *Parser) next() item {
	p.pos++
	if p.pos == len(p.tokens) {
		return p.tokens[p.pos-1]
	}
	// log.Printf("next token: %v", p.tokens[p.pos])
	return p.tokens[p.pos]
}

func (p *Parser) backup() {
	if p.pos > 0 {
		p.pos--
	}
}

func (p *Parser) peek() item {
	next := p.next()
	p.backup()
	return next
}

func itemNumberToInt(i item) (int, error) {
	r, err := strconv.Atoi(i.val)
	return int(r), err
}

func itemNumberToFloat64(i item) (float64, error) {
	r, err := strconv.ParseFloat(i.val, 64)
	return r, err
}
