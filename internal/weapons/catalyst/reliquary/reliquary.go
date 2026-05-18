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
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ReliquaryOfTruth, NewWeapon)
}

const (
	secretOfLiesKey = "reliquary-secretoflies"
	moonOfTruthKey  = "reliquary-moonoftruth"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// CRIT Rate is increased by 8%/10%/12%/14%/16%. When the equipping character
// unleashes an Elemental Skill, they gain the Secret of Lies effect: Elemental
// Mastery is increased by 80/100/120/140/160 for 12s. When the equipping character
// deals Lunar-Bloom DMG to an opponent, they gain the Moon of Truth effect:
// CRIT DMG is increased by 24%/30%/36%/42%/48% for 4s. When both the Secret of
// Lies and Moon of Truth effects are active at the same time, the results of
// both effects will be increased by 50%.
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

	emBuff := 60 + 20*float64(r)
	m2 := make([]float64, attributes.EndStatType)
	secretoflies := func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(secretOfLiesKey, 12*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				m2[attributes.EM] = emBuff
				if char.StatModIsActive(moonOfTruthKey) {
					m2[attributes.EM] = emBuff * 1.5
				}
				return m2
			},
		})
	}

	cdBuff := 0.18 + 0.06*float64(r)
	m3 := make([]float64, attributes.EndStatType)
	moonoftruth := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(moonOfTruthKey, 4*60),
			AffectedStat: attributes.CD,
			Amount: func() []float64 {
				m3[attributes.CD] = cdBuff
				if char.StatModIsActive(secretOfLiesKey) {
					m3[attributes.CD] = cdBuff * 1.5
				}
				return m3
			},
		})
	}

	c.Events.Subscribe(event.OnSkill, secretoflies, fmt.Sprintf("reliquary-secretoflies-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnEnemyDamage, moonoftruth, fmt.Sprintf("reliquary-moonoftruth-%v", char.Base.Key.String()))

	return w, nil
}
