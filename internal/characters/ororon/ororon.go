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

const (
	superJumpBeginFrames = 28 + 15 // Jump Frames + Jump->SuperJump Frames

	jumpNsDelay        = 49 // From swap ui gray to nightsoul state emblem
	jumpStamDrainDelay = 5
	jumpStamDrainAmt   = 75
	jumpStamReqAmt     = 1

	maxJumpFrames      = 162                       // From swap ui gray to glider wings appear
	plungeCancelFrames = superJumpBeginFrames + 18 // From start of jump animation to plunge animation start
	fallCancelFrames   = superJumpBeginFrames + 46 // From From start of jump animation to UI changes from gliding to standard UI

	fallFrames = 44 // From fall animation start to swap icon un-gray.
)

func init() {
	core.RegisterCharFunc(keys.Ororon, NewChar)

	// Hold Jump
	jumpHoldFrames = make([][]int, 2)
	// Hold Jump -> X
	jumpHoldFrames[0] = frames.InitAbilSlice(60 * 10) // set to very high number for most abilities
	jumpHoldFrames[0][action.ActionHighPlunge] = plungeCancelFrames
	// Fall -> X
	jumpHoldFrames[1] = frames.InitAbilSlice(fallFrames)
	jumpHoldFrames[1][action.ActionAttack] = 45
	jumpHoldFrames[1][action.ActionAim] = 46
	jumpHoldFrames[1][action.ActionSkill] = 45
	jumpHoldFrames[1][action.ActionBurst] = 46
	jumpHoldFrames[1][action.ActionDash] = 46
	jumpHoldFrames[1][action.ActionJump] = 47
	jumpHoldFrames[1][action.ActionWalk] = 47
	jumpHoldFrames[1][action.ActionSwap] = 44
}

type char struct {
	*tmpl.Character
	nightsoulState     *nightsoul.State
	particlesGenerated bool
	c2Bonus            []float64
	c6stacks           *stacks.MultipleRefreshNoRemove
	c6bonus            []float64
	jmpSrc             int
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
