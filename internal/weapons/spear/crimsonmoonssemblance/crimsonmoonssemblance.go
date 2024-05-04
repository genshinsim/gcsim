package crimsonmoonssemblance

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

const (
	icdKey      = "crimsonmoonssemblance-icd"
	icdDuration = 14 * 60
)

func init() {
	core.RegisterWeaponFunc(keys.CrimsonMoonsSemblance, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	refine := p.Refine

	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}

		if c.Player.Active() != char.Index {
			return false
		}

		if char.StatusIsActive(icdKey) {
			return false
		}

		char.AddStatus(icdKey, icdDuration, true)
		char.ModifyHPDebtByRatio(0.25)

		return false
	}, fmt.Sprintf("crimsonmoonssemblance-hit-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		maxhp := char.MaxHP()
		m := make([]float64, attributes.EndStatType)

		if char.Index != index {
			return false
		}

		if char.CurrentHPDebt() > 0 {
			m[attributes.DmgP] += 0.08 + 0.04*float64(refine)
		}

		if char.CurrentHPDebt() >= 0.3*maxhp {
			m[attributes.DmgP] += 0.16 + 0.08*float64(refine)
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("crimsonmoonssemblance-bonus", -1),
			AffectedStat: attributes.DmgP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("crimsonmoonssemblance-hp-debt-%v", char.Base.Key.String()))

	return w, nil
}
