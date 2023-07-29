package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type GoldenMajesty struct {
	Index int
}

func (b *GoldenMajesty) SetIndex(idx int) { b.Index = idx }
func (b *GoldenMajesty) Init() error      { return nil }

func NewGoldenMajesty(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &GoldenMajesty{}
	r := p.Refine

	const buffKey = "golden-majesty"
	const icdKey = "golden-majesty-icd"

	shd := .15 + float64(r)*.05
	atkbuff := 0.03 + 0.01*float64(r)
	key := fmt.Sprintf("golden-majesty-%v", char.Base.Key.String())
	c.Player.Shields.AddShieldBonusMod(key, -1, func() (float64, bool) { return shd, false })

	stacks := 0
	m := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 18, true)

		if !char.StatModIsActive(buffKey) {
			stacks = 0
		}

		stacks++
		if stacks > 5 {
			stacks = 5
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 60*8),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				m[attributes.ATKP] = atkbuff * float64(stacks)
				if char.Index == c.Player.Active() && c.Player.Shields.PlayerIsShielded() {
					m[attributes.ATKP] *= 2
				}
				return m, true
			},
		})
		return false
	}, key)

	return w, nil
}
