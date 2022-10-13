package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const skillBName = "Mortuary Rite"

var (
	skillHitmark  = 21
	skillBHitmark = 28
	skillFrames   []int
	skillBFrames  []int
	skillA1Frames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 42

	// burst frames
	skillBFrames = frames.InitAbilSlice(34)
	skillBFrames[action.ActionDash] = 30
	skillBFrames[action.ActionJump] = 31
	skillBFrames[action.ActionSwap] = 33

	// skill has different frame data with endseers
	skillA1Frames = frames.InitAbilSlice(35)
	skillA1Frames[action.ActionDash] = 30
	skillA1Frames[action.ActionJump] = 31
	skillA1Frames[action.ActionSwap] = 33
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// check if Q is up for different E
	if c.StatusIsActive(burstKey) {
		return c.skillB() // SkillB is Mortuary Rite (skill during burst)
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Secret Rite: Chasmic Soulfarer",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0,
		CanBeDefenseHalted: false,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		skillHitmark,
		skillHitmark,
	)

	c.Core.QueueParticle("cyno", 3, attributes.Electro, skillHitmark+c.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 450, 17)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillB() action.ActionInfo {
	c.tryBurstPPSlide(skillBHitmark)

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               skillBName,
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skillB[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.03 * 60,
		CanBeDefenseHalted: false,
	}

	if !c.StatusIsActive(a1Key) { // check for endseer buff
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy),
			skillBHitmark,
			skillBHitmark,
		)
	} else {
		// apply the extra damage on skill
		c.judiscation()
		if c.Base.Cons >= 1 && c.StatusIsActive(c1Key) {
			c.c1()
		}
		if c.Base.Cons >= 6 { // constellation 6 giving 4 stacks on judication
			c.c6stacks += 4
			c.AddStatus("cyno-c6", 480, true) // 8s*60
			if c.c6stacks > 8 {
				c.c6stacks = 8
			}
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy),
			skillBHitmark,
			skillBHitmark,
		)
		// Apply the extra hit
		ai.Abil = "Duststalker Bolt"
		ai.Mult = 1.0
		ai.FlatDmg = c.Stat(attributes.EM) * 2.5 // this is the A4
		ai.ICDTag = combat.ICDTagCynoBolt
		ai.ICDGroup = combat.ICDGroupCynoBolt
		ai.IsDeployable = true

		// 3 instances
		// TODO: timing (frames) of each instance
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
				skillBHitmark,
				skillBHitmark,
			)
		}

	}
	if c.burstExtension < 2 { // burst can only be extended 2 times per burst cycle (up to 18s, 10s base and +4 each time)
		c.ExtendStatus(burstKey, 240) // 4s*60
		c.burstExtension++
	}

	var count float64 = 1 // 33% of generating 2 on furry form
	if c.Core.Rand.Float64() < .33 {
		count++
	}
	c.Core.QueueParticle("cyno", count, attributes.Electro, skillBHitmark+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 180, 26)

	f := skillBFrames
	if c.StatusIsActive(a1Key) {
		f = skillA1Frames
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(f),
		AnimationLength: skillBFrames[action.InvalidAction],
		CanQueueAfter:   skillBFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

// TODO:I am applying the skill dmg bonus this way to ensure that the skill gets the dmg bonus even if endseer expires mid cast (this may not be neccesary)
func (c *char) judiscation() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.35
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("cyno-a1-dmg", 60), // 1 second should be enough to be applied... TODO: this is scuff af
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			// cyno's a1 only buffs the skill dmg, not bolts. ugly hack, but gets the job done
			if atk.Info.Abil != skillBName {
				return nil, false
			}
			return m, true
		},
	})
}
