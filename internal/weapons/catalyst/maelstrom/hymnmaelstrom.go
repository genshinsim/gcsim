package maelstrom

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	reactionBuffKey      = "maelstrom-reaction-boost"
	reactionBuffDuration = 5 * 60
)

type Weapon struct {
	Index int
}

func maelstromAtkBoost(maxHP float64, refine, count int) float64 {
	multiplier := 0.75 + float64(refine)*0.25
	intervals := max(0, min(int((maxHP-40000)/1000), 20))
	boost := 0.004 * float64(intervals) * multiplier * float64(count)
	return min(boost, 0.08*multiplier)
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Increases Healing Bonus by 4%.
//
//	When performing healing, the equipping character gains the "Vatsamonga's Vatic Vintage"
//
// effect, which increases Max HP by 4% as well as increases the currently active party
// member's ATK by 0.4% for every 1,000 Max HP the equipping character has over 40,000.
// A maximum of 8% ATK can be gained in this way. This effect lasts 10s, max 3 stacks.
// When a nearby party member triggers a Frozen or Stellar Swirl reaction, the aforementioned
// Max HP and ATK boosts will be further increased by 75% for the next 5s.
// The aforementioned effects can still trigger even when the equipping character is
// not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	multiplier := 0.75 + float64(p.Refine)*0.25
	index := 0
	hpboost := 0.04 * multiplier

	atkcap := 0.08 * multiplier
	stackKey := []string{
		"maelstrom-stack-1",
		"maelstrom-stack-2",
		"maelstrom-stack-3",
	}
	atkBuffKey := fmt.Sprintf("maelstrom-ATK-%v", char.Base.Key.String())
	m := make([]float64, attributes.EndStatType)

	updateAtkBuff := func(target *character.CharWrapper) {
		duration := 0
		count := 0
		for _, key := range stackKey {
			remaining := char.StatusDuration(key)
			if remaining > 0 {
				count++
				duration = max(duration, remaining)
			}
		}
		if count == 0 {
			target.DeleteStatMod(atkBuffKey)
			return
		}

		target.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(atkBuffKey, duration),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				count := 0
				for _, key := range stackKey {
					if char.StatusIsActive(key) {
						count++
					}
				}
				boost := 1.0
				if char.StatusIsActive(reactionBuffKey) {
					boost = 1.75
				}
				m[attributes.ATKP] = min(maelstromAtkBoost(char.MaxHP(), p.Refine, count)*boost, atkcap*boost)
				return m
			},
		})
	}

	mheal := make([]float64, attributes.EndStatType)
	mheal[attributes.Heal] = 0.04 * multiplier
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("maelstrom-heal-bonus", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() []float64 {
			return mheal
		},
	})

	onBoostReaction := func(args ...any) {
		char.AddStatus(reactionBuffKey, reactionBuffDuration, true)
	}
	reactionSubKey := fmt.Sprintf("maelstrom-reaction-%v", char.Base.Key.String())
	c.Events.Subscribe(event.OnFrozen, onBoostReaction, reactionSubKey)
	c.Events.Subscribe(event.OnStellarSwirl, onBoostReaction, reactionSubKey)

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		c.Player.ByIndex(prev).DeleteStatMod(atkBuffKey)
		updateAtkBuff(c.Player.ByIndex(next))
	}, fmt.Sprintf("maelstrom-swap-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...any) {
		source := args[0].(*info.HealInfo)
		amount := args[2].(float64)
		// check if the wearer realises the heal
		if source.Caller != char.Index() {
			return
		}
		if amount <= 0 {
			return
		}
		// add status to char given index
		char.AddStatus(stackKey[index], 600, true)
		// update buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("maelstrom-HP", 600),
			AffectedStat: attributes.HPP,
			Amount: func() []float64 {
				count := 0
				for _, v := range stackKey {
					if char.StatusIsActive(v) {
						count++
					}
				}
				boost := 1.0
				if char.StatusIsActive(reactionBuffKey) {
					boost = 1.75
				}
				m[attributes.HPP] = hpboost * float64(count) * boost
				return m
			},
		})
		updateAtkBuff(c.Player.ActiveChar())
		index++
		if index == 3 {
			index = 0
		}
	}, fmt.Sprintf("maelstrom-%v", char.Base.Key.String()))

	return w, nil
}
