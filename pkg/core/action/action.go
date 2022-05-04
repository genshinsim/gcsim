//Package action describes the valid actions that any character may take
package action

type ActionInfo struct {
	Frames          func(next Action) int
	AnimationLength int
	CanQueueAfter   int
	Post            int
	State           AnimationState
}

type Action int

const (
	InvalidAction Action = iota
	ActionSkill
	ActionBurst
	ActionAttack
	ActionCharge
	ActionHighPlunge
	ActionLowPlunge
	ActionAim
	ActionDash
	ActionJump
	//following action have to implementations
	ActionSwap
	ActionWalk
	ActionWait // character should stand around and wait
	EndActionType
	//these are only used for frames purposes and that's why it's after end
	ActionSkillHoldFramesOnly
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
	"dash",
	"jump",
	"swap",
	"walk",
	"wait",
}

func (a Action) String() string {
	return astr[a]
}

type AnimationState int

const (
	Idle AnimationState = iota
	NormalAttackState
	ChargeAttackState
	PlungeAttackState
	SkillState
	BurstState
	AimState
	DashState
)
