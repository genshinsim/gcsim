package ast

import (
	"fmt"
	"testing"
)

func TestFields(t *testing.T) {
	input := `if .status.field > 0 { print("hi") };`

	l := lex(input)
	for n := l.nextItem(); n.Typ != itemEOF; n = l.nextItem() {
		fmt.Println(n)
	}
}
func TestBasicToken(t *testing.T) {
	input := `
	let y = fn(x) {
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
	
	for x = 0; x < 5; x = x + 1 {
		let i = y(x);
	}

	-1
	1
	-
	-a
	`

	expected := []Token{
		// function
		{Typ: keywordLet, Val: "let"},
		{Typ: itemIdentifier, Val: "y"},
		{Typ: itemAssign, Val: "="},
		{Typ: keywordFn, Val: "fn"},
		{Typ: itemLeftParen, Val: "("},
		{Typ: itemIdentifier, Val: "x"},
		// {typ: typeNum, Val: "num"},
		{Typ: itemRightParen, Val: ")"},
		// {typ: typeNum, Val: "num"}
		{Typ: itemLeftBrace, Val: "{"},
		{Typ: keywordReturn, Val: "return"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: ItemPlus, Val: "+"},
		{Typ: itemNumber, Val: "1"},
		{Typ: itemTerminateLine, Val: ";"},
		{Typ: itemRightBrace, Val: "}"},
		// variable
		{Typ: keywordLet, Val: "let"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemAssign, Val: "="},
		{Typ: itemNumber, Val: "5"},
		{Typ: itemTerminateLine, Val: ";"},
		// label
		{Typ: keywordLabel, Val: "label"},
		{Typ: itemIdentifier, Val: "A"},
		{Typ: itemColon, Val: ":"},
		// while loop
		{Typ: keywordWhile, Val: "while"},
		{Typ: itemLeftBrace, Val: "{"},
		// comment
		// {typ: itemComment, Val: "comment"},
		// function call
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemAssign, Val: "="},
		{Typ: itemIdentifier, Val: "y"},
		{Typ: itemLeftParen, Val: "("},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemRightParen, Val: ")"},
		{Typ: itemTerminateLine, Val: ";"},
		// if statement
		{Typ: keywordIf, Val: "if"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: OpGreaterThan, Val: ">"},
		{Typ: itemNumber, Val: "10"},
		{Typ: itemLeftBrace, Val: "{"},
		// break
		{Typ: keywordBreak, Val: "break"},
		{Typ: itemIdentifier, Val: "A"},
		{Typ: itemTerminateLine, Val: ";"},
		// end if
		{Typ: itemRightBrace, Val: "}"},
		// comment
		// {typ: itemComment, Val: "comment"},
		// switch
		{Typ: keywordSwitch, Val: "switch"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemLeftBrace, Val: "{"},
		// case
		{Typ: keywordCase, Val: "case"},
		{Typ: itemNumber, Val: "1"},
		{Typ: itemColon, Val: ":"},
		{Typ: keywordFallthrough, Val: "fallthrough"},
		{Typ: itemTerminateLine, Val: ";"},
		// case
		{Typ: keywordCase, Val: "case"},
		{Typ: itemNumber, Val: "2"},
		{Typ: itemColon, Val: ":"},
		{Typ: keywordFallthrough, Val: "fallthrough"},
		{Typ: itemTerminateLine, Val: ";"},
		// case
		{Typ: keywordCase, Val: "case"},
		{Typ: itemNumber, Val: "3"},
		{Typ: itemColon, Val: ":"},
		{Typ: keywordBreak, Val: "break"},
		{Typ: itemIdentifier, Val: "A"},
		{Typ: itemTerminateLine, Val: ";"},
		// end switch
		{Typ: itemRightBrace, Val: "}"},
		// end while
		{Typ: itemRightBrace, Val: "}"},
		// for loop
		{Typ: keywordFor, Val: "for"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemAssign, Val: "="},
		{Typ: itemNumber, Val: "0"},
		{Typ: itemTerminateLine, Val: ";"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: OpLessThan, Val: "<"},
		{Typ: itemNumber, Val: "5"},
		{Typ: itemTerminateLine, Val: ";"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemAssign, Val: "="},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: ItemPlus, Val: "+"},
		{Typ: itemNumber, Val: "1"},
		{Typ: itemLeftBrace, Val: "{"},
		// body
		{Typ: keywordLet, Val: "let"},
		{Typ: itemIdentifier, Val: "i"},
		{Typ: itemAssign, Val: "="},
		{Typ: itemIdentifier, Val: "y"},
		{Typ: itemLeftParen, Val: "("},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: itemRightParen, Val: ")"},
		{Typ: itemTerminateLine, Val: ";"},
		// end for
		{Typ: itemRightBrace, Val: "}"},
		// misc tests
		{Typ: itemNumber, Val: "-1"},
		{Typ: itemNumber, Val: "1"},
		{Typ: ItemMinus, Val: "-"},
		{Typ: ItemMinus, Val: "-"},
		{Typ: itemIdentifier, Val: "a"},
	}

	l := lex(input)
	i := 0
	for n := l.nextItem(); n.Typ != itemEOF; n = l.nextItem() {
		if expected[i].Typ != n.Typ && expected[i].Val != n.Val {
			t.Errorf("expected %v got %v", expected[i], n)
		}
		if i < len(expected)-1 {
			i++
		}
	}
}

func TestElseSpace(t *testing.T) {
	input := `
	if x > 1{}else{}
	`

	expected := []Token{
		// function
		{Typ: keywordIf, Val: "if"},
		{Typ: itemIdentifier, Val: "x"},
		{Typ: OpGreaterThan, Val: ">"},
		{Typ: itemNumber, Val: "1"},
		{Typ: itemLeftBrace, Val: "{"},
		{Typ: itemRightBrace, Val: "}"},
		{Typ: keywordElse, Val: "else"},
		{Typ: itemLeftBrace, Val: "{"},
		{Typ: itemRightBrace, Val: "}"},
	}

	l := lex(input)
	i := 0
	for n := l.nextItem(); n.Typ != itemEOF; n = l.nextItem() {
		if expected[i].Typ != n.Typ && expected[i].Val != n.Val {
			t.Errorf("expected %v got %v", expected[i], n)
			t.FailNow()
		}
		if i < len(expected)-1 {
			i++
		}
	}
}
