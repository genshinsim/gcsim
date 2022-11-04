package sayu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillPressFrames []int
var skillShortHoldFrames []int
var skillHoldFrames []int

const skillPressDoTHitmark = 7 // 1 DoT tick
const skillPressCDStart = 14
const skillPressKickHitmark = 25

const skillShortHoldCDStart = 30
const skillShortHoldKickHitmark = 48

const skillHoldCDStart = 648
const skillHoldKickHitmark = 667

func init() {
	// Tap E
	skillPressFrames = frames.InitAbilSlice(44) // Tap E -> N1/D
	skillPressFrames[action.ActionBurst] = 43   // Tap E -> Q
	skillPressFrames[action.ActionJump] = 45    // Tap E -> J
	skillPressFrames[action.ActionSwap] = 42    // Tap E -> Swap

	// Short Hold E
	skillShortHoldFrames = frames.InitAbilSlice(89) // Short Hold E -> N1/Q/D/J
	skillShortHoldFrames[action.ActionSwap] = 87    // Short Hold E -> Swap

	// Hold E
	skillHoldFrames = frames.InitAbilSlice(709) // Hold E -> N1/Q
	skillHoldFrames[action.ActionDash] = 708    // Hold E -> J
	skillHoldFrames[action.ActionJump] = 708    // Hold E -> J
	skillHoldFrames[action.ActionSwap] = 706    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	short_hold := p["short_hold"]
	if p["short_hold"] != 0 {
		short_hold = 1
	}
	if short_hold == 1 {
		return c.skillShortHold(p)
	}

	hold := p["hold"]
	if hold > 0 {
		if hold > 600 { // 10s
			hold = 600
		}
		return c.skillHold(p, hold)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {
	c.c2Bonus = 0.033

	// Fuufuu Windwheel DMG
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fuufuu Windwheel (DoT Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArtAnemo,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 0.1), skillPressDoTHitmark)

	// Fuufuu Whirlwind Kick Press DMG
	ai = combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Press)",
		AttackTag:        combat.AttackTagElementalArt,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillPressEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap = c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 0.5), skillPressKickHitmark)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, skillPressKickHitmark+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 6*60, skillPressCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillShortHold(p map[string]int) action.ActionInfo {
	c.eDuration = c.Core.F + skillShortHoldKickHitmark
	c.c2Bonus = .0

	c.eAbsorb = attributes.NoElement
	c.eAbsorbTag = combat.ICDTagNone
	c.absorbCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 0.1)

	// 1 DoT Tick
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, 1), 18)

	c.Core.Tasks.Add(func() {
		c.Core.QueueAttackEvent(d, 0)

		if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
			c.c2Bonus += 0.033
			c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", glog.LogCharacterEvent, c.Index).
				Write("dmg bonus%", c.c2Bonus)
		}
	}, 18)
	c.Core.QueueParticle("sayu-skill-hold", 1, attributes.Anemo, 18+c.ParticleDelay)

	// Fuufuu Whirlwind Kick Hold DMG
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Hold)",
		AttackTag:        combat.AttackTagElementalArtHold,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillHoldEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 0.5), skillShortHoldKickHitmark)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, skillShortHoldKickHitmark+c.ParticleDelay)

	// 6.2s cooldown
	c.SetCDWithDelay(action.ActionSkill, 372, skillShortHoldCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillShortHoldFrames),
		AnimationLength: skillShortHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillShortHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.ActionInfo {
	c.eDuration = c.Core.F + (skillHoldKickHitmark - 600) + duration
	c.c2Bonus = .0

	c.eAbsorb = attributes.NoElement
	c.eAbsorbTag = combat.ICDTagNone
	c.absorbCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 0.1)

	// ticks
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, int(duration/12)), 18)

	for i := 0; i <= duration; i += 30 { // 1 tick for sure
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackEvent(d, 0)

			if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
				c.c2Bonus += 0.033
				c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", glog.LogCharacterEvent, c.Index).
					Write("dmg bonus%", c.c2Bonus)
			}
		}, 18+i)

		if i%180 == 0 { // 3s
			c.Core.QueueParticle("sayu-skill-hold", 1, attributes.Anemo, 18+i+c.ParticleDelay)
		}
	}

	// Fuufuu Whirlwind Kick Hold DMG
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Hold)",
		AttackTag:        combat.AttackTagElementalArtHold,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillHoldEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 0.5), (skillHoldKickHitmark-600)+duration)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, (skillHoldKickHitmark-600)+duration+c.ParticleDelay)

	// +2 frames for not proc the sacrificial by "Yoohoo Art: Fuuin Dash (Elemental DMG)"
	c.SetCDWithDelay(action.ActionSkill, int(6*60+float64(duration)*0.5), (skillHoldCDStart-600)+duration+2)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillHoldFrames[next] - 600 + duration },
		AnimationLength: skillHoldFrames[action.InvalidAction] - 600 + duration,
		CanQueueAfter:   skillHoldFrames[action.ActionSwap] - 600 + duration, // earliest cancel
		State:           action.SkillState,
	}
}

// TODO: is this helper needed?
func (c *char) createSkillHoldSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Windwheel (DoT Hold)",
		AttackTag:        combat.AttackTagElementalArt,
		ICDTag:           combat.ICDTagElementalArtAnemo,
		ICDGroup:         combat.ICDGroupDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.01 * 60,
		HitlagFactor:     0.05,
		IsDeployable:     true,
	}
	snap := c.Snapshot(&ai)

	return (&combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewCircleHit(c.Core.Combat.Player(), 0.5),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
}

func (c *char) absorbCheck(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}

		c.eAbsorb = c.Core.Combat.AbsorbCheck(c.absorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if c.eAbsorb != attributes.NoElement {
			switch c.eAbsorb {
			case attributes.Pyro:
				c.eAbsorbTag = combat.ICDTagElementalArtPyro
			case attributes.Hydro:
				c.eAbsorbTag = combat.ICDTagElementalArtHydro
			case attributes.Electro:
				c.eAbsorbTag = combat.ICDTagElementalArtElectro
			case attributes.Cryo:
				c.eAbsorbTag = combat.ICDTagElementalArtCryo
			}
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"sayu absorbed ", c.eAbsorb.String(),
			)
			return
		}
		c.Core.Tasks.Add(c.absorbCheck(src, count+1, max), 12)
	}
}

func (c *char) rollAbsorb() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		if atk.Info.Element != attributes.Anemo || c.eAbsorb == attributes.NoElement {
			return false
		}
		if c.Core.F > c.eDuration {
			return false
		}

		switch atk.Info.AttackTag {
		case combat.AttackTagElementalArt:
			// DoT Elemental DMG
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Fuufuu Windwheel Elemental (Elemental DoT Hold)",
				AttackTag:  combat.AttackTagElementalArtHold,
				ICDTag:     c.eAbsorbTag,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    c.eAbsorb,
				Durability: 25,
				Mult:       skillAbsorb[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1), 1, 1)
		case combat.AttackTagElementalArtHold:
			// Kick Elemental DMG
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Fuufuu Whirlwind Elemental (Elemental Kick Hold)",
				AttackTag:  combat.AttackTagElementalArt,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    c.eAbsorb,
				Durability: 25,
				Mult:       skillAbsorbEnd[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1), 1, 1)
		}

		return false
	}, "sayu-absorb-check")
}
