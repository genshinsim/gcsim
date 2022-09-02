package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var SkillACDStarts = []int{30, 31}
var SkillBCDStarts = []int{52, 52}

var SkillAHitmarks = []int{32, 33}
var SkillBHitmarks = []int{55, 55}

var SkillAFrames [][]int
var SkillBFrames [][]int

const ()

func init() {
	// Tap E
	SkillAFrames = make([][]int, 2)

	// outside of Q
	SkillAFrames[0] = frames.InitAbilSlice(74) // Tap E -> Swap
	SkillAFrames[0][action.ActionAttack] = 70  // Tap E -> N1
	SkillAFrames[0][action.ActionBurst] = 69   // Tap E -> Q
	SkillAFrames[0][action.ActionDash] = 31    // Tap E -> D
	SkillAFrames[0][action.ActionJump] = 31    // Tap E -> J

	// inside of Q
	SkillAFrames[1] = frames.InitAbilSlice(76) // Tap E -> Swap
	SkillAFrames[1][action.ActionSwap] = 75    // Tap E -> N1
	SkillAFrames[1][action.ActionDash] = 32    // Tap E -> D
	SkillAFrames[1][action.ActionJump] = 32    // Tap E -> J

	// Hold E
	SkillBFrames = make([][]int, 2)

	// outside of Q
	SkillBFrames[0] = frames.InitAbilSlice(103) // Hold E -> Q
	SkillBFrames[0][action.ActionAttack] = 102  // Hold E -> N1
	SkillBFrames[0][action.ActionDash] = 52     // Hold E -> D
	SkillBFrames[0][action.ActionJump] = 52     // Hold E -> J
	SkillBFrames[0][action.ActionSwap] = 91     // Hold E -> Swap

	// inside of Q
	SkillBFrames[1] = frames.InitAbilSlice(96) // Hold E -> N1
	SkillBFrames[1][action.ActionDash] = 53    // Hold E -> D
	SkillBFrames[1][action.ActionJump] = 52    // Hold E -> J
	SkillBFrames[1][action.ActionSwap] = 88    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// check if Q is up for different E frames
	burstActive := 0
	if c.StatusIsActive(burstKey) {
		burstActive = 1
		return c.SkillB(burstActive) //SkillB is Mortuary Rite (skill during burst)
	}

	return c.SkillA(burstActive) //SkillA is normal, non burst boosted skill
}

func (c *char) SkillA(burstActive int) action.ActionInfo {
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
		SkillAHitmarks[burstActive],
		SkillAHitmarks[burstActive],
	)

	c.Core.QueueParticle("cyno", 3, attributes.Electro, SkillAHitmarks[burstActive]+c.ParticleDelay)
	cd := 7.5 * 60
	c.SetCDWithDelay(action.ActionSkill, int(cd), SkillACDStarts[burstActive])

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(SkillAFrames[burstActive]),
		AnimationLength: SkillAFrames[burstActive][action.InvalidAction],
		CanQueueAfter:   SkillAFrames[burstActive][action.ActionDash], // earliest cancel is 1f before SkillAHitmark
		State:           action.SkillState,
	}
}

func (c *char) SkillB(burstActive int) action.ActionInfo {
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
			SkillBHitmarks[burstActive],
			SkillBHitmarks[burstActive],
		)
	} else {
		//apply the extra damage on skill
		c.judiscation()

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy),
			SkillBHitmarks[burstActive],
			SkillBHitmarks[burstActive],
		)
		//Apply the extra hit
		ai.Abil = "Duststalker Bolt"
		ai.Mult = 0.5
		ai.FlatDmg = c.Stat(attributes.EM) * 2.5 //this is the A4
		//3 instances
		//TODO: timing (frames) of each instance
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
				SkillBHitmarks[burstActive],
				SkillBHitmarks[burstActive],
			)
		}

	}

	c.ExtendStatus(burstKey, 4)

	var count float64 = 1 //33% of generating 2 on furry form
	if c.Core.Rand.Float64() < .33 {
		count++
	}
	c.Core.QueueParticle("cyno", count, attributes.Electro, SkillAHitmarks[burstActive]+c.ParticleDelay)

	cd := 3 * 60
	c.SetCDWithDelay(action.ActionSkill, int(cd), SkillACDStarts[burstActive])

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(SkillBFrames[burstActive]),
		AnimationLength: SkillBFrames[burstActive][action.InvalidAction],
		CanQueueAfter:   SkillBFrames[burstActive][action.ActionJump], // earliest cancel is 3f before SkillBHitmark
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
