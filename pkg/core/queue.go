package core

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type QueueHandler interface {
	Next() ([]ActionItem, error)
	SetActionList(a []Action)
}

type QueueCtrl struct {
	core *Core
	prio []Action
}

func NewQueueCtr(c *Core) *QueueCtrl {
	return &QueueCtrl{
		core: c,
	}
}

func (q *QueueCtrl) SetActionList(a []Action) {
	q.prio = a
}

func (q *QueueCtrl) logSkipped(a Action, reason string, keysAndValue ...interface{}) {
	if q.core.Flags.LogDebug {
		//build exec str
		var sb strings.Builder
		for _, v := range a.Exec {
			sb.WriteString(v.Typ.String())
			sb.WriteString(",")
		}
		str := sb.String()
		if len(str) > 0 {
			str = str[:len(str)-1]
		}
		items := []interface{}{
			"frame", q.core.F,
			"event", LogQueueEvent,
			"failed", true,
			"reason", reason,
			"target", a.Target,
			"exec", str,
			"raw", a.Raw,
		}
		items = append(items, keysAndValue...)
		q.core.Log.Debugw(
			"skip",
			items...,
		)
	}
}

func (q *QueueCtrl) Next() ([]ActionItem, error) {
	var r []ActionItem
	f := q.core.F
	active := q.core.Chars[q.core.ActiveChar].Key()
next:
	for i, v := range q.prio {
		char, ok := q.core.CharByName(v.Target)
		if !ok {
			continue next
		}
		//check if disabled
		if v.Disabled {
			// q.core.Log.Debugw("queue not rdy; disabled", "frame", f, "event", LogQueueEvent, "raw", v.Raw)
			q.logSkipped(v, "disabled")
			continue next
		}
		//check if still locked
		if v.ActionLock > f-v.Last && v.Last != -1 {
			// q.core.Log.Debugw("queue not rdy; on action lock", "frame", f, "event", LogQueueEvent, "raw", v.Raw)
			q.logSkipped(v, "locked", "last", v.Last, "lock_for", v.ActionLock)
			continue next
		}
		//check active char
		if v.ActiveCond != 0 {
			if v.ActiveCond != active {
				// q.core.Log.Debugw("queue not rdy; char not active", "frame", f, "event", LogQueueEvent, "active", active, "cond", v.ActiveCond, "raw", v.Raw)
				q.logSkipped(v, "inactive", "active", active, "cond", v.ActiveCond)
				continue next
			}
		}
		//check if char requested is even alive
		//check if actor is alive first, if not return 0 and call it a day
		if char.HP() <= 0 {
			// q.core.Log.Debugw("queue not rdy; char dead", "frame", f, "event", LogQueueEvent, "character", v.Target, "hp", char.HP(), "raw", v.Raw)
			q.logSkipped(v, "dead", "hp", char.HP())
			continue next
		}

		//check if we need to swap for this, and if so is swapcd = 0
		if v.Target != active {
			if q.core.SwapCD > 0 {
				// q.core.Log.Debugw("queue not rdy; swap on cd", "frame", f, "event", LogQueueEvent, "swap_cd", q.core.SwapCD, "raw", v.Raw)
				q.logSkipped(v, "swap cd", "swap_cd", q.core.SwapCD)
				continue next
			}
		}

		ready := false

		switch {
		case v.IsSeq && v.IsStrict:
			ready = true
			for _, a := range v.Exec {
				ready = ready && char.ActionReady(a.Typ, a.Param)
			}
		case v.IsSeq:
			if v.Pos >= len(v.Exec) {
				ready = false
			} else {
				ready = char.ActionReady(v.Exec[v.Pos].Typ, v.Exec[v.Pos].Param)
			}
		default:
			ready = char.ActionReady(v.Exec[0].Typ, v.Exec[0].Param)
		}

		if !ready {
			// q.core.Log.Debugw("queue not rdy; actions not rdy", "frame", f, "event", LogQueueEvent, "raw", v.Raw)
			q.logSkipped(v, "not rdy")
			continue next
		}

		//walk the tree
		if v.Conditions != nil {
			ok, err := q.evalTree(v.Conditions)
			if err != nil {
				return nil, err
			}
			if !ok {
				// q.core.Log.Debugw("queue not rdy; conditions not met", "frame", f, "event", LogQueueEvent, "condition", v.Conditions, "raw", v.Raw)
				q.logSkipped(v, "cond failed", "condition", v.Conditions)
				continue next
			}
		}

		//add this point ability is ready and we can queue
		//if active char is not current, then add swap first to queue
		if active != v.Target {
			r = append(r, ActionItem{
				Target: v.Target,
				Typ:    ActionSwap,
			})
		}

		//if it's execute once, disable it for future
		if v.Once {
			q.prio[i].Disabled = true

		}
		q.prio[i].Last = f //TODO: check this doesnt bug out since we're queuing actions

		//queue up swap lock
		if v.SwapLock > 0 {
			r = append(r, ActionItem{
				Typ:      ActionSwapLock,
				SwapLock: v.SwapLock,
			})
		}

		//queue up the abilities
		l := 1
		switch {
		case v.IsSeq && v.IsStrict:
			r = append(r, v.Exec...)
			l = len(v.Exec)
		case v.IsSeq:
			r = append(r, v.Exec[v.Pos])
			v.Pos++
		default:
			r = append(r, v.Exec[0])
		}

		//check for any cancel actions
		switch v.PostAction {
		case ActionDash:
			r = append(r, ActionItem{
				Typ: ActionDash,
			})
			l++
		case ActionJump:
			r = append(r, ActionItem{
				Typ: ActionJump,
			})
			l++
		}
		//check for any force swaps at the end
		if v.SwapTo != 0 {
			if _, ok := q.core.CharByName(v.SwapTo); ok {
				r = append(r, ActionItem{
					Target: v.SwapTo,
					Typ:    ActionSwap,
				})
				l++
			}
		}
		q.core.Log.Debugw(
			"item queued",
			"frame", f,
			"event", LogQueueEvent,
			"name", v.Name,
			"target", v.Target,
			"is seq", v.IsSeq,
			"strict", v.IsStrict,
			"exec", v.Exec,
			"once", v.Once,
			"post", v.PostAction.String(),
			"swap_to", v.SwapTo,
			"raw", v.Raw,
		)

		return r, nil
	}
	return nil, nil // no item to add
}

func (q *QueueCtrl) evalTree(node *ExprTreeNode) (bool, error) {
	//recursively evaluate tree nodes
	if node.IsLeaf {
		r, err := q.evalCond(node.Expr)
		// s.Log.Debugw("evaluating leaf node", "frame", s.F, "event", LogQueueEvent, "result", r, "node", node)
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
	// s.Log.Debugw("evaluating tree node", "frame", s.F, "event", LogQueueEvent, "left val", left, "right val", right, "node", node)
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

func (q *QueueCtrl) evalCond(c Condition) (bool, error) {

	switch c.Fields[0] {
	case ".debuff":
		return q.evalDebuff(c)
	case ".element":
		return q.evalElement(c)
	case ".cd":
		return q.evalCD(c)
	case ".energy":
		return q.evalEnergy(c)
	case ".status":
		return q.evalStatus(c)
	case ".tags":
		return q.evalTags(c)
	case ".stam":
		return q.evalStam(c)
	case ".ready":
		return q.evalAbilReady(c)
	}
	return false, nil
}

func (q *QueueCtrl) evalStam(c Condition) (bool, error) {
	return compInt(c.Op, int(q.core.Stam), c.Value), nil
}

func (q *QueueCtrl) evalAbilReady(c Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval abil: unexpected short field, expected at least 3")
	}
	cs := strings.TrimPrefix(c.Fields[2], ".")
	key := keys.CharNameToKey[cs]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, nil
	}
	a := strings.TrimPrefix(c.Fields[1], ".")
	val := c.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	ready := 0
	switch a {
	case "burst":
		if char.ActionReady(ActionBurst, nil) {
			ready = 1
		}
	case "skill":
		if char.ActionReady(ActionSkill, nil) {
			ready = 1
		}
	default:
		return false, nil
	}
	return ready == val, nil

}

func (q *QueueCtrl) evalDebuff(c Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval debuff: unexpected short field, expected at least 3")
	}
	t := strings.TrimPrefix(c.Fields[1], ".")

	val := c.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	active := 0
	d := strings.TrimPrefix(c.Fields[2], ".")
	//expecting the value to be either 0 or not 0; 0 for false

	switch t {
	case "res":
		if q.core.Combat.TargetHasResMod(d, DefaultTargetIndex) {
			active = 1
		}
	case "def":
		if q.core.Combat.TargetHasDefMod(d, DefaultTargetIndex) {
			active = 1
		}
	default:
		return false, nil
	}

	return compInt(c.Op, active, val), nil
}

func (q *QueueCtrl) evalElement(c Condition) (bool, error) {
	if len(c.Fields) < 2 {
		return false, errors.New("eval element: unexpected short field, expected at least 2")
	}
	ele := strings.TrimPrefix(c.Fields[1], ".")
	//expecting the value to be either 0 or not 0; 0 for false
	val := c.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	active := 0
	e := StringToEle(ele)
	if e == UnknownElement {
		return false, nil
	}

	if q.core.Combat.TargetHasElement(e, 0) {
		active = 1
	}
	return compInt(c.Op, active, val), nil
}

func (q *QueueCtrl) evalCD(c Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval cd: unexpected short field, expected at least 3")
	}
	//check target is valid
	name := strings.TrimPrefix(c.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval cd: invalid char in condition")
	}
	var cd int
	switch c.Fields[2] {
	case ".skill":
		cd = char.Cooldown(ActionSkill)
	case ".burst":
		cd = char.Cooldown(ActionBurst)
	default:
		return false, nil
	}
	//check vs the conditions
	return compInt(c.Op, cd, c.Value), nil
}

func (q *QueueCtrl) evalEnergy(c Condition) (bool, error) {
	if len(c.Fields) < 2 {
		return false, errors.New("eval energy: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval energy: invalid char in condition")
	}
	e := char.CurrentEnergy()
	return compFloat(c.Op, e, float64(c.Value)), nil
}

func (q *QueueCtrl) evalStatus(c Condition) (bool, error) {
	if len(c.Fields) < 2 {
		return false, errors.New("eval status: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	status := q.core.Status.Duration(name)
	// q.core.Log.Debugw("queue status check", "frame", q.core.F, "event", LogQueueEvent, "status", name, "val", status, "expected", c.Value, "op", c.Op)
	return compInt(c.Op, status, c.Value), nil

}

func (q *QueueCtrl) evalTags(c Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(c.Fields[2], ".")
	v := char.Tag(tag)
	q.core.Log.Debugw("evaluating tags", "frame", q.core.F, "event", LogQueueEvent, "char", char.CharIndex(), "targ", tag, "val", v)
	return compInt(c.Op, v, c.Value), nil
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
