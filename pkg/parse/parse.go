package parse

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"

	"github.com/genshinsim/gsim/pkg/def"
)

var actionKeys = map[string]def.ActionType{
	"sequence":        def.ActionSequence,
	"sequence_strict": def.ActionSequenceStrict,
	"reset_sequence":  def.ActionSequenceReset,
	"skill":           def.ActionSkill,
	"burst":           def.ActionBurst,
	"attack":          def.ActionAttack,
	"charge":          def.ActionCharge,
	"high_plunge":     def.ActionHighPlunge,
	"low_lunge":       def.ActionLowPlunge,
	"aim":             def.ActionAim,
	"dash":            def.ActionDash,
	"jump":            def.ActionJump,
	"swap":            def.ActionSwap,
}

var eleKeys = map[string]def.EleType{
	"pyro":            def.Pyro,
	"hydro":           def.Hydro,
	"cryo":            def.Cryo,
	"electro":         def.Electro,
	"geo":             def.Geo,
	"anemo":           def.Anemo,
	"dendro":          def.Dendro,
	"physical":        def.Physical,
	"frozen":          def.Frozen,
	"electro-charged": def.EC,
	"":                def.NoElement,
}

type Parser struct {
	l      *lexer
	tokens []item
	pos    int //current position

	//results
	result *def.Config
	chars  map[string]*def.CharacterProfile
}

type parseFn func(*Parser) (parseFn, error)

func New(name, input string) *Parser {
	p := &Parser{}
	p.l = lex(name, input)
	p.pos = -1
	return p
}

func (p *Parser) Parse() (def.Config, error) {
	//initialize
	var err error

	p.result = &def.Config{}
	p.chars = make(map[string]*def.CharacterProfile)

	//default run options
	p.result.Mode.Duration = 90
	p.result.Mode.Iteration = 1000
	p.result.Mode.Workers = 20

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

	if p.result.Mode.HP > 0 {
		p.result.Mode.HPMode = true
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

	//options mode=average iteration=5000 duration=90 simhp=0 workers=24;
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
			case itemAverage:
				p.result.Mode.Average = true
			case itemSingle:
				p.result.Mode.Average = false
			default:
				return nil, fmt.Errorf("bad token %v parsing options, expecting average or single. line %v", n.val, p.tokens)
			}
		case itemIterations:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.Mode.Iteration, err = itemNumberToInt(n)
			}
		case itemDuration:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.Mode.Duration, err = itemNumberToInt(n)
			}
		case itemSimHP:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.Mode.HP, err = itemNumberToFloat64(n)
			}
		case itemWorkers:
			n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
			if err == nil {
				p.result.Mode.Workers, err = itemNumberToInt(n)
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
