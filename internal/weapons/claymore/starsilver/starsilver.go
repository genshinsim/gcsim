package starsilver

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
	core.RegisterWeaponFunc(keys.SnowTombedStarsilver, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := 0.65 + float64(r)*0.15
	mc := 1.6 + float64(r)*0.4
	prob := 0.5 + float64(r)*0.1

	const icdKey = "starsilver-icd"

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < prob {
			char.AddStatus(icdKey, 600, true)
			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Starsilver Proc",
				AttackTag:  attacks.AttackTagWeaponSkill,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       m,
			}
			if t.AuraContains(attributes.Cryo, attributes.Frozen) {
				ai.Mult = mc
			}
			c.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3), 0, 1)
		}
		return false
	}, fmt.Sprintf("starsilver-%v", char.Base.Key.String()))
	return w, nil
}
