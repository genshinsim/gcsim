package urakumisugiri

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
	core.RegisterWeaponFunc(keys.UrakuMisugiri, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey = "urakumisugiri-increase-buff"
)

// Normal Attack DMG is increased by 16/20/24/28/32% and Elemental Skill DMG is increased by 24/30/36/42/48%.
// After a nearby active character deals Geo DMG, the aforementioned effects increase by 100% for 15s.
// Additionally, the wielder's DEF is increased by 20/25/30/35/40%.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	normalDmg := 0.12 + 0.04*float64(r)
	skillDmg := 0.18 + 0.06*float64(r)
	defIncrease := 0.15 + 0.05*float64(r)

	mNormal := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("urakumisugiri-na", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal {
				return nil, false
			}

			mNormal[attributes.DmgP] = normalDmg
			if char.StatusIsActive(buffKey) {
				mNormal[attributes.DmgP] *= 2
			}
			return mNormal, true
		},
	})

	mSkill := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("urakumisugiri-skill", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil, false
			}

			mSkill[attributes.DmgP] = skillDmg
			if char.StatusIsActive(buffKey) {
				mSkill[attributes.DmgP] *= 2
			}
			return mSkill, true
		},
	})

	mDef := make([]float64, attributes.EndStatType)
	mDef[attributes.DEFP] = defIncrease
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("urakumisugiri-def", -1),
		AffectedStat: attributes.DEFP,
		Amount: func() ([]float64, bool) {
			return mDef, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Element != attributes.Geo {
			return false
		}
		char.AddStatus(buffKey, 900, true)
		return false
	}, fmt.Sprintf("urakumisugiri-%v", char.Base.Key.String()))

	return w, nil
}
