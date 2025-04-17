package mizuki

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1SwirlKey                    = "mizuki-a1-swirl-%v"
	a1ICDKey                      = "mizuki-a1-icd"
	a1ICD                         = 0.3 * 60
	a4Duration                    = 4 * 60
	a4Key                         = "mizuki-a4"
	a4EMBuff                      = 100
	dreamDrifterExtensions        = 2
	dreamDrifterDurationExtension = 2.5 * 60
)

// When Yumemizuki Mizuki triggers Swirl while in her Dreamdrifter state, Dreamdrifter's duration increases by 2.5s.
// This effect can trigger once every 0.3s for a maximum of 2 times per Dreamdrifter state.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	swirlFunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)

		// Mizuki should trigger the swirl
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

		// ICD
		if c.StatusIsActive(a1ICDKey) {
			return false
		}

		c.AddStatus(a1ICDKey, a1ICD, false)

		c.ExtendStatus(dreamDrifterStateKey, dreamDrifterDurationExtension)

		c.dreamDrifterExtensionsRemaining--

		return false
	}

	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlFunc, fmt.Sprintf(a1SwirlKey, attributes.Pyro))
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlFunc, fmt.Sprintf(a1SwirlKey, attributes.Hydro))
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlFunc, fmt.Sprintf(a1SwirlKey, attributes.Electro))
	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlFunc, fmt.Sprintf(a1SwirlKey, attributes.Cryo))
}

// While Yumemizuki Mizuki is in the Dreamdrifter state, when other nearby party members hit opponents with
// Pyro, Hydro, Cryo, or Electro attacks, her Elemental Mastery will increase by 100 for 4s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.EM] = a4EMBuff

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

		// is there any ICD?
		element := atk.Info.Element

		// Only when enemy is hit by Pyro, Hydro, Cryo, Electro
		if element == attributes.Electro || element == attributes.Hydro || element == attributes.Pyro || element == attributes.Cryo {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(a4Key, a4Duration), // 4s
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return c.a4Buff, true
				},
			})
		}

		return false
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, hitFunc, a4Key)
}
