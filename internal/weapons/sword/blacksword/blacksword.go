package black

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheBlackSword, NewWeapon)
}

// Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
// regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.15 + 0.05*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("blacksword", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			return val, true
		},
	})

	const icdKey = "black-sword-icd"
	heal := 0.5 + .1*float64(r)
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if crit {
			c.Player.Heal(player.HealInfo{
				Caller:  char.Index,
				Target:  c.Player.Active(),
				Message: "The Black Sword",
				Src:     heal * (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]),
				Bonus:   char.Stat(attributes.Heal),
			})
			// trigger cd
			char.AddStatus(icdKey, 300, true) // every 5s
		}
		return false
	}, fmt.Sprintf("black-sword-%v", char.Base.Key.String()))
	return w, nil
}
