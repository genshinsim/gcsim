package parser

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestErrorInvalidWeaponName(t *testing.T) {
	input := `xingqiu add weapon="notarealweapon" refine=1 lvl=1/1;`
	file := ast.NewFile()
	p := New(file, input)
	_, _, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid weapon name")
	}
	t.Logf("got expected error: %v", err)
}

func TestErrorDuplicateFnParam(t *testing.T) {
	input := `fn f(a, a) { }`
	file := ast.NewFile()
	p := New(file, input)
	_, _, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for duplicate fn param")
	}
	t.Logf("got expected error: %v", err)
}

func TestErrorNonNumberMapParam(t *testing.T) {
	input := `xingqiu add weapon="harbingerofdawn" refine=1 lvl=1/1 +params=[x=y];`
	file := ast.NewFile()
	p := New(file, input)
	_, _, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for non-number map value in add params")
	}
	t.Logf("got expected error: %v", err)
}

func TestErrorSyntax(t *testing.T) {
	input := `let x = ;`
	file := ast.NewFile()
	p := New(file, input)
	_, _, err := p.Parse()
	if err == nil {
		t.Fatal("expected syntax error")
	}
	t.Logf("got expected error: %v", err)
}

func TestErrorActionStartLine(t *testing.T) {
	input := `xingqiu attack; skill`
	file := ast.NewFile()
	p := New(file, input)
	_, _, err := p.Parse()
	if err == nil {
		t.Fatal("expected error: line starting with action key")
	}
	t.Logf("got expected error: %v", err)
}
