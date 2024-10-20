package mountainbracingbolt

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.MountainBracingBolt, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	baseBuffKey  = "mountain-bracing-bolt-base"
	otherBuffKey = "mountain-bracing-bolt-other"
)

// Decreases Climbing Stamina Consumption by 15% and increases Elemental Skill DMG by 12%.
// Also, after other nearby party members use Elemental Skills,
// the equipping character's Elemental Skill DMG will also increase by 12% for 8s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.09 + float64(r)*0.03

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(baseBuffKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil, false
			}
			return m, true
		},
	})

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() == char.Index {
			return false
		}
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(otherBuffKey, 8*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
					return nil, false
				}
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("mountain-bracing-bolt-%v", char.Base.Key.String()))

	return w, nil
}
