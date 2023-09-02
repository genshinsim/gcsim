package eval

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalBasicBlockStmt(t *testing.T) {
	n := &ast.BlockStmt{
		List: []ast.Node{
			&ast.BinaryExpr{
				Left: &ast.NumberLit{
					IntVal: 5,
				},
				Right: &ast.NumberLit{
					IntVal: 4,
				},
				Op: ast.Token{
					Typ: ast.ItemMinus,
				},
			},
			&ast.BinaryExpr{
				Left: &ast.NumberLit{
					IntVal: 5,
				},
				Right: &ast.NumberLit{
					IntVal: 2,
				},
				Op: ast.Token{
					Typ: ast.ItemMinus,
				},
			},
		},
	}

	e := evalFromNode(n)
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.evalNext(nil)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(val)
	}
	if !done {
		t.Error("expected node to be done, got false")
	}
	v, ok := val.(*number)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	if v.ival != 3 {
		t.Errorf("expected result to be %v, got %v", 3, v.ival)
	}
}

func TestEvalBlockWithReturnStmt(t *testing.T) {
	n := &ast.BlockStmt{
		List: []ast.Node{
			&ast.NumberLit{
				IntVal: -1,
			},
			&ast.StringLit{
				Value: "should never get here",
			},
			&ast.ReturnStmt{
				Val: &ast.BinaryExpr{
					Left: &ast.NumberLit{
						IntVal: 5,
					},
					Right: &ast.NumberLit{
						IntVal: 4,
					},
					Op: ast.Token{
						Typ: ast.ItemMinus,
					},
				},
			},
			&ast.StringLit{
				Value: "should never get here",
			},
		},
	}

	e := evalFromNode(n)
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.evalNext(nil)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(val)
	}
	if !done {
		t.Error("expected node to be done, got false")
	}
	v, ok := val.(*retval)
	if !ok {
		t.Errorf("res is not a retval, got %v", val.Typ())
	}
	amt, ok := v.res.(*number)
	if !ok {
		t.Errorf("retval is not a number, got %v", v.res.Typ())
	}
	if amt.ival != 1 {
		t.Errorf("expected result to be %v, got %v", 1, amt.ival)
	}
}
