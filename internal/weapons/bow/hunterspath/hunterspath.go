package hunterspath

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

const (
	buffKey = "hunterspath-tireless-hunt"
	icdKey  = "hunterspath-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.HuntersPath, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	// Gain 12% All Elemental DMG Bonus. Obtain the Tireless Hunt effect after
	// hitting an opponent with a Charged Attack. This effect increases Charged
	// Attack DMG by 160% of Elemental Mastery. This effect will be removed after
	// 12 Charged Attacks or 10s. Only 1 instance of Tireless Hunt can be gained
	// every 12s.
	w := &Weapon{}
	r := p.Refine

	dmgBonus := 0.09 + 0.03*float64(r)
	val := make([]float64, attributes.EndStatType)
	for i := attributes.PyroP; i <= attributes.DendroP; i++ {
		val[i] = dmgBonus
	}
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("hunterspath-dmg-bonus", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	caBoost := 1.2 + 0.4*float64(r)
	procCount := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		// The buff is a ping dependent action, we're assuming the first hit won't
		// have extra damage.
		if !char.StatusIsActive(icdKey) {
			char.AddStatus(buffKey, 600, true)
			char.AddStatus(icdKey, 720, true)
			procCount = 12
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
		c.Log.NewEvent("hunterspath proc dmg add", glog.LogPreDamageMod, char.Index).
			Write("base_added_dmg", baseDmgAdd).
			Write("remaining_stacks", procCount)
		return false
	}, fmt.Sprintf("hunterspath-%v", char.Base.Key.String()))

	return w, nil
}
