package echoesoftheheart

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

	emBuff := make([]float64, attributes.EndStatType)
	emBuff[attributes.EM] = 45 + float64(r)*15

	onReaction := func(args ...any) {
		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("echoes-of-the-heart-em", 12*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return emBuff
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
			Base: modifier.NewBaseWithHitlag("echoes-of-the-heart-stellar", 12*60),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag.IsStellar() {
					return stellarBuff
				}
				return 0
			},
		})
	}

	for evt := event.ReactionEventStartDelim + 1; evt < event.ReactionEventEndDelim; evt++ {
		c.Events.Subscribe(evt, onReaction, "echoes-of-the-heart-on-reaction")
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, "echoes-of-the-heart-on-stellar")
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, "echoes-of-the-heart-on-stellar")

	return w, nil
}
