//Package action describes the valid actions that any character may take
package action

//TODO: add a sync.Pool here to save some memory allocs
type ActionInfo struct {
	Frames              func(next Action) int
	AnimationLength     int
	CanQueueAfter       int
	State               AnimationState
	FramePausedOnHitlag func() bool
	OnRemoved           func()
	//hidden stuff
	queued       []queuedAction
	cachedFrames [EndActionType]int
	timePassed   float64
}

type queuedAction struct {
	f     func()
	delay float64
}

func (a *ActionInfo) CacheFrames() {
	for i := range a.cachedFrames {
		a.cachedFrames[i] = a.Frames(Action(i))
	}
}

func (a *ActionInfo) QueueAction(f func(), delay int) {
	a.queued = append(a.queued, queuedAction{f: f, delay: float64(delay)})
}

func (a *ActionInfo) CanQueueNext() bool {
	return a.timePassed >= float64(a.CanQueueAfter)
}

func (a *ActionInfo) CanUse(next Action) bool {
	//can't use anything if we're frozen
	if a.FramePausedOnHitlag != nil && a.FramePausedOnHitlag() {
		return false
	}
	return a.timePassed >= float64(a.cachedFrames[next])
}

func (a *ActionInfo) AnimationState() AnimationState {
	return a.State
}

func (a *ActionInfo) Tick() bool {
	//time only goes on if either not hitlag function, or not paused
	if a.FramePausedOnHitlag == nil || !a.FramePausedOnHitlag() {
		a.timePassed++
	}

	//execute all action such that timePassed > delay, and then remove from
	//slice
	n := -1
	for i := range a.queued {
		if a.queued[i].delay <= a.timePassed {
			a.queued[i].f()
		} else {
			n = i
			break
		}
	}
	if n == -1 {
		a.queued = nil
	} else {
		a.queued = a.queued[n:]
	}

	//check if animation is over
	if a.timePassed > float64(a.AnimationLength) {
		//handle remove
		if a.OnRemoved != nil {
			a.OnRemoved()
		}
		return true
	}

	return false
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
	JumpState
	WalkState
)
