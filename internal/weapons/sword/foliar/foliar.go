package foliar

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
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffKey = "foliar-whitemoon-bristle"
	icdKey  = "foliar-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.FoliarIncision, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//CRIT Rate is increased by 4%. 
	//After Normal Attacks deal Elemental DMG, the Foliar Incision effect will be obtained, 
	//increasing DMG dealt by Normal Attacks and Elemental Skills by 120% of Elemental Mastery. 
	//This effect will disappear after 28 DMG instances or 12s. You can obtain Foliar Incision once every 12s.
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.03 + float64(r)*0.01
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("foliar-crit-rate", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	caBoost := 0.09 + 0.3*float64(r)
	procCount := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if !(atk.Info.AttackTag == combat.AttackTagNormal || atk.Info.AttackTag == combat.AttackTagElementalArt || atk.Info.AttackTag == combat.AttackTagElementalArtHold) {
			return false
		}
		// The buff is a ping dependent action, we're assuming the first hit won't
		// have extra damage.
		if !char.StatusIsActive(icdKey) {
			char.AddStatus(buffKey, 600, true)
			char.AddStatus(icdKey, 720, true)
			procCount = 28
			return false
		}
		if !char.StatusIsActive(buffKey) {
			return false
		}
		baseDmgAdd := char.Stat(attributes.EM) * caBoost
		atk.Info.FlatDmg += baseDmgAdd
		procCount -= 1
		if procCount <= 0 {
			char.DeleteStatus(buffKey)
		}
		c.Log.NewEvent("foliarincision proc dmg add", glog.LogPreDamageMod, char.Index).
			Write("base_added_dmg", baseDmgAdd).
			Write("remaining_stacks", procCount)
		return false
	}, fmt.Sprintf("foliarincision-%v", char.Base.Key.String()))

	return w, nil
}