package dehya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const (
	a1ReductionKey      = "dehya-a1-reduction"
	a1ReductionDuration = 6 * 60
	a1ReductionMult     = 0.6
	a4ICDKey            = "dehya-a4-icd"
	a4HealMsg           = "Stalwart and True (A4)"
	a4Threshold         = 0.4
	a4InitialHealRatio  = 0.2
	a4DoTHealRatio      = 0.06
	a4DoTHealInterval   = 2 * 60
	a4ICD               = 20 * 60
)

// Within 6 seconds after Dehya retrieves the Fiery Sanctum field through Molten Inferno: Ranging Flame
// or Leonine Bite, she will take 60% less DMG when receiving DMG from Redmane's Blood.
// This effect can be triggered once every 2s.
func (c *char) a1Reduction() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(a1ReductionKey, a1ReductionDuration, true)
}

// TODO: interrupt res part of a1 is not implemented
// Additionally, within 9s after Dehya unleashes Molten Inferno: Indomitable Flame,
// she will grant all party members the Gold-Forged Form state.
// This state will further increase a character's resistance to interruption when
// they are within the Fiery Sanctum field. Gold-Forged Form can be activated once every 18s.

// When her HP is less than 40%, Dehya will recover 20% of her Max HP
// and will restore 6% of her Max HP every 2s for the next 10s.
// This effect can be triggered once every 20s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	// TODO: should also check once every 1s but this is good enough...
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		if c.CurrentHPRatio() >= a4Threshold {
			return false
		}
		if c.StatusIsActive(a4ICDKey) {
			return false
		}
		c.AddStatus(a4ICDKey, a4ICD, true)
		// 20% HP Part
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: a4HealMsg,
			Src:     a4InitialHealRatio * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
		// 6% every 2s for 10s part (5 times)
		c.QueueCharTask(c.a4DotHeal(0), a4DoTHealInterval)
		return false
	}, "hutao-c6")
}

// don't need source because it is impossible for multiple a4s to be up at the same time
func (c *char) a4DotHeal(count int) func() {
	return func() {
		if count == 5 {
			return
		}
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: a4HealMsg,
			Src:     a4DoTHealRatio * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
		c.QueueCharTask(c.a4DotHeal(count+1), a4DoTHealInterval)
	}
}
