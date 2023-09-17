package kirara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillPressFrames     []int
	skillShortHoldFrames []int
	skillHoldFrames      []int
)

const (
	skillPressHitmark = 14

	skillShortHoldCDStart     = 27
	skillShortHoldKickHitmark = 48

	skillHoldCDStart     = 614
	skillHoldKickHitmark = 636

	kickParticleICDKey = "kirara-kick-particle-icd"
	rollParticleICDKey = "kirara-roll-particle-icd"
)

func init() {
	// Tap E
	skillPressFrames = frames.InitAbilSlice(38) // E -> Walk
	skillPressFrames[action.ActionAttack] = 34
	skillPressFrames[action.ActionSkill] = 34
	skillPressFrames[action.ActionBurst] = 34
	skillPressFrames[action.ActionDash] = 35
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionSwap] = 33

	// Short Hold E
	skillShortHoldFrames = frames.InitAbilSlice(79) // Short Hold E -> Walk
	skillShortHoldFrames[action.ActionAttack] = 72
	skillShortHoldFrames[action.ActionSkill] = 75
	skillShortHoldFrames[action.ActionBurst] = 74
	skillShortHoldFrames[action.ActionDash] = 74
	skillShortHoldFrames[action.ActionJump] = 74
	skillShortHoldFrames[action.ActionSwap] = 72

	// Hold E
	skillHoldFrames = frames.InitAbilSlice(668) // Hold E -> Walk
	skillHoldFrames[action.ActionAttack] = 659
	skillHoldFrames[action.ActionSkill] = 662
	skillHoldFrames[action.ActionBurst] = 663
	skillHoldFrames[action.ActionDash] = 662
	skillHoldFrames[action.ActionJump] = 663
	skillHoldFrames[action.ActionSwap] = 663
}

func (c *char) Skill(p map[string]int) action.Info {
	shortHold := p["short_hold"]
	if p["short_hold"] != 0 {
		shortHold = 1
	}
	if shortHold == 1 {
		return c.skillShortHold()
	}

	hold := p["hold"]
	if hold > 0 {
		if hold > 10*60 {
			hold = 10 * 60
		}
		return c.skillHold(hold)
	}
	return c.skillPress()
}

func (c *char) kickParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(kickParticleICDKey) {
		return
	}
	c.AddStatus(kickParticleICDKey, 0.3*60, true)
	c.Core.QueueParticle("kirara-kick", 3, attributes.Dendro, c.ParticleDelay)
}

func (c *char) rollParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(rollParticleICDKey) {
		return
	}
	c.AddStatus(rollParticleICDKey, 4*60, true)
	c.Core.QueueParticle("kirara-roll", 1, attributes.Dendro, c.ParticleDelay)
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Tail-Flicking Flying Kick",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 2.5, 270),
		skillPressHitmark,
		skillPressHitmark,
		c.kickParticleCB,
	)

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.QueueCharTask(c.generateSkillShield, skillPressHitmark)
	c.SetCDWithDelay(action.ActionSkill, 8*60, 12)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillShortHold() action.Info {
	// 1 tick
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(func() {
		// pattern shouldn't snapshot on attack event creation because the skill follows the player
		d.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
		c.Core.QueueAttackEvent(d, 0)
	}, 17)

	// Flipclaw Strike DMG
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Flipclaw Strike",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               flipclawDmg[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.1 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 2.6),
		skillShortHoldKickHitmark,
		c.kickParticleCB,
	)

	if c.Base.Ascension >= 1 {
		c.QueueCharTask(c.a1, skillShortHoldCDStart-1)
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.QueueCharTask(c.generateSkillShield, 14)
	c.SetCDWithDelay(action.ActionSkill, 8.2*60, skillShortHoldCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillShortHoldFrames),
		AnimationLength: skillShortHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillShortHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(duration int) action.Info {
	// ticks
	d := c.createSkillHoldSnapshot()

	for i := 16; i <= duration+12; i += 0.5 * 60 {
		c.Core.Tasks.Add(func() {
			// pattern shouldn't snapshot on attack event creation because the skill follows the player
			d.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
			c.Core.QueueAttackEvent(d, 0)
		}, i)
	}

	// Flipclaw Strike DMG
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Flipclaw Strike",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               flipclawDmg[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.1 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 2.6),
		(skillHoldKickHitmark-600)+duration,
		c.kickParticleCB,
	)

	if c.Base.Ascension >= 1 {
		c.QueueCharTask(c.a1, (skillHoldCDStart-600)+duration)
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	cd := 8*60 + duration/30*12
	c.QueueCharTask(c.generateSkillShield, 14)
	c.SetCDWithDelay(action.ActionSkill, cd, (skillHoldCDStart-600)+duration)

	return action.Info{
		Frames:          func(next action.Action) int { return skillHoldFrames[next] - 600 + duration },
		AnimationLength: skillHoldFrames[action.InvalidAction] - 600 + duration,
		CanQueueAfter:   skillHoldFrames[action.ActionAttack] - 600 + duration, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) createSkillHoldSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Urgent Neko Parcel",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagElementalArtHold,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       catDmg[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	// pattern shouldn't snapshot on attack event creation because the skill follows the player
	ae := combat.AttackEvent{
		Info:        ai,
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	}
	ae.Callbacks = append(ae.Callbacks, c.rollParticleCB)
	if c.Base.Ascension >= 1 {
		ae.Callbacks = append(ae.Callbacks, c.a1StackGain)
	}
	return &ae
}

func (c *char) generateSkillShield() {
	c.genShield("Shield of Safe Transport", c.shieldHP())

	player, ok := c.Core.Combat.Player().(*avatar.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}
	player.ApplySelfInfusion(attributes.Dendro, 25, 0.1*60)
}
