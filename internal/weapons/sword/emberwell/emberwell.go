package emberwell

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.12 + float64(r)*0.04

	onReaction := func(args ...any) {
		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("emberwell-em", 12*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return atkBuff
			},
		})
	}

	stellarBuff := 0.12 + float64(r)*0.04

	onStellar := func(args ...any) {
		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != char.Index() {
			return
		}

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag("emberwell-stellar", 12*60),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag.IsStellar() {
					return stellarBuff
				}
				return 0
			},
		})
	}

	for evt := event.ReactionEventStartDelim + 1; evt < event.ReactionEventEndDelim; evt++ {
		c.Events.Subscribe(evt, onReaction, "emberwell-on-reaction")
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, "emberwell-stellar")
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, "emberwell-stellar")

	return w, nil
}
