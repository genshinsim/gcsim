package kagura

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
	core.RegisterWeaponFunc(keys.KagurasVerity, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Gains the Kagura Dance effect when using an Elemental Skill, causing the
	// Elemental Skill DMG of the character wielding this weapon to increase by
	// 12% for 16s. Max 3 stacks. This character will gain 12% All Elemental DMG
	// Bonus when they possess 3 stacks.
	w := &Weapon{}
	r := p.Refine
	stacks := 0
	key := fmt.Sprintf("kaguradance-%v", char.Base.Key.String())
	dmg := 0.12 + 0.03*float64(r-1)
	val := make([]float64, attributes.EndStatType)
	const stackKey = "kaguras-verity-stacks"
	stackDuration := 960 // 16s * 60

	//TODO: this used to be on postskill. make sure nothing broke here
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if !char.StatusIsActive(stackKey) {
			// reset stacks back to 0
			stacks = 0
		}
		char.AddStatus(stackKey, stackDuration, true)
		if stacks < 3 {
			stacks++
		}
		// bonus ele damage if 3 stacks
		if stacks == 3 {
			val[attributes.PyroP] = dmg
			val[attributes.HydroP] = dmg
			val[attributes.CryoP] = dmg
			val[attributes.ElectroP] = dmg
			val[attributes.AnemoP] = dmg
			val[attributes.GeoP] = dmg
			val[attributes.PhyP] = dmg
			val[attributes.DendroP] = dmg
		} else {
			// clean stacks ele dmg% otherwise
			val[attributes.PyroP] = 0
			val[attributes.HydroP] = 0
			val[attributes.CryoP] = 0
			val[attributes.ElectroP] = 0
			val[attributes.AnemoP] = 0
			val[attributes.GeoP] = 0
			val[attributes.PhyP] = 0
			val[attributes.DendroP] = 0
		}
		// add mod for duration, override last
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("kaguras-verity", stackDuration),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.ActorIndex != char.Index {
					return nil, false
				}
				if atk.Info.AttackTag == attacks.AttackTagElementalArt || atk.Info.AttackTag == attacks.AttackTagElementalArtHold {
					val[attributes.DmgP] = dmg * float64(stacks)
				} else {
					val[attributes.DmgP] = 0
				}
				return val, true
			},
		})
		return false
	}, key)

	return w, nil
}
