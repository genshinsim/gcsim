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

	caFunc := func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.Player.Active() != char.Index() {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("dawning-frost-ca", 10*64),
			Amount: func() ([]float64, bool) {
				return emBuffCa, true
			},
		})
		return false
	}

	skillFunc := func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.Player.Active() != char.Index() {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt &&
			atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("dawning-frost-skill", 10*60),
			Amount: func() ([]float64, bool) {
				return emBuffSkill, true
			},
		})
		return false
	}

	c.Events.Subscribe(event.OnEnemyDamage, skillFunc,
		fmt.Sprintf("dawning-frost-%v-skill-hit", char.Base.Key.String()))
	c.Events.Subscribe(event.OnEnemyDamage, caFunc,
		fmt.Sprintf("dawning-frost-%v-ca-hit", char.Base.Key.String()))

	return w, nil
}
