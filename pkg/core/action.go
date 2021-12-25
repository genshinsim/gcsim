package core

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type CommandHandler interface {
	Exec(n Command) (frames int, done bool, err error) //return frames, if executed, any errors
}

type CommandType int

const (
	CommandTypeAction CommandType = iota
	CommandTypeWait
	CommandTypeNoSwap
)

//Command is what gets executed by the sim.
type Command interface {
	Type() CommandType
}

type CmdWaitType int

const (
	CmdWaitTypeTimed CmdWaitType = iota
	CmdWaitTypeParticle
)

type CmdWait struct {
	For    CmdWaitType
	Max    int //cannot be 0 if type is timed
	Fields string
	Val    int
	Op     string
}

func (c *CmdWait) Type() CommandType { return CommandTypeWait }

type CmdNoSwap struct {
	Val int
}

func (c *CmdNoSwap) Type() CommandType { return CommandTypeNoSwap }

type ActionBlockType int

const (
	ActionBlockTypeSequence ActionBlockType = iota
	ActionBlockTypeWait
	ActionBlockTypeChain
)

type ActionBlock struct {
	Label string //label for this block
	Type  ActionBlockType
	//sequence is only relevant to ActionBlockTypeSequence
	Sequence     []ActionItem
	SequenceChar keys.Char

	ChainSequences []ActionBlock

	//conditions
	Conditions *ExprTreeNode //conditions to be met
	OnField    bool          //if true then can only use if char is on field; sequence only
	Needs      string        //previous queued action block label must match this
	Limit      int           //number of times this action block can be queued
	Timeout    int           //the action block cannot be used again for x frames

	//options
	SwapTo            keys.Char //character to swap to after this block
	SwapLock          int       //must stay on current char for x frames
	Try               bool      //if true then drop rest of queue if any action is not ready
	TryDropIfNotReady bool      //if false will keep trying next action; other wise drop sequence. Only if Try is set to true

	//tracking
	NumQueued int //number of times this action block has been queued
	LastUsed  int //last time this action block was queued

	//options related to wait
	Wait CmdWait
}

type ActionItem struct {
	Typ    ActionType
	Param  map[string]int
	Target keys.Char
}

func (a *ActionItem) Type() CommandType { return CommandTypeAction }

type ActionType int

const (
	InvalidAction ActionType = iota
	ActionSkill
	ActionBurst
	ActionAttack
	ActionCharge
	ActionHighPlunge
	ActionLowPlunge
	ActionAim
	ActionCancellable // delim cancellable action
	ActionDash
	ActionJump
	ActionSwap
	EndActionType
)

var astr = []string{
	"invalid",
	"skill",
	"burst",
	"attack",
	"charge",
	"high_plunge",
	"low_plunge",
	"aim",
	"",
	"dash",
	"jump",
}

func (a ActionType) String() string {
	return astr[a]
}

type ExprTreeNode struct {
	Left   *ExprTreeNode
	Right  *ExprTreeNode
	IsLeaf bool
	Op     string //&& || ( )
	Expr   Condition
}

type Condition struct {
	Fields []string
	Op     string
	Value  int
}

func (c Condition) String() {
	var sb strings.Builder
	for _, v := range c.Fields {
		sb.WriteString(v)
	}
	sb.WriteString(c.Op)
}

type ActionCtrl struct {
	core *Core
}

func NewActionCtrl(c *Core) *ActionCtrl {
	return &ActionCtrl{
		core: c,
	}
}

func (a *ActionCtrl) Exec(n Command) (int, bool, error) {
	switch v := n.(type) {
	case *ActionItem:
		return a.execAction(v)
	case *CmdWait:
		return a.execWait(v)
	case *CmdNoSwap:
		return a.execNoSwap(v)
	}
	return 0, false, errors.New("unrecognized command")
}

func (a *ActionCtrl) execWait(n *CmdWait) (int, bool, error) {
	//if condition not met, skip this frame

	return 0, true, nil
}

func (a *ActionCtrl) execNoSwap(n *CmdNoSwap) (int, bool, error) {
	a.core.SwapCD += n.Val
	a.core.Log.Debugw(
		"locked swap",
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"dur", n.Val,
		"cd", a.core.SwapCD,
	)
	return 0, true, nil
}

func (a *ActionCtrl) execAction(n *ActionItem) (int, bool, error) {
	c := a.core.Chars[a.core.ActiveChar]
	f := 0
	done := true
	a.core.Log.Debugw(
		"attempting to execute "+n.Typ.String(),
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_pre", a.core.SwapCD,
		"stam_pre", a.core.Stam,
	)

	//do one last ready check
	if !c.ActionReady(n.Typ, n.Param) {
		a.core.Log.Warnw("queued action is not ready, should not happen; skipping frame", "frame", a.core.F, "event", LogSimEvent)
		return 0, false, nil
	}
	switch n.Typ {
	case ActionSkill:
		f, done = a.execActionItem(n, PreSkill, PostSkill, SkillState, true, c.Skill)
	case ActionBurst:
		f, done = a.execActionItem(n, PreBurst, PostBurst, BurstState, true, c.Burst)
	case ActionAttack:
		f, done = a.execActionItem(n, PreAttack, PostAttack, NormalAttackState, false, c.Attack)
	case ActionCharge:
		req := a.core.StamPercentMod(ActionCharge) * c.ActionStam(ActionCharge, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: charge attack", "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f, done = a.execActionItem(n, PreChargeAttack, PostChargeAttack, ChargeAttackState, true, c.ChargeAttack)
		a.core.Events.Emit(OnStamUse, ActionCharge)
	case ActionHighPlunge:
		f, done = a.execActionItem(n, PrePlunge, PostPlunge, PlungeAttackState, true, c.HighPlungeAttack)
	case ActionLowPlunge:
		f, done = a.execActionItem(n, PrePlunge, PostPlunge, PlungeAttackState, true, c.LowPlungeAttack)
	case ActionAim:
		f, done = a.execActionItem(n, PreAimShoot, PostAimShoot, AimState, true, c.Aimed)
	case ActionDash:
		req := a.core.StamPercentMod(ActionDash) * c.ActionStam(ActionDash, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: dash", "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f, done = a.execActionItem(n, PreDash, PostDash, DashState, true, c.Aimed)
		a.core.Events.Emit(OnStamUse, ActionDash)
	case ActionJump:
		f = JumpFrames
		a.core.ResetAllNormalCounter()
	case ActionSwap:
		if a.core.SwapCD > 0 {
			a.core.Log.Warnw("swap on cd", "cd", a.core.SwapCD, "frame", a.core.F, "event", LogActionEvent)
			return 0, false, nil
		}
		f = a.core.Swap(n.Target)
		a.core.ClearState()
	}

	a.core.Log.Debugw(
		"executed "+n.Typ.String(),
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_post", a.core.SwapCD,
		"stam_post", a.core.Stam,
		"animation", f,
	)

	a.core.LastAction = *n

	return f, done, nil
}

func (a *ActionCtrl) execActionItem(
	n *ActionItem,
	pre, post EventType,
	state AnimationState,
	reset bool,
	abil func(map[string]int) (int, int),
) (int, bool) {
	a.core.Events.Emit(pre)
	f, l := abil(n.Param)
	a.core.SetState(state, l)
	if reset {
		a.core.ResetAllNormalCounter()
	}
	a.core.Tasks.Add(func() {
		a.core.Events.Emit(post, f)
	}, f)
	return f, true
}
