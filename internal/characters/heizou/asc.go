package heizou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Shikanoin Heizou activates a Swirl reaction while on the field,
// he will gain 1 Declension stack for Heartstopper Strike.
// This effect can be triggered once every 0.1s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	const a1IcdKey = "heizou-a1-icd"
	swirlCB := func() func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			if c.StatusIsActive(a1IcdKey) {
				return false
			}
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			if c.Core.Player.Active() != c.Index {
				return false
			}
			switch atk.Info.AttackTag {
			case combat.AttackTagSwirlPyro:
			case combat.AttackTagSwirlHydro:
			case combat.AttackTagSwirlElectro:
			case combat.AttackTagSwirlCryo:
			default:
				return false
			}
			//icd is triggered regardless if stacks are maxed or not
			c.AddStatus(a1IcdKey, 6, true)
			c.addDecStack()
			return false
		}
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, swirlCB(), "heizou-a1")
}

// After Shikanoin Heizou's Heartstopper Strike hits an opponent,
// increases all party members' (excluding Shikanoin Heizou) Elemental Mastery by 80 for 10s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	dur := 60 * 10
	for i, char := range c.Core.Player.Chars() {
		if i == c.Index {
			continue //nothing for heizou
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("heizou-a4", dur),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return c.a4Buff, true
			},
		})
	}
	c.Core.Log.NewEvent("heizou a4 triggered", glog.LogCharacterEvent, c.Index).Write("em snapshot", c.a4Buff[attributes.EM]).Write("expiry", c.Core.F+dur)
}
