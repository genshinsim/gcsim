package parse

import (
	"testing"
)

func TestBasicToken(t *testing.T) {
	input := `
	let y = fn(x num) num {
		return x + 1;
	}
	let x = 5;
	label A:
	while {
		#comment
		x = y(x);
		if x > 10 {
			break A;
		}
		//comment
		switch x {
		case 1:
			fallthrough;
		case 2:
			fallthrough;
		case 3:
			break A;
		}
	}
	`

	expected := []Token{
		//function
		{typ: keywordLet, val: "let"},
		{typ: itemIdentifier, val: "y"},
		{typ: itemAssign, val: "="},
		{typ: keywordFn, val: "fn"},
		{typ: itemLeftParen, val: "("},
		{typ: itemIdentifier, val: "x"},
		{typ: keywordNum, val: "num"},
		{typ: itemRightParen, val: ")"},
		{typ: keywordNum, val: "num"},
		{typ: itemLeftBrace, val: "{"},
		{typ: keywordReturn, val: "return"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemPlus, val: "+"},
		{typ: itemNumber, val: "1"},
		{typ: itemTerminateLine, val: ";"},
		{typ: itemRightBrace, val: "}"},
		//variable
		{typ: keywordLet, val: "let"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemAssign, val: "="},
		{typ: itemNumber, val: "5"},
		{typ: itemTerminateLine, val: ";"},
		//label
		{typ: keywordLabel, val: "label"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemColon, val: ":"},
		//while loop
		{typ: keywordWhile, val: "while"},
		{typ: itemLeftBrace, val: "{"},
		//comment
		{typ: itemComment, val: "comment"},
		//function call
		{typ: itemIdentifier, val: "x"},
		{typ: itemAssign, val: "="},
		{typ: itemIdentifier, val: "y"},
		{typ: itemLeftParen, val: "("},
		{typ: itemIdentifier, val: "x"},
		{typ: itemRightParen, val: ")"},
		{typ: itemTerminateLine, val: ";"},
		//if statement
		{typ: keywordIf, val: "if"},
		{typ: itemIdentifier, val: "x"},
		{typ: OpGreaterThan, val: ">"},
		{typ: itemNumber, val: "10"},
		{typ: itemLeftBrace, val: "{"},
		//break
		{typ: keywordBreak, val: "break"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemTerminateLine, val: ";"},
		//end if
		{typ: itemRightBrace, val: "}"},
		//comment
		{typ: itemComment, val: "comment"},
		//switch
		{typ: keywordSwitch, val: "switch"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemLeftBrace, val: "{"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "1"},
		{typ: itemColon, val: ":"},
		{typ: keywordFallthrough, val: "fallthrough"},
		{typ: itemTerminateLine, val: ";"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "2"},
		{typ: itemColon, val: ":"},
		{typ: keywordFallthrough, val: "fallthrough"},
		{typ: itemTerminateLine, val: ";"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "3"},
		{typ: itemColon, val: ":"},
		{typ: keywordBreak, val: "break"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemTerminateLine, val: ";"},
		//end switch
		{typ: itemRightBrace, val: "}"},
		//end while
		{typ: itemRightBrace, val: "}"},
	}

	l := lex("test", input)
	i := 0
	for n := l.nextItem(); n.typ != itemEOF; n = l.nextItem() {
		if expected[i].typ != n.typ && expected[i].val != n.val {
			t.Errorf("expected %v got %v", expected[i], n)
		}
		if i < len(expected)-1 {
			i++
		}
	}
}

func testActionToken(t *testing.T) {

	input := `
	let y = fn(x num) num {
		return x + 1;
	}
	let x = 5;
	label A:
	while {
		## comment
		x = y(x);
		if x > 10 {
			break A;
		}
		// comment
		switch x {
		case 1:
			fallthrough;
		case 2:
			fallthrough;
		case 3:
			break A;
		}
	}
	`

	expected := []Token{
		//function
		{typ: keywordLet, val: "let"},
		{typ: itemIdentifier, val: "y"},
		{typ: itemAssign, val: "="},
		{typ: keywordFn, val: "fn"},
		{typ: itemLeftParen, val: "("},
		{typ: itemIdentifier, val: "x"},
		{typ: keywordNum, val: "num"},
		{typ: itemRightParen, val: ")"},
		{typ: keywordNum, val: "num"},
		{typ: itemLeftBrace, val: "{"},
		{typ: keywordReturn, val: "return"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemPlus, val: "+"},
		{typ: itemNumber, val: "1"},
		{typ: itemTerminateLine, val: ";"},
		{typ: itemRightBrace, val: "}"},
		//variable
		{typ: keywordLet, val: "let"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemAssign, val: "="},
		{typ: itemNumber, val: "5"},
		{typ: itemTerminateLine, val: ";"},
		//label
		{typ: keywordLabel, val: "label"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemColon, val: ":"},
		//while loop
		{typ: keywordWhile, val: "while"},
		{typ: itemLeftBrace, val: "{"},
		//function call
		{typ: itemIdentifier, val: "x"},
		{typ: itemAssign, val: "="},
		{typ: itemIdentifier, val: "y"},
		{typ: itemLeftParen, val: "("},
		{typ: itemIdentifier, val: "x"},
		{typ: itemRightParen, val: ")"},
		{typ: itemTerminateLine, val: ";"},
		//if statement
		{typ: keywordIf, val: "if"},
		{typ: itemIdentifier, val: "x"},
		{typ: OpGreaterThan, val: ">"},
		{typ: itemNumber, val: "10"},
		{typ: itemLeftBrace, val: "{"},
		//break
		{typ: keywordBreak, val: "break"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemTerminateLine, val: ";"},
		//end if
		{typ: itemRightBrace, val: "}"},
		//switch
		{typ: keywordSwitch, val: "switch"},
		{typ: itemIdentifier, val: "x"},
		{typ: itemLeftBrace, val: "{"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "1"},
		{typ: itemColon, val: ":"},
		{typ: keywordFallthrough, val: "fallthrough"},
		{typ: itemTerminateLine, val: ";"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "2"},
		{typ: itemColon, val: ":"},
		{typ: keywordFallthrough, val: "fallthrough"},
		{typ: itemTerminateLine, val: ";"},
		//case
		{typ: keywordCase, val: "case"},
		{typ: itemNumber, val: "3"},
		{typ: itemColon, val: ":"},
		{typ: keywordBreak, val: "break"},
		{typ: itemIdentifier, val: "A"},
		{typ: itemTerminateLine, val: ";"},
		//end switch
		{typ: itemRightBrace, val: "}"},
		//end while
		{typ: itemRightBrace, val: "}"},
	}

	l := lex("test", input)
	i := 0
	for n := l.nextItem(); n.typ != itemEOF; n = l.nextItem() {
		if expected[i].typ != n.typ || expected[i].val != n.val {
			t.Errorf("expected %v got %v", expected[i], n)
		}
		if i < len(expected)-1 {
			i++
		}
	}

}
