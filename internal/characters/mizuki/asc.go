package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1ICDKey                      = "mizuki-a1-icd"
	a1ICD                         = 0.3 * 60
	a4Duration                    = 4 * 60
	a4EMBuff                      = 100
	a4BuffKey                     = "mizuki-a4-buff"
	dreamDrifterExtensions        = 2
	dreamDrifterDurationExtension = 2.5 * 60
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	swirlfunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		// Only when dream drifter is active
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return false
		}

		// Max 2 extensions per E
		if c.dreamDrifterExtensionsRemaining <= 0 {
			return false
		}

		// Check ICD
		if c.StatusIsActive(a1ICDKey) {
			return false
		}

		// set the ICD
		c.AddStatus(a1ICDKey, a1ICD, false)

		c.ExtendStatus(dreamDrifterStateKey, dreamDrifterDurationExtension)

		for _, char := range c.Core.Player.Chars() {
			char.ExtendStatus(dreamDrifterSwirlBuffKey, dreamDrifterDurationExtension)
		}

		c.dreamDrifterExtensionsRemaining--

		return false
	}

	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc, "mizuki-a1-pyro-swirl")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc, "mizuki-a1-hydro-swirl")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc, "mizuki-a1-electro-swirl")
	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc, "mizuki-a1-cryo-swirl")
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	a4Buff := make([]float64, attributes.EndStatType)
	a4Buff[attributes.EM] = a4EMBuff

	hitFunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		// Only when others attack
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex == c.Index {
			return false
		}

		// Only when dream drifter is active
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return false
		}

		element := atk.Info.Element

		// Only when enemy is hit by Pyro, Hydro, Cryo, Electro
		if element == attributes.Electro || element == attributes.Hydro || element == attributes.Pyro || element == attributes.Cryo {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(a4BuffKey, a4Duration), // 4s
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return a4Buff, true
				},
			})
		}

		return false
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, hitFunc, "mizuki-a4-hit")
}
