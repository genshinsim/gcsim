package parse

import (
	"testing"
)

func TestBasicToken(t *testing.T) {
	input := `
	func y(x num) num {
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
	
	-1
	1
	-
	-a
	`

	expected := []Token{
		//function
		{typ: keywordFunc, Val: "func"},
		{typ: itemIdentifier, Val: "y"},
		{typ: itemLeftParen, Val: "("},
		{typ: itemIdentifier, Val: "x"},
		{typ: keywordNum, Val: "num"},
		{typ: itemRightParen, Val: ")"},
		{typ: keywordNum, Val: "num"},
		{typ: itemLeftBrace, Val: "{"},
		{typ: keywordReturn, Val: "return"},
		{typ: itemIdentifier, Val: "x"},
		{typ: itemPlus, Val: "+"},
		{typ: itemNumber, Val: "1"},
		{typ: itemTerminateLine, Val: ";"},
		{typ: itemRightBrace, Val: "}"},
		//variable
		{typ: keywordLet, Val: "let"},
		{typ: itemIdentifier, Val: "x"},
		{typ: itemAssign, Val: "="},
		{typ: itemNumber, Val: "5"},
		{typ: itemTerminateLine, Val: ";"},
		//label
		{typ: keywordLabel, Val: "label"},
		{typ: itemIdentifier, Val: "A"},
		{typ: itemColon, Val: ":"},
		//while loop
		{typ: keywordWhile, Val: "while"},
		{typ: itemLeftBrace, Val: "{"},
		//comment
		{typ: itemComment, Val: "comment"},
		//function call
		{typ: itemIdentifier, Val: "x"},
		{typ: itemAssign, Val: "="},
		{typ: itemIdentifier, Val: "y"},
		{typ: itemLeftParen, Val: "("},
		{typ: itemIdentifier, Val: "x"},
		{typ: itemRightParen, Val: ")"},
		{typ: itemTerminateLine, Val: ";"},
		//if statement
		{typ: keywordIf, Val: "if"},
		{typ: itemIdentifier, Val: "x"},
		{typ: OpGreaterThan, Val: ">"},
		{typ: itemNumber, Val: "10"},
		{typ: itemLeftBrace, Val: "{"},
		//break
		{typ: keywordBreak, Val: "break"},
		{typ: itemIdentifier, Val: "A"},
		{typ: itemTerminateLine, Val: ";"},
		//end if
		{typ: itemRightBrace, Val: "}"},
		//comment
		{typ: itemComment, Val: "comment"},
		//switch
		{typ: keywordSwitch, Val: "switch"},
		{typ: itemIdentifier, Val: "x"},
		{typ: itemLeftBrace, Val: "{"},
		//case
		{typ: keywordCase, Val: "case"},
		{typ: itemNumber, Val: "1"},
		{typ: itemColon, Val: ":"},
		{typ: keywordFallthrough, Val: "fallthrough"},
		{typ: itemTerminateLine, Val: ";"},
		//case
		{typ: keywordCase, Val: "case"},
		{typ: itemNumber, Val: "2"},
		{typ: itemColon, Val: ":"},
		{typ: keywordFallthrough, Val: "fallthrough"},
		{typ: itemTerminateLine, Val: ";"},
		//case
		{typ: keywordCase, Val: "case"},
		{typ: itemNumber, Val: "3"},
		{typ: itemColon, Val: ":"},
		{typ: keywordBreak, Val: "break"},
		{typ: itemIdentifier, Val: "A"},
		{typ: itemTerminateLine, Val: ";"},
		//end switch
		{typ: itemRightBrace, Val: "}"},
		//end while
		{typ: itemRightBrace, Val: "}"},
		//misc tests
		{typ: itemNumber, Val: "-1"},
		{typ: itemNumber, Val: "1"},
		{typ: itemMinus, Val: "-"},
		{typ: itemMinus, Val: "-"},
		{typ: itemIdentifier, Val: "a"},
	}

	l := lex(input)
	i := 0
	for n := l.nextItem(); n.typ != itemEOF; n = l.nextItem() {
		if expected[i].typ != n.typ && expected[i].Val != n.Val {
			t.Errorf("expected %v got %v", expected[i], n)
		}
		if i < len(expected)-1 {
			i++
		}
	}
}
