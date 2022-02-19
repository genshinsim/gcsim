package core

import (
	"reflect"
	"testing"
)

func TestCloneExprTree(t *testing.T) {
	tree := ExprTreeNode{
		Left: &ExprTreeNode{
			IsLeaf: true,
			Expr: Condition{
				Fields: []string{".test"},
				Op:     "=",
				Value:  1,
			},
		},
		Right: &ExprTreeNode{
			IsLeaf: true,
			Expr: Condition{
				Fields: []string{".test"},
				Op:     "=",
				Value:  1,
			},
		},
		Op: "||",
	}

	clone := tree.Clone()

	if !reflect.DeepEqual(tree, *clone) {
		t.Fail()
	}

	//check they are not same reference
	if clone.Left == tree.Left {
		t.Fail()
	}

	if clone.Right == tree.Right {
		t.Fail()
	}

	if &clone.Left.Expr.Fields == &tree.Left.Expr.Fields {
		t.Fail()
	}

	if &clone.Right.Expr.Fields == &tree.Right.Expr.Fields {
		t.Fail()
	}

}
