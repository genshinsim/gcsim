package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) Burst(p map[string]int) action.ActionInfo {
	f, a := c.ActionFrames(action.ActionBurst, p)

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

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, 82)

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

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(310/18)), "kaz-absorb-check", 10)

	//from kisa's count: ticks starts at 147, + 117 gap each roughly; 5 ticks total
	for i := 0; i < 5; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			if c.qInfuse != attributes.NoElement {
				aiAbsorb.Element = c.qInfuse
				c.Core.QueueAttackWithSnap(aiAbsorb, snapAbsorb, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			}
		}, "kazuha-burst-tick", 147+117*i)
	}

	//reset skill cd
	if c.Base.Cons > 0 {
		c.ResetActionCooldown(action.ActionSkill)
	}

	//add em to kazuha even if off-field
	//add em to all char, but only activate if char is active
	if c.Base.Cons >= 2 {
		// TODO: Lasts while Q field is on stage is ambiguous.
		// Does it apply to Kazuha's initial hit?
		// Not sure when it lasts from and until
		// For consistency with how it was previously done, assume that it lasts from button press to the last tick
		val := make([]float64, attributes.EndStatType)
		val[attributes.EM] = 200
		for _, char := range c.Core.Chars {
			this := char
			char.AddMod(core.CharStatMod{
				Key:    "kazuha-c2",
				Expiry: c.Core.F + 147 + 117*5,
				Amount: func() ([]float64, bool) {
					switch this.CharIndex() {
					case c.Core.ActiveChar, c.CharIndex():
						return val, true
					}
					return nil, false
				},
			})
		}
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    attributes.Anemo,
			Tags:   []combat.AttackTag{combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCDWithDelay(action.ActionBurst, 15*60, 7)
	c.ConsumeEnergy(7)
	return f, a
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
		c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), "kaz-q-absorb-check", 18)
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
		c.Core.Tasks.Add(c.absorbCheckA1(src, count+1, max), "kaz-a1-absorb-check", 6)
	}
}
