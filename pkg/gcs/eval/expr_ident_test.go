package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalBasicIdent(t *testing.T) {
	n := &ast.Ident{
		Value: "print",
	}
	env := NewEnv(nil)
	var o Obj = &strval{str: "test"}
	env.put("print", &o)

	val, err := runEvalReturnResWhenDone(evalFromNode(n), env)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v, ok := val.(*strval)
	if !ok {
		t.Errorf("res is not bfuncval, got %v", val.Typ())
	}
	if v.str != "test" {
		t.Errorf("expected string to be test, got %v", v.str)
	}
}
