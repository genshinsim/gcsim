package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4Key           = "wanderer-a4"
	a4IcdKey        = "wanderer-a4-icd"
	a1ElectroKey    = "wanderer-a1-electro"
	a1ElectroIcdKey = "wanderer-a1-electro-icd"
	a1PyroKey       = "wanderer-a1-pyro"
	a1CryoKey       = "wanderer-a1-cryo"
)

func (c *char) a4CB(a combat.AttackCB) {
	if !c.StatusIsActive(skillKey) || c.StatusIsActive(a4Key) || c.StatusIsActive(a4IcdKey) {
		return
	}

	c.AddStatus(a4IcdKey, 6, true)

	if c.Core.Rand.Float64() > c.a4Prob {
		c.a4Prob += 0.12
		return
	}

	c.Core.Log.NewEvent("wanderer-a4 proc'd", glog.LogCharacterEvent, c.Index).
		Write("probability", c.a4Prob)

	c.a4Prob = 0.16

	c.AddStatus(a4Key, 20*60, true)

	if c.Core.Player.CurrentState() == action.DashState {
		c.a4()
		return
	}
}

func (c *char) a4() bool {
	if c.StatusIsActive(a4Key) {
		c.DeleteStatus(a4Key)

		a4Mult := 0.35

		if c.StatusIsActive("wanderer-c1-atkspd") {
			a4Mult = 0.6
		}

		a4Info := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Gales of Reverie",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagWandererA4,
			ICDGroup:   combat.ICDGroupWandererA4,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       a4Mult,
		}

		for i := 0; i < 4; i++ {
			c.Core.QueueAttack(a4Info, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1),
				a4Release[i], a4Release[i]+a4Hitmark)
		}

		return true

	}

	return false
}

func (c *char) absorbCheckA1() {
	if len(c.a1ValidBuffs) <= c.a1MaxAbsorb {
		return
	}

	a1AbsorbCheckLocation := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	absorbCheck := c.Core.Combat.AbsorbCheck(a1AbsorbCheckLocation, c.a1ValidBuffs...)

	if absorbCheck != attributes.NoElement {
		c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
			"wanderer a1 absorbed ", absorbCheck.String(),
		)
		c.deleteFromValidBuffs(absorbCheck)
		c.addA1Buff(absorbCheck)
		if c.Base.Cons >= 4 && len(c.a1ValidBuffs) == 3 {
			chosenElement := c.a1ValidBuffs[c.Core.Rand.Intn(len(c.a1ValidBuffs))]
			c.addA1Buff(chosenElement)
			c.deleteFromValidBuffs(chosenElement)
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"wanderer c4 applied a1 ", chosenElement.String(),
			)
		}
	}

	//otherwise queue up
	delay := 6
	if c.skydwellerPoints*6 > delay {
		c.Core.Tasks.Add(c.absorbCheckA1, delay)
	}

}

// Buffs, need to be manually removed when state is ending
func (c *char) addA1Buff(absorbCheck attributes.Element) {
	switch absorbCheck {

	case attributes.Hydro:
		c.maxSkydwellerPoints += 20
		c.skydwellerPoints += 20

	case attributes.Pyro:
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.3
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1PyroKey, 1200),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

	case attributes.Cryo:
		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.2
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1CryoKey, 1200),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

	case attributes.Electro:
		c.AddStatus(a1ElectroKey, 1200, true)
	}
}

func (c *char) a1ElectroCB(cb combat.AttackCB) {
	if !c.StatusIsActive(a1ElectroKey) {
		return
	}
	if c.StatusIsActive(a1ElectroIcdKey) {
		return
	}
	c.AddStatus(a1ElectroIcdKey, 12, true)
	c.AddEnergy("wanderer-a1-electro-energy", 0.8)
}

func (c *char) deleteFromValidBuffs(ele attributes.Element) {
	var results []attributes.Element
	for _, e := range c.a1ValidBuffs {
		if e != ele {
			results = append(results, e)
		}
	}
	c.a1ValidBuffs = results
}
