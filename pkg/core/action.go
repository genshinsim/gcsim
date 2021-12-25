package core

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Action struct {
	Name   string
	Target keys.Char //either character or a sequence name

	Exec     []ActionItem //if len > 1 then it's a sequence
	IsSeq    bool         // is this a sequence
	IsStrict bool         //strict sequence?
	Once     bool         // is this an opener (execute once)
	Disabled bool         // is this action disabled
	Pos      int          //current position in execution, default 0

	Last       int //last time this was executed (frame)
	ActionLock int //how many frames this action is locked from executing again

	ActiveCond keys.Char
	SwapTo     keys.Char
	SwapLock   int
	PostAction ActionType

	Wait int //wait till this frame, only for calc

	Conditions *ExprTreeNode //conditions to be met

	Raw []string //raw action in string
}

type ActionItem struct {
	Typ            ActionType
	Param          map[string]int
	Target         keys.Char
	SwapLock       int  //used for swaplock
	FramesOverride bool //true if using custom frames
	Frames         int  //frames if overridden
}

type ActionType int

const (
	ActionSequence ActionType = iota
	ActionSequenceStrict
	ActionDelimiter
	ActionSequenceReset
	ActionSkill
	ActionBurst
	ActionAttack
	ActionCharge
	ActionHighPlunge
	ActionLowPlunge
	ActionSpecialProc
	ActionAim
	ActionSwap
	ActionSwapLock    //force swap lock
	ActionCancellable // delim cancellable action
	ActionDash
	ActionJump
	ActionOtherEvents //delim for other events
	ActionHurt        //damage characters
	//delim
	EndActionType
)

var astr = []string{
	"sequence",
	"sequence_strict",
	"",
	"reset_sequence",
	"skill",
	"burst",
	"attack",
	"charge",
	"high_plunge",
	"low_plunge",
	"proc",
	"aim",
	"swap",
	"swaplock",
	"",
	"dash",
	"jump",
	"",
	"hurt",
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

type ActionHandler interface {
	Exec(n ActionItem) (int, bool, error) //return frames, if executed, any errors
}

type ActionCtrl struct {
	core *Core
}

func NewActionCtrl(c *Core) *ActionCtrl {
	return &ActionCtrl{
		core: c,
	}
}

func (a *ActionCtrl) Exec(n ActionItem) (int, bool, error) {

	c := a.core.Chars[a.core.ActiveChar]
	f := 0
	l := 0

	a.core.Log.Debugw(
		"attempting to execute "+n.Typ.String(),
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_pre", a.core.SwapCD,
		"stam_pre", a.core.Stam,
		"animation", f,
	)

	//do one last ready check
	if !c.ActionReady(n.Typ, n.Param) {
		a.core.Log.Warnw("queued action is not ready, should not happen; skipping frame", "frame", a.core.F, "event", LogSimEvent)
		return 0, false, nil
	}
	switch n.Typ {
	case ActionSwapLock:
		a.core.SwapCD += n.SwapLock
		a.core.Log.Debugw(
			"locked swap",
			"frame", a.core.F,
			"event", LogActionEvent,
			"char", a.core.ActiveChar,
			"dur", n.SwapLock,
			"cd", a.core.SwapCD,
		)
		return 0, true, nil
	case ActionSkill:
		a.core.Events.Emit(PreSkill)
		f, l = c.Skill(n.Param)
		a.core.SetState(SkillState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(PostSkill, f)
	case ActionBurst:
		a.core.Events.Emit(PreBurst)
		f, l = c.Burst(n.Param)
		a.core.SetState(BurstState, l)
		a.core.ResetAllNormalCounter()
		a.core.Tasks.Add(func() {
			a.core.Events.Emit(PostBurst, f)
		}, f)
	case ActionAttack:
		a.core.Events.Emit(PreAttack)
		f, l = c.Attack(n.Param)
		a.core.SetState(NormalAttackState, l)
		a.core.Events.Emit(PostAttack, f)
	case ActionCharge:
		req := a.core.StamPercentMod(ActionCharge) * c.ActionStam(ActionCharge, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: charge attack", "have", a.core.Stam)
			return 0, false, nil
		}

		a.core.Stam -= req
		a.core.Events.Emit(PreChargeAttack)
		f, l = c.ChargeAttack(n.Param)
		a.core.SetState(ChargeAttackState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(PostChargeAttack, f)
		a.core.Events.Emit(OnStamUse, ActionCharge)

	case ActionHighPlunge:
		a.core.Events.Emit(PrePlunge)
		f, l = c.HighPlungeAttack(n.Param)
		a.core.SetState(PlungeAttackState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(PostPlunge, f)
	case ActionLowPlunge:
		a.core.Events.Emit(PrePlunge)
		f, l = c.LowPlungeAttack(n.Param)
		a.core.SetState(PlungeAttackState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(PostPlunge, f)
	case ActionAim:
		a.core.Events.Emit(PreAimShoot)
		f, l = c.Aimed(n.Param)
		a.core.SetState(AimState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(PostAimShoot, f)
	case ActionSwap:
		if a.core.SwapCD > 0 {
			a.core.Log.Warnw("swap on cd", "cd", a.core.SwapCD, "frame", a.core.F, "event", LogActionEvent)
			return 0, false, nil
		}
		f = a.core.Swap(n.Target)
		a.core.ClearState()
	case ActionCancellable:
	case ActionDash:
		//check if enough req
		req := a.core.StamPercentMod(ActionDash) * c.ActionStam(ActionDash, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: dash", "have", a.core.Stam)
			return 0, false, nil
		}

		a.core.Stam -= req
		f, l = c.Dash(n.Param)
		a.core.SetState(DashState, l)
		a.core.ResetAllNormalCounter()
		a.core.Events.Emit(OnDash)
		a.core.Events.Emit(OnStamUse, ActionDash)

	case ActionJump:
		f = JumpFrames
		a.core.ResetAllNormalCounter()
	}

	// s.Log.Infof("[%v] %v executing %v", s.Frame(), s.ActiveChar, a.Action)
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

	a.core.LastAction = n
	a.core.LastAction.Target = c.Key()

	return f, true, nil
}
