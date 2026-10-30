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
	buffKey  = "heretics-molten-blade-atk"
	icdKey   = "heretics-molten-blade-icd"
	gleamKey = "heretics-molten-blade-gleam"
	maxDist  = 7
)

type Weapon struct {
	Index    int
	char     *character.CharWrapper
	src      int
	dist     float64
	buff     []float64
	minBonus float64
	maxBonus float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	r := p.Refine
	w := &Weapon{
		char:     char,
		buff:     make([]float64, attributes.EndStatType),
		minBonus: 0.135 + 0.045*float64(r),
		maxBonus: 0.27 + 0.09*float64(r),
	}

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() != char.Index() || char.StatusIsActive(icdKey) {
			return
		}

		char.AddStatus(icdKey, 14*60, true)
		char.AddStatus(gleamKey, 14*60, true)
		w.src = c.F
		w.dist = 0
		w.buff[attributes.ATKP] = w.minBonus
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 14*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				if !char.StatusIsActive(gleamKey) {
					return nil
				}
				return w.buff
			},
		})
		char.QueueCharTask(w.tick(w.src), 60)
	}, "heretics-molten-blade-on-skill-"+char.Base.Key.String())

	c.Events.Subscribe(event.OnMovement, func(args ...any) {
		if !char.StatusIsActive(gleamKey) || c.Player.Active() != char.Index() {
			return
		}
		if d := args[0].(float64); d > 0 {
			w.dist += d
		}
	}, "heretics-molten-blade-movement-"+char.Base.Key.String())

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		if args[0].(int) == char.Index() {
			char.DeleteStatMod(buffKey)
			char.DeleteStatus(gleamKey)
			w.src = -1
		}
	}, "heretics-molten-blade-swap-"+char.Base.Key.String())

	return w, nil
}

func (w *Weapon) tick(src int) func() {
	return func() {
		if w.src != src || !w.char.StatusIsActive(gleamKey) {
			return
		}

		// Min and max were verified in-game and values in between are continuous;
		// linear scaling up to 7m/s is an assumption.
		w.buff[attributes.ATKP] = w.minBonus + (w.maxBonus-w.minBonus)*min(w.dist/maxDist, 1)
		w.dist = 0
		w.char.QueueCharTask(w.tick(src), 60)
	}
}
