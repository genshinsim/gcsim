package core

import "strings"

type Querer interface {
	Next(active string) ([]ActionItem, error)
}

type Action struct {
	Name   string
	Target string //either character or a sequence name

	Exec     []ActionItem //if len > 1 then it's a sequence
	IsSeq    bool         // is this a sequence
	IsStrict bool         //strict sequence?
	Once     bool         // is this an opener (execute once)
	Disabled bool         // is this action disabled
	Pos      int          //current position in execution, default 0

	Last       int //last time this was executed (frame)
	ActionLock int //how many frames this action is locked from executing again

	ActiveCond string
	SwapTo     string
	SwapLock   int
	PostAction ActionType

	Conditions *ExprTreeNode //conditions to be met

	Raw []string //raw action in string
}

type ActionItem struct {
	Typ      ActionType
	Param    map[string]int
	Target   string
	SwapLock int //used for swaplock
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
