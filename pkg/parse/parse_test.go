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

func TestParseIf(t *testing.T) {
	s := `
	if x > y {
		1 + 1;
		2 + 2;
		3 + 3;
	} else {
		//do stuff
		c + d;
	}
	`
	p := New(s)
	res, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(res.Program.String())
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

const cfg = `
	let y = fn(a, b) {
		a + b;
	}
	let x = 0;
	while x < 10 {
		//x = y(x, 1);
		//do loopy stuff
	}
`

func TestCfg(t *testing.T) {
	p := New(cfg)
	fmt.Printf("parsing:\n %v\n", cfg)
	res, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("output:")
	fmt.Println(res.Program.String())
}
