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

// based on sayu/thoma frames
// TODO: update frames
var (
	skillPressFrames []int
	skillHoldFrames  []int
)

const (
	skillPressHitmark = 11

	skillHoldCDStart     = 648
	skillHoldKickHitmark = 667

	kickParticleICDKey = "kirara-kick-particle-icd"
	rollParticleICDKey = "kirara-roll-particle-icd"
)

func init() {
	// Tap E
	skillPressFrames = frames.InitAbilSlice(46)
	skillPressFrames[action.ActionDash] = 32
	skillPressFrames[action.ActionJump] = 32
	skillPressFrames[action.ActionSwap] = 44

	// Hold E
	skillHoldFrames = frames.InitAbilSlice(709) // Hold E -> N1/Q
	skillHoldFrames[action.ActionDash] = 708    // Hold E -> J
	skillHoldFrames[action.ActionJump] = 708    // Hold E -> J
	skillHoldFrames[action.ActionSwap] = 706    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold > 0 {
		if hold > 10*60 {
			hold = 10 * 60
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

func (c *char) skillPress(p map[string]int) action.ActionInfo {
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
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: 1}, 3, 270),
		skillPressHitmark,
		skillPressHitmark,
		c.kickParticleCB,
	)

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.QueueCharTask(c.generateSkillShield, 7*2)
	c.SetCDWithDelay(action.ActionSkill, 8*60, 7*2)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash],
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.ActionInfo {
	// ticks
	d := c.createSkillHoldSnapshot()

	for i := 0; i <= duration; i += 20 * 2 { // 1 tick for sure
		c.Core.Tasks.Add(func() {
			// pattern shouldn't snapshot on attack event creation because the skill follows the player
			d.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
			c.Core.QueueAttackEvent(d, 0)
		}, 18+i)
	}

	// Flipclaw Strike DMG
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Flipclaw Strike",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Durability:       25,
		Mult:             flipclawDmg[c.TalentLvlSkill()],
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

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.QueueCharTask(c.generateSkillShield, 7*2)
	c.SetCDWithDelay(action.ActionSkill, int(8*60+float64(duration)*0.4), (skillHoldCDStart-600)+duration)

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
		Abil:             "Urgent Neko Parcel",
		AttackTag:        attacks.AttackTagElementalArtHold,
		ICDTag:           attacks.ICDTagElementalArt,
		ICDGroup:         attacks.ICDGroupDefault, // TODO: Crash?
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Durability:       25,
		Mult:             catDmg[c.TalentLvlSkill()],
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
	if c.Base.Ascension >= 1 {
		ae.Callbacks = append(ae.Callbacks, c.a1)
	}
	return &ae
}

func (c *char) generateSkillShield() {
	c.genShield("Shield of Safe Transport", c.shieldHP())

	player, ok := c.Core.Combat.Player().(*avatar.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}
	player.ApplySelfInfusion(attributes.Dendro, 25, 30)
}
