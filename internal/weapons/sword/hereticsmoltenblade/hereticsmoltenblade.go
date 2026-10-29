package hereticsmoltenblade

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffKey = "heretics-molten-blade-atk"
	icdKey  = "heretics-molten-blade-icd"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dist, ok := p.Params["dist"]
	if !ok {
		dist = 7
	}
	if dist < 0 {
		dist = 0
	}
	if dist > 7 {
		dist = 7
	}

	minBonus := 0.135 + 0.045*float64(r)
	maxBonus := 0.27 + 0.09*float64(r)
	// Min and max were verified in-game and values in between are continuous;
	// linear scaling up to 7m/s is an assumption.
	bonus := minBonus + (maxBonus-minBonus)*float64(dist)/7
	buff := make([]float64, attributes.EndStatType)
	buff[attributes.ATKP] = bonus

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() != char.Index() || char.StatusIsActive(icdKey) {
			return
		}

		char.AddStatus(icdKey, 14*60, true)
		// Distance is constant, so apply the bonus immediately for the full 14s.
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 14*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return buff
			},
		})
	}, "heretics-molten-blade-on-skill-"+char.Base.Key.String())

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		if args[0].(int) == char.Index() {
			char.DeleteStatMod(buffKey)
		}
	}, "heretics-molten-blade-swap-"+char.Base.Key.String())

	return w, nil
}
