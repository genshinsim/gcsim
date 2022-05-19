package parse

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestParseIdent(t *testing.T) {
	p := New("x;")

	al, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	prettyPrint(al)
}

func TestOrderPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1+2*3;",
			"(1 + (2 * 3));",
		},
		{
			"1+2+3;",
			`((1 + 2) + 3);`,
		},
		{
			"1 * 2 + 3;",
			"((1 * 2) + 3);",
		},
		{
			"a * b + c;",
			`((a * b) + c);`,
		},
		{
			"-a * b;",
			"((-a) * b);",
		},
		{
			"a - b;",
			"(a - b);",
		},
		{
			"!-a;",
			"(!(-a));",
		},
		{
			"(1+2)*3;",
			"((1 + 2) * 3);",
		},
	}

	for _, test := range tests {
		p := New(test.input)
		res, err := p.Parse()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// prettyPrint(res)
		actual := res.Program.String()
		//strip \n
		actual = strings.TrimSuffix(actual, "\n")
		if actual != test.expected {
			t.Errorf("expected=%q, got %q", test.expected, actual)
		}

	}
}

func prettyPrint(body interface{}) {
	b, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
