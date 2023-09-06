package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func mapExprEval(n *ast.MapExpr) evalNode {
	return &mapExprEvalNode{
		root: n,
	}
}

type mapFieldWrapper struct {
	key  string
	node evalNode
}

type mapExprEvalNode struct {
	root   *ast.MapExpr
	stack  []mapFieldWrapper
	fields map[string]Obj
}

func (m *mapExprEvalNode) nextAction(env *Env) (Obj, bool, error) {
	if len(m.root.Fields) == 0 {
		return &mapval{}, true, nil
	}

	if m.fields == nil {
		m.fields = make(map[string]Obj)
		m.stack = make([]mapFieldWrapper, 0, len(m.root.Fields))
		//order really doesn't matter here
		for k, v := range m.root.Fields {
			m.stack = append(m.stack, mapFieldWrapper{
				key:  k,
				node: evalFromExpr(v),
			})
		}
	}

	for len(m.stack) > 0 {
		idx := len(m.stack) - 1
		res, done, err := m.stack[idx].node.nextAction(env)
		if err != nil {
			return nil, false, err
		}
		if done {
			m.fields[m.stack[idx].key] = res
			m.stack = m.stack[:idx]
		}
		if res.Typ() == typAction {
			return res, false, nil //done is false b/c the node is not done yet
		}
	}

	return &mapval{
		fields: m.fields,
	}, true, nil
}
