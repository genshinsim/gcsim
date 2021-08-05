package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

type Queuer struct {
	s    def.Sim
	prio []def.Action
	log  *zap.SugaredLogger
}

func New(s def.Sim, prio []def.Action, log *zap.SugaredLogger) *Queuer {
	q := &Queuer{
		s:    s,
		prio: prio,
		log:  log,
	}
	return q
}

func (q *Queuer) Next(active string) ([]def.ActionItem, error) {
	var r []def.ActionItem
	f := q.s.Frame()
next:
	for i, v := range q.prio {
		char, ok := q.s.CharByName(v.Target)
		if !ok {
			continue next
		}
		//check if disabled
		if v.Disabled {
			q.log.Debugw("queue not rdy; disabled", "frame", f, "event", def.LogQueueEvent, "raw", v.Raw)
			continue next
		}
		//check if still locked
		if v.ActionLock > f-v.Last && v.Last != -1 {
			q.log.Debugw("queue not rdy; on action lock", "frame", f, "event", def.LogQueueEvent, "last", v.Last, "lock_for", v.ActionLock, "raw", v.Raw)
			continue next
		}
		//check active char
		if v.ActiveCond != "" {
			if v.ActiveCond != active {
				q.log.Debugw("queue not rdy; char not active", "frame", f, "event", def.LogQueueEvent, "active", active, "cond", v.ActiveCond, "raw", v.Raw)
				continue next
			}
		}
		//check if char requested is even alive
		//check if actor is alive first, if not return 0 and call it a day
		if char.HP() <= 0 {
			q.log.Debugw("queue not rdy; char dead", "frame", f, "event", def.LogQueueEvent, "character", v.Target, "hp", char.HP(), "raw", v.Raw)
			continue next
		}

		//check if we need to swap for this, and if so is swapcd = 0
		if v.Target != active {
			if q.s.SwapCD() > 0 {
				q.log.Debugw("queue not rdy; swap on cd", "frame", f, "event", def.LogQueueEvent, "swap_cd", q.s.SwapCD(), "raw", v.Raw)
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
			q.log.Debugw("queue not rdy; actions not rdy", "frame", f, "event", def.LogQueueEvent, "raw", v.Raw)
			continue next
		}

		//walk the tree
		if v.Conditions != nil {
			ok, err := q.evalTree(v.Conditions)
			if err != nil {
				return nil, err
			}
			if !ok {
				q.log.Debugw("queue not rdy; conditions not met", "frame", f, "event", def.LogQueueEvent, "condition", v.Conditions, "raw", v.Raw)
				continue next
			}
		}

		//add this point ability is ready and we can queue
		//if active char is not current, then add swap first to queue
		if active != v.Target {
			r = append(r, def.ActionItem{
				Target: v.Target,
				Typ:    def.ActionSwap,
			})
		}

		//if it's execute once, disable it for future
		if v.Once {
			q.prio[i].Disabled = true

		}
		q.prio[i].Last = f //TODO: check this doesnt bug out since we're queuing actions

		//queue up swap lock
		if v.SwapLock > 0 {
			r = append(r, def.ActionItem{
				Typ:      def.ActionSwapLock,
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
		case def.ActionDash:
			r = append(r, def.ActionItem{
				Typ: def.ActionDash,
			})
			l++
		case def.ActionJump:
			r = append(r, def.ActionItem{
				Typ: def.ActionJump,
			})
			l++
		}
		//check for any force swaps at the end
		if v.SwapTo != "" {
			if _, ok := q.s.CharByName(v.SwapTo); ok {
				r = append(r, def.ActionItem{
					Target: v.SwapTo,
					Typ:    def.ActionSwap,
				})
				l++
			}
		}
		q.log.Debugw(
			"item queued",
			"frame", f,
			"event", def.LogQueueEvent,
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

func (q *Queuer) evalTree(node *def.ExprTreeNode) (bool, error) {
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

func (q *Queuer) evalCond(c def.Condition) (bool, error) {

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
	}
	return false, nil
}

func (q *Queuer) evalStam(c def.Condition) (bool, error) {
	return compInt(c.Op, int(q.s.Stam()), c.Value), nil
}

func (q *Queuer) evalDebuff(c def.Condition) (bool, error) {
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
		if q.s.TargetHasResMod(d, 0) {
			active = 1
		}
	case "def":
		if q.s.TargetHasDefMod(d, 0) {
			active = 1
		}
	default:
		return false, nil
	}

	return compInt(c.Op, active, val), nil
}

func (q *Queuer) evalElement(c def.Condition) (bool, error) {
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
	e := def.StringToEle(ele)
	if e == def.UnknownElement {
		return false, nil
	}

	if q.s.TargetHasElement(e, 0) {
		active = 1
	}
	return compInt(c.Op, active, val), nil
}

func (q *Queuer) evalCD(c def.Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval cd: unexpected short field, expected at least 3")
	}
	//check target is valid
	name := strings.TrimPrefix(c.Fields[1], ".")
	char, ok := q.s.CharByName(name)
	if !ok {
		return false, errors.New("eval cd: invalid char in condition")
	}
	var cd int
	switch c.Fields[2] {
	case ".skill":
		cd = char.Cooldown(def.ActionSkill)
	case ".burst":
		cd = char.Cooldown(def.ActionBurst)
	default:
		return false, nil
	}
	//check vs the conditions
	return compInt(c.Op, cd, c.Value), nil
}

func (q *Queuer) evalEnergy(c def.Condition) (bool, error) {
	if len(c.Fields) < 2 {
		return false, errors.New("eval energy: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	char, ok := q.s.CharByName(name)
	if !ok {
		return false, errors.New("eval energy: invalid char in condition")
	}
	e := char.CurrentEnergy()
	return compFloat(c.Op, e, float64(c.Value)), nil
}

func (q *Queuer) evalStatus(c def.Condition) (bool, error) {
	if len(c.Fields) < 2 {
		return false, errors.New("eval status: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	status := q.s.Status(name)
	q.log.Debugw("queue status check", "frame", q.s.Frame(), "event", def.LogQueueEvent, "status", name, "val", status, "expected", c.Value, "op", c.Op)
	return compInt(c.Op, status, c.Value), nil

}

func (q *Queuer) evalTags(c def.Condition) (bool, error) {
	if len(c.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(c.Fields[1], ".")
	char, ok := q.s.CharByName(name)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(c.Fields[2], ".")
	v := char.Tag(tag)
	q.log.Debugw("evaluating tags", "frame", q.s.Frame(), "event", def.LogQueueEvent, "char", char.CharIndex(), "targ", tag, "val", v)
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
