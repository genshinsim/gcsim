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

const (
	icdKey      = "cinnabar-icd"
	durationKey = "cinnabar-buff-active"
)

// Elemental Skill DMG is increased by 40% of DEF. The effect will be triggered
// no more than once every 1.5s and will be cleared 0.1s after the Elemental
// Skill deals DMG.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	defPer := .3 + float64(r)*.1
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		//don't do anything if we're in icd period
		if char.StatusIsActive(icdKey) {
			return false
		}
		//otherwise if this is first time proc'ing, set the duration and queue
		//task to set icd
		if !char.StatusIsActive(durationKey) {
			//TODO: we're assuming icd starts after the effect
			char.QueueCharTask(func() {
				char.AddStatus(icdKey, 90, false) //icd lasts for 1.5s
			}, 6) //icd starts 6 frames after
			char.AddStatus(durationKey, 6, false)
		}
		damageAdd := (char.Base.Def*(1+char.Stat(attributes.DEFP)) + char.Stat(attributes.DEF)) * defPer
		atk.Info.FlatDmg += damageAdd

		c.Log.NewEvent("Cinnabar Spindle proc dmg add", glog.LogPreDamageMod, char.Index).
			Write("damage_added", damageAdd)
		return false
	}, fmt.Sprintf("cinnabar-%v", char.Base.Key.String()))

	return w, nil
}
