package redhorn

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.RedhornStonethresher, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//DEF is increased by 28%. Normal and Charged Attack DMG is increased by 40% of DEF.
	w := &Weapon{}
	r := p.Refine

	defBoost := .21 + 0.07*float64(r)
	val := make([]float64, attributes.EndStatType)
	val[attributes.DEFP] = defBoost
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("redhorn-stonethrasher-def-boost", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	nacaBoost := .3 + .1*float64(r)
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if !(atk.Info.AttackTag == attacks.AttackTagNormal || atk.Info.AttackTag == attacks.AttackTagExtra) {
			return false
		}
		baseDmgAdd := (char.Base.Def*(1+char.Stat(attributes.DEFP)) + char.Stat(attributes.DEF)) * nacaBoost
		atk.Info.FlatDmg += baseDmgAdd
		c.Log.NewEvent("Redhorn proc dmg add", glog.LogPreDamageMod, char.Index).
			Write("base_added_dmg", baseDmgAdd)
		return false
	}, fmt.Sprintf("redhorn-%v", char.Base.Key.String()))

	return w, nil
}
