package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
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
	effectLastProc := -9999

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.F < effectLastProc+15*60 {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			effectLastProc = c.F
			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Prototype Archaic Proc",
				AttackTag:  combat.AttackTagWeaponSkill,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: combat.StrikeTypeDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       atk,
			}
			c.QueueAttack(ai, combat.NewDefCircHit(.6, false, combat.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Base.Key.String()))
	return w, nil
}
