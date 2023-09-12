package sayu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillPressFrames     []int
	skillShortHoldFrames []int
	skillHoldFrames      []int
)

const (
	skillPressDoTHitmark  = 7 // 1 DoT tick
	skillPressCDStart     = 14
	skillPressKickHitmark = 25

	skillShortHoldCDStart     = 30
	skillShortHoldKickHitmark = 48

	skillHoldCDStart     = 648
	skillHoldKickHitmark = 667

	kickParticleICDKey = "sayu-kick-particle-icd"
	rollParticleICDKey = "sayu-roll-particle-icd"
)

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

func (c *char) Skill(p map[string]int) action.Info {
	shortHold := p["short_hold"]
	if p["short_hold"] != 0 {
		shortHold = 1
	}
	if shortHold == 1 {
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

func (c *char) kickParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(kickParticleICDKey) {
		return
	}
	c.AddStatus(kickParticleICDKey, 0.5*60, true)
	c.Core.QueueParticle("sayu-kick", 2, attributes.Anemo, c.ParticleDelay)
}

func (c *char) rollParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(rollParticleICDKey) {
		return
	}
	c.AddStatus(rollParticleICDKey, 3*60, true)
	c.Core.QueueParticle("sayu-roll", 1, attributes.Anemo, c.ParticleDelay)
}

func (c *char) skillPress(p map[string]int) action.Info {
	c.c2Bonus = 0.033

	// Fuufuu Windwheel DMG
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fuufuu Windwheel (DoT Press)",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagElementalArtAnemo,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		skillPressDoTHitmark,
	)

	// Fuufuu Whirlwind Kick Press DMG
	ai = combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Press)",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillPressEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap = c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 2.5),
		skillPressKickHitmark,
		c.kickParticleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 6*60, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillShortHold(p map[string]int) action.Info {
	c.eDuration = c.Core.F + skillShortHoldKickHitmark
	c.c2Bonus = .0

	c.eAbsorb = attributes.NoElement
	c.eAbsorbTag = attacks.ICDTagNone
	c.absorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.2)

	// 1 DoT Tick
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, 1), 18)

	c.Core.Tasks.Add(func() {
		// pattern shouldn't snapshot on attack event creation because the skill follows the player
		d.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
		c.Core.QueueAttackEvent(d, 0)

		if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
			c.c2Bonus += 0.033
			c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", glog.LogCharacterEvent, c.Index).
				Write("dmg bonus%", c.c2Bonus)
		}
	}, 18)

	// Fuufuu Whirlwind Kick Hold DMG
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Hold)",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillHoldEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 3),
		skillShortHoldKickHitmark,
		c.kickParticleCB,
	)

	// 6.2s cooldown
	c.SetCDWithDelay(action.ActionSkill, 372, skillShortHoldCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillShortHoldFrames),
		AnimationLength: skillShortHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillShortHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.Info {
	c.eDuration = c.Core.F + (skillHoldKickHitmark - 600) + duration
	c.c2Bonus = .0

	c.eAbsorb = attributes.NoElement
	c.eAbsorbTag = attacks.ICDTagNone
	c.absorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.2)

	// ticks
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, int(duration/12)), 18)

	for i := 0; i <= duration; i += 30 { // 1 tick for sure
		c.Core.Tasks.Add(func() {
			// pattern shouldn't snapshot on attack event creation because the skill follows the player
			d.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
			c.Core.QueueAttackEvent(d, 0)

			if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
				c.c2Bonus += 0.033
				c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", glog.LogCharacterEvent, c.Index).
					Write("dmg bonus%", c.c2Bonus)
			}
		}, 18+i)
	}

	// Fuufuu Whirlwind Kick Hold DMG
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fuufuu Whirlwind (Kick Hold)",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillHoldEnd[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 3),
		(skillHoldKickHitmark-600)+duration,
		c.kickParticleCB,
	)

	// +2 frames for not proc the sacrificial by "Yoohoo Art: Fuuin Dash (Elemental DMG)"
	c.SetCDWithDelay(action.ActionSkill, int(6*60+float64(duration)*0.5), (skillHoldCDStart-600)+duration+2)

	return action.Info{
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
		AttackTag:        attacks.AttackTagElementalArtHold,
		ICDTag:           attacks.ICDTagElementalArtAnemo,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.01 * 60,
		HitlagFactor:     0.05,
		IsDeployable:     true,
	}
	snap := c.Snapshot(&ai)
	// pattern shouldn't snapshot on attack event creation because the skill follows the player
	ae := combat.AttackEvent{
		Info:        ai,
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	}
	ae.Callbacks = append(ae.Callbacks, c.rollParticleCB)
	return &ae
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
				c.eAbsorbTag = attacks.ICDTagElementalArtPyro
			case attributes.Hydro:
				c.eAbsorbTag = attacks.ICDTagElementalArtHydro
			case attributes.Electro:
				c.eAbsorbTag = attacks.ICDTagElementalArtElectro
			case attributes.Cryo:
				c.eAbsorbTag = attacks.ICDTagElementalArtCryo
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
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		if atk.Info.Element != attributes.Anemo || c.eAbsorb == attributes.NoElement {
			return false
		}
		if c.Core.F > c.eDuration {
			return false
		}

		switch atk.Info.AttackTag {
		// DoT always has ElementalArtHold tag
		case attacks.AttackTagElementalArtHold:
			// DoT Elemental DMG
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Fuufuu Windwheel Elemental (Elemental DoT Hold)",
				AttackTag:  attacks.AttackTagElementalArtHold,
				ICDTag:     c.eAbsorbTag,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    c.eAbsorb,
				Durability: 25,
				Mult:       skillAbsorb[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(e.Key()), 1, 1)
		// Kick always has ElementalArt tag
		case attacks.AttackTagElementalArt:
			// Kick Elemental DMG
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Fuufuu Whirlwind Elemental (Elemental Kick Hold)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    c.eAbsorb,
				Durability: 25,
				Mult:       skillAbsorbEnd[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(
				ai,
				combat.NewSingleTargetHit(e.Key()),
				1,
				1,
			)
		}

		return false
	}, "sayu-absorb-check")
}
