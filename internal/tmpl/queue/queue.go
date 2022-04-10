package queue

import (
	"errors"
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Queuer struct {
	core       *core.Core
	pq         []core.ActionBlock
	labels     map[string]int
	prevQueued int //index of previously queued action
}

func NewQueuer(c *core.Core) *Queuer {
	return &Queuer{
		core:   c,
		labels: make(map[string]int),
	}
}

func (q *Queuer) SetActionList(a []core.ActionBlock) error {
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
	q.core.Log.NewEvent(
		"priority queued set",
		core.LogQueueEvent,
		-1,
		"pq", q.pq,
	)
	return nil
}

func (q *Queuer) Next() (next []core.Command, dropIfNotReady bool, err error) {
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
			q.core.Log.NewEvent(
				"item queued",
				core.LogQueueEvent,
				-1,
				"queued", next,
				"full", q.pq[i],
			)

			if v.Type == core.ActionBlockTypeWait {
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

func (q *Queuer) createQueueFromBlock(a core.ActionBlock) []core.Command {
	//set tracking info
	// a.NumQueued++
	// a.LastQueued = q.core.F

	var res []core.Command

	switch a.Type {
	case core.ActionBlockTypeWait:
		return q.createWaitCommand(a)
	case core.ActionBlockTypeChain:
		return q.createQueueFromChain(a)
	case core.ActionBlockTypeSequence:
		//check first if we need to swap char for this sequence
		if q.core.Chars[q.core.ActiveChar].Key() != a.SequenceChar {
			res = append(res, &core.ActionItem{
				Typ:    core.ActionSwap,
				Target: a.SequenceChar,
			})
		}
		return append(res, q.createQueueFromSequence(a)...)
	default:
		//unknown type
		return nil
	}
}

func (c *Queuer) createWaitCommand(a core.ActionBlock) []core.Command {
	//we can either wait for particles or wait for some status
	v := a.Wait
	return []core.Command{
		&v,
	}
}

func (q *Queuer) createQueueFromChain(a core.ActionBlock) []core.Command {
	var res []core.Command

	//add lock out if any
	if a.SwapLock > 0 {
		res = append(res, &core.CmdNoSwap{
			Val: a.SwapLock,
		})
	}

	active := q.core.Chars[q.core.ActiveChar].Key()
	//add up sequences for each subchain
	for i := 0; i < len(a.ChainSequences); i++ {
		//swap to this char if not currently active; only if v is a sequence command
		if a.ChainSequences[i].Type == core.ActionBlockTypeSequence && active != a.ChainSequences[i].SequenceChar {
			q.core.Log.NewEvent(
				"adding swap before sequence",
				core.LogQueueEvent,
				-1,
				"active", active,
				"next", a.ChainSequences[i].SequenceChar,
				"full", a.ChainSequences[i],
			)
			res = append(res, &core.ActionItem{
				Typ:    core.ActionSwap,
				Target: a.ChainSequences[i].SequenceChar,
			})
			active = a.ChainSequences[i].SequenceChar
		}
		//append
		res = append(res, q.createQueueFromSequence(a.ChainSequences[i])...)
	}

	//if swapto, add to end of sequence
	if a.SwapTo > core.NoChar {
		res = append(res, &core.ActionItem{
			Typ:    core.ActionSwap,
			Target: a.SwapTo,
		})
	}

	return res
}

func (q *Queuer) createQueueFromSequence(a core.ActionBlock) []core.Command {
	var res []core.Command

	//add lock out if any
	if a.SwapLock > 0 {
		res = append(res, &core.CmdNoSwap{
			Val: a.SwapLock,
		})
	}

	//check type of queue
	switch a.Type {
	case core.ActionBlockTypeSequence:
		//add abilities to the res
		for i := 0; i < len(a.Sequence); i++ {
			res = append(res, &a.Sequence[i])
		}
	case core.ActionBlockTypeWait:
		res = append(res, &a.Wait)
	case core.ActionBlockTypeResetLimit:
		res = append(res, core.CmdResetLimit{})
	}

	//if swapto, add to end of sequence
	if a.SwapTo > core.NoChar {
		res = append(res, &core.ActionItem{
			Typ:    core.ActionSwap,
			Target: a.SwapTo,
		})
	}

	return res
}

func (q *Queuer) blockUseable(a core.ActionBlock) (bool, error) {
	// wait blocks are always useable
	// chain blocks are useable if every sequence is useable
	// sequence useable if conditions are met + all abil are ready
	switch a.Type {
	case core.ActionBlockTypeWait:
		return true, nil
	case core.ActionBlockTypeChain:
		return q.chainUseable(a)
	case core.ActionBlockTypeSequence:
		return q.sequenceUseable(a)
	case core.ActionBlockTypeResetLimit:
		for i := range q.pq {
			if q.pq[i].Limit > 0 {
				q.pq[i].NumQueued = 0
			}
		}
		q.core.Log.NewEvent(
			"reset limits",
			core.LogQueueEvent,
			-1,
		)
		return true, nil
	// Add catch cases for calc mode blocks for clean error message purposes
	case core.ActionBlockTypeCalcRestart:
		return false, errors.New("invalid restart keyword detected in action priority mode list - did you mean to use sequential mode?")
	case core.ActionBlockTypeCalcWait:
		return false, errors.New("invalid wait keyword detected in action priority mode list - did you mean to use sequential mode?")
	default:
		//unknown type
		return false, errors.New("unknown action block type")
	}
}

func (q *Queuer) logSkipped(a core.ActionBlock, reason string, keysAndValue ...interface{}) {
	if q.core.Flags.LogDebug {
		//build exec str
		var sb strings.Builder
		switch a.Type {
		case core.ActionBlockTypeSequence:
			for _, v := range a.Sequence {
				sb.WriteString(v.Typ.String())
				sb.WriteString(",")
			}
		case core.ActionBlockTypeChain:
			sb.WriteString("chain,")
		case core.ActionBlockTypeWait:
			sb.WriteString("wait,")
		}

		str := sb.String()
		if len(str) > 0 {
			str = str[:len(str)-1]
		}
		q.core.Log.NewEvent(
			"skip",
			core.LogQueueEvent,
			-1,
			"failed", true,
			"reason", reason,
			"exec", str,
			"raw", a,
		).Write(keysAndValue...)

	}
}

func (q *Queuer) chainUseable(a core.ActionBlock) (bool, error) {
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
		case core.ActionBlockTypeSequence:
			ok, err := q.sequenceUseable(v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		case core.ActionBlockTypeChain:
			return false, errors.New("invalid nested chain in chain")
		}

	}
	//check conditions

	return true, nil
}

func (q *Queuer) sequenceUseable(a core.ActionBlock) (bool, error) {
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
