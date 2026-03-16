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
	basicChargeWindup        = 20
	basicChargeHitmark       = 44
	phantasmPostAttackCancel = 69
	phantasmPostChargeCancel = 65
	slitherMaxDuration       = 150
	slitherActivationFrames  = 60
	slitherMinCancelFrames   = 24
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

func phantasmChargeCanQueueAfter(phaseStartOffset int) int {
	return phaseStartOffset + 1
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
	return 0
}

func (c *char) ActionStamina(a action.Action, p map[string]int) action.StaminaSpec {
	if a != action.ActionCharge {
		timing := action.StaminaConsumeOnExec
		if a == action.ActionDash {
			timing = action.StaminaConsumeByAbility
		}
		cost := c.Character.ActionStam(a, p)
		return action.StaminaSpec{
			Requirement: cost,
			Consume:     cost,
			Timing:      timing,
		}
	}
	if c.canTriggerPhantasm() {
		return action.StaminaSpec{
			Timing: action.StaminaConsumeByAbility,
		}
	}
	// Entering Slither is a hard prerequisite for any non-Phantasm charged-attack
	// route. Based on current in-game observation, Nefer needs roughly one second
	// worth of Slither stamina available to begin that route, but this threshold is
	// not consumed up front. Actual stamina consumption still happens only through
	// the Slither tick drain once Slither has non-zero duration.
	return action.StaminaSpec{
		Requirement: c.slitherActivationThreshold(),
		Timing:      action.StaminaConsumeByAbility,
	}
}

func (c *char) slitherActivationThreshold() float64 {
	req := slitherStamTickCost * slitherActivationFrames
	req *= 1 + c.Core.Player.StamPercentMod(action.ActionCharge)
	if req < 0 {
		return 0
	}
	return req
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
			case action.ActionAttack:
				return basicChargeWindup + phantasmPostAttackCancel
			case action.ActionCharge:
				return basicChargeWindup + phantasmPostChargeCancel
			default:
				return basicChargeWindup + phantasmAnimationLength
			}
		},
		AnimationLength: basicChargeWindup + phantasmAnimationLength,
		CanQueueAfter:   basicChargeWindup + 1,
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
				case action.ActionAttack:
					return phaseStartOffset + phantasmPostAttackCancel
				case action.ActionCharge:
					return phaseStartOffset + phantasmPostChargeCancel
				default:
					return phaseStartOffset + phantasmAnimationLength
				}
			}
			if c.chargeRoute.releaseStartFrame > 0 {
				phaseStartOffset := c.chargeRoute.releaseStartFrame - src
				return phaseStartOffset + chargeFrames[next]
			}
			switch next {
			case action.ActionCharge:
				return windup + slitherDuration + chargeFrames[next]
			case action.ActionAttack, action.ActionSkill, action.ActionBurst, action.ActionDash, action.ActionJump, action.ActionSwap, action.ActionWalk:
				return windup + slitherMinCancelFrames
			default:
				return windup + slitherDuration + chargeFrames[action.InvalidAction]
			}
		},
		AnimationLength: windup + slitherDuration + chargeFrames[action.InvalidAction],
		CanQueueAfter:   windup + 1,
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
	phaseStartOffset := c.chargeRoute.phantasmStartFrame - src
	c.Core.Player.SetActionLength(c.chargeRoute.phantasmEndFrame-src, phantasmChargeCanQueueAfter(phaseStartOffset))
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
	offset := c.chargeRoute.releaseStartFrame - src
	c.Core.Player.SetActionLength(offset+chargeFrames[action.InvalidAction], offset+basicChargeCanQueueAfter())
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
		Mult:       charge[0][c.TalentLvlAttack()],
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
		c.Core.Player.ConsumeVerdantDew(1)
		c.absorbSeedsOfDeceit()
	}, consumeFrame)

	shadeScaleBonus := c.c1ShadeScaleBonus()
	phantasmVeilMultiplier := 1 + c.phantasmVeilBonus()

	neferHit1 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Phantasm Performance (Nefer 1)",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm[0][c.TalentLvlSkill()] * phantasmVeilMultiplier,
		FlatDmg:    c.Stat(attributes.EM) * phantasm[1][c.TalentLvlSkill()] * phantasmVeilMultiplier,
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
		Mult:       phantasm[2][c.TalentLvlSkill()] * phantasmVeilMultiplier,
		FlatDmg:    c.Stat(attributes.EM) * phantasm[3][c.TalentLvlSkill()] * phantasmVeilMultiplier,
	}
	if c.Base.Cons >= 6 {
		neferHit2.AttackTag = attacks.AttackTagDirectLunarBloom
		neferHit2.Durability = 0
		neferHit2.UseEM = true
		neferHit2.IgnoreDefPercent = 1
		neferHit2.Mult = c6PhantasmHit2EM * phantasmVeilMultiplier
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
		Mult:             (phantasm[4][c.TalentLvlSkill()] + shadeScaleBonus) * phantasmVeilMultiplier,
	}
	shadeHit2 := shadeHit1
	shadeHit2.Abil = "Phantasm Performance (Shade 2)"
	shadeHit2.Mult = (phantasm[5][c.TalentLvlSkill()] + shadeScaleBonus) * phantasmVeilMultiplier
	shadeHit3 := shadeHit1
	shadeHit3.Abil = "Phantasm Performance (Shade 3)"
	shadeHit3.Mult = (phantasm[6][c.TalentLvlSkill()] + shadeScaleBonus) * phantasmVeilMultiplier

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(neferHit1, ap, phantasmHit1, phantasmHit1)
	c.Core.QueueAttack(shadeHit1, ap, phantasmHit2, phantasmHit2)
	c.Core.QueueAttack(shadeHit2, ap, phantasmHit3, phantasmHit3)
	c.Core.QueueAttack(neferHit2, ap, phantasmHit4, phantasmHit4)
	c.Core.QueueAttack(shadeHit3, ap, phantasmHit5, phantasmHit5)
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
			Mult:             c6PhantasmEndEM * phantasmVeilMultiplier,
		}
		c.Core.QueueAttack(c6EndHit, ap, phantasmAnimationLength, phantasmAnimationLength)
	}
}
