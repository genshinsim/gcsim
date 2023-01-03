package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4Key = "wanderer-a4"
)

func (c *char) makeA4Callback() func(cb combat.AttackCB) {
	return func(a combat.AttackCB) {
		if !c.StatusIsActive(skillKey) || c.StatusIsActive(a4Key) {
			return
		}

		if c.Core.Rand.Float64() > c.a4Prob {
			c.a4Prob += 0.12
			return
		}

		c.Core.Log.NewEvent("wanderer-a4 proc'd", glog.LogCharacterEvent, c.Index).
			Write("probability", c.a4Prob)

		c.a4Prob = 0.16

		if c.Core.Player.CurrentState() == action.DashState {
			c.a4()
			return
		}

		c.AddStatus(a4Key, 20*60, true)
	}
}

func (c *char) a4() {
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
			c.Core.QueueAttack(a4Info, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 1),
				a4Release[i], a4Release[i]+a4Hitmark)
		}
	}
}

func (c *char) absorbCheckA1(src int) func() {
	return func() {

		if len(c.a1ValidBuffs) <= c.a1MaxAbsorb {
			return
		}

		a1AbsorbCheckLocation := combat.NewCircleHit(c.Core.Combat.Player(), 5)

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

			if 4-len(c.a1ValidBuffs) >= c.a1MaxAbsorb {
				return
			}
		}
		//otherwise queue up
		// TODO: determine delay

		delay := 6

		if c.skydwellerPoints*6 > delay {
			c.Core.Tasks.Add(c.absorbCheckA1(src), delay)
		}

	}
}

// Buffs, need to be manually removed when state is ending
func (c *char) addA1Buff(absorbCheck attributes.Element) {
	switch absorbCheck {

	case attributes.Hydro:
		c.maxSkydwellerPoints += 20
		c.skydwellerPoints += 20

	case attributes.Pyro:
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("wanderer-a1-pyro", 1200),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				m := make([]float64, attributes.EndStatType)
				m[attributes.ATKP] = 0.3
				return m, true
			},
		})

	case attributes.Cryo:
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("wanderer-a1-cryo", 1200),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				m := make([]float64, attributes.EndStatType)
				m[attributes.CR] = 0.2
				return m, true
			},
		})

	case attributes.Electro:
		c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
			if c.StatusIsActive("wanderer-a1-electro-icd") {
				return false
			}
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index || atk.Info.AttackTag != combat.AttackTagNormal {
				return false
			}
			if c.Core.Player.Active() != c.Index {
				return false
			}

			c.AddStatus("wanderer-a1-electro-icd", 12, true)
			c.AddEnergy("wanderer-a1-electro-energy", 0.8)
			return false
		}, "wanderer-a1-electro")
	}
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
