package bladeofatonement

import (
	"fmt"

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

	emBuff := make([]float64, attributes.EndStatType)
	emBuff[attributes.EM] = 48 + float64(r)*16

	onReaction := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("blade-of-atonement-em", 12*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return emBuff
			},
		})
	}

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.12 + float64(r)*0.04

	onStellar := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("blade-of-atonement-atk", 12*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return atkBuff
			},
		})
	}

	for evt := event.ReactionEventStartDelim + 1; evt < event.ReactionEventEndDelim; evt++ {
		c.Events.Subscribe(evt, onReaction, "blade-of-atonement-on-reaction")
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, fmt.Sprintf("blade-of-atonement-on-stellar-conduct-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, fmt.Sprintf("blade-of-atonement-on-stellar-swirl-%v", char.Base.Key.String()))

	return w, nil
}
