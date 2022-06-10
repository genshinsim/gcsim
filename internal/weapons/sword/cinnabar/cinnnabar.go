package cinnabar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.CinnabarSpindle, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Elemental Skill DMG is increased by 40% of DEF. The effect will be triggered
// no more than once every 1.5s and will be cleared 0.1s after the Elemental
// Skill deals DMG.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	effectICDExpiry := 0
	effectDurationExpiry := 0
	effectLastProc := 0
	defPer := .3 + float64(r)*.1
	c.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		if effectDurationExpiry < c.F && c.F <= effectICDExpiry {
			return false
		}
		damageAdd := (atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[attributes.DEFP]) + atk.Snapshot.Stats[attributes.DEF]) * defPer
		atk.Info.FlatDmg += damageAdd

		c.Log.NewEvent("Cinnabar Spindle proc dmg add", glog.LogPreDamageMod, char.Index, "damage_added", damageAdd, "lastproc", effectLastProc, "effect_ends_at", effectDurationExpiry, "effect_icd_ends_at", effectICDExpiry)

		// TODO: Assumes that the ICD starts after the last duration ends
		effectICDExpiry = c.F + 6 + 90

		// Only want to update the last proc and duration if we're not within the currently active period
		if !(effectLastProc < c.F && c.F <= effectDurationExpiry) {
			effectLastProc = c.F
			effectDurationExpiry = c.F + 6
		}

		return false
	}, fmt.Sprintf("cinnabar-%v", char.Base.Name))

	return w, nil
}
