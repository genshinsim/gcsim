package heizou

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// For 5s after Shikanoin Heizou takes the field, his Normal Attack SPD is increased by 15%.
// He also gains 1 Declension stack for Heartstopper Strike. This effect can be triggered once every 10s.
func (c *char) c1() {
	const c1Icd = "heizou-c1-icd"
	// No log value saved as stat mod already shows up in debug view
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.StatusIsActive(c1Icd) {
			return false
		}
		next := args[1].(int)
		if next != c.Index {
			return false
		}
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("heizou-c1", 300), //5s
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				return c.c1buff, true
			},
		})
		c.addDecStack()
		c.AddStatus(c1Icd, 600, true)
		return false
	}, "heizou enter")

}

// The first Windmuster Iris explosion in each Windmuster Kick will regenerate 9 Elemental Energy for Shikanoin Heizou.
// Every subsequent explosion in that Windmuster Kick will each regenerate an additional 1.5 Energy for Heizou.
// One Windmuster Kick can regenerate a total of 13.5 Energy for Heizou in this manner.
func (c *char) c4(i int) {
	switch i {
	case 1:
		c.AddEnergy("heizou c4", 9.0)
	case 2, 3, 4:
		c.AddEnergy("heizou c4", 1.5)
	}
}

// Each Declension stack will increase the CRIT Rate of the Heartstopper Strike unleashed by 4%.
// When Heizou possesses Conviction, this Heartstopper Strike's CRIT DMG is increased by 32%.
func (c *char) c6() {
	val := make([]float64, attributes.EndStatType)

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("heizou-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}
			if c.decStack == 0 {
				return nil, false
			}
			val[attributes.CR] = 0.04 * float64(c.decStack)
			if c.decStack == 4 {
				val[attributes.CD] = 0.32
			} else {
				val[attributes.CD] = 0
			}
			return val, true
		},
	})
}
