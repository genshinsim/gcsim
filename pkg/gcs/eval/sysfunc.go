package eval

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (e *Evaluator) initSysFuncs(env *Env) {
	// std funcs
	e.addSysFunc("print", e.print, env)
}

func (e *Evaluator) addSysFunc(name string, f systemFunc, env *Env) {
	var obj Obj = &bfuncval{
		Body: f,
		Env:  NewEnv(env),
	}
	env.varMap[name] = &obj
}

func (e *Evaluator) print(args []Obj, env *Env) (Obj, error) {
	//concat all args
	var sb strings.Builder
	for _, v := range args {
		sb.WriteString(v.Inspect())
	}
	if e.Core != nil {
		e.Core.Log.NewEvent(sb.String(), glog.LogUserEvent, -1)
	} else {
		fmt.Println(sb.String())
	}
	return &null{}, nil
}
