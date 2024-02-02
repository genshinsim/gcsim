package ast

import "testing"

func TestLetTyping(t *testing.T) {
	parseAndPrint(
		`let x = 1;`,
		t,
	)
	parseAndPrint(
		`let x number = 1;`,
		t,
	)
	parseAndPrint(
		`let x fn() = fn() { print("hi"); };`,
		t,
	)
	parseAndPrint(
		`let x fn() string = fn() string { print("hi"); };`,
		t,
	)
}
