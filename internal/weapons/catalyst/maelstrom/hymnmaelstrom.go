package maelstrom

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	reactionBuffKey      = "maelstrom-reaction-boost"
	reactionBuffDuration = 5 * 60
	stackDuration        = 10 * 60
)

type Weapon struct {
	Index        int
	char         *character.CharWrapper
	stackTracker *stacks.MultipleRefreshNoRemove
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
	w := &Weapon{
		char:         char,
		stackTracker: stacks.NewMultipleRefreshNoRemove(3, char.QueueCharTask, &c.F),
	}
	multiplier := 0.75 + float64(p.Refine)*0.25
	hpboost := 0.04 * multiplier
	atkcap := 0.08 * multiplier

	atkBuffKey := fmt.Sprintf("maelstrom-ATK-%v", char.Base.Key.String())
	m := make([]float64, attributes.EndStatType)

	updateAtkBuff := func(target *character.CharWrapper) {
		if w.stackTracker.Count() == 0 {
			target.DeleteStatMod(atkBuffKey)
			return
		}

		target.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(atkBuffKey, stackDuration),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				boost := 1.0
				if char.StatusIsActive(reactionBuffKey) {
					boost = 1.75
				}
				m[attributes.ATKP] = min(maelstromAtkBoost(char.MaxHP(), p.Refine, w.stackTracker.Count())*boost, atkcap*boost)
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
		w.stackTracker.Add(stackDuration)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("maelstrom-HP", stackDuration),
			AffectedStat: attributes.HPP,
			Amount: func() []float64 {
				boost := 1.0
				if char.StatusIsActive(reactionBuffKey) {
					boost = 1.75
				}
				m[attributes.HPP] = hpboost * float64(w.stackTracker.Count()) * boost
				return m
			},
		})
		updateAtkBuff(c.Player.ActiveChar())
	}, fmt.Sprintf("maelstrom-%v", char.Base.Key.String()))

	return w, nil
}
