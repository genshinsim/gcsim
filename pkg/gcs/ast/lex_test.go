package ast

import (
	"fmt"
	"testing"
)

func TestFields(t *testing.T) {
	input := `if .status.field > 0 { print("hi") };`

	l := NewLexer(input)
	for n := l.NextItem(); n.Typ != ItemEOF; n = l.NextItem() {
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
		{Typ: KeywordLet, Val: "let"},
		{Typ: ItemIdentifier, Val: "y"},
		{Typ: ItemAssign, Val: "="},
		{Typ: KeywordFn, Val: "fn"},
		{Typ: ItemLeftParen, Val: "("},
		{Typ: ItemIdentifier, Val: "x"},
		// {typ: typeNum, Val: "num"},
		{Typ: ItemRightParen, Val: ")"},
		// {typ: typeNum, Val: "num"}
		{Typ: ItemLeftBrace, Val: "{"},
		{Typ: KeywordReturn, Val: "return"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemPlus, Val: "+"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemTerminateLine, Val: ";"},
		{Typ: ItemRightBrace, Val: "}"},
		// variable
		{Typ: KeywordLet, Val: "let"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemAssign, Val: "="},
		{Typ: ItemNumber, Val: "5"},
		{Typ: ItemTerminateLine, Val: ";"},
		// label
		{Typ: KeywordLabel, Val: "label"},
		{Typ: ItemIdentifier, Val: "A"},
		{Typ: ItemColon, Val: ":"},
		// while loop
		{Typ: KeywordWhile, Val: "while"},
		{Typ: ItemLeftBrace, Val: "{"},
		// comment
		// {typ: itemComment, Val: "comment"},
		// function call
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemAssign, Val: "="},
		{Typ: ItemIdentifier, Val: "y"},
		{Typ: ItemLeftParen, Val: "("},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemRightParen, Val: ")"},
		{Typ: ItemTerminateLine, Val: ";"},
		// if statement
		{Typ: KeywordIf, Val: "if"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: OpGreaterThan, Val: ">"},
		{Typ: ItemNumber, Val: "10"},
		{Typ: ItemLeftBrace, Val: "{"},
		// break
		{Typ: KeywordBreak, Val: "break"},
		{Typ: ItemIdentifier, Val: "A"},
		{Typ: ItemTerminateLine, Val: ";"},
		// end if
		{Typ: ItemRightBrace, Val: "}"},
		// comment
		// {typ: itemComment, Val: "comment"},
		// switch
		{Typ: KeywordSwitch, Val: "switch"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemLeftBrace, Val: "{"},
		// case
		{Typ: KeywordCase, Val: "case"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemColon, Val: ":"},
		{Typ: KeywordFallthrough, Val: "fallthrough"},
		{Typ: ItemTerminateLine, Val: ";"},
		// case
		{Typ: KeywordCase, Val: "case"},
		{Typ: ItemNumber, Val: "2"},
		{Typ: ItemColon, Val: ":"},
		{Typ: KeywordFallthrough, Val: "fallthrough"},
		{Typ: ItemTerminateLine, Val: ";"},
		// case
		{Typ: KeywordCase, Val: "case"},
		{Typ: ItemNumber, Val: "3"},
		{Typ: ItemColon, Val: ":"},
		{Typ: KeywordBreak, Val: "break"},
		{Typ: ItemIdentifier, Val: "A"},
		{Typ: ItemTerminateLine, Val: ";"},
		// end switch
		{Typ: ItemRightBrace, Val: "}"},
		// end while
		{Typ: ItemRightBrace, Val: "}"},
		// for loop
		{Typ: KeywordFor, Val: "for"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemAssign, Val: "="},
		{Typ: ItemNumber, Val: "0"},
		{Typ: ItemTerminateLine, Val: ";"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: OpLessThan, Val: "<"},
		{Typ: ItemNumber, Val: "5"},
		{Typ: ItemTerminateLine, Val: ";"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemAssign, Val: "="},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemPlus, Val: "+"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemLeftBrace, Val: "{"},
		// body
		{Typ: KeywordLet, Val: "let"},
		{Typ: ItemIdentifier, Val: "i"},
		{Typ: ItemAssign, Val: "="},
		{Typ: ItemIdentifier, Val: "y"},
		{Typ: ItemLeftParen, Val: "("},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: ItemRightParen, Val: ")"},
		{Typ: ItemTerminateLine, Val: ";"},
		// end for
		{Typ: ItemRightBrace, Val: "}"},
		// misc tests
		{Typ: ItemMinus, Val: "-"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemMinus, Val: "-"},
		{Typ: ItemMinus, Val: "-"},
		{Typ: ItemIdentifier, Val: "a"},
	}

	l := NewLexer(input)
	i := 0
	for n := l.NextItem(); n.Typ != ItemEOF; n = l.NextItem() {
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
		{Typ: KeywordIf, Val: "if"},
		{Typ: ItemIdentifier, Val: "x"},
		{Typ: OpGreaterThan, Val: ">"},
		{Typ: ItemNumber, Val: "1"},
		{Typ: ItemLeftBrace, Val: "{"},
		{Typ: ItemRightBrace, Val: "}"},
		{Typ: KeywordElse, Val: "else"},
		{Typ: ItemLeftBrace, Val: "{"},
		{Typ: ItemRightBrace, Val: "}"},
	}

	l := NewLexer(input)
	i := 0
	for n := l.NextItem(); n.Typ != ItemEOF; n = l.NextItem() {
		if expected[i].Typ != n.Typ && expected[i].Val != n.Val {
			t.Errorf("expected %v got %v", expected[i], n)
			t.FailNow()
		}
		if i < len(expected)-1 {
			i++
		}
	}
}
