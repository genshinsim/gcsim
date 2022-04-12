package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalTree(node *core.ExprTreeNode) (bool, error) {
	//recursively evaluate tree nodes
	if node.IsLeaf {
		r, err := q.evalCond(node.Expr)
		// s.Log.Debugw("evaluating leaf node", LogQueueEvent, "result", r, "node", node)
		return r, err
	}
	//so this is a node, then we want to evalute the left and right
	//and then apply operator on both and return that
	left, err := q.evalTree(node.Left)
	if err != nil {
		return false, err
	}
	right, err := q.evalTree(node.Right)
	if err != nil {
		return false, err
	}
	// s.Log.Debugw("evaluating tree node", LogQueueEvent, "left val", left, "right val", right, "node", node)
	switch node.Op {
	case "||":
		return left || right, nil
	case "&&":
		return left && right, nil
	default:
		//if unrecognized op then return false
		return false, nil
	}

}

func (q *Queuer) evalCond(cond core.Condition) (bool, error) {

	switch cond.Fields[0] {
	case ".debuff":
		return q.evalDebuff(cond)
	case ".element":
		return q.evalElement(cond)
	case ".cd":
		return q.evalCD(cond)
	case ".energy":
		return q.evalEnergy(cond)
	case ".status":
		return q.evalStatus(cond)
	case ".tags":
		return q.evalTags(cond)
	case ".stam":
		return q.evalStam(cond)
	case ".ready":
		return q.evalAbilReady(cond)
	case ".mods":
		return q.evalCharMods(cond)
	case ".infusion":
		return q.evalInfusion(cond)
	case ".construct":
		return q.evalConstruct(cond)
	case ".normal":
		return q.evalNormal(cond)
	}
	return false, nil
}

func (q *Queuer) evalStam(cond core.Condition) (bool, error) {
	return compInt(cond.Op, int(q.core.Stam), cond.Value), nil
}

func (q *Queuer) evalAbilReady(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval abil: unexpected short field, expected at least 3")
	}
	cs := strings.TrimPrefix(cond.Fields[2], ".")
	key := core.CharNameToKey[cs]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, nil
	}
	a := strings.TrimPrefix(cond.Fields[1], ".")
	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	ready := 0
	switch a {
	case "burst":
		if char.ActionReady(core.ActionBurst, nil) {
			ready = 1
		}
	case "skill":
		if char.ActionReady(core.ActionSkill, nil) {
			ready = 1
		}
	default:
		return false, nil
	}
	return ready == val, nil

}

func compFloat(op string, a, b float64) bool {
	switch op {
	case "==":
		return a == b
	case "!=":
		return a != b
	case ">":
		return a > b
	case ">=":
		return a >= b
	case "<":
		return a < b
	case "<=":
		return a <= b
	}
	return false
}

func compInt(op string, a, b int) bool {
	switch op {
	case "==":
		return a == b
	case "!=":
		return a != b
	case ">":
		return a > b
	case ">=":
		return a >= b
	case "<":
		return a < b
	case "<=":
		return a <= b
	}
	return false
}
