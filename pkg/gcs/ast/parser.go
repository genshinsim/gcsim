package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

type Parser struct {
	lex *lexer
	res *ActionList

	//other information tracked as we parse
	chars          map[keys.Char]*character.CharacterProfile
	charOrder      []keys.Char
	currentCharKey keys.Char

	//lookahead
	token []Token
	pos   int

	//parseFn
	prefixParseFns map[TokenType]func() Expr
	infixParseFns  map[TokenType]func(Expr) Expr
}
type ActionList struct {
	Targets     []enemy.EnemyProfile         `json:"targets"`
	PlayerPos   core.Coord                   `json:"player_initial_pos"`
	Characters  []character.CharacterProfile `json:"characters"`
	InitialChar keys.Char                    `json:"initial"`
	Program     *BlockStmt
	Energy      EnergySettings
	Settings    SimulatorSettings
}

type EnergySettings struct {
	Every  float64
	Amount int
	Lambda float64
}

type SimulatorSettings struct {
	Duration   float64
	DamageMode bool
	//other stuff
	NumberOfWorkers int // how many workers to run the simulation
	Iterations      int // how many iterations to run
	Delays          Delays
}

type Delays struct {
	Skill  int
	Burst  int
	Attack int
	Charge int
	Aim    int
	Dash   int
	Jump   int
	Swap   int
}

func (c *ActionList) Copy() *ActionList {

	r := *c

	r.Targets = make([]enemy.EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters = make([]character.CharacterProfile, len(c.Characters))
	for i, v := range c.Characters {
		r.Characters[i] = v.Clone()
	}

	r.Program = c.Program.CopyBlock()
	return &r
}

func (a *ActionList) PrettyPrint() string {
	prettyJson, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(prettyJson)
}

type parseFn func(*Parser) (parseFn, error)

func New(input string) *Parser {
	p := &Parser{
		chars:          make(map[keys.Char]*character.CharacterProfile),
		prefixParseFns: make(map[TokenType]func() Expr),
		infixParseFns:  make(map[TokenType]func(Expr) Expr),
		token:          make([]Token, 0, 20),
		pos:            -1,
	}
	p.lex = lex(input)
	p.res = &ActionList{
		Program: newBlockStmt(0),
	}
	//expr functions
	p.prefixParseFns[itemIdentifier] = p.parseIdent
	p.prefixParseFns[itemNumber] = p.parseNumber
	p.prefixParseFns[itemString] = p.parseString
	p.prefixParseFns[LogicNot] = p.parseUnaryExpr
	p.prefixParseFns[ItemMinus] = p.parseUnaryExpr
	p.prefixParseFns[itemLeftParen] = p.parseParen
	p.infixParseFns[LogicAnd] = p.parseBinaryExpr
	p.infixParseFns[LogicOr] = p.parseBinaryExpr
	p.infixParseFns[ItemPlus] = p.parseBinaryExpr
	p.infixParseFns[ItemMinus] = p.parseBinaryExpr
	p.infixParseFns[ItemForwardSlash] = p.parseBinaryExpr
	p.infixParseFns[ItemAsterisk] = p.parseBinaryExpr
	p.infixParseFns[OpEqual] = p.parseBinaryExpr
	p.infixParseFns[OpNotEqual] = p.parseBinaryExpr
	p.infixParseFns[OpLessThan] = p.parseBinaryExpr
	p.infixParseFns[OpLessThanOrEqual] = p.parseBinaryExpr
	p.infixParseFns[OpGreaterThan] = p.parseBinaryExpr
	p.infixParseFns[OpGreaterThanOrEqual] = p.parseBinaryExpr
	p.infixParseFns[itemLeftParen] = p.parseCall

	return p
}

// consume returns err if next token does not match expected
// otherwise return next token and nil error
func (p *Parser) consume(i TokenType) (Token, error) {
	n := p.next()
	if n.Typ != i {
		return n, errors.New("unexpected token")
	}
	return n, nil
}

// next returns the next token.
func (p *Parser) next() Token {
	p.pos++
	if p.pos == len(p.token) {
		//grab more from the stream
		n := p.lex.nextItem()
		p.token = append(p.token, n)
	}
	return p.token[p.pos]
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.pos--
	//no op if at beginning
	if p.pos < -1 {
		p.pos = -1
	}
}

// peek returns but does not consume the next token.
func (p *Parser) peek() Token {
	n := p.next()
	p.backup()
	return n
}

func (p *Parser) acceptSeqReturnLast(items ...TokenType) (Token, error) {
	var n Token
	for _, v := range items {
		n = p.next()
		if n.Typ != v {
			_, file, no, _ := runtime.Caller(1)
			return n, fmt.Errorf("(%s#%d) expecting %v, got token %v", file, no, v, n)
		}
	}
	return n, nil
}

func itemNumberToInt(i Token) (int, error) {
	r, err := strconv.Atoi(i.Val)
	return int(r), err
}

func itemNumberToFloat64(i Token) (float64, error) {
	r, err := strconv.ParseFloat(i.Val, 64)
	return r, err
}
