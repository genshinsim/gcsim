package songofthevigil

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const energyICDKey = "song-of-the-vigil-energy-icd"

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	onReaction := func(args ...any) {
		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != char.Index() {
			return
		}

		if char.StatusIsActive(energyICDKey) {
			return
		}
		char.AddStatus(energyICDKey, 9*60, true)
		char.AddEnergy("song-of-the-vigil-energy", 4)
	}

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.15 + float64(r)*0.05

	onStellar := func(args ...any) {
		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("song-of-the-vigil-atk", 12*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return atkBuff
			},
		})
	}

	for evt := event.ReactionEventStartDelim + 1; evt < event.ReactionEventEndDelim; evt++ {
		c.Events.Subscribe(evt, onReaction, "song-of-the-vigil-on-reaction")
	}

	c.Events.Subscribe(event.OnStellarConduct, onStellar, "song-of-the-vigil-on-stellar")
	c.Events.Subscribe(event.OnStellarSwirl, onStellar, "song-of-the-vigil-on-stellar")

	return w, nil
}
