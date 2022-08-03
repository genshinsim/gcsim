package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstAnimation = 100
const burstHitmark = 82
const burstFirstTick = 140

func init() {
	burstFrames = frames.InitAbilSlice(burstAnimation)
	burstFrames[action.ActionAttack] = 92
	burstFrames[action.ActionSkill] = 92
	burstFrames[action.ActionDash] = 92
	burstFrames[action.ActionJump] = 92
	burstFrames[action.ActionSwap] = 95
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.qInfuse = attributes.NoElement
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Kazuha Slash",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Anemo,
		Durability:         50,
		Mult:               burstSlash[c.TalentLvlBurst()],
		HitlagHaltFrames:   0.05 * 60,
		HitlagFactor:       0.05,
		CanBeDefenseHalted: false,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy), 0, burstHitmark)

	//apply dot and check for absorb
	ai.Abil = "Kazuha Slash (Dot)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.Durability = 25
	// no more hitlag after initial slash
	ai.HitlagHaltFrames = 0
	snap := c.Snapshot(&ai)

	aiAbsorb := ai
	aiAbsorb.Abil = "Kazuha Slash (Absorb Dot)"
	aiAbsorb.Mult = burstEleDot[c.TalentLvlBurst()]
	aiAbsorb.Element = attributes.NoElement
	snapAbsorb := c.Snapshot(&aiAbsorb)

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(310/18)), 10)

	//from kisa's count: ticks starts at 147, + 117 gap each roughly; 5 ticks total
	//updated to 140 based on koli's count: https://docs.google.com/spreadsheets/d/1uEbP13O548-w_nGxFPGsf5jqj1qGD3pqFZ_AiV4w3ww/edit#gid=775340159

	// make sure that this task gets executed:
	// - after q initial hit hitlag happened
	// - before kazuha can get affected by any more hitlag
	c.QueueCharTask(func() {
		for i := 0; i < 5; i++ {
			c.Core.Tasks.Add(func() {
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 0)
				if c.qInfuse != attributes.NoElement {
					aiAbsorb.Element = c.qInfuse
					c.Core.QueueAttackWithSnap(aiAbsorb, snapAbsorb, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 0)
				}
			}, (burstFirstTick-(burstHitmark+5))+117*i)
		}

		//add em to kazuha even if off-field
		//add em to all char, but only activate if char is active
		//The Autumn Whirlwind field created by Kazuha Slash has the following effects:
		//• Increases Kaedehara Kazuha's own Elemental Mastery by 200 for its duration.
		//• Increases the Elemental Mastery of characters within the field by 200.
		//The Elemental Mastery-increasing effects of this Constellation do not stack.
		if c.Base.Cons >= 2 {
			// TODO: Lasts while Q field is on stage is ambiguous.
			// Does it apply to Kazuha's initial hit?
			// Not sure when it lasts from and until
			// For consistency with how it was previously done, assume that it lasts from button press to the last tick
			for _, char := range c.Core.Player.Chars() {
				this := char
				//use non hitlag since it's from the field?
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("kazuha-c2", (burstFirstTick-(burstHitmark+5))+117*5),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						switch this.Index {
						case c.Core.Player.Active(), c.Index:
							return c.c2buff, true
						}
						return nil, false
					},
				})
			}
		}
	}, burstHitmark+5)

	//reset skill cd
	if c.Base.Cons > 0 {
		c.ResetActionCooldown(action.ActionSkill)
	}

	if c.Base.Cons == 6 {
		c.c6()
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstAnimation,
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.qInfuse != attributes.NoElement {
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), 18)
	}
}
