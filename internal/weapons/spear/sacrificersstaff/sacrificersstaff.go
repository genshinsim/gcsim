package sacrificersstaff

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

const buffKey = "sacrificersstaff-buff"

func init() {
	core.RegisterWeaponFunc(keys.SacrificersStaff, NewWeapon)
}

type Weapon struct {
	Index  int
	buff   []float64
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// For 6s after an Elemental Skill hits an opponent, ATK is increased by 8% and Energy Recharge is increased by 6%. Max 3 stacks. This effect can be triggered even when the equipping character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	w.buff = make([]float64, attributes.EndStatType)
	atkBuff := 0.08 + 0.02*float64(r)
	erBuff := 0.06 + 0.015*float64(r)

	// add on hit effect
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		w.stacks = max(w.stacks+1, 3)
		w.buff[attributes.ATKP] = atkBuff * float64(w.stacks)
		w.buff[attributes.ER] = erBuff * float64(w.stacks)

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 300),
			Amount: func() ([]float64, bool) {
				return w.buff, true
			},
		})
		return false
	}, fmt.Sprintf("sacrificersstaff-%v", char.Base.Key.String()))

	return w, nil
}
