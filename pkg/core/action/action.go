// Package action describes the valid actions that any character may take
package action

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// TODO: add a sync.Pool here to save some memory allocs
type ActionInfo struct {
	Frames              func(next Action) int `json:"-"`
	AnimationLength     int
	CanQueueAfter       int
	State               AnimationState
	FramePausedOnHitlag func() bool               `json:"-"`
	OnRemoved           func(next AnimationState) `json:"-"`
	//following are exposed only so we can log it properly
	CachedFrames         [EndActionType]int //TODO: consider removing the cache frames and instead cache the frames function instead
	TimePassed           float64
	NormalizedTimePassed float64
	UseNormalizedTime    func(next Action) bool
	//hidden stuff
	queued []queuedAction
}

// ActionEval represents a sim action
type ActionEval struct {
	Char   keys.Char
	Action Action
	Param  map[string]int
}

// Evaluator provides method for getting next action
type Evaluator interface {
	Continue()
	NextAction() (*ActionEval, error)
}

type queuedAction struct {
	f     func()
	delay float64
}

func (a *ActionInfo) CacheFrames() {
	for i := range a.CachedFrames {
		a.CachedFrames[i] = a.Frames(Action(i))
	}
}

func (a *ActionInfo) QueueAction(f func(), delay int) {
	a.queued = append(a.queued, queuedAction{f: f, delay: float64(delay)})
}

func (a *ActionInfo) CanQueueNext() bool {
	return a.TimePassed >= float64(a.CanQueueAfter)
}

func (a *ActionInfo) CanUse(next Action) bool {
	if a.UseNormalizedTime != nil && a.UseNormalizedTime(next) {
		return a.NormalizedTimePassed >= float64(a.CachedFrames[next])
	}
	//can't use anything if we're frozen
	if a.FramePausedOnHitlag != nil && a.FramePausedOnHitlag() {
		return false
	}
	return a.TimePassed >= float64(a.CachedFrames[next])
}

func (a *ActionInfo) AnimationState() AnimationState {
	return a.State
}

func (a *ActionInfo) Tick() bool {
	a.NormalizedTimePassed++ //this always increments
	//time only goes on if either not hitlag function, or not paused
	if a.FramePausedOnHitlag == nil || !a.FramePausedOnHitlag() {
		a.TimePassed++
	}

	//execute all action such that timePassed > delay, and then remove from
	//slice
	if a.queued != nil {
		n := 0
		for i := 0; i < len(a.queued); i++ {
			if a.queued[i].delay <= a.TimePassed {
				a.queued[i].f()
			} else {
				a.queued[n] = a.queued[i]
				n++
			}
		}
		a.queued = a.queued[:n]
	}

	//check if animation is over
	if a.TimePassed > float64(a.AnimationLength) {
		//handle remove
		if a.OnRemoved != nil {
			a.OnRemoved(Idle)
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

func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(astr[a])
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range astr {
		if v == s {
			*a = Action(i)
			return nil
		}
	}
	return errors.New("unrecognized action")
}

func StringToAction(s string) Action {
	for i, v := range astr {
		if v == s {
			return Action(i)
		}
	}
	return InvalidAction
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
	SwapState
)

var statestr = []string{
	"idle",
	"normal",
	"charge",
	"plunge",
	"skill",
	"burst",
	"aim",
	"dash",
	"jump",
	"walk",
	"swap",
}

func (a AnimationState) String() string {
	return statestr[a]
}
