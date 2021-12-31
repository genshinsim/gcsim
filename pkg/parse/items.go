package parse

import "fmt"

// item represents a token or text string returned from the scanner.
type item struct {
	typ  ItemType // The type of this item.
	pos  Pos      // The starting position, in bytes, of this item in the input string.
	val  string   // The value of this item.
	line int      // The line number at the start of this item.
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

	itemEqual         // equals ('=') introducing an assignment
	itemComma         // coma (,) used to break up list of ident
	itemTerminateLine // \n to denote end of a line
	itemEOF
	itemField            // alphanumeric identifier starting with '.'
	itemIdentifier       // alphanumeric identifier not starting with '.'
	itemVariable         // variable starting with '$', such as '$' or  '$1' or '$hello'
	itemNumber           // simple number
	itemForwardSlash     // '/'
	itemLeftParen        // '('
	itemRightParen       // ')'
	itemLeftSquareParen  // '['
	itemRightSquareParen // ']'
	itemColon            // ':'
	itemPlus             // '+'
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
	itemKeyword // used only to delimit the keywords
	// these are command keywords
	itemOptions    // options
	itemChain      // chain
	itemWaitFor    // wait_for
	itemWait       // wait (calc mode only)
	itemRestart    // restart (calc mode only)
	itemResetLimit // reset_limit
	itemHurt       // hurt
	itemEnergy     // energy
	itemActive     // active
	itemTarget     // target
	// these are options related keywords
	itemDebug      // debug
	itemIterations // iteration
	itemDuration   // duration
	itemWorkers    // workers
	itemERCalc     // er_calc
	// these are special char related keywords

	itemAdd     // add
	itemChar    // char
	itemStats   // stats
	itemWeapon  // weapon
	itemSet     // set
	itemLvl     // lvl
	itemRefine  // refine
	itemCons    // cons
	itemTalent  // talent
	itemStartHP // start_hp
	itemCount   // count
	itemParams  // params
	itemLabel   // label
	itemUntil   // until
	// these are flags

	itemSwapLock // swap_lock
	itemIf       // if
	itemSwap     // swap_to
	itemOnField  // is_onfield
	itemLimit    // limit
	itemTry      // try
	itemTimeout  // timeout
	itemNeeds    // needs
	// these are wait specific key words

	itemValue  // value
	itemMax    // max
	itemFiller // filler
	// these are energy related flags

	itemInterval // interval
	itemAmount   // amount
	itemOnce     // once
	itemEvery    // every

	// these are hurt related flags

	itemEle //ele

	// these are related to target setting

	itemResist

	// stat types after the rest

	itemKeys
	itemStatKey      // stats: def%, def, etc..
	itemElementKey   // elements: pyro, hydro, etc..
	itemCharacterKey // characters: albedo, amber, etc..
	itemActionKey    // actions: skill, burst, attack, charge, etc...
)
