package reliquary

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ReliquaryOfTruth, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.06 + 0.02*float64(r)

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("reliquary-cr", -1),
		Amount: func() []float64 {
			return m
		},
	})

	secretoflies := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 60 + 20*float64(r)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("reliquary-secretoflies", 10*60),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})

		bothActive(char, r)
	}

	moonoftruth := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.18 + 0.06*float64(r)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("reliquary-moonoftruth", 15*60),
			Extra:        true,
			AffectedStat: attributes.CR,
			Amount: func() []float64 {
				return m
			},
		})

		bothActive(char, r)
	}

	c.Events.Subscribe(event.OnSkill, secretoflies, fmt.Sprintf("reliquary-secretoflies-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnEnemyHit, moonoftruth, fmt.Sprintf("reliquary-moonoftruth-%v", char.Base.Key.String()))

	return w, nil
}

func bothActive(char *character.CharWrapper, r int) {
	if !char.StatusIsActive("reliquary-secretoflies") || !char.StatusIsActive("reliquary-moonoftruth") {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 30 + 10*float64(r)
	m[attributes.CR] = 0.09 + 0.03*float64(r)

	duration := min(char.StatusDuration("reliquary-secretoflies"), char.StatusDuration("reliquary-moonoftruth"))

	char.AddStatMod(character.StatMod{
		Base:  modifier.NewBaseWithHitlag("reliquary-both", duration),
		Extra: true,
		Amount: func() []float64 {
			return m
		},
	})
}
