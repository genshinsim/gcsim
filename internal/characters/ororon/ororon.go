package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/model"
)

var jumpHoldFrames [][]int

// TODO: find real Frame delays
const (
	jumpNsDelay        = 10
	jumpNsDuration     = 18
	jumpStamDrainDelay = 5
	jumpStamDrainAmt   = 75
	jumpStamReqAmt     = 1 // TODO: Find real value

	maxJumpFrames   = 60
	minCancelFrames = 15 // assume is the same as minPlungeFrames.
	// minPlungeFrames = 17
	jumpNoStamFallDelayFrames = maxJumpFrames // If ororon has 0 stam, fall cancel takes longer.

	fallFrames = 60 // Time it takes from cancelling high jump to hitting the ground.

	// TODO: How to prevent stamina from regenerating until allowed?
	fallStamResumeDelay = 60 // Time it takes stamina to start regenerating again after landing from fall.
)

func init() {
	core.RegisterCharFunc(keys.Ororon, NewChar)

	// Hold Jump
	jumpHoldFrames = make([][]int, 2)
	// Hold Jump -> X
	jumpHoldFrames[0] = frames.InitAbilSlice(60 * 10) // set to very high number for most abilities
	jumpHoldFrames[0][action.ActionHighPlunge] = minCancelFrames
	// Fall -> X
	jumpHoldFrames[1] = frames.InitAbilSlice(fallFrames)
	jumpHoldFrames[1][action.ActionAttack] = 158
	jumpHoldFrames[1][action.ActionBurst] = 159
	jumpHoldFrames[1][action.ActionSwap] = 155
}

type char struct {
	*tmpl.Character
	nightsoulState     *nightsoul.State
	particlesGenerated bool
	c2Bonus            []float64
	c6stacks           *stacks.MultipleRefreshNoRemove
	c6bonus            []float64
	jmpSrc             int
	allowFallFrame     int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 80

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Init()
	c.a4Init()
	c.c1Init()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 14
	}
	return c.Character.AnimationStartDelay(k)
}

// TODO: Should this return stam used or just stam required to start?
func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionJump && p["hold"] != 0 {
		return 75
	}
	return c.Character.ActionStam(a, p)
}

// Ororon is in NS if either he has it through high jump or if he has it through his ascention.
// Each has independent duration, so must be checked for in parallel.
func (c *char) StatusIsActive(key string) bool {
	if key == nightsoul.NightsoulBlessingStatus {
		return (c.Character.StatusIsActive(nightsoul.NightsoulBlessingStatus) ||
			c.Character.StatusIsActive(jumpNsStatusTag))
	}
	return c.Character.StatusIsActive(key)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if a1 window is active is on-field
	if a == action.ActionJump && p["hold"] != 0 {
		if c.Core.Player.Stam < jumpStamReqAmt {
			return false, action.InsufficientStamina
		}
	}
	return c.Character.ActionReady(a, p)
}
