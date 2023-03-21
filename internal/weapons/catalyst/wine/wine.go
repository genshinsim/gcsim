package wine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.WineAndSong, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Hitting an opponent with a Normal Attack decreases the Stamina consumption
	//of Sprint or Alternate Sprint by 14% for 5s. Additionally, using a Sprint
	//or Alternate Sprint ability increases ATK by 20% for 5s.

	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = .15 + float64(r)*.05
	stamReduction := .12 + float64(r)*.02
	key := fmt.Sprintf("wineandsong-%v", char.Base.Key.String())
	c.Events.Subscribe(event.OnDash, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("wineandsong", 60*5),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, key)

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if c.Player.Active() != char.Index {
			return false
		}
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		c.Player.AddStamPercentMod(key, 300, func(a action.Action) (float64, bool) {
			if a == action.ActionDash {
				return -stamReduction, false
			}
			return 0, false
		})
		return false
	}, key)

	return w, nil
}
