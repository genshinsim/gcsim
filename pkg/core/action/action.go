//Package action describes the valid actions that any character may take
package action

type ActionInfo struct {
	Frames func(next Action) int
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
