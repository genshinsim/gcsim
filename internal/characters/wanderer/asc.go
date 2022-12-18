package wanderer

import (
	"fmt"
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

func (c *char) a4Init() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		// if e is not active or a4 already active
		if !c.StatusIsActive(skillKey) || c.StatusIsActive(a4Key) {
			return false
		}

		atk := args[1].(*combat.AttackEvent)

		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}

		// If this is not a normal or charged attack then ignore
		if !(atk.Info.AttackTag == combat.AttackTagNormal || atk.Info.AttackTag == combat.AttackTagExtra) {
			return false
		}

		if c.Core.Rand.Float64() > c.a4Prob {
			c.a4Prob += 0.12
			return false
		}

		c.AddStatus(a4Key, -1, true)

		c.Core.Log.NewEvent("wanderer-a4 proc'd", glog.LogCharacterEvent, c.Index).
			Write("probability", c.a4Prob)

		c.a4Prob = 0.16

		return false
	}, fmt.Sprintf("wanderer-a4"))
}

func (c *char) a4Activation() {
	// TODO
}

func (c *char) absorbCheckA1(src int) func() {
	return func() {

		absorbCheck := c.Core.Combat.AbsorbCheck(c.a1AbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if absorbCheck != attributes.NoElement && c.checkIfA1BuffExists(absorbCheck) {
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"wanderer a1 absorbed ", absorbCheck.String(),
			)
			c.a1Buffs = append(c.a1Buffs, absorbCheck)

			c.addA1Buff(absorbCheck)

			maxAbsorb := 2

			if c.Base.Cons >= 4 {
				maxAbsorb = 3

				if len(c.a1Buffs) < maxAbsorb {
					validElements := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}
					var possibleElements []attributes.Element
					for _, e := range validElements {
						if c.checkIfA1BuffExists(e) {
							possibleElements = append(possibleElements, e)
						}
					}
					chosenElement := possibleElements[c.Core.Rand.Intn(len(possibleElements))]
					c.addA1Buff(chosenElement)
					c.Core.Log.NewEvent("wanderer c4 applied", glog.LogCharacterEvent, c.Index)
				}

			}

			if len(c.a1Buffs) >= maxAbsorb {
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
			if atk.Info.ActorIndex != c.Index {
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

func (c *char) checkIfA1BuffExists(ele attributes.Element) bool {
	for _, e := range c.a1Buffs {
		if e == ele {
			return true
		}
	}
	return false
}
