package queue

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type ActionBlockType int

const (
	ActionBlockTypeSequence ActionBlockType = iota
	ActionBlockTypeWait
	ActionBlockTypeChain
	ActionBlockTypeResetLimit
	ActionBlockTypeCalcWait
	ActionBlockTypeCalcRestart
)

type ActionBlock struct {
	Label string //label for this block
	Type  ActionBlockType
	//sequence is only relevant to ActionBlockTypeSequence
	Sequence     []ActionItem
	SequenceChar keys.CharKey

	ChainSequences []ActionBlock

	//conditions
	Conditions *ExprTreeNode //conditions to be met
	OnField    bool          //if true then can only use if char is on field; sequence only
	Needs      string        //previous queued action block label must match this
	Limit      int           //number of times this action block can be queued
	Timeout    int           //the action block cannot be used again for x frames

	//options
	SwapTo            keys.CharKey //character to swap to after this block
	SwapLock          int          //must stay on current char for x frames
	Try               bool         //if true then drop rest of queue if any action is not ready
	TryDropIfNotReady bool         //if false will keep trying next action; other wise drop sequence. Only if Try is set to true

	//tracking
	NumQueued  int //number of times this action block has been queued
	LastQueued int //last time this action block was queued

	//options related to wait
	Wait     CmdWait
	CalcWait CmdCalcWait
}

func (a *ActionBlock) Clone() ActionBlock {
	next := *a
	//always check fo rnil since some of these fields may not exist
	//all the slices have to be manually cloned
	if a.Sequence != nil {
		next.Sequence = make([]ActionItem, len(a.Sequence))
		for i, v := range a.Sequence {
			next.Sequence[i] = v.Clone()
		}
	}
	//clone conditions
	if a.Conditions != nil {
		next.Conditions = a.Conditions.Clone()
	}
	next.Wait = a.Wait.Clone()

	//clone chain sequence last
	if a.ChainSequences != nil {
		next.ChainSequences = make([]ActionBlock, len(a.ChainSequences))
		for i, v := range a.ChainSequences {
			next.ChainSequences[i] = v.Clone()
		}
	}

	return next
}

type ActionItem struct {
	Typ    action.Action
	Param  map[string]int
	Target keys.CharKey
}

func (a *ActionItem) Type() CommandType { return CommandTypeAction }

func (a *ActionItem) Clone() ActionItem {
	next := *a

	if a.Param != nil {
		next.Param = make(map[string]int, len(a.Param))
		for k, v := range a.Param {
			next.Param[k] = v
		}
	}

	return next
}
