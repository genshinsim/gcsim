package kokomi

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Passive 2 - permanently modify stats for +25% healing bonus and -100% CR
func (c *char) passive() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = .25
	m[attributes.CR] = -1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("kokomi-passive", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// If Sangonomiya Kokomi's own Bake-Kurage is on the field when she uses Nereid's Ascension, the Bake-Kurage's duration will be refreshed.
//
// - checks for ascension level in burst.go to avoid queuing this up only to fail the ascension level check
func (c *char) a1() {
	if c.Core.Status.Duration("kokomiskill") <= 0 {
		return
	}
	// +1 to avoid same frame expiry issues with skill tick
	c.Core.Status.Add("kokomiskill", 12*60+1)
}

// While donning the Ceremonial Garment created by Nereid's Ascension, the Normal and Charged Attack DMG Bonus
// Sangonomiya Kokomi gains based on her Max HP will receive a further increase based on 15% of her Healing Bonus.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("kokomiburst") == 0 {
			return false
		}

		a4Bonus := c.Stat(attributes.Heal) * 0.15 * c.MaxHP()
		atk.Info.FlatDmg += a4Bonus

		return false
	}, "kokomi-a4")
}
