package flameforgedinsight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	icdKey  = "flame-forged-insight-icd"
	buffKey = "flame-forged-insight-buff"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	energy := 9 + float64(r)*3
	em := 45 + float64(r)*15

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = em

	key := fmt.Sprintf("flameforgedinsight-%v", char.Base.Key.String())

	reactionProc := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if char.StatusIsActive(icdKey) {
			return
		}

		char.AddStatus(icdKey, 15*60, true)
		char.AddEnergy("flameforgedinsight", energy)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 15*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}

	c.Events.Subscribe(event.OnElectroCharged, reactionProc, key)
	c.Events.Subscribe(event.OnLunarCharged, reactionProc, key)
	c.Events.Subscribe(event.OnBloom, reactionProc, key)
	c.Events.Subscribe(event.OnLunarBloom, reactionProc, key)

	c.Events.Subscribe(event.OnCrystallizeHydro, reactionProc, key)
	c.Events.Subscribe(event.OnCrystallizePyro, reactionProc, key)
	c.Events.Subscribe(event.OnCrystallizeElectro, reactionProc, key)
	c.Events.Subscribe(event.OnCrystallizeCryo, reactionProc, key)
	c.Events.Subscribe(event.OnLunarCrystallize, reactionProc, key)

	return w, nil
}
