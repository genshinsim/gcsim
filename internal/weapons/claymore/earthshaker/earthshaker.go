package earthshaker

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
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.EarthShaker, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After a party member triggers a Pyro-related reaction,
// the equipping character's Elemental Skill DMG is increased by 32% for 8s.
// This effect can be triggered even when the triggering party member is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	amt := 0.12 + float64(r)*0.04
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = amt

	buffSkill := func() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("earth-shaker", 8*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
					return nil, false
				}
				return m, true
			},
		})
	}

	buffSkillNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}
		buffSkill()
		return false
	}

	charKey := char.Base.Key.String()
	c.Events.Subscribe(event.OnOverload, buffSkillNoGadget, fmt.Sprintf("earth-shaker-overload-%s", charKey))
	c.Events.Subscribe(event.OnVaporize, buffSkillNoGadget, fmt.Sprintf("earth-shaker-vaporize-%s", charKey))
	c.Events.Subscribe(event.OnMelt, buffSkillNoGadget, fmt.Sprintf("earth-shaker-melt-%s", charKey))
	c.Events.Subscribe(event.OnSwirlPyro, buffSkillNoGadget, fmt.Sprintf("earth-shaker-pyro-swirl-%s", charKey))
	c.Events.Subscribe(event.OnCrystallizePyro, buffSkillNoGadget, fmt.Sprintf("earth-shaker-pyro-crystallize-%s", charKey))
	c.Events.Subscribe(event.OnBurning, buffSkillNoGadget, fmt.Sprintf("earth-shaker-burning-%s", charKey))
	c.Events.Subscribe(event.OnBurgeon, buffSkillNoGadget, fmt.Sprintf("earth-shaker-burgeon-%s", charKey))

	return w, nil
}
