package parser

import (
	"fmt"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/validation"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

// ---------------------------------------------------------------------------
// Types shared with generated parser (used in action code blocks)
// ---------------------------------------------------------------------------

type mapEntry struct {
	key string
	val ast.Expr
}

type callPart struct {
	lparenPos int
	args      []ast.Expr
}

type paramDef struct {
	name string
	pos  int
	typ  ast.ExprType
}

type actionItem struct {
	actionName string
	actionKey  action.Action
	pos        int
	params     *ast.MapExpr
	count      int
}

type charNameInfo struct {
	name string
	key  keys.Char
	pos  int
}

type energyItem struct {
	kind      string
	ival      int64
	startIval int64
	endIval   int64
}

type hurtItem struct {
	kind      string
	ival      int64
	startIval int64
	endIval   int64
	minF      float64
	maxF      float64
	ele       attributes.Element
}

type optionItem struct {
	kind string
	text string
	bval bool
	ival int64
	fval float64
}

type targetItem struct {
	kind      string
	intVal    int64
	floatVal  float64
	ele       attributes.Element
	x, y      float64
	text      string
	paramList []paramPair
}

type paramPair struct {
	key string
	val ast.Expr
	pos int
}

type addStatItem struct {
	kind  string
	key   attributes.Stat
	val   float64
	label string
}

type charDetailData struct {
	kind      string
	intVal    int64
	intVal2   int64
	intVal3   int64
	mapVal  map[string]int
}

type addWeaponItem struct {
	kind    string
	intVal  int64
	intVal2 int64
	mapVal  map[string]int
}

type addSetItem struct {
	kind    string
	intVal  int64
	mapVal  map[string]int
}

type randomStatItem struct {
	kind     string
	intVal   int64
	statName string
}

type forParts struct {
	Init any
	Cond ast.Expr
	Post any
	Body *ast.BlockStmt
}

// ---------------------------------------------------------------------------
// Error interfaces (mirror json/errors.go pattern)
// ---------------------------------------------------------------------------

type ErrorLister interface {
	Errors() []error
}

func (e errList) Errors() []error { return e }

type ParserError interface {
	Error() string
	InnerError() error
	Pos() (line, col, offset int)
	Expected() []string
}

func (p *parserError) InnerError() error {
	return p.Inner
}

func (p *parserError) Pos() (line, col, offset int) {
	return p.pos.line, p.pos.col, p.pos.offset
}

func (p *parserError) Expected() []string {
	return p.expected
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func (p *Parser) errf(off int, format string, args ...any) {
	panic(ast.NewErrorf(p.file.Position(ast.Pos(off)), format, args...))
}

func (p *Parser) err(off int, msg string) {
	panic(ast.NewErrorf(p.file.Position(ast.Pos(off)), "%s", msg))
}

// ---------------------------------------------------------------------------
// Keyword / identifier classification
// ---------------------------------------------------------------------------

var reservedWords map[string]bool

func init() {
	reservedWords = make(map[string]bool)
	for _, w := range ast.Keywords() {
		reservedWords[w] = true
	}
	reservedWords["true"] = true
	reservedWords["false"] = true
}

func isPlainIdent(word string) bool {
	if reservedWords[word] {
		return false
	}
	if _, ok := ast.StatKeys[word]; ok {
		return false
	}
	if _, ok := ast.EleKeys[word]; ok {
		return false
	}
	if _, ok := shortcut.CharNameToKey[word]; ok {
		return false
	}
	if _, ok := ast.ActionKeys[word]; ok {
		return false
	}
	return true
}

func isCharName(word string) bool {
	_, ok := shortcut.CharNameToKey[word]
	return ok
}

func isActionName(word string) bool {
	_, ok := ast.ActionKeys[word]
	return ok
}

func isStatName(word string) bool {
	_, ok := ast.StatKeys[word]
	return ok
}

func isEleName(word string) bool {
	_, ok := ast.EleKeys[word]
	return ok
}

// ---------------------------------------------------------------------------
// Token building helpers
// ---------------------------------------------------------------------------

func identTok(c *current, off int, val string) ast.Token {
	p := c.globalStore["parser"].(*Parser)
	return ast.Token{
		Typ:  ast.ItemIdentifier,
		Pos:  ast.Pos(off),
		Val:  val,
		Line: p.file.Position(ast.Pos(off)).Line,
	}
}

func opTok(c *current, typ ast.TokenType, off int, val string) ast.Token {
	p := c.globalStore["parser"].(*Parser)
	return ast.Token{
		Typ:  typ,
		Pos:  ast.Pos(off),
		Val:  val,
		Line: p.file.Position(ast.Pos(off)).Line,
	}
}

// ---------------------------------------------------------------------------
// Expression helpers
// ---------------------------------------------------------------------------

func buildBinaryExpr(p *Parser, left ast.Expr, opTok ast.Token, right ast.Expr) ast.Expr {
	expr := &ast.BinaryExpr{
		Pos:   opTok.Pos,
		Op:    opTok,
		Left:  left,
		Right: right,
	}
	if p.constantFolding {
		folded, err := foldConstants(expr)
		if err != nil {
			p.err(int(opTok.Pos), err.Error())
		}
		return folded
	}
	return expr
}

func buildUnaryExpr(p *Parser, opTok ast.Token, right ast.Expr) ast.Expr {
	expr := &ast.UnaryExpr{
		Pos:   opTok.Pos,
		Op:    opTok,
		Right: right,
	}
	if p.constantFolding {
		folded, err := foldConstants(expr)
		if err != nil {
			p.err(int(opTok.Pos), err.Error())
		}
		return folded
	}
	return expr
}

func constFloat(p *Parser, e ast.Expr) float64 {
	switch v := e.(type) {
	case *ast.NumberLit:
		if v.IsFloat {
			return v.FloatVal
		}
		return float64(v.IntVal)
	default:
		p.errf(int(e.Position()), "expecting number, got %v", e.String())
		return 0
	}
}

func constFloatFold(p *Parser, e ast.Expr) float64 {
	var folded ast.Expr
	if p.constantFolding {
		var err error
		folded, err = foldConstants(e)
		if err != nil {
			p.err(int(e.Position()), err.Error())
		}
	} else {
		folded = e
	}
	return constFloat(p, folded)
}

func parseIntNumber(p *Parser, off int, s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		p.errf(off, "cannot parse %v to int", s)
	}
	return v
}

func buildNumberLit(p *Parser, off int, text string) *ast.NumberLit {
	num := &ast.NumberLit{Pos: ast.Pos(off)}
	iv, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		num.IntVal = iv
		num.FloatVal = float64(iv)
	} else {
		fv, err := strconv.ParseFloat(text, 64)
		if err != nil {
			p.errf(off, "cannot parse %v to number", text)
		}
		num.IsFloat = true
		num.FloatVal = fv
	}
	return num
}

func buildBoolLit(off int, val string) *ast.NumberLit {
	num := &ast.NumberLit{Pos: ast.Pos(off)}
	switch val {
	case ast.TrueVal:
		num.IntVal = 1
		num.FloatVal = 1
	case ast.FalseVal:
		num.IntVal = 0
		num.FloatVal = 0
	}
	return num
}

func stripQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// ---------------------------------------------------------------------------
// Action block building
// ---------------------------------------------------------------------------

func buildActionBlock(p *Parser, cn charNameInfo, items []actionItem) *ast.BlockStmt {
	var actions []*ast.CallExpr
	charKey := cn.key
	charPos := cn.pos
	for _, item := range items {
		mapExpr := item.params
		if mapExpr == nil {
			mapExpr = &ast.MapExpr{Pos: ast.Pos(item.pos)}
		}
		m := mapExpr.Fields
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		err := validation.ValidateCharParamKeys(charKey, item.actionKey, keys)
		if err != nil {
			p.errf(charPos, "character %v: %v", charKey, err)
		}
		expr := &ast.CallExpr{
			Pos: ast.Pos(charPos),
			Fun: &ast.Ident{
				Pos:   ast.Pos(item.pos),
				Value: "execute_action",
			},
			Args: []ast.Expr{
				&ast.NumberLit{
					Pos:      ast.Pos(charPos),
					IntVal:   int64(charKey),
					FloatVal: float64(charKey),
				},
				&ast.NumberLit{
					Pos:      ast.Pos(item.pos),
					IntVal:   int64(item.actionKey),
					FloatVal: float64(item.actionKey),
				},
				mapExpr,
			},
		}
		for range item.count {
			actions = append(actions, expr)
		}
	}
	b := ast.NewBlockStmt(ast.Pos(charPos))
	for _, v := range actions {
		b.Append(v)
	}
	return b
}

// ---------------------------------------------------------------------------
// Character helpers
// ---------------------------------------------------------------------------

func (p *Parser) ensureChar(key keys.Char) {
	if _, ok := p.chars[key]; !ok {
		r := info.CharacterProfile{}
		r.Base.Key = key
		r.Stats = make([]float64, attributes.EndStatType)
		r.StatsByLabel = make(map[string][]float64)
		r.Params = make(map[string]int)
		r.Sets = make(map[keys.Set]int)
		r.SetParams = make(map[keys.Set]map[string]int)
		r.Weapon.Params = make(map[string]int)
		r.Base.Element = keys.CharKeyToEle[key]
		p.chars[key] = &r
		p.charOrder = append(p.charOrder, key)
	}
}

// ---------------------------------------------------------------------------
// Target helpers
// ---------------------------------------------------------------------------

var allElements = []attributes.Element{
	attributes.Electro, attributes.Cryo, attributes.Hydro,
	attributes.Physical, attributes.Pyro, attributes.Geo,
	attributes.Dendro, attributes.Anemo,
}

func (p *Parser) acceptOptionalTargetParamsFromBody(pairs []paramPair) enemy.TargetParams {
	result := enemy.TargetParams{
		HpMultiplier: 0.0,
		Particles:    true,
	}
	for _, pp := range pairs {
		switch pp.key {
		case "hp_mult":
			result.HpMultiplier = constFloatFold(p, pp.val)
		case "particles":
			v := int(constFloatFold(p, pp.val))
			result.Particles = v != 0
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Energy/Hurt helpers (accept []any for use from grammar)
// ---------------------------------------------------------------------------

func parseEnergyIntervalOnce(p *Parser, items []any) {
	p.res.EnergySettings.Active = true
	p.res.EnergySettings.Once = true
	for _, it := range items {
		it = unwrapItem(it)
		e := it.(energyItem)
		switch e.kind {
		case "interval":
			p.res.EnergySettings.Start = int(e.ival)
		case "amount":
			p.res.EnergySettings.Amount = int(e.ival)
		}
	}
}

func parseEnergyIntervalEvery(p *Parser, items []any) {
	p.res.EnergySettings.Active = true
	p.res.EnergySettings.Once = false
	for _, it := range items {
		it = unwrapItem(it)
		e := it.(energyItem)
		switch e.kind {
		case "interval":
			p.res.EnergySettings.Start = int(e.startIval)
			p.res.EnergySettings.End = int(e.endIval)
		case "amount":
			p.res.EnergySettings.Amount = int(e.ival)
		}
	}
}

func parseHurtOnce(p *Parser, items []any) {
	p.res.HurtSettings.Active = true
	p.res.HurtSettings.Once = true
	for _, it := range items {
		it = unwrapItem(it)
		e := it.(hurtItem)
		switch e.kind {
		case "interval":
			p.res.HurtSettings.Start = int(e.ival)
		case "amount":
			p.res.HurtSettings.Min = e.minF
			p.res.HurtSettings.Max = e.maxF
		case "element":
			p.res.HurtSettings.Element = e.ele
		}
	}
}

func parseHurtEvery(p *Parser, items []any) {
	p.res.HurtSettings.Active = true
	p.res.HurtSettings.Once = false
	for _, it := range items {
		it = unwrapItem(it)
		e := it.(hurtItem)
		switch e.kind {
		case "interval":
			p.res.HurtSettings.Start = int(e.startIval)
			p.res.HurtSettings.End = int(e.endIval)
		case "amount":
			p.res.HurtSettings.Min = e.minF
			p.res.HurtSettings.Max = e.maxF
		case "element":
			p.res.HurtSettings.Element = e.ele
		}
	}
}

// ---------------------------------------------------------------------------
// Options helpers
// ---------------------------------------------------------------------------

func applyOptions(p *Parser, items []any) {
	for _, it := range items {
		it = unwrapItem(it)
		e := it.(optionItem)
		switch e.kind {
		case "debug":

		case "defhalt":
			p.res.Settings.DefHalt = e.bval
		case "hitlag":
			p.res.Settings.EnableHitlag = e.bval
		case "iteration":
			p.res.Settings.Iterations = int(e.ival)
		case "duration":
			p.res.Settings.Duration = e.fval
		case "workers":
			p.res.Settings.NumberOfWorkers = int(e.ival)
		case "mode":

		case "swap_delay":
			p.res.Settings.Delays.Swap = int(e.ival)
		case "attack_delay":
			p.res.Settings.Delays.Attack = int(e.ival)
		case "charge_delay":
			p.res.Settings.Delays.Charge = int(e.ival)
		case "skill_delay":
			p.res.Settings.Delays.Skill = int(e.ival)
		case "burst_delay":
			p.res.Settings.Delays.Burst = int(e.ival)
		case "jump_delay":
			p.res.Settings.Delays.Jump = int(e.ival)
		case "dash_delay":
			p.res.Settings.Delays.Dash = int(e.ival)
		case "aim_delay":
			p.res.Settings.Delays.Aim = int(e.ival)
		case "frame_defaults":
			if e.text == "human" {
				p.res.Settings.Delays.Swap = 8
				p.res.Settings.Delays.Attack = 5
				p.res.Settings.Delays.Charge = 5
				p.res.Settings.Delays.Skill = 5
				p.res.Settings.Delays.Burst = 5
				p.res.Settings.Delays.Dash = 5
				p.res.Settings.Delays.Jump = 5
				p.res.Settings.Delays.Aim = 5
			} else {
				p.err(0, fmt.Sprintf("unrecognized option for frame_defaults specified: %v", e.text))
			}
		case "ignore_burst_energy":
			p.res.Settings.IgnoreBurstEnergy = e.bval
		default:
			p.err(0, fmt.Sprintf("unrecognized option specified: %v", e.kind))
		}
	}
}

// ---------------------------------------------------------------------------
// Target application helper
// ---------------------------------------------------------------------------

func applyTargets(p *Parser, items []any) {
	var unwrapped []any
	for _, it := range items {
		unwrapped = append(unwrapped, unwrapItem(it))
	}
	items = unwrapped
	r := info.EnemyProfile{}
	r.Resist = make(map[attributes.Element]float64)
	r.ParticleElement = attributes.NoElement
	for _, it := range items {
		e := it.(targetItem)
		switch e.kind {
		case "lvl":
			r.Level = int(e.intVal)
		case "hp":
			r.HP = e.floatVal
			p.res.Settings.DamageMode = true
			r.Modified = true
		case "resist":
			for _, elem := range allElements {
				r.Resist[elem] += e.floatVal
			}
			r.Modified = true
		case "element":
			r.Resist[e.ele] += e.floatVal
			r.Modified = true
		case "pos":
			r.Pos.X = e.x
			r.Pos.Y = e.y
		case "radius":
			r.Pos.R = e.floatVal
		case "type":
			params := p.acceptOptionalTargetParamsFromBody(e.paramList)
			err := enemy.ConfigureTarget(&r, e.text, params)
			if err != nil {
				p.err(0, err.Error())
			}
			p.res.Settings.DamageMode = true
		case "freeze_resist":
			r.FreezeResist = e.floatVal
			r.Modified = true
		case "particle_threshold":
			r.ParticleDropThreshold = e.floatVal
			r.ParticleDrops = nil
			r.ParticleElement = attributes.NoElement
			r.Modified = true
		case "particle_drop_count":
			r.ParticleDropCount = e.floatVal
			r.Modified = true
		case "particle_element":
			r.ParticleElement = e.ele
			r.Modified = true
		}
	}
	p.res.Targets = append(p.res.Targets, r)
}

// ---------------------------------------------------------------------------
// FnCore param list duplicate check
// ---------------------------------------------------------------------------

func acceptIntMap(p *Parser, off int, m *ast.MapExpr) map[string]int {
	r := make(map[string]int)
	for k, v := range m.Fields {
		switch vt := v.(type) {
		case *ast.NumberLit:
			r[k] = int(vt.IntVal)
		default:
			p.errf(off, "expected number in the map, got %v", v.String())
		}
	}
	return r
}

func unwrapItem(it any) any {
	if sl, ok := it.([]any); ok && len(sl) == 2 {
		return sl[1]
	}
	return it
}

func applyCharDetails(p *Parser, key keys.Char, items []any) {
	c := p.chars[key]
	for _, it := range items {
		d := unwrapItem(it).(charDetailData)
		switch d.kind {
		case "lvl":
			c.Base.Level = int(d.intVal)
			c.Base.MaxLevel = int(d.intVal2)
		case "cons":
			c.Base.Cons = int(d.intVal)
		case "talent":
			c.Talents.Attack = int(d.intVal)
			c.Talents.Skill = int(d.intVal2)
			c.Talents.Burst = int(d.intVal3)
		case "params":
			c.Params = d.mapVal
		}
	}
}

func applyCharAddWeapon(p *Parser, key keys.Char, name string, items []any) {
	var unwrapped []any
	for _, it := range items {
		unwrapped = append(unwrapped, unwrapItem(it))
	}
	items = unwrapped
	c := p.chars[key]
	weaponName := stripQuotes(name)
	label, ok := shortcut.WeaponNameToKey[weaponName]
	if !ok {
		p.errf(0, "invalid weapon %v", weaponName)
	}
	c.Weapon.Key = label
	c.Weapon.Name = c.Weapon.Key.String()
	for _, it := range items {
		d := it.(addWeaponItem)
		switch d.kind {
		case "lvl":
			c.Weapon.Level = int(d.intVal)
			c.Weapon.MaxLevel = int(d.intVal2)
		case "refine":
			c.Weapon.Refine = int(d.intVal)
		case "params":
			c.Weapon.Params = d.mapVal
		}
	}
}

func applyCharAddSet(p *Parser, key keys.Char, name string, items []any) {
	var unwrapped []any
	for _, it := range items {
		unwrapped = append(unwrapped, unwrapItem(it))
	}
	items = unwrapped
	c := p.chars[key]
	setName := stripQuotes(name)
	label, ok := shortcut.SetNameToKey[setName]
	if !ok {
		p.errf(0, "invalid set %v", setName)
	}
	for _, it := range items {
		d := it.(addSetItem)
		switch d.kind {
		case "count":
			c.Sets[label] = int(d.intVal)
		case "params":
			c.SetParams[label] = d.mapVal
		}
	}
}

func applyCharAddStats(p *Parser, key keys.Char, items []any) {
	var unwrapped []any
	for _, it := range items {
		unwrapped = append(unwrapped, unwrapItem(it))
	}
	items = unwrapped
	c := p.chars[key]
	line := make([]float64, attributes.EndStatType)
	var keyLabel string
	for _, it := range items {
		d := it.(addStatItem)
		switch d.kind {
		case "stat":
			line[d.key] += d.val
		case "label":
			keyLabel = d.label
		}
	}
	m, ok := c.StatsByLabel[keyLabel]
	if !ok {
		m = make([]float64, attributes.EndStatType)
	}
	for i, v := range line {
		c.Stats[i] += v
		m[i] += v
	}
	c.StatsByLabel[keyLabel] = m
}

func applyCharAddRandomStats(p *Parser, key keys.Char, items []any) {
	rs := &info.RandomSubstats{Rarity: 5}
	for _, it := range items {
		it = unwrapItem(it)
		d := it.(randomStatItem)
		switch d.kind {
		case "rarity":
			rs.Rarity = int(d.intVal)
		case "sand":
			rs.Sand = ast.StatKeys[d.statName]
		case "goblet":
			rs.Goblet = ast.StatKeys[d.statName]
		case "circlet":
			rs.Circlet = ast.StatKeys[d.statName]
		}
	}
	if err := rs.Validate(); err != nil {
		p.err(0, err.Error())
	}
	p.chars[key].RandomSubstats = rs
}

func checkDuplicateParams(p *Parser, args []*ast.Ident) {
	chk := make(map[string]bool)
	for _, v := range args {
		if _, ok := chk[v.Value]; ok {
			p.errf(int(v.Pos), "fn contains duplicated param name %v", v.Value)
		}
		chk[v.Value] = true
	}
}
