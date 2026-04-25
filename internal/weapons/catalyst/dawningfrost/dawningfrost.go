package dawningfrost

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
	core.RegisterWeaponFunc(keys.DawningFrost, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// For 10s after a Charged Attack hits an opponent, Elemental Mastery is increased by 72.
// For 10s after an Elemental Skill hits an opponent, Elemental Mastery is increased by 48.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	emBuffCa := make([]float64, attributes.EndStatType)
	emBuffCa[attributes.EM] = 54 + float64(r)*18

	emBuffSkill := make([]float64, attributes.EndStatType)
	emBuffSkill[attributes.EM] = 36 + float64(r)*12

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if c.Player.Active() != char.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagExtra:
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("dawning-frost-ca", 10*60),
				Amount: func() []float64 {
					return emBuffCa
				},
			})
		case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("dawning-frost-skill", 10*60),
				Amount: func() []float64 {
					return emBuffSkill
				},
			})
		}
	}, fmt.Sprintf("dawning-frost-%v-ondamage", char.Base.Key.String()))

	return w, nil
}
