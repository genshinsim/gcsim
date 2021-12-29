package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Pos int

func (p Pos) Position() Pos {
	return p
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
	l.items <- item{
		typ:  t,
		pos:  l.start,
		val:  l.input[l.start:l.pos],
		line: l.startLine,
	}
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
	case r == ':':
		l.emit(itemColon)
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
			l.emit(itemEqual)
		}
	case r == ',':
		l.emit(itemComma)
	case r == '+':
		//check if next item is a number or not; if number lexNumber
		//otherwise it's a + sign
		n := l.next()
		if isNumeric(n) {
			//back up twice
			l.backup()
			l.backup()
			return lexNumber
		}
		//otherwise it's a plus sign
		l.backup()
		l.emit(itemPlus)
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
	case r == '/':
		l.emit(itemForwardSlash)
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
				l.emit(checkIdentifier(word))
			}
			break Loop
		}
	}
	return lexText
}

func checkIdentifier(word string) ItemType {
	if _, ok := statKeys[word]; ok {
		return itemStatKey
	}
	if _, ok := eleKeys[word]; ok {
		return itemElementKey
	}
	if _, ok := keys.CharNameToKey[word]; ok {
		return itemCharacterKey
	}
	if _, ok := actionKeys[word]; ok {
		return itemActionKey
	}
	return itemIdentifier
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
	return r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '%'
}

// is Numeric reports whether r is a digit
func isNumeric(r rune) bool {
	return unicode.IsDigit(r)
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
