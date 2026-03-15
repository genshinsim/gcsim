package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	basicChargeWindup       = 20
	basicChargeHitmark      = 44
	slitherMinCancelFrames  = 24
	slitherMoveInterval     = 1
	slitherMoveDistance     = 0.1
	slitherStamInterval     = 1
	slitherStamTickCost     = 18.15 / 60.0
	phantasmAnimationLength = 106
	phantasmRecoverFrame    = 89
	phantasmConsumeDewFrame = 29
	phantasmHit1            = 30
	phantasmHit2            = 35
	phantasmHit3            = 43
	phantasmHit4            = 44
	phantasmHit5            = 45
)

var chargeFrames []int

func init() {
	chargeFrames = frames.InitAbilSlice(72)
	chargeFrames[action.ActionAttack] = 49
	chargeFrames[action.ActionSkill] = 48
	chargeFrames[action.ActionBurst] = 48
	chargeFrames[action.ActionDash] = 48
	chargeFrames[action.ActionJump] = 48
	chargeFrames[action.ActionSwap] = 48
	chargeFrames[action.ActionWalk] = 71
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a != action.ActionCharge {
		return c.Character.ActionStam(a, p)
	}
	if c.canTriggerPhantasm() {
		return 0
	}
	if c.StatusIsActive(shadowDanceKey) {
		return 25
	}
	return 50
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.canTriggerPhantasm() {
		return c.specialChargeAttack()
	}
	c.clearPhantasmChargeLoop()
	return c.basicChargeAttack()
}

func (c *char) specialChargeAttack() (action.Info, error) {
	continuing := c.phantasmChargeHoldActive() && c.Core.Player.CurrentState() == action.ChargeAttackState
	remaining := max(c.StatusDuration(shadowDanceKey), slitherMinCancelFrames)
	windup := 0
	if !continuing {
		windup = basicChargeWindup
		c.AddStatus(phantasmChargeHoldKey, -1, false)
		c.phantasmChargeSrc = c.Core.F
		c.phantasmEndFrame = 0
		c.QueueCharTask(c.phantasmChargeLoopTask(c.phantasmChargeSrc), windup)
	}

	return action.Info{
		Frames: func(next action.Action) int {
			phaseFrames := slitherMinCancelFrames
			phantasmActive := c.phantasmActive()
			if phantasmActive {
				phaseFrames = max(c.phantasmEndFrame-c.Core.F, 1)
			}
			switch next {
			case action.ActionCharge:
				if phantasmActive {
					return windup + phaseFrames
				}
				return windup + remaining
			case action.ActionAttack, action.ActionSkill, action.ActionBurst, action.ActionDash, action.ActionJump, action.ActionSwap, action.ActionWalk:
				return windup + phaseFrames
			default:
				return windup + remaining
			}
		},
		AnimationLength: windup + remaining,
		CanQueueAfter:   windup + 1,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			if next != action.ChargeAttackState {
				c.clearPhantasmChargeLoop()
			}
		},
	}, nil
}

func (c *char) phantasmChargeLoopTask(src int) func() {
	return func() {
		if c.chargeLoopInterrupted(src) {
			return
		}
		if c.phantasmActive() {
			return
		}
		if c.canTriggerPhantasm() {
			c.triggerPhantasmFromLoop(src)
			return
		}
		c.startSlitherPhase(src)
	}
}

func (c *char) startSlitherPhase(src int) {
	if !c.chargeLoopActive(src) {
		return
	}
	if c.slitherActive() {
		return
	}
	c.AddStatus(slitherKey, -1, false)
	c.slitherSrc = src
	c.QueueCharTask(c.slitherTickTask(src), 0)
}

func (c *char) slitherTickTask(src int) func() {
	return func() {
		if c.slitherLoopInterrupted(src) {
			return
		}
		if c.canTriggerPhantasm() {
			c.triggerPhantasmFromLoop(src)
			return
		}
		if c.phantasmActive() {
			c.clearSlither()
			return
		}

		req := slitherStamTickCost * (1 + c.Core.Player.StamPercentMod(action.ActionCharge))
		if req < 0 {
			req = 0
		}
		if c.Core.Player.Stam < req {
			c.clearSlither()
			return
		}

		player := c.Core.Combat.Player()
		target := c.Core.Combat.PrimaryTarget()
		if target != nil {
			player.SetDirection(target.Pos())
		}
		nextPos := info.CalcOffsetPoint(player.Pos(), info.Point{Y: slitherMoveDistance}, player.Direction())
		c.Core.Combat.SetPlayerPos(nextPos)
		c.Core.Player.UseStam(req, action.ActionWait)
		c.absorbSeedsOfDeceit()
		c.QueueCharTask(c.slitherTickTask(src), max(slitherMoveInterval, slitherStamInterval))
	}
}

func (c *char) triggerPhantasmFromLoop(src int) {
	if !c.chargeLoopActive(src) || !c.canTriggerPhantasm() {
		return
	}

	c.clearSlither()
	c.phantasmCharges--
	c.phantasmEndFrame = c.Core.F + phantasmAnimationLength
	c.Core.Player.SetActionLength(c.phantasmEndFrame-src, c.phantasmEndFrame-src)
	c.queuePhantasmPerformance(src)
	c.QueueCharTask(func() {
		if !c.chargeLoopActive(src) {
			return
		}
		c.phantasmEndFrame = 0
	}, phantasmAnimationLength)
}

func (c *char) basicChargeAttack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       charge[0][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: -2}, 3, 9),
		basicChargeHitmark,
		basicChargeHitmark,
	)
	c.QueueCharTask(c.absorbSeedsOfDeceit, basicChargeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) queuePhantasmPerformance(src int) {
	consumeFrame := phantasmConsumeDewFrame
	c.QueueCharTask(func() {
		if c.phantasmChargeSrc != src {
			return
		}
		c.Core.Player.ConsumeVerdantDew(1)
		c.absorbSeedsOfDeceit()
	}, consumeFrame)

	neferHit1 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Phantasm Performance (Nefer 1)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm[0][c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * phantasm[1][c.TalentLvlSkill()],
	}
	neferHit2 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Phantasm Performance (Nefer 2)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm[2][c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * phantasm[3][c.TalentLvlSkill()],
	}
	shadeHit1 := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Phantasm Performance (Shade 1)",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		UseEM:            true,
		IgnoreDefPercent: 1,
		Mult:             phantasm[4][c.TalentLvlSkill()],
	}
	shadeHit2 := shadeHit1
	shadeHit2.Abil = "Phantasm Performance (Shade 2)"
	shadeHit2.Mult = phantasm[5][c.TalentLvlSkill()]
	shadeHit3 := shadeHit1
	shadeHit3.Abil = "Phantasm Performance (Shade 3)"
	shadeHit3.Mult = phantasm[6][c.TalentLvlSkill()]

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(neferHit1, ap, phantasmHit1, phantasmHit1)
	c.Core.QueueAttack(shadeHit1, ap, phantasmHit2, phantasmHit2)
	c.Core.QueueAttack(shadeHit2, ap, phantasmHit3, phantasmHit3)
	c.Core.QueueAttack(neferHit2, ap, phantasmHit4, phantasmHit4)
	c.Core.QueueAttack(shadeHit3, ap, phantasmHit5, phantasmHit5)
}
