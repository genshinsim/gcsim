package parser

import (
	"strings"
	"testing"
)

func TestConstantFolding(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1+2*3;",
			"7",
		},
		{
			"1+2+3;",
			"6",
		},
		{
			"1 * 2 + 3;",
			"5",
		},
		{
			"(1+2)*3;",
			"9",
		},
		{
			"1==2 && 3!=4;",
			"0",
		},
		{
			"1 && 0 || 1+2 == 3;",
			"1",
		},
	}

	for _, test := range tests {
		p := New(test.input)
		_, prog, err := p.Parse()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// prettyPrint(res)
		actual := prog.String()
		// strip \n
		actual = strings.TrimSuffix(actual, "\n")
		if actual != test.expected {
			t.Errorf("expected=%q, got %q", test.expected, actual)
		}
	}
}
