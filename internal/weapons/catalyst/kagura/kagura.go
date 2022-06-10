package kagura

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.KagurasVerity, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	stacks := 0
	var ctick = func(char *character.CharWrapper, c *core.Core) func() {
		return func() {
			if c.Status.Duration("kaguradance-"+char.Base.Name) <= 0 {
				stacks = 0
				return
			}
		}
	}

	c.Events.Subscribe(event.PostSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		c.Status.Add("kaguradance-"+char.Base.Name, 16*60)
		if stacks < 3 {
			stacks++
		}
		c.Tasks.Add(ctick(char, c), 16*60)
		return false

	}, fmt.Sprintf("kaguradance-%v", char.Base.Name))

	mod := float64(r - 1)
	char.AddAttackMod("kagurasverity",

		-1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.ActorIndex != char.Index {
				return nil, false
			}
			val := make([]float64, attributes.EndStatType)
			if stacks == 3 {
				val[attributes.PyroP] = 0.12 + 0.03*mod
				val[attributes.HydroP] = 0.12 + 0.03*mod
				val[attributes.CryoP] = 0.12 + 0.03*mod
				val[attributes.ElectroP] = 0.12 + 0.03*mod
				val[attributes.AnemoP] = 0.12 + 0.03*mod
				val[attributes.GeoP] = 0.12 + 0.03*mod
				val[attributes.PhyP] = 0.12 + 0.03*mod
				val[attributes.DendroP] = 0.12 + 0.03*mod
			}
			if atk.Info.AttackTag == combat.AttackTagElementalArt || atk.Info.AttackTag == combat.AttackTagElementalArtHold {
				val[attributes.DmgP] = float64(stacks) * (0.12 + mod*0.03)
			}

			return val, true
		})

	return w, nil

}
