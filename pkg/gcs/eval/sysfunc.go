package eval

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/conditional"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (e *Evaluator) initSysFuncs() {
	// std funcs
	e.addSysFunc("print", e.print)
	e.addSysFunc("execute_action", e.executeAction)
	e.addSysFunc("evaluate_conditional", e.evaluateConditions)
}

func (e *Evaluator) addSysFunc(name string, f systemFunc) {
	var obj Obj = &bfuncval{
		Body: f,
	}
	e.env.put(name, &obj)
}

func (e *Evaluator) print(args []Obj) (Obj, error) {
	// concat all args
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

func (e *Evaluator) evaluateConditions(args []Obj) (Obj, error) {
	// expecting args to be all strings
	var vals []string
	for _, v := range args {
		str, ok := v.(*strval)
		if !ok {
			return nil, fmt.Errorf("system error; expecting str for conditional args, got %v", v.Typ())
		}
		vals = append(vals, str.str)
	}
	r, err := conditional.Eval(e.Core, vals)
	if err != nil {
		return nil, err
	}

	num := &number{}
	switch v := r.(type) {
	case bool:
		if v {
			num.ival = 1
		}
	case int:
		num.ival = int64(v)
	case int64:
		num.ival = v
	case float64:
		num.fval = v
		num.isFloat = true
	default:
		return nil, fmt.Errorf("field condition '.%v' does not evaluate to a number, got %v", strings.Join(vals, "."), v)
	}
	return num, nil
}

func (e *Evaluator) executeAction(args []Obj) (Obj, error) {
	// execute_action(char, action, params)
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
