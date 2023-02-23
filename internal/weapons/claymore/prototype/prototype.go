package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.PrototypeArchaic, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := 1.8 + float64(r)*0.6
	const icdKey = "prototype-archaic-icd"

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal && ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			char.AddStatus(icdKey, 15*60, true)
			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Prototype Archaic Proc",
				AttackTag:  attacks.AttackTagWeaponSkill,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       atk,
			}
			trg := args[0].(combat.Target)
			c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Base.Key.String()))
	return w, nil
}
