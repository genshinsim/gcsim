package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ  ItemType // The type of this item.
	pos  Pos      // The starting position, in bytes, of this item in the input string.
	val  string   // The value of this item.
	line int      // The line number at the start of this item.
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ == itemTerminateLine:
		return "End Line"
	case i.typ > itemCompareOp && i.typ < itemKeyword:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
		// case len(i.val) > 10:
		// 	return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// ItemType identifies the type of lex items.
type ItemType int

const (
	itemError ItemType = iota // error occurred; value is text of error
	itemBool                  // boolean constant

	itemAssign        // equals ('=') introducing an assignment
	itemDeclare       // colon-equals (':=') introducing a declaration
	itemAddToList     // plus-equals ('+=') introducing add to list
	itemComma         // coma (,) used to break up list of ident
	itemTerminateLine // \n to denote end of a line
	itemEOF
	itemField            // alphanumeric identifier starting with '.'
	itemIdentifier       // alphanumeric identifier not starting with '.'
	itemVariable         // variable starting with '$', such as '$' or  '$1' or '$hello'
	itemNumber           // simple number
	itemLeftParen        // '('
	itemRightParen       // ')'
	itemLeftSquareParen  // '['
	itemRightSquareParen // ']'
	itemString           // string, including quotes
	// following is logic operator
	itemLogicOP // used only to delimit logical operation
	LogicAnd    // && keyword
	LogicOr     // || keyword
	// following is comparison operator
	itemCompareOp        // used only to delimi comparison operators
	OpEqual              // == keyword
	OpNotEqual           // != keyword
	OpGreaterThan        // > keyword
	OpGreaterThanOrEqual // >= keyword
	OpLessThan           // < keyword
	OpLessThanOrEqual    // <= keyword
	// Keywords appear after all the rest.
	itemKeyword    // used only to delimit the keywords
	itemDot        // the cursor, spelled '.'
	itemAction     // action keyword
	itemChar       // char
	itemStats      // stats
	itemWeapon     // weapon
	itemArt        // art
	itemHurt       // hurt
	itemLvl        // lvl
	itemCons       // cons
	itemTalent     // talent
	itemRefine     // refine
	itemParam      // param
	itemLabel      // label
	itemCount      // count
	itemEle        // ele
	itemTarget     // char keyword
	itemExec       // exec keyword
	itemLock       // lock keyword
	itemIf         // if keyword
	itemSwap       // swap keyword
	itemPost       // trail keyword
	itemActive     // active keyword
	itemInterval   // interval keyword
	itemAmount     // amount keyword
	itemOnce       // once keyword
	itemEvery      // every keyword
	itemActionLock // actionlock keyword
	itemStartHP    // starthp keyword
	// these are configuration options
	itemOptions    // option
	itemMode       // mode
	itemSingle     // single
	itemAverage    // average
	itemIterations // iteration
	itemDuration   // duration
	itemSimHP      // simhp
	itemWorkers    // workers
	// stat types after the rest
	statKeyword  // delimit stats
	statDEFP     // def%
	statDEF      // def
	statHP       // hp
	statHPP      // hp%
	statATK      // atk
	statATKP     // atk%
	statER       // er
	statEM       // em
	statCR       // cr
	statCD       // cd
	statHeal     // heal
	statPyroP    // pyro%
	statHydroP   // hydro%
	statCryoP    // cryo%
	statElectroP // electro%
	statAnemoP   // anemo%
	statGeoP     // geo%
	statPhyP     // phys%
	statDendroP  // dendro%
	eleTypeKeyword
	elePyro     // pyro
	eleHydro    // hydro
	eleCryo     // cryo
	eleElectro  // electro
	eleGeo      // geo
	eleAnemo    // anemo
	eleDendro   // dendro
	elePhysical // physical

)

var key = map[string]ItemType{
	".": itemDot,
	//config related
	"options":   itemOptions,
	"mode":      itemMode,
	"single":    itemSingle,
	"average":   itemAverage,
	"iteration": itemIterations,
	"duration":  itemDuration,
	"simhp":     itemSimHP,
	"workers":   itemWorkers,
	//action related
	"actions":    itemAction,
	"char":       itemChar,
	"stats":      itemStats,
	"weapon":     itemWeapon,
	"art":        itemArt,
	"hurt":       itemHurt,
	"lvl":        itemLvl,
	"cons":       itemCons,
	"talent":     itemTalent,
	"refine":     itemRefine,
	"param":      itemParam,
	"label":      itemLabel,
	"count":      itemCount,
	"ele":        itemEle,
	"target":     itemTarget,
	"exec":       itemExec,
	"lock":       itemLock,
	"if":         itemIf,
	"swap":       itemSwap,
	"post":       itemPost,
	"active":     itemActive,
	"interval":   itemInterval,
	"amount":     itemAmount,
	"once":       itemOnce,
	"every":      itemEvery,
	"actionlock": itemActionLock,
	"starthp":    itemStartHP,
	//stats
	"def%":     statDEFP,
	"def":      statDEF,
	"hp":       statHP,
	"hp%":      statHPP,
	"atk":      statATK,
	"atk%":     statATKP,
	"er":       statER,
	"em":       statEM,
	"cr":       statCR,
	"cd":       statCD,
	"heal":     statHeal,
	"pyro%":    statPyroP,
	"hydro%":   statHydroP,
	"cryo%":    statCryoP,
	"electro%": statElectroP,
	"anemo%":   statAnemoP,
	"geo%":     statGeoP,
	"phys%":    statPhyP,
	"dendro%":  statDendroP,
	//element types
	"pyro":     elePyro,
	"hydro":    eleHydro,
	"cryo":     eleCryo,
	"electro":  eleElectro,
	"geo":      eleGeo,
	"anemo":    eleAnemo,
	"dendro":   eleDendro,
	"physical": elePhysical,
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name         string    // the name of the input; used only for error reports
	input        string    // the string being scanned
	pos          Pos       // current position in the input
	start        Pos       // start position of this item
	width        Pos       // width of last rune read from input
	items        chan item // channel of scanned items
	line         int       // 1+number of newlines seen
	startLine    int       // start line of this item
	parenDepth   int
	sqParenDepth int
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t ItemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], l.startLine}
	l.start = l.pos
	l.startLine = l.line
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), l.startLine}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	return <-l.items
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:      name,
		input:     input,
		items:     make(chan item),
		line:      1,
		startLine: 1,
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// lexText scans until an opening action delimiter, "{{".
func lexText(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	// n := l.peek()
	// log.Printf("lexText next is %c\n", n)
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case r == ';':
		l.emit(itemTerminateLine)
	case isSpace(r):
		l.ignore()
	case r == '#':
		return lexComment
	case r == '=':
		n := l.next()
		if n == '=' {
			l.emit(OpEqual)
		} else {
			l.backup()
			l.emit(itemAssign)
		}
	case r == ':':
		if l.next() != '=' {
			return l.errorf("expected :=")
		}
		l.emit(itemDeclare)
	case r == ',':
		l.emit(itemComma)
	case r == '+':
		n := l.next()
		if n == '=' {
			l.emit(itemAddToList)
		} else {
			l.backup()
			return lexNumber
		}
	case r == '.':
		// special look-ahead for ".field" so we don't break l.backup().
		if l.pos < Pos(len(l.input)) {
			r := l.input[l.pos]
			if r < '0' || '9' < r {
				return lexField
			}
		}
		fallthrough // '.' can start a number.
	case r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case r == '>':
		if n := l.next(); n == '=' {
			l.emit(OpGreaterThanOrEqual)
		} else {
			l.backup()
			l.emit(OpGreaterThan)
		}
	case r == '<':
		switch n := l.next(); n {
		case '=':
			l.emit(OpLessThanOrEqual)
		case '>':
			l.emit(OpNotEqual)
		default:
			l.backup()
			l.emit(OpLessThan)
		}
	case r == '|':
		if n := l.next(); n == '|' {
			l.emit(LogicOr)
		} else {
			return l.errorf("unrecognized character in action: %#U", r)
		}
	case r == '"':
		return lexQuote
	case r == '&':
		if n := l.next(); n == '&' {
			l.emit(LogicAnd)
		} else {
			return l.errorf("unrecognized character in action: %#U", r)
		}
	case r == '(':
		l.emit(itemLeftParen)
		l.parenDepth++
	case r == ')':
		l.emit(itemRightParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
	case r == '[':
		l.emit(itemLeftSquareParen)
		l.sqParenDepth++
	case r == ']':
		l.emit(itemRightSquareParen)
		l.sqParenDepth--
		if l.sqParenDepth < 0 {
			return l.errorf("unexpected right sq paren %#U", r)
		}
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexText
}

func lexComment(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '\n':
			l.backup()
			return lexText
		default:
			// absorb
		}
	}

}

// lexField scans a field: .Alphanumeric.
// The . has been scanned.
func lexField(l *lexer) stateFn {
	return lexFieldOrVariable(l, itemField)
}

// lexVariable scans a Variable: $Alphanumeric.
// The $ has been scanned.
func lexVariable(l *lexer) stateFn {
	if l.atTerminator() { // Nothing interesting follows -> "$".
		l.emit(itemVariable)
		return lexText
	}
	return lexFieldOrVariable(l, itemVariable)
}

// lexVariable scans a field or variable: [.$]Alphanumeric.
// The . or $ has been scanned.
func lexFieldOrVariable(l *lexer, typ ItemType) stateFn {
	if l.atTerminator() { // Nothing interesting follows -> "." or "$".
		if typ == itemVariable {
			l.emit(itemVariable)
		} else {
			l.emit(itemDot)
		}
		return lexText
	}
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	if !l.atTerminator() {
		return l.errorf("bad character %#U", r)
	}
	l.emit(typ)
	return lexText
}

func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexText
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > itemKeyword:
				l.emit(key[word])
			case word[0] == '.':
				l.emit(itemField)
			case word == "true", word == "false":
				l.emit(itemBool)
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}
	return lexText
}

func lexNumber(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")

	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumber)

	return lexText
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '%'
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(', '+', '=', '>', '<', '&', '!', ';', '[', ']':
		return true
	}
	return false
}
