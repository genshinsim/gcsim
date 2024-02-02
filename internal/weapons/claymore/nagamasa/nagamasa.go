package nagamasa

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
	core.RegisterWeaponFunc(keys.KatsuragikiriNagamasa, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Increases Elemental Skill DMG by 6%. After Elemental Skill hits an
	// opponent, the character loses 3 Energy but regenerates 3 Energy every 2s
	// for the next 6s. This effect can occur once every 10s. Can be triggered
	// even when the character is not on the field.0
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	base := 0.045 + float64(r)*0.015
	regen := 2.5 + float64(r)*0.5

	m[attributes.DmgP] = base
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("nagamasa-skill-dmg-buff", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagElementalArt || atk.Info.AttackTag == attacks.AttackTagElementalArtHold {
				return m, true
			}
			return nil, false
		},
	})

	const icdKey = "nagamasa-icd"
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 600, true)
		char.AddEnergy("nagamasa", -3)
		for i := 120; i <= 360; i += 120 {
			// use char queue for hitlag
			char.QueueCharTask(func() {
				char.AddEnergy("nagamasa", regen)
			}, i)
		}
		return false
	}, fmt.Sprintf("nagamasa-%v", char.Base.Key.String()))

	return w, nil
}
