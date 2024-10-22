package sturdybone

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	buffKey = "sturdy-bone"
)

func init() {
	core.RegisterWeaponFunc(keys.SturdyBone, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Sprint or Alternate Sprint Stamina Consumption decreased by 15%. Additionally,
// after using Sprint or Alternate Sprint, Normal Attack DMG is increased by 32% of ATK.
// This effect expires after triggering 18 times or 7s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	c.Player.AddStamPercentMod(buffKey, -1, func(a action.Action) (float64, bool) {
		return -0.15, false
	})

	naDmg := 0.12 + 0.4*float64(r)
	c.Events.Subscribe(event.OnDash, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		char.SetTag(buffKey, 18)
		char.AddStatus(buffKey, 7*60, true)

		return false
	}, fmt.Sprintf("sturdybone-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		if char.Tags[buffKey] == 0 {
			return false
		}
		if !char.StatusIsActive(buffKey) {
			return false
		}

		dmgAdded := char.Stat(attributes.ATK) * naDmg
		atk.Info.FlatDmg += dmgAdded
		c.Log.NewEvent("sturdy bone buff", glog.LogPreDamageMod, char.Index).
			Write("damage_added", dmgAdded).
			Write("remaining_stacks", char.Tags[buffKey])

		return false
	}, "sturdy-bone")
	return w, nil
}
