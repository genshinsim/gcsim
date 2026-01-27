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
	chargeHitmark           = 73
	chargeReleaseFrame      = 67
	deerTransformationFrame = 49
	deerStatusKey           = "lauma-spirit-envoy"
)

func init() {
	// deer state frames
	deerFrames = frames.InitAbilSlice(64) // skill
	deerFrames[action.ActionAttack] = 61
	deerFrames[action.ActionCharge] = 63
	// deerFrames[action.ActionSkillHoldFramesOnly] = 63
	deerFrames[action.ActionBurst] = 63
	deerFrames[action.ActionSkill] = 63
	deerFrames[action.ActionSwap] = deerTransformationFrame
	deerFrames[action.ActionJump] = deerTransformationFrame

	// charge attack frames
	chargeFrames = frames.InitAbilSlice(68) // CA -> swap
	chargeFrames[action.ActionAttack] = chargeReleaseFrame
	chargeFrames[action.ActionSkill] = chargeReleaseFrame
	chargeFrames[action.ActionBurst] = chargeReleaseFrame
	chargeFrames[action.ActionDash] = chargeReleaseFrame
	chargeFrames[action.ActionJump] = chargeReleaseFrame
	chargeFrames[action.ActionSwap] = chargeReleaseFrame
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// by default enters different state where she doesnt hit enemies and consumes 25 stamina per second
	if c.deerStateReady && !c.StatusIsActive(deerStatusKey) {
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
			2.8,
			8,
		),
		chargeFrames[action.InvalidAction],
		chargeHitmark,
	)

	return action.Info{
		Frames: func(next action.Action) int {
			if c.deerStateReady && next == action.ActionCharge {
				return chargeReleaseFrame
			}
			return chargeFrames[next]
		},
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) enterDeerState() (action.Info, error) {
	c.QueueCharTask(func() {
		c.deerStateStaminaBleed()
		c.AddStatus(deerStatusKey, 10*60, true)
	}, deerTransformationFrame)

	return action.Info{
		Frames:          frames.NewAbilFunc(deerFrames),
		AnimationLength: deerFrames[action.InvalidAction],
		CanQueueAfter:   deerFrames[action.ActionSwap],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) chargeInit() {
	c.Core.Events.Subscribe(event.OnActionExec, func(args ...any) bool {
		if !c.StatusIsActive(deerStatusKey) {
			return false
		}

		if c.Core.Player.Active() != c.Index() {
			return false
		}

		a := args[1].(action.Action)
		if a == action.ActionJump || a == action.ActionWalk {
			return false
		}
		c.endDeerState()
		return false
	}, "lauma-exit-deer-state")

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) bool {
		if !c.StatusIsActive(deerStatusKey) {
			return false
		}

		prev := args[0].(int)
		if prev != c.Index() {
			return false
		}
		c.endDeerState()

		return false
	}, "lauma-exit-deer-state-swap")
}

func (c *char) deerStateStaminaBleed() {
	if !c.StatusIsActive(deerStatusKey) {
		return
	}
	staminaCost := 25.0 / 60.0
	if c.Core.Player.Stam < staminaCost {
		c.endDeerState()
	}
	c.Core.Player.UseStam(staminaCost, action.ActionWait)
	c.Core.Tasks.Add(c.deerStateStaminaBleed, 1)
}

func (c *char) endDeerState() {
	c.DeleteStatus(deerStatusKey)
	cd := int(4 * 60 * c.a4SpiritEnvoyCooldownReduction())
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), "spirit envoy cooldown triggered").
			Write("type", "charge").
			Write("expiry", c.Core.F+cd).
			Write("original_cd", c.Core.F+cd)
	}

	c.Core.Tasks.Add(func() {
		c.deerStateReady = true
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), "spirit envoy cooldown ready").
				Write("type", "charge")
		}
	}, cd)
}
