package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type QueueHandler interface {
	//returns a sequence of 1 or more commands to execute,
	//whether or not to drop sequence if any is not ready, and any error
	Next() (queue []Command, dropIfFailed bool, err error)
	SetActionList(pq []ActionBlock) error
}

type Queuer struct {
	core       *Core
	pq         []ActionBlock
	labels     map[string]int
	prevQueued int //index of previously queued action
}

func NewQueuer(c *Core) *Queuer {
	return &Queuer{
		core:   c,
		labels: make(map[string]int),
	}
}

func (q *Queuer) SetActionList(a []ActionBlock) error {
	//set labels
	for i, v := range a {
		if v.Label == "" {
			continue
		}
		if _, ok := q.labels[v.Label]; ok {
			return fmt.Errorf("duplicated label in action list: %v", v.Label)
		}
		q.labels[v.Label] = i
	}
	q.pq = a
	q.core.Log.Debugw(
		"priority queued set",
		"frame", q.core.F,
		"event", LogQueueEvent,
		"pq", q.pq,
	)
	return nil
}

func (q *Queuer) Next() (next []Command, dropIfNotReady bool, err error) {
	// from the action block we need to build the command list
	var ok bool
	for i, v := range q.pq {
		//find the first item on prior queue that's useable
		ok, err = q.blockUseable(v)
		if err != nil {
			return
		}
		if ok {
			q.pq[i].NumQueued++
			q.pq[i].LastQueued = q.core.F
			next = q.createQueueFromBlock(v)
			q.core.Log.Debugw(
				"item queued",
				"frame", q.core.F,
				"event", LogQueueEvent,
				"full", q.pq[i],
				"queued", next,
			)

			if v.Type == ActionBlockTypeWait {
				return
			}
			//check if /try is set if this is a sequence or chain
			if v.Try {
				dropIfNotReady = v.TryDropIfNotReady
			}
			return
		}
	}
	//if we hit here then that means no action is ready
	//we should log this
	return
}

func (q *Queuer) createQueueFromBlock(a ActionBlock) []Command {
	//set tracking info
	// a.NumQueued++
	// a.LastQueued = q.core.F

	var res []Command

	switch a.Type {
	case ActionBlockTypeWait:
		return q.createWaitCommand(a)
	case ActionBlockTypeChain:
		return q.createQueueFromChain(a)
	case ActionBlockTypeSequence:
		//check first if we need to swap char for this sequence
		if q.core.Chars[q.core.ActiveChar].Key() != a.SequenceChar {
			res = append(res, &ActionItem{
				Typ:    ActionSwap,
				Target: a.SequenceChar,
			})
		}
		return append(res, q.createQueueFromSequence(a)...)
	default:
		//unknown type
		return nil
	}
}

func (c *Queuer) createWaitCommand(a ActionBlock) []Command {
	//we can either wait for particles or wait for some status
	v := a.Wait
	return []Command{
		&v,
	}
}

func (q *Queuer) createQueueFromChain(a ActionBlock) []Command {
	var res []Command

	active := q.core.Chars[q.core.ActiveChar].Key()
	//add up sequences for each subchain
	for i := 0; i < len(a.ChainSequences); i++ {
		//swap to this char if not currently active; only if v is a sequence command
		if a.ChainSequences[i].Type == ActionBlockTypeSequence && active != a.ChainSequences[i].SequenceChar {
			q.core.Log.Debugw(
				"adding swap before sequence",
				"frame", q.core.F,
				"event", LogQueueEvent,
				"active", active,
				"next", a.ChainSequences[i].SequenceChar,
				"full", a.ChainSequences[i],
			)
			res = append(res, &ActionItem{
				Typ:    ActionSwap,
				Target: a.ChainSequences[i].SequenceChar,
			})
		}
		//append
		res = append(res, q.createQueueFromSequence(a.ChainSequences[i])...)
	}

	return res
}

func (q *Queuer) createQueueFromSequence(a ActionBlock) []Command {
	var res []Command

	//add lock out if any
	if a.SwapLock > 0 {
		res = append(res, &CmdNoSwap{
			Val: a.SwapLock,
		})
	}

	//check type of queue
	switch a.Type {
	case ActionBlockTypeSequence:
		//add abilities to the res
		for i := 0; i < len(a.Sequence); i++ {
			res = append(res, &a.Sequence[i])
		}
	case ActionBlockTypeWait:
		res = append(res, &a.Wait)
	case ActionBlockTypeResetLimit:
		res = append(res, CmdResetLimit{})
	}

	//if swapto, add to end of sequence
	if a.SwapTo > keys.NoChar {
		res = append(res, &ActionItem{
			Typ:    ActionSwap,
			Target: a.SwapTo,
		})
	}

	return res
}

func (q *Queuer) blockUseable(a ActionBlock) (bool, error) {
	// wait blocks are always useable
	// chain blocks are useable if every sequence is useable
	// sequence useable if conditions are met + all abil are ready
	switch a.Type {
	case ActionBlockTypeWait:
		return true, nil
	case ActionBlockTypeChain:
		return q.chainUseable(a)
	case ActionBlockTypeSequence:
		return q.sequenceUseable(a)
	case ActionBlockTypeResetLimit:
		for i := range q.pq {
			if q.pq[i].Limit > 0 {
				q.pq[i].NumQueued = 0
			}
		}
		return true, nil
	default:
		//unknown type
		return false, errors.New("unknown action block type")
	}
}

func (q *Queuer) logSkipped(a ActionBlock, reason string, keysAndValue ...interface{}) {
	if q.core.Flags.LogDebug {
		//build exec str
		var sb strings.Builder
		switch a.Type {
		case ActionBlockTypeSequence:
			for _, v := range a.Sequence {
				sb.WriteString(v.Typ.String())
				sb.WriteString(",")
			}
		case ActionBlockTypeChain:
			sb.WriteString("chain,")
		case ActionBlockTypeWait:
			sb.WriteString("wait,")
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
			"exec", str,
			"raw", a,
		}
		items = append(items, keysAndValue...)
		q.core.Log.Debugw(
			"skip",
			items...,
		)
	}
}

func (q *Queuer) chainUseable(a ActionBlock) (bool, error) {
	// a chain is useable if by all sequences in it are useable
	if len(a.ChainSequences) == 0 {
		q.logSkipped(a, "invalid chain - length 0")
		return false, nil
	}
	//check easy stuff first
	//can't exceed limit
	if a.Limit > 0 && a.NumQueued >= a.Limit {
		q.logSkipped(a, "over limit", "limit", a.Limit, "count", a.NumQueued)
		return false, nil
	}
	//can't be timed out
	if a.Timeout > 0 && q.core.F-a.LastQueued < a.Timeout {
		q.logSkipped(a, "still in timeout", "timeout", a.Timeout)
		return false, nil
	}
	//check needs
	if a.Needs != "" {
		needs, ok := q.labels[a.Needs]
		if !ok {
			q.logSkipped(a, "need does not exist", "needs", a.Needs)
			return false, nil
		}
		if needs != q.prevQueued {
			q.logSkipped(a, "needs not last executed", "needs", a.Needs, "needs_index", needs, "prev_index", q.prevQueued)
			return false, nil
		}
	}
	//check the tree
	if a.Conditions != nil {
		ok, err := q.evalTree(a.Conditions)
		if err != nil {
			return false, err
		}
		if !ok {
			q.logSkipped(a, "conditions not met")
			return false, nil
		}
	}

	// if try is set, only check the first action of the first sequence
	if a.Try {
		return q.sequenceUseable(a.ChainSequences[0])
	}
	//otherwise check every sequence
	for _, v := range a.ChainSequences {
		//check type, either sequence or wait or reset
		//no nested chains!!
		switch v.Type {
		case ActionBlockTypeSequence:
			ok, err := q.sequenceUseable(v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		case ActionBlockTypeChain:
			return false, errors.New("invalid nested chain in chain")
		}

	}
	//check conditions

	return true, nil
}

func (q *Queuer) sequenceUseable(a ActionBlock) (bool, error) {
	/**
	for a sequence to be useable we need to check
		- active char
		- onfield
		- need label
		- limit count
		- timeout
		- conditions
		- abil ready and/or with try
	**/
	//forget it if sequence is blank for whatever reason
	if len(a.Sequence) == 0 {
		q.logSkipped(a, "invalid sequence - length 0")
		return false, nil
	}
	//check easy stuff first
	//can't exceed limit
	if a.Limit > 0 && a.NumQueued >= a.Limit {
		q.logSkipped(a, "over limit", "limit", a.Limit, "count", a.NumQueued)
		return false, nil
	}
	//can't be timed out
	if a.Timeout > 0 && q.core.F-a.LastQueued < a.Timeout {
		q.logSkipped(a, "still in timeout", "timeout", a.Timeout)
		return false, nil
	}
	//check needs
	if a.Needs != "" {
		needs, ok := q.labels[a.Needs]
		if !ok {
			q.logSkipped(a, "need does not exist", "needs", a.Needs)
			return false, nil
		}
		if needs != q.prevQueued {
			q.logSkipped(a, "needs not last executed", "needs", a.Needs, "needs_index", needs, "prev_index", q.prevQueued)
			return false, nil
		}
	}

	//make sure sequence refers to valid char just in case
	charPos, ok := q.core.CharPos[a.SequenceChar]
	if !ok {
		return false, errors.New("invalid character in action list " + a.SequenceChar.String())
	}
	//check if swap required, and if so check to make sure swapcd ==0
	if q.core.ActiveChar != charPos {
		//if we need to be on field then forget it
		if a.OnField {
			q.logSkipped(a, "not on field")
			return false, nil
		}
		//other wise check swap
		if q.core.SwapCD > 0 {
			q.logSkipped(a, "swap on cd", "active", q.core.ActiveChar, "charPos", charPos)
			return false, nil
		}

	}

	char := q.core.Chars[charPos]
	//make sure char is alive
	if char.HP() <= 0 {
		q.logSkipped(a, "char is dead")
		return false, nil
	}

	//check the tree
	if a.Conditions != nil {
		ok, err := q.evalTree(a.Conditions)
		if err != nil {
			return false, err
		}
		if !ok {
			q.logSkipped(a, "conditions not met")
			return false, nil
		}
	}

	//finally check if abilities are ready
	//if try is set the only the first ability has to be useable
	if a.Try {
		rdy := char.ActionReady(a.Sequence[0].Typ, a.Sequence[0].Param)
		if !rdy {
			q.logSkipped(a, "first action not ready in try")
		}
		return rdy, nil
	}

	//check each ability now
	for _, s := range a.Sequence {
		if !char.ActionReady(s.Typ, s.Param) {
			q.logSkipped(a, "action not ready", "failed_at", s.Typ.String())
			return false, nil
		}
	}

	//at this point all checks passed
	return true, nil
}

func (q *Queuer) evalTree(node *ExprTreeNode) (bool, error) {
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

func (q *Queuer) evalCond(cond Condition) (bool, error) {

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
	}
	return false, nil
}

func (q *Queuer) evalStam(cond Condition) (bool, error) {
	return compInt(cond.Op, int(q.core.Stam), cond.Value), nil
}

func (q *Queuer) evalAbilReady(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval abil: unexpected short field, expected at least 3")
	}
	cs := strings.TrimPrefix(cond.Fields[2], ".")
	key := keys.CharNameToKey[cs]
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

func (q *Queuer) evalDebuff(cond Condition) (bool, error) {
	//.debuff.res.1.name
	if len(cond.Fields) < 4 {
		return false, errors.New("eval debuff: unexpected short field, expected at least 3")
	}
	typ := strings.TrimPrefix(cond.Fields[1], ".")
	trg := strings.TrimPrefix(cond.Fields[2], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		return false, errors.New("eval debuff: expected int for target, got " + trg)
	}

	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	active := 0
	d := strings.TrimPrefix(cond.Fields[3], ".")
	//expecting the value to be either 0 or not 0; 0 for false

	switch typ {
	case "res":
		if q.core.Combat.TargetHasResMod(d, int(tid)) {
			active = 1
		}
	case "def":
		if q.core.Combat.TargetHasDefMod(d, int(tid)) {
			active = 1
		}
	default:
		return false, nil
	}

	return compInt(cond.Op, active, val), nil
}

func (q *Queuer) evalElement(cond Condition) (bool, error) {
	//.element.1.pyro
	if len(cond.Fields) < 3 {
		return false, errors.New("eval element: unexpected short field, expected at least 2")
	}
	trg := strings.TrimPrefix(cond.Fields[1], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		return false, errors.New("eval element: expected int for target, got " + trg)
	}

	ele := strings.TrimPrefix(cond.Fields[2], ".")
	//expecting the value to be either 0 or not 0; 0 for false
	val := cond.Value
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

	if q.core.Combat.TargetHasElement(e, int(tid)) {
		active = 1
	}
	return compInt(cond.Op, active, val), nil
}

func (q *Queuer) evalCD(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval cd: unexpected short field, expected at least 3")
	}
	//check target is valid
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval cd: invalid char in condition")
	}
	var cd int
	switch cond.Fields[2] {
	case ".skill":
		cd = char.Cooldown(ActionSkill)
	case ".burst":
		cd = char.Cooldown(ActionBurst)
	default:
		return false, nil
	}
	//check vs the conditions
	return compInt(cond.Op, cd, cond.Value), nil
}

func (q *Queuer) evalEnergy(cond Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval energy: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval energy: invalid char in condition")
	}
	e := char.CurrentEnergy()
	return compFloat(cond.Op, e, float64(cond.Value)), nil
}

func (q *Queuer) evalStatus(cond Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval status: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	status := q.core.Status.Duration(name)
	// q.core.Log.Debugw("queue status check", "frame", q.core.F, "event", LogQueueEvent, "status", name, "val", status, "expected", c.Value, "op", c.Op)
	return compInt(cond.Op, status, cond.Value), nil

}

func (q *Queuer) evalTags(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(cond.Fields[2], ".")
	v := char.Tag(tag)
	q.core.Log.Debugw("evaluating tags", "frame", q.core.F, "event", LogQueueEvent, "char", char.CharIndex(), "targ", tag, "val", v)
	return compInt(cond.Op, v, cond.Value), nil
}

func (q *Queuer) evalCharMods(cond Condition) (bool, error) {
	//.mods.bennett.buff==1
	if len(cond.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(cond.Fields[2], ".")
	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	q.core.Log.Debugw("evaluating mods", "frame", q.core.F, "event", LogQueueEvent, "char", char.CharIndex(), "mod", tag)
	return char.ModIsActive(tag) == (val == 1), nil
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
