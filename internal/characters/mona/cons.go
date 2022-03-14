package mona

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

//When a Normal Attack hits, there is a 20% chance that it will be automatically followed by a Charged Attack.
//This effect can only occur once every 5s.
func (c *char) c2cb(a core.AttackCB) {
	if c.Core.Rand.Float64() > .2 {
		return
	}
	if c.c2icd > c.Core.Frame {
		return
	}
	c.c2icd = c.Core.Frame + 300 //every 5 seconds
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  coretype.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, coretype.TargettableEnemy), 0, 0)
}

//When any party member attacks an opponent affected by an Omen, their CRIT Rate is increased by 15%.
func (c *char) c4() {
	val := make([]float64, core.EndStatType)
	val[core.CR] = 0.15
	for _, char := range c.Core.Chars {
		char.AddPreDamageMod(coretype.PreDamageMod{
			Key:    "mona-c4",
			Expiry: -1,
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				//ignore if omen or bubble not present
				if t.GetTag(bubbleKey) < c.Core.Frame && t.GetTag(omenKey) < c.Core.Frame {
					return nil, false
				}
				return val, true
			},
		})
	}
}
