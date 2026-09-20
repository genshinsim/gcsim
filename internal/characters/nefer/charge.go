package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	basicChargeWindup        = 20
	basicChargeHitmark       = 44
	slitherMaxDuration      = 150
	slitherActivationFrames = 60
	slitherMinCancelFrames  = 24
	slitherMoveInterval      = 1
	slitherMoveDistance      = 0.1
	slitherStamInterval      = 1
	slitherStamTickCost      = 18.15 / 60.0
	phantasmAnimationLength  = 106
	phantasmConsumeDewFrame  = 9
	phantasmHit1             = 10
	phantasmHit2             = 15
	phantasmHit3             = 23
	phantasmHit4             = 24
	phantasmHit5             = 25
	c6PhantasmHit2EM         = 0.85
	c6PhantasmEndEM          = 1.20
)

var chargeFrames []int

func basicChargeCanQueueAfter() int {
	return min(chargeFrames[action.ActionAttack], chargeFrames[action.ActionSwap])
}

func init() {
	chargeFrames = frames.InitAbilSlice(72)
	chargeFrames[action.ActionAttack] = 29
	chargeFrames[action.ActionSkill] = 48
	chargeFrames[action.ActionBurst] = 48
	chargeFrames[action.ActionDash] = 48
	chargeFrames[action.ActionJump] = 48
	chargeFrames[action.ActionSwap] = 28
	chargeFrames[action.ActionWalk] = 71
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a != action.ActionCharge {
		return c.Character.ActionStam(a, p)
	}
	if c.canTriggerPhantasm() {
		return 0
	}
	return slitherStamTickCost * slitherActivationFrames
}

// ChargeAttack uses `hold` as an explicit hold duration in frames.
//
// Route split:
//   - hold=0: tap. If Phantasm is already available, execute the direct Phantasm CA.
//     Otherwise execute the no-hold ordinary CA route.
//   - hold>0: hold for exactly that many frames, capped to 150. If Phantasm is
//     already available, the result is identical to tap. Otherwise enter Slither
//     first. Slither then resolves into either Phantasm CA as soon as the Phantasm
//     condition becomes true or ordinary CA when the hold route ends.
//
// Ordinary CA is modeled as the no-hold Slither-release path. Workbook ordinary
// CA rows are currently interpreted as button-to-button timings for that full
// route rather than as release-only timings after a separate CA windup. The
// branch does not yet know how many of those frames belong to the embedded
// Slither-entry segment, so it currently assumes that contribution is 0f.
//
// The shared player ReadyCheck therefore gates non-Phantasm charge startup on
// the Slither-entry threshold: if there is not enough stamina to enter Slither,
// the ordinary CA route never begins.
func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if p["hold"] != 0 {
		return c.heldChargeAttack(p)
	}
	if c.canTriggerPhantasm() {
		return c.tapPhantasmChargeAttack()
	}
	c.Core.Player.RestoreStam(c.Core.Player.AbilStamCost(c.Index(), action.ActionCharge, p))
	c.clearPhantasmChargeLoop()
	return c.basicChargeAttack()
}

func heldChargeDuration(p map[string]int) int {
	return min(max(p["hold"], 1), slitherMaxDuration)
}

func (c *char) tapPhantasmChargeAttack() (action.Info, error) {
	c.clearPhantasmChargeLoop()
	src := c.Core.F
	c.chargeRoute.src = src
	c.QueueCharTask(func() {
		c.startPhantasmPhase(src)
	}, basicChargeWindup)

	return action.Info{
		Frames: func(next action.Action) int {
			switch next {
			case action.ActionBurst, action.ActionDash:
				return basicChargeWindup + phantasmHit1
			default:
				return basicChargeWindup + phantasmAnimationLength
			}
		},
		AnimationLength: basicChargeWindup + phantasmAnimationLength,
		CanQueueAfter:   1,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			if c.chargeRoute.src != src {
				return
			}
			c.clearPhantasmChargeLoop()
		},
	}, nil
}

func (c *char) heldChargeAttack(p map[string]int) (action.Info, error) {
	// If Phantasm is already available, hold and tap intentionally collapse to the
	// same direct Phantasm route. The only difference between tap and hold is the
	// Slither prefix used while waiting for either Phantasm or the ordinary release.
	if c.canTriggerPhantasm() {
		return c.tapPhantasmChargeAttack()
	}

	c.Core.Player.RestoreStam(c.Core.Player.AbilStamCost(c.Index(), action.ActionCharge, p))
	c.clearPhantasmChargeLoop()
	src := c.Core.F
	windup := basicChargeWindup
	slitherDuration := heldChargeDuration(p)
	c.chargeRoute.src = src
	c.QueueCharTask(func() {
		if c.chargeRouteInterrupted(src) {
			return
		}
		c.startSlitherPhase(src)
	}, windup)
	c.QueueCharTask(func() {
		c.finishHeldCharge(src)
	}, windup+slitherDuration)

	return action.Info{
		Frames: func(next action.Action) int {
			if c.phantasmActive() {
				phaseStartOffset := c.chargeRoute.phantasmStartFrame - src
				switch next {
				case action.ActionBurst, action.ActionDash:
					return phaseStartOffset + phantasmHit1
				default:
					return phaseStartOffset + phantasmAnimationLength
				}
			}
			if c.chargeRoute.releaseStartFrame > 0 {
				phaseStartOffset := c.chargeRoute.releaseStartFrame - src
				return phaseStartOffset + chargeFrames[next]
			}
			switch next {
			case action.ActionBurst, action.ActionDash:
				return windup + slitherMinCancelFrames
			case action.ActionCharge:
				return windup + slitherDuration + chargeFrames[next]
			default:
				return windup + slitherDuration + chargeFrames[action.InvalidAction]
			}
		},
		AnimationLength: windup + slitherDuration + chargeFrames[action.InvalidAction],
		CanQueueAfter:   1,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			if next != action.ChargeAttackState {
				c.clearPhantasmChargeLoop()
			}
		},
	}, nil
}

func (c *char) startSlitherPhase(src int) {
	if !c.chargeRouteActive(src) {
		return
	}
	if c.slitherActive() {
		return
	}
	c.AddStatus(slitherKey, -1, false)
	c.chargeRoute.slitherSrc = src
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

// finishHeldCharge resolves the non-Phantasm hold path into the same ordinary
// CA release used by the tap route. At this point the only difference from an
// immediate ordinary CA is the Slither prefix that already happened.
func (c *char) finishHeldCharge(src int) {
	if c.chargeRouteInterrupted(src) || c.phantasmActive() || c.chargeRoute.releaseStartFrame > 0 {
		return
	}
	c.clearSlither()
	c.startBasicChargeRelease(src)
}

func (c *char) triggerPhantasmFromLoop(src int) {
	if !c.chargeRouteActive(src) || !c.canTriggerPhantasm() {
		return
	}

	c.clearSlither()
	c.startPhantasmPhase(src)
}

func (c *char) startPhantasmPhase(src int) {
	if c.chargeRoute.src != src || !c.canTriggerPhantasm() {
		return
	}
	c.phantasmCharges--
	c.chargeRoute.phantasmStartFrame = c.Core.F
	c.chargeRoute.phantasmEndFrame = c.Core.F + phantasmAnimationLength
	c.queuePhantasmPerformance(src)
	c.QueueCharTask(func() {
		if c.chargeRoute.src != src {
			return
		}
		c.chargeRoute.phantasmStartFrame = 0
		c.chargeRoute.phantasmEndFrame = 0
	}, phantasmAnimationLength)
}

func (c *char) startBasicChargeRelease(src int) {
	if c.chargeRoute.src != src {
		return
	}
	c.chargeRoute.releaseStartFrame = c.Core.F
	c.queueBasicChargeRelease()
}

func (c *char) basicChargeAttack() (action.Info, error) {
	c.clearPhantasmChargeLoop()
	c.queueBasicChargeRelease()

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   basicChargeCanQueueAfter(),
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) queueBasicChargeRelease() {
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
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: -2}, 3, 9),
		basicChargeHitmark,
		basicChargeHitmark,
	)
	c.QueueCharTask(c.absorbSeedsOfDeceit, basicChargeHitmark)
}

func (c *char) queuePhantasmPerformance(src int) {
	consumeFrame := phantasmConsumeDewFrame
	c.QueueCharTask(func() {
		if c.chargeRoute.src != src {
			return
		}
		c.Core.Player.ConsumeDew(1)
		c.absorbSeedsOfDeceit()
	}, consumeFrame)

	shadeScaleBonus := c.c1ShadeScaleBonus()

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
	if c.Base.Cons >= 6 {
		neferHit2.AttackTag = attacks.AttackTagDirectLunarBloom
		neferHit2.Durability = 0
		neferHit2.UseEM = true
		neferHit2.IgnoreDefPercent = 1
		neferHit2.Mult = c6PhantasmHit2EM
		neferHit2.FlatDmg = 0
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
		Mult:             phantasm[4][c.TalentLvlSkill()] + shadeScaleBonus,
	}
	shadeHit2 := shadeHit1
	shadeHit2.Abil = "Phantasm Performance (Shade 2)"
	shadeHit2.Mult = phantasm[5][c.TalentLvlSkill()] + shadeScaleBonus
	shadeHit3 := shadeHit1
	shadeHit3.Abil = "Phantasm Performance (Shade 3)"
	shadeHit3.Mult = phantasm[6][c.TalentLvlSkill()] + shadeScaleBonus

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.queuePhantasmHit(src, neferHit1, ap, phantasmHit1, true)
	c.queuePhantasmHit(src, shadeHit1, ap, phantasmHit2, false)
	c.queuePhantasmHit(src, shadeHit2, ap, phantasmHit3, false)
	c.queuePhantasmHit(src, neferHit2, ap, phantasmHit4, false)
	c.queuePhantasmHit(src, shadeHit3, ap, phantasmHit5, false)
	if c.Base.Cons >= 6 {
		c6EndHit := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Phantasm Performance (C6 End Hit)",
			AttackTag:        attacks.AttackTagDirectLunarBloom,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Dendro,
			UseEM:            true,
			IgnoreDefPercent: 1,
			Mult:             c6PhantasmEndEM,
		}
		c.queuePhantasmHit(src, c6EndHit, ap, phantasmAnimationLength, false)
	}
}

func (c *char) queuePhantasmHit(
	src int,
	ai info.AttackInfo,
	ap info.AttackPattern,
	delay int,
	first bool,
) {
	c.QueueCharTask(func() {
		if first {
			if !c.chargeRouteActive(src) {
				return
			}
			c.phantasmCommitSrc = src
			c.phantasmVeilMultiplier = 1 + c.phantasmVeilBonus()
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("nefer phantasm veil snapshot", glog.LogCharacterEvent, c.Index()).
					Write("ability", ai.Abil).
					Write("stacks", c.currentVeilStacks()).
					Write("bonus", c.phantasmVeilMultiplier-1).
					Write("multiplier", c.phantasmVeilMultiplier)
			}
		} else if c.phantasmCommitSrc != src {
			return
		}
		ai.Mult *= c.phantasmVeilMultiplier
		ai.FlatDmg *= c.phantasmVeilMultiplier
		c.Core.QueueAttackEvent(&info.AttackEvent{
			Info:        ai,
			Pattern:     ap,
			Snapshot:    c.Snapshot(&ai),
			SourceFrame: src,
		}, 0)
	}, delay)
}
