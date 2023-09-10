package ast

import "fmt"

// Token represents a token or text string returned from the scanner.
type Token struct {
	Typ  TokenType // The type of this item.
	pos  Pos       // The starting position, in bytes, of this item in the input string.
	Val  string    // The value of this item.
	line int       // The line number at the start of this item.
}

func (i Token) String() string {
	switch {
	case i.Typ == itemEOF:
		return "EOF"
	case i.Typ == itemError:
		return i.Val
	case i.Typ == itemTerminateLine:
		return ";"
	case i.Typ > itemTerminateLine && i.Typ < itemKeyword:
		return i.Val
	case i.Typ > itemKeyword:
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
	ItemPlus             // '+'
	ItemMinus            // '-'
	ItemAsterisk         // '*'
	ItemForwardSlash     // '/'
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
	keywordReturn      // return
	keywordFor         // for
	// Keywords after this are specific to Genshin (i.e. not generic scripting keywords)
	// These are special char related keywords
	keywordOptions           // options
	keywordAdd               // add
	keywordChar              // char
	keywordStats             // stats
	keywordWeapon            // weapon
	keywordSet               // set
	keywordLvl               // lvl
	keywordRefine            // refine
	keywordCons              // cons
	keywordTalent            // talent
	keywordCount             // count
	keywordParams            // params
	keywordLabel             // label
	keywordUntil             // until
	keywordActive            // active
	keywordTarget            // target
	keywordResist            // resist
	keywordEnergy            // energy
	keywordParticleThreshold // particle_threshold
	keywordParticleDropCount // particle_drop_count
	keywordHurt              // hurt

	// Keywords specific to gcsim appears after this
	itemKeys
	itemStatKey      // stats: def%, def, etc..
	itemElementKey   // elements: pyro, hydro, etc..
	itemCharacterKey // characters: albedo, amber, etc..
	itemActionKey    // actions: skill, burst, attack, charge, etc...
)
