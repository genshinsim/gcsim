package lauma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	deerFrames   []int
	chargeFrames []int
)

const (
	chargeHitmark = 73
	deerStatusKey = "lauma-deer-state"
)

func init() {
	// deer state frames
	deerFrames = frames.InitAbilSlice(64) // skill
	deerFrames[action.ActionAttack] = 61
	deerFrames[action.ActionCharge] = 63
	// deerFrames[action.ActionSkillHoldFramesOnly] = 63
	deerFrames[action.ActionBurst] = 63
	deerFrames[action.ActionSkill] = 63
	deerFrames[action.ActionSwap] = 0

	// charge attack frames
	chargeFrames = frames.InitAbilSlice(68) // CA -> swap
	chargeFrames[action.ActionAttack] = 67
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// by default enters different state where she doesnt hit enemies and consumes 25 stamina per second
	if c.deerStateReady {
		c.deerStateReady = false
		return c.enterDeerState()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: 0},
			8,
			3,
		),
		chargeFrames[action.InvalidAction],
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) enterDeerState() (action.Info, error) {
	c.deerStateStaminaBleed()
	c.AddStatus(deerStatusKey, 10*60, true)

	c.endDeerStateCondition()

	return action.Info{
		Frames:          frames.NewAbilFunc(deerFrames),
		AnimationLength: deerFrames[action.InvalidAction],
		CanQueueAfter:   deerFrames[action.ActionSwap],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) endDeerStateCondition() {
	c.Core.Events.Subscribe(event.OnActionExec, func(args ...any) bool {
		if !c.StatusIsActive(deerStatusKey) {
			return true
		}

		a := args[1].(action.Action)
		if a == action.ActionJump {
			return false
		}
		c.endDeerState()
		return true
	}, "lauma-exit-deer-state")
}

func (c *char) deerStateStaminaBleed() {
	if !c.StatusIsActive(deerStatusKey) {
		return
	}
	staminaCost := 25.0 / 60.0
	if c.Core.Player.Stam < staminaCost {
		c.endDeerState()
	}
	c.Core.Player.RestoreStam(-staminaCost)
	c.Core.Tasks.Add(c.deerStateStaminaBleed, 1)
}

func (c *char) endDeerState() {
	c.DeleteStatus(deerStatusKey)
	c.Core.Tasks.Add(func() {
		c.deerStateReady = true
		c.Core.Log.NewEvent("deer state refreshed", glog.LogCharacterEvent, c.Index())
	}, 4*60)
}
