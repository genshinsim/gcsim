package dragonspine

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
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.DragonspineSpear, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	prob := 0.5 + float64(r)*0.1

	const icdKey = "dragonspine-spear-icd"
	icd := 600 // 10s *60

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
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
		if c.Rand.Float64() < prob {
			char.AddStatus(icdKey, icd, true)
			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Dragonspine Proc",
				AttackTag:  attacks.AttackTagWeaponSkill,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       atk,
			}
			if t.AuraContains(attributes.Cryo, attributes.Frozen) {
				ai.Mult = atkc
			}
			c.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3), 0, 1)
		}
		return false
	}, fmt.Sprintf("dragonspine-%v", char.Base.Key.String()))
	return w, nil
}
