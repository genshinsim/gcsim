package parse

import "fmt"

// Token represents a token or text string returned from the scanner.
type Token struct {
	typ  tokenType // The type of this item.
	pos  Pos       // The starting position, in bytes, of this item in the input string.
	val  string    // The value of this item.
	line int       // The line number at the start of this item.
}

func (i Token) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ == itemTerminateLine:
		return ";"
	case i.typ > itemCompareOp && i.typ < itemKeyword:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
		// case len(i.val) > 10:
		// 	return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// tokenType identifies the type of lex items.
type tokenType int

const (
	itemError tokenType = iota // error occurred; value is text of error
	itemBool                   // boolean constant

	itemComment       // comment text
	itemAssign        // equals ('=') introducing an assignment
	itemComma         // coma (,) used to break up list of ident
	itemTerminateLine // \n to denote end of a line
	itemEOF
	itemField            // alphanumeric identifier starting with '.'
	itemIdentifier       // alphanumeric identifier not starting with '.'
	itemNumber           // simple number
	itemLeftParen        // '('
	itemRightParen       // ')'
	itemLeftSquareParen  // '['
	itemRightSquareParen // ']'
	itemLeftBrace        // '{'
	itemRightBrace       // '}'
	itemColon            // ':'
	itemPlus             // '+'
	itemMinus            // '-'
	itemMultiply         // '*'
	itemDivide           // '/'
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
	itemDot              // the cursor, spelled '.'
	// Keywords appear after all the rest.
	itemKeyword        // used only to delimit the keywords
	keywordLet         // let
	keywordWhile       // while
	keywordIf          // if
	keywordElse        // else
	keywordFn          // fn
	keywordSwitch      // switch
	keywordCase        // case
	keywordBreak       // break
	keywordFallthrough // fallthrough
	keywordLabel       // label
	keywordNum         // num
	keywordReturn      // return

	// Keywords specific to gcsim appears after this
	itemKeys
	itemStatKey      // stats: def%, def, etc..
	itemElementKey   // elements: pyro, hydro, etc..
	itemCharacterKey // characters: albedo, amber, etc..
	itemActionKey    // actions: skill, burst, attack, charge, etc...
)
