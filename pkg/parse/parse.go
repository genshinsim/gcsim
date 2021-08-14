package parse

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"

	"github.com/genshinsim/gsim/pkg/core"
)

var actionKeys = map[string]core.ActionType{
	"sequence":        core.ActionSequence,
	"sequence_strict": core.ActionSequenceStrict,
	"reset_sequence":  core.ActionSequenceReset,
	"skill":           core.ActionSkill,
	"burst":           core.ActionBurst,
	"attack":          core.ActionAttack,
	"charge":          core.ActionCharge,
	"high_plunge":     core.ActionHighPlunge,
	"low_lunge":       core.ActionLowPlunge,
	"aim":             core.ActionAim,
	"dash":            core.ActionDash,
	"jump":            core.ActionJump,
	"swap":            core.ActionSwap,
}

var eleKeys = map[string]core.EleType{
	"pyro":            core.Pyro,
	"hydro":           core.Hydro,
	"cryo":            core.Cryo,
	"electro":         core.Electro,
	"geo":             core.Geo,
	"anemo":           core.Anemo,
	"dendro":          core.Dendro,
	"physical":        core.Physical,
	"frozen":          core.Frozen,
	"electro-charged": core.EC,
	"":                core.NoElement,
}

type Parser struct {
	l      *lexer
	tokens []item
	pos    int //current position

	//results
	result *core.Config
	chars  map[string]*core.CharacterProfile
}

type parseFn func(*Parser) (parseFn, error)

func New(name, input string) *Parser {
	p := &Parser{}
	p.l = lex(name, input)
	p.pos = -1
	return p
}

func (p *Parser) Parse() (core.Config, error) {
	//initialize
	var err error

	p.result = &core.Config{}
	p.chars = make(map[string]*core.CharacterProfile)

	//default run options
	p.result.RunOptions.Duration = 90
	p.result.RunOptions.Iteration = 1000
	p.result.RunOptions.Workers = 20

	state := parseRows
	for state != nil && err == nil {
		state, err = state(p)
		if err != nil {
			return *p.result, err
		}
	}

	if err != nil {
		return *p.result, err
	}

	keys := make([]string, 0, len(p.chars))
	for k := range p.chars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		p.result.Characters.Profile = append(p.result.Characters.Profile, *p.chars[k])
	}

	//check target hp

	for i, v := range p.result.Targets {
		if p.result.RunOptions.DamageMode && v.HP <= 0 {
			p.result.Targets[i].HP = 1000000 //default 1 mil hp
		} else if !p.result.RunOptions.DamageMode {
			p.result.Targets[i].HP = 0 //make sure its 0 if not running hp mode
		}
	}

	return *p.result, nil
}

func parseRows(p *Parser) (parseFn, error) {
	p.tokens = make([]item, 0, 20)
	p.pos = -1

	//consume the entire line
	for n := p.l.nextItem(); n.typ != itemEOF; n = p.l.nextItem() {
		if n.typ == itemError {
			return nil, errors.New(n.val)
		}
		p.tokens = append(p.tokens, n)
		if n.typ == itemTerminateLine {
			break
		}
	}

	if len(p.tokens) == 0 {
		return nil, nil //at end of the file
	}

	n := p.next()

	//check if we're parsing options
	if n.typ == itemOptions {
		return parseOptions, nil
	}

	x, err := p.consume(itemAddToList)
	if err != nil {
		return nil, fmt.Errorf("<parse row> expecting += but got %v; line %v", x, p.tokens)
	}

	//check what the first token of the line i
	switch n.typ {
	case itemLabel: //profile label
		return parseLabel, nil
	case itemAction: //action
		return parseAction, nil
	case itemChar: //char basic
		return parseChar, nil
	case itemStats: //char stats
		return parseStats, nil
	case itemWeapon: //weapon data
		return parseWeapon, nil
	case itemArt: //artifact sets
		return parseArtifacts, nil
	case itemHurt: //hurt events
		return parseHurtEvent, nil
	case itemTarget: //enemy related
		return parseTarget, nil
	case itemActive: //active char
		return parseActiveChar, nil
	default:
		return nil, fmt.Errorf("<parse row> invalid token at start of line: %v", p.tokens)
	}

}

func parseOptions(p *Parser) (parseFn, error) {
	var err error

	//options mode=damage debug=true iteration=5000 duration=90 workers=24;
	for n := p.next(); n.typ != itemEOF; n = p.next() {
		switch n.typ {
		case itemMode:
			//expect avg or single
			_, err = p.consume(itemAssign)
			if err != nil {
				return nil, err
			}
			item := p.next()
			switch item.typ {
			case itemDamage:
				p.result.RunOptions.DamageMode = true
			case itemTime:
				p.result.RunOptions.DamageMode = false
			default:
				return nil, fmt.Errorf("bad token %v parsing options, expecting average or single. line %v", n.val, p.tokens)
			}
		case itemDebug:
			n, err = p.acceptSeqReturnLast(itemAssign, itemBool)
			if err == nil {
				p.result.RunOptions.Debug = n.val == "true"
			}
		case itemIterations:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.RunOptions.Iteration, err = itemNumberToInt(n)
			}
		case itemDuration:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.RunOptions.Duration, err = itemNumberToInt(n)
			}
		case itemWorkers:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.RunOptions.Workers, err = itemNumberToInt(n)
			}
		case itemTerminateLine:
			return parseRows, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing options")
}

func parseLabel(p *Parser) (parseFn, error) {
	ident, err := p.consume(itemIdentifier)
	if err != nil {
		return nil, err
	}
	p.result.Label = ident.val
	n, err := p.consume(itemTerminateLine)
	if err != nil {
		return nil, fmt.Errorf("bad token %v parsing label, expecting ;. line %v", n.val, p.tokens)
	}
	return parseRows, nil
}

func itemNumberToInt(i item) (int, error) {
	r, err := strconv.Atoi(i.val)
	return int(r), err
}

func itemNumberToFloat64(i item) (float64, error) {
	r, err := strconv.ParseFloat(i.val, 64)
	return r, err
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

func tokensToStringArray(tokens []item) []string {
	r := make([]string, 0, len(tokens))
	for _, v := range tokens {
		r = append(r, v.val)
	}
	return r
}
