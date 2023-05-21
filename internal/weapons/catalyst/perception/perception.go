package perception

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
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.EyeOfPerception, NewWeapon)
}

type Weapon struct {
	Index int
	ai    combat.AttackInfo
	snap  combat.Snapshot
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const bounceKey = "eye-of-perception-bounce"

func (w *Weapon) chain(count int, c *core.Core, char *character.CharWrapper) func(a combat.AttackCB) {
	if count == 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		// check target is an enemey
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		// shouldn't proc more than one chain if multiple enemies are hit
		if done {
			return
		}
		done = true

		next := c.Combat.ClosestEnemyWithinArea(
			combat.NewCircleHitOnTarget(t, nil, 8),
			func(e combat.Enemy) bool {
				return !e.StatusIsActive(bounceKey)
			},
		)
		if next != nil {
			next.AddStatus(bounceKey, 36, true)
			c.QueueAttackWithSnap(w.ai, w.snap, combat.NewCircleHitOnTarget(next, nil, 0.6), 10, w.chain(count+1, c, char))
		}
	}
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "perception-icd"
	cd := (13 - r) * 60
	dmg := 2.1 * float64(r) * 0.3

	w.ai = combat.AttackInfo{
		ActorIndex: char.Index,
		Abil:       "Eye of Preception Proc",
		AttackTag:  attacks.AttackTagWeaponSkill,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 100,
		Mult:       dmg,
	}

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagNormal && ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, cd, true)
		w.snap = char.Snapshot(&w.ai)

		enemy := c.Combat.ClosestEnemyWithinArea(
			combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 8),
			func(e combat.Enemy) bool {
				return !e.StatusIsActive(bounceKey)
			},
		)
		if enemy != nil {
			enemy.AddStatus(bounceKey, 36, true)
			c.QueueAttackWithSnap(w.ai, w.snap, combat.NewCircleHitOnTarget(enemy, nil, 0.6), 10, w.chain(0, c, char))
		}

		return false
	}, fmt.Sprintf("perception-%v", char.Base.Key.String()))

	return w, nil
}
