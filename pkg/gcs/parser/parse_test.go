package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestOrderPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1+2*3;",
			"(1 + (2 * 3))",
		},
		{
			"1+2+3;",
			`((1 + 2) + 3)`,
		},
		{
			"1 * 2 + 3;",
			"((1 * 2) + 3)",
		},
		{
			"a * b + c;",
			`((a * b) + c)`,
		},
		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"a - b;",
			"(a - b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"(1+2)*3;",
			"((1 + 2) * 3)",
		},
		{
			"1==2 && 3!=4;",
			"((1 == 2) && (3 != 4))",
		},
		{
			"1 && 0 || 1+2 == 3;",
			"((1 && 0) || ((1 + 2) == 3))",
		},
	}

	for _, test := range tests {
		file := ast.NewFile()
		p := New(file, test.input)
		p.constantFolding = false
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

// func prettyPrint(body interface{}) {
// 	b, err := json.MarshalIndent(body, "", "\t")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(b))
// }

const cfg = `
	active xingqiu;
	switch a {
	case 1:
		1+1;
		fallthrough;
	case 2:
		2+2;
		break;
	default:
		3+3;
	}
	fn y(a, b) {
		let c = a + b;
		return c;
	}
	let x = 0;
	while x < 10 {
		x = y(x, 1);
		//do loopy stuff
		xingqiu skill;
		xingqiu attack:4;
		if x > 0 {
			continue;
		} else {
			break;
		}
	}
	for x = 0; x < 5; x = x + 1 {
		let i = y(x);
	}
`

func TestCfg(t *testing.T) {
	file := ast.NewFile()
	p := New(file, cfg)
	fmt.Printf("parsing:\n %v\n", cfg)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("output:")
	fmt.Println(prog.String())
}

const charaction = `
xingqiu attack[randomparam=2]:4,skill;
xingqiu burst[orbital=0];
active xingqiu;
`

func TestCharAction(t *testing.T) {
	file := ast.NewFile()
	p := New(file, charaction)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("output:")
	fmt.Println(prog.String())
}

const charstats = `
raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="favoniuslance" refine=3 lvl=90/90;
raiden add set="tenacityofthemillelith" count=4;
raiden add stats hp=4780 atk=311.0 er=0.5180 cr=0.3110 electro%=0.4660;
raiden add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

xingqiu char lvl=90/90 cons=6 talent=9,9,9;
xingqiu add weapon="harbingerofdawn" refine=5 lvl=90/90;
xingqiu add set="noblesseoblige" count=4;
xingqiu add stats hp=4780 atk=311.0 atk%=0.4660 cr=0.3110 hydro%=0.4660;
xingqiu add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.3306 em=39.64 cr=0.2648 cd=0.7944 ;																						

bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="thealleyflash" refine=1 lvl=90/90;
bennett add set="instructor" count=4;
bennett add stats hp=3571 atk=232.0 em=187.0 cr=0.2320 pyro%=0.3480;
bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.2979 cd=0.4634 ;																						

xiangling char lvl=90/90 cons=6 talent=9,9,9;
xiangling add weapon="thecatch" refine=5 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;
xiangling add stats hp=4780 atk=311.0 em=187.0 cr=0.3110 pyro%=0.4660;
xiangling add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944;

active raiden;

energy every interval=480,720 amount=1;
`

func TestCharAdd(t *testing.T) {
	file := ast.NewFile()
	p := New(file, charstats)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	spew.Config.Dump(prog)
}

func TestField(t *testing.T) {
	file := ast.NewFile()
	p := New(file, `if .status.field > 0 { print("hi"); }`)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	spew.Config.Dump(prog)
}

func TestActionStartLine(t *testing.T) {
	file := ast.NewFile()
	p := New(file, `xingqiu attack; skill`)
	_, prog, err := p.Parse()
	if err == nil {
		t.Errorf("xingqiu attack; skill parsed incorrectly without error")
		t.FailNow()
	}
	spew.Config.Dump(prog)
}

func parseAndPrint(s string, t *testing.T) {
	file := ast.NewFile()
	p := New(file, s)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("output:")
	fmt.Println(prog.String())
}

func TestCharFieldRanges(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"valid refine", `xiangling add weapon="thecatch" refine=5 lvl=90/90;`, false},
		{"refine too low", `xiangling add weapon="thecatch" refine=0 lvl=90/90;`, true},
		{"refine too high", `xiangling add weapon="thecatch" refine=6 lvl=90/90;`, true},
		// refine is serialized as an int32, without a range check this wraps to 5
		{"refine overflowing int32", `xiangling add weapon="thecatch" refine=4294967301 lvl=90/90;`, true},
		{"valid cons", `xiangling char lvl=90/90 cons=6 talent=9,9,9;`, false},
		{"cons too high", `xiangling char lvl=90/90 cons=7 talent=9,9,9;`, true},
		{"valid talent", `xiangling char lvl=90/90 cons=0 talent=1,10,10;`, false},
		{"attack talent too low", `xiangling char lvl=90/90 cons=0 talent=0,9,9;`, true},
		{"skill talent too high", `xiangling char lvl=90/90 cons=0 talent=9,11,9;`, true},
		{"burst talent too high", `xiangling char lvl=90/90 cons=0 talent=9,9,11;`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := ast.NewFile()
			p := New(file, tt.input)
			_, _, err := p.Parse()
			if tt.expectErr && err == nil {
				t.Errorf("%v: expected an error, got none", tt.input)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("%v: expected no error, got %v", tt.input, err)
			}
		})
	}
}
