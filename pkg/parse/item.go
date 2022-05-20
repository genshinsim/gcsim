package parse

import "fmt"

// Token represents a token or text string returned from the scanner.
type Token struct {
	typ  TokenType // The type of this item.
	pos  Pos       // The starting position, in bytes, of this item in the input string.
	Val  string    // The value of this item.
	line int       // The line number at the start of this item.
}

func (i Token) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.Val
	case i.typ == itemTerminateLine:
		return ";"
	case i.typ > itemTerminateLine && i.typ < itemKeyword:
		return i.Val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.Val)
		// case len(i.val) > 10:
		// 	return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// TokenType identifies the type of lex items.
type TokenType int

const (
	itemError TokenType = iota // error occurred; value is text of error

	itemEOF
	itemTerminateLine    // \n to denote end of a line
	itemAssign           // equals ('=') introducing an assignment
	itemComma            // coma (,) used to break up list of ident
	itemLeftParen        // '('
	itemRightParen       // ')'
	itemLeftSquareParen  // '['
	itemRightSquareParen // ']'
	itemLeftBrace        // '{'
	itemRightBrace       // '}'
	itemColon            // ':'
	itemPlus             // '+'
	itemMinus            // '-'
	itemAsterisk         // '*'
	itemSlash            // '/'
	// following is logic operator
	itemLogicOP // used only to delimit logical operation
	LogicNot    // !
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
	// item types
	itemTypes
	itemField      // alphanumeric identifier starting with '.'
	itemIdentifier // alphanumeric identifier not starting with '.'
	itemNumber     // simple number
	itemBool       // boolean
	itemString     // string, including quotes
	// Keywords appear after all the rest.
	itemKeyword        // used only to delimit the keywords
	keywordLet         // let
	keywordWhile       // while
	keywordIf          // if
	keywordElse        // else
	keywordFn          // fn
	keywordSwitch      // switch
	keywordCase        // case
	keywordDefault     // default
	keywordBreak       // break
	keywordContinue    // continue
	keywordFallthrough // fallthrough
	keywordLabel       // label
	keywordReturn      // return

	// Keywords specific to gcsim appears after this
	itemKeys
	itemStatKey      // stats: def%, def, etc..
	itemElementKey   // elements: pyro, hydro, etc..
	itemCharacterKey // characters: albedo, amber, etc..
	itemActionKey    // actions: skill, burst, attack, charge, etc...
)
