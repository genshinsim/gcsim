package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var SkillACDStarts = 32
var SkillBCDStarts = 32

var SkillAHitmarks = 32
var SkillBHitmarks = 32

var SkillAFrames []int
var SkillBFrames []int

const ()

func init() {
	// Tap E
	SkillAFrames = make([]int, 2)

	// outside of Q
	SkillAFrames = frames.InitAbilSlice(37)

	// Furry E
	SkillBFrames = make([]int, 2)

	// inside of Q
	SkillBFrames = frames.InitAbilSlice(37) // Hold E -> N1

}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// check if Q is up for different E
	if c.StatusIsActive(burstKey) {
		return c.SkillB() //SkillB is Mortuary Rite (skill during burst)
	}

	return c.SkillA() //SkillA is normal, non burst boosted skill
}

func (c *char) SkillA() action.ActionInfo {
	//TODO: Adjust the attack frame values (this ones are source: i made them the fuck up)
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Secret Rite: Chasmic Soulfarer",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skillA[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.1 * 60,
		HitlagFactor:       0.03,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		SkillAHitmarks,
		SkillAHitmarks,
	)

	c.Core.QueueParticle("cyno", 3, attributes.Electro, SkillAHitmarks+c.ParticleDelay)
	cd := 7.5 * 60
	c.SetCDWithDelay(action.ActionSkill, int(cd), SkillACDStarts)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(SkillAFrames),
		AnimationLength: SkillAFrames[action.InvalidAction],
		CanQueueAfter:   SkillAFrames[action.ActionDash], // earliest cancel is 1f before SkillAHitmark
		State:           action.SkillState,
	}
}

func (c *char) SkillB() action.ActionInfo {
	//TODO: Adjust the attack frame values (this ones are source: i made them the fuck up)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mortuary Rite",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillB[c.TalentLvlSkill()],
	}

	if !c.StatusIsActive(a4key) { //check for endseer buff
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy),
			SkillBHitmarks,
			SkillBHitmarks,
		)
	} else {
		//apply the extra damage on skill
		c.judiscation()

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy),
			SkillBHitmarks,
			SkillBHitmarks,
		)
		//Apply the extra hit
		ai.Abil = "Duststalker Bolt"
		ai.Mult = 1.0
		ai.FlatDmg = c.Stat(attributes.EM) * 2.5 //this is the A4
		//3 instances
		//TODO: timing (frames) of each instance
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
				SkillBHitmarks,
				SkillBHitmarks,
			)
		}

	}

	c.ExtendStatus(burstKey, 4)

	var count float64 = 1 //33% of generating 2 on furry form
	if c.Core.Rand.Float64() < .33 {
		count++
	}
	c.Core.QueueParticle("cyno", count, attributes.Electro, SkillBHitmarks+c.ParticleDelay)

	cd := 3 * 60
	c.SetCDWithDelay(action.ActionSkill, int(cd), SkillBCDStarts)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(SkillBFrames),
		AnimationLength: SkillBFrames[action.InvalidAction],
		CanQueueAfter:   SkillBFrames[action.ActionJump], // earliest cancel is 3f before SkillBHitmark
		State:           action.SkillState,
	}
}

// TODO:I am applying the skill dmg bonus this way to ensure that the skill gets the dmg bonus even if endseer expires mid cast (this may not be neccesary)
func (c *char) judiscation() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.35
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("cyno-a1", 60), //1 second should be enough to be applied... TODO: this is scuff af
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
				return nil, false
			}
			return m, true
		},
	})

}
