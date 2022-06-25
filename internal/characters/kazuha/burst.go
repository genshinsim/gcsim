package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstAnimation = 100
const burstHitmark = 82

func init() {
	burstFrames = frames.InitAbilSlice(burstAnimation)
	burstFrames[action.ActionAttack] = 95
	burstFrames[action.ActionSkill] = 95
	burstFrames[action.ActionSwap] = 95
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.qInfuse = attributes.NoElement
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kazuha Slash",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       burstSlash[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, burstHitmark)

	//apply dot and check for absorb
	ai.Abil = "Kazuha Slash (Dot)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.Durability = 25
	snap := c.Snapshot(&ai)

	aiAbsorb := ai
	aiAbsorb.Abil = "Kazuha Slash (Absorb Dot)"
	aiAbsorb.Mult = burstEleDot[c.TalentLvlBurst()]
	aiAbsorb.Element = attributes.NoElement
	snapAbsorb := c.Snapshot(&aiAbsorb)

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(310/18)), 10)

	//from kisa's count: ticks starts at 147, + 117 gap each roughly; 5 ticks total
	for i := 0; i < 5; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			if c.qInfuse != attributes.NoElement {
				aiAbsorb.Element = c.qInfuse
				c.Core.QueueAttackWithSnap(aiAbsorb, snapAbsorb, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			}
		}, 147+117*i)
	}

	//reset skill cd
	if c.Base.Cons > 0 {
		c.ResetActionCooldown(action.ActionSkill)
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
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 200
		for _, char := range c.Core.Player.Chars() {
			this := char
			char.AddStatMod("kazuha-c2", 147+117*5, attributes.EM, func() ([]float64, bool) {
				switch this.Index {
				case c.Core.Player.Active(), c.Index:
					return m, true
				}
				return nil, false
			})
		}
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + burstAnimation + 300
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			"kazuha-c6-infusion",
			attributes.Anemo,
			burstAnimation+300,
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		)
	}

	c.SetCDWithDelay(action.ActionBurst, 15*60, 7)
	c.ConsumeEnergy(7)

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

func (c *char) absorbCheckA1(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.a1Ele = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.a1Ele != attributes.NoElement {
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"kazuha a1 infused ", c.a1Ele.String(),
			)
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckA1(src, count+1, max), 6)
	}
}
