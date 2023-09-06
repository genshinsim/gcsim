package eval

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (e *Evaluator) initSysFuncs(env *Env) {
	// std funcs
	e.addSysFunc("print", e.print, env)
	e.addSysFunc("execute_action", e.executeAction, env)
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

func (e *Evaluator) executeAction(args []Obj, env *Env) (Obj, error) {
	//execute_action(char, action, params)
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of params for execute_action, expected 3 got %v", len(args))
	}

	// char
	charId := args[0]
	if charId.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument char should evaluate to a number, got %v", charId.Inspect())
	}
	char := charId.(*number)

	// action
	actionId := args[1]
	if actionId.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument action should evaluate to a number, got %v", actionId.Inspect())
	}
	ac := actionId.(*number)

	// params
	paramsMap := args[2]
	if paramsMap.Typ() != typMap {
		return nil, fmt.Errorf("execute_action argument params should evaluate to a map, got %v", paramsMap.Inspect())
	}

	p := paramsMap.(*mapval)
	params := make(map[string]int)
	for k, v := range p.fields {
		if v.Typ() != typNum {
			return nil, fmt.Errorf("map params should evaluate to a number, got %v", v.Inspect())
		}
		params[k] = int(v.(*number).ival)
	}

	charKey := keys.Char(char.ival)
	actionKey := action.Action(ac.ival)
	return &actionval{
		char:   charKey,
		action: actionKey,
		param:  params,
	}, nil
}
