package black

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheBlackSword, NewWeapon)
}

//Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
//regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.15 + 0.05*float64(r)
	char.AddAttackMod(character.AttackMod{Base: modifier.NewBase("blacksword", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return nil, false
		}
		return val, true
	}})

	last := 0
	heal := 0.5 + .1*float64(r)
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if crit && (c.F-last >= 300 || last == 0) {
			c.Player.Heal(player.HealInfo{
				Caller:  char.Index,
				Target:  c.Player.Active(),
				Message: "The Black Sword",
				Src:     heal * (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]),
				Bonus:   char.Stat(attributes.Heal),
			})
			//trigger cd
			last = c.F
		}
		return false
	}, fmt.Sprintf("black-sword-%v", char.Base.Key.String()))
	return w, nil
}
