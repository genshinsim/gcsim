package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames [][]int
var aimedA1Frames []int

var aimedHitmarks = []int{15 - 12, 15, 86}

const aimedA1Hitmark = 50

func init() {
	// outside of E status
	aimedFrames = make([][]int, 3)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(25 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(25)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot
	aimedFrames[2] = frames.InitAbilSlice(96)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// Fully-Charged Aimed Shot (Crowfeather)
	aimedA1Frames = frames.InitAbilSlice(60)
	aimedA1Frames[action.ActionDash] = aimedA1Hitmark
	aimedA1Frames[action.ActionJump] = aimedA1Hitmark
}

// Aimed charge attack damage queue generator
// Additionally handles crowfeather state, E skill damage, and A4
// A4 effect is: When Tengu Juurai: Ambush hits opponents, Kujou Sara will restore 1.2 Energy to all party members for every 100% Energy Recharge she has. This effect can be triggered once every 3s.
// Has two parameters, "travel", used to set the number of frames that the arrow is in the air (default = 10)
// weak_point, used to determine if an arrow is hitting a weak point (default = 1 for true)
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 2
	}
	if hold > 2 {
		hold = 2
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Electro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	var a action.ActionInfo

	// A1:
	// While in the Crowfeather Cover state provided by Tengu Stormcall, Aimed Shot charge times are decreased by 60%.
	if c.Core.Status.Duration(coverKey) > 0 && hold == 2 {
		ai.Abil += " (A1)"
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedA1Frames),
			AnimationLength: aimedA1Frames[action.InvalidAction],
			CanQueueAfter:   aimedA1Hitmark,
			State:           action.AimState,
		}
	} else {
		if hold < 2 {
			ai.Abil = "Aimed Shot"
			if hold == 0 {
				ai.Abil += " (ARCC)"
			}
			ai.AttackTag = combat.AttackTagExtra
			ai.Element = attributes.Physical
			ai.Mult = aim[c.TalentLvlAttack()]
		}
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedFrames[hold]),
			AnimationLength: aimedFrames[hold][action.InvalidAction],
			CanQueueAfter:   aimedHitmarks[hold],
			State:           action.AimState,
		}
	}

	c.Core.QueueAttack(ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
		a.CanQueueAfter,
		a.CanQueueAfter+travel,
	)

	if c.Core.Status.Duration(coverKey) > 0 && hold == 2 {
		// Cover state handling - drops crowfeather, which explodes after 1.5 seconds
		// Not sure what kind of strike type this is
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}

		//TODO: snapshot?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
			aimedA1Hitmark,
			aimedA1Hitmark+travel+90,
			c.a4,
		)
		c.attackBuff(aimedA1Hitmark + travel + 90)

		// Particles are emitted after the ambush thing hits
		c.Core.QueueParticle("sara", 3, attributes.Electro, aimedA1Hitmark+travel+90+c.ParticleDelay)

		c.Core.Status.Delete(coverKey)
	}

	return a
}
