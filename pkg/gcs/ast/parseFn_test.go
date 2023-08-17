package ast

import (
	"fmt"
	"testing"
)

const fntest = `
active bennett;
fn y(x) {
    print(x);
    return x +1;
}
let xy = fn(a number, b number) number {
	return a+b;
};

let z = f(2);

print(z);

print(xy(1,2));
`

func TestFnCall(t *testing.T) {
	p := New(fntest)
	_, prog, err := p.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("output:")
	fmt.Println(prog.String())
}

func TestFnTyping(t *testing.T) {
	parseAndPrint(
		`fn z() { print("hi"); }`,
		t,
	)

	parseAndPrint(
		`fn z(a number) { print("hi"); }`,
		t,
	)

	parseAndPrint(
		`fn z(a number) string { print("hi"); }`,
		t,
	)

	parseAndPrint(
		`fn z(a fn(number) number) string { print("hi"); } `,
		t,
	)

	parseAndPrint(
		`fn z(a fn() string) fn() { print("hi"); } `,
		t,
	)

}

func TestAnonFn(t *testing.T) {
	parseAndPrint(
		`let x = fn() { return 1; }() + 2;`,
		t,
	)
}
