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

func (c *Queuer) SetActionList(a []ActionBlock) error {
	//set labels
	for i, v := range a {
		if _, ok := c.labels[v.Label]; ok {
			return fmt.Errorf("duplicated label in action list: %v", v.Label)
		}
		c.labels[v.Label] = i
	}
	c.pq = a
	return nil
}

func (c *Queuer) Next() (next []Command, dropIfNotReady bool, err error) {
	// from the action block we need to build the command list
	var ok bool
	for _, v := range c.pq {
		//find the first item on prior queue that's useable
		ok, err = c.blockUseable(v)
		if err != nil {
			return
		}
		if ok {
			next = c.createQueueFromBlock(v)
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

func (c *Queuer) createQueueFromBlock(a ActionBlock) []Command {
	//set tracking info
	a.NumQueued++
	a.LastQueued = c.core.F

	var res []Command

	switch a.Type {
	case ActionBlockTypeWait:
		return c.createWaitCommand(a)
	case ActionBlockTypeChain:
		return c.createQueueFromChain(a)
	case ActionBlockTypeSequence:
		//check first if we need to swap char for this sequence
		if c.core.Chars[c.core.ActiveChar].Key() != a.SequenceChar {
			res = append(res, &ActionItem{
				Typ:    ActionSwap,
				Target: a.SequenceChar,
			})
		}
		return append(res, c.createQueueFromSequence(a)...)
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

func (c *Queuer) createQueueFromChain(a ActionBlock) []Command {
	var res []Command

	active := c.core.Chars[c.core.ActiveChar].Key()
	//add up sequences for each subchain
	for _, v := range a.ChainSequences {
		//swap to this char if not currently active
		if active != v.SequenceChar {
			res = append(res, &ActionItem{
				Typ:    ActionSwap,
				Target: v.SequenceChar,
			})
		}
		//append
		res = append(res, c.createQueueFromSequence(a)...)
	}

	return res
}

func (c *Queuer) createQueueFromSequence(a ActionBlock) []Command {
	var res []Command

	//add lock out if any
	if a.SwapLock > 0 {
		res = append(res, &CmdNoSwap{
			Val: a.SwapLock,
		})
	}

	//add abilities to the res
	for i := 0; i < len(a.Sequence); i++ {
		res = append(res, &a.Sequence[i])
	}
	// for _, v := range a.Sequence {
	// 	res = append(res, &v)
	// }

	//if swapto, add to end of sequence
	if a.SwapTo > keys.NoChar {
		res = append(res, &ActionItem{
			Typ:    ActionSwap,
			Target: a.SwapTo,
		})
	}

	return res
}

func (c *Queuer) blockUseable(a ActionBlock) (bool, error) {
	// wait blocks are always useable
	// chain blocks are useable if every sequence is useable
	// sequence useable if conditions are met + all abil are ready
	switch a.Type {
	case ActionBlockTypeWait:
		return true, nil
	case ActionBlockTypeChain:
		return c.chainUseable(a)
	case ActionBlockTypeSequence:
		return c.sequenceUseable(a)
	default:
		//unknown type
		return false, errors.New("unknown action block type")
	}
}

func (c *Queuer) chainUseable(a ActionBlock) (bool, error) {
	// a chain is useable if by all sequences in it are useable
	if len(a.ChainSequences) == 0 {
		return false, nil
	}
	// if try is set, only check the first action of the first sequence
	if a.Try {
		return c.sequenceUseable(a.ChainSequences[0])
	}
	//otherwise check every sequence
	for _, v := range a.ChainSequences {
		ok, err := c.sequenceUseable(v)
		if err != nil {
			return false, err
		}
		if !ok {
			return ok, nil
		}
	}

	return true, nil
}

func (c *Queuer) sequenceUseable(a ActionBlock) (bool, error) {
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
		return false, nil
	}
	//check easy stuff first
	//can't exceed limit
	if a.NumQueued >= a.Limit {
		return false, nil
	}
	//can't be timed out
	if c.core.F-a.LastQueued < a.Timeout {
		return false, nil
	}
	//check needs
	needs, ok := c.labels[a.Label]
	if !ok {
		return false, nil
	}
	if needs != c.prevQueued {
		return false, nil
	}

	//make sure sequence refers to valid char just in case
	charPos, ok := c.core.CharPos[a.SequenceChar]
	if !ok {
		return false, errors.New("invalid character in action list " + a.SequenceChar.String())
	}
	//check if swap required, and if so check to make sure swapcd ==0
	if c.core.ActiveChar != charPos {
		//if we need to be on field then forget it
		if a.OnField {
			return false, nil
		}
		//other wise check swap
		if c.core.SwapCD > 0 {
			return false, nil
		}

	}

	char := c.core.Chars[charPos]
	//make sure char is alive
	if char.HP() <= 0 {
		return false, nil
	}

	//check the tree
	if a.Conditions != nil {
		ok, err := c.evalTree(a.Conditions)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	//finally check if abilities are ready
	//if try is set the only the first ability has to be useable
	if a.Try {
		return char.ActionReady(a.Sequence[0].Typ, a.Sequence[0].Param), nil
	}

	//check each ability now
	for _, s := range a.Sequence {
		if !char.ActionReady(s.Typ, s.Param) {
			return false, nil
		}
	}

	//at this point all checks passed
	return true, nil
}

func (c *Queuer) evalTree(node *ExprTreeNode) (bool, error) {
	//recursively evaluate tree nodes
	if node.IsLeaf {
		r, err := c.evalCond(node.Expr)
		// s.Log.Debugw("evaluating leaf node", "frame", s.F, "event", LogQueueEvent, "result", r, "node", node)
		return r, err
	}
	//so this is a node, then we want to evalute the left and right
	//and then apply operator on both and return that
	left, err := c.evalTree(node.Left)
	if err != nil {
		return false, err
	}
	right, err := c.evalTree(node.Right)
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

func (c *Queuer) evalCond(cond Condition) (bool, error) {

	switch cond.Fields[0] {
	case ".debuff":
		return c.evalDebuff(cond)
	case ".element":
		return c.evalElement(cond)
	case ".cd":
		return c.evalCD(cond)
	case ".energy":
		return c.evalEnergy(cond)
	case ".status":
		return c.evalStatus(cond)
	case ".tags":
		return c.evalTags(cond)
	case ".stam":
		return c.evalStam(cond)
	case ".ready":
		return c.evalAbilReady(cond)
	case ".mods":
		return c.evalCharMods(cond)
	}
	return false, nil
}

func (c *Queuer) evalStam(cond Condition) (bool, error) {
	return compInt(cond.Op, int(c.core.Stam), cond.Value), nil
}

func (c *Queuer) evalAbilReady(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval abil: unexpected short field, expected at least 3")
	}
	cs := strings.TrimPrefix(cond.Fields[2], ".")
	key := keys.CharNameToKey[cs]
	char, ok := c.core.CharByName(key)
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

func (c *Queuer) evalDebuff(cond Condition) (bool, error) {
	//.debuff.res.1.name
	if len(cond.Fields) < 4 {
		return false, errors.New("eval debuff: unexpected short field, expected at least 3")
	}
	typ := strings.TrimPrefix(cond.Fields[1], ".")
	trg := strings.TrimPrefix(cond.Fields[2], ".")
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
		if c.core.Combat.TargetHasResMod(d, int(tid)) {
			active = 1
		}
	case "def":
		if c.core.Combat.TargetHasDefMod(d, int(tid)) {
			active = 1
		}
	default:
		return false, nil
	}

	return compInt(cond.Op, active, val), nil
}

func (c *Queuer) evalElement(cond Condition) (bool, error) {
	//.element.1.pyro
	if len(cond.Fields) < 3 {
		return false, errors.New("eval element: unexpected short field, expected at least 2")
	}
	trg := strings.TrimPrefix(cond.Fields[1], ".")
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

	if c.core.Combat.TargetHasElement(e, int(tid)) {
		active = 1
	}
	return compInt(cond.Op, active, val), nil
}

func (c *Queuer) evalCD(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval cd: unexpected short field, expected at least 3")
	}
	//check target is valid
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := c.core.CharByName(key)
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

func (c *Queuer) evalEnergy(cond Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval energy: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := c.core.CharByName(key)
	if !ok {
		return false, errors.New("eval energy: invalid char in condition")
	}
	e := char.CurrentEnergy()
	return compFloat(cond.Op, e, float64(cond.Value)), nil
}

func (c *Queuer) evalStatus(cond Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval status: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	status := c.core.Status.Duration(name)
	// q.core.Log.Debugw("queue status check", "frame", q.core.F, "event", LogQueueEvent, "status", name, "val", status, "expected", c.Value, "op", c.Op)
	return compInt(cond.Op, status, cond.Value), nil

}

func (c *Queuer) evalTags(cond Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := c.core.CharByName(key)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(cond.Fields[2], ".")
	v := char.Tag(tag)
	c.core.Log.Debugw("evaluating tags", "frame", c.core.F, "event", LogQueueEvent, "char", char.CharIndex(), "targ", tag, "val", v)
	return compInt(cond.Op, v, cond.Value), nil
}

func (c *Queuer) evalCharMods(cond Condition) (bool, error) {
	//.mods.bennett.buff==1
	if len(cond.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := keys.CharNameToKey[name]
	char, ok := c.core.CharByName(key)
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
	c.core.Log.Debugw("evaluating mods", "frame", c.core.F, "event", LogQueueEvent, "char", char.CharIndex(), "mod", tag)
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
