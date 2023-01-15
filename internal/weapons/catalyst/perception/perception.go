package perception

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
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

const bounceKey = "eye-of-preception-bounce"

func (w *Weapon) chain(count int, c *core.Core, char *character.CharWrapper) func(a combat.AttackCB) {
	if count == 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		//check target is an enemey
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		t.SetTag(bounceKey, c.F+36)
		trgs := c.Combat.EnemyByDistance(a.Target.Shape().Pos(), a.Target.Key())
		next := -1
		for _, v := range trgs {
			trg, ok := c.Combat.Enemy(v).(*enemy.Enemy)
			if !ok {
				continue
			}
			if trg.GetTag(bounceKey) < c.F {
				next = v
				break
			}
		}

		if next == -1 {
			return
		}

		cb := w.chain(count+1, c, char)
		c.QueueAttackWithSnap(w.ai, w.snap, combat.NewCircleHitOnTarget(c.Combat.Enemy(next), nil, 0.6), 10, cb)
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
		AttackTag:  combat.AttackTagWeaponSkill,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 100,
		Mult:       dmg,
	}

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, cd, true)

		cb := w.chain(0, c, char)
		w.snap = char.Snapshot(&w.ai)
		c.QueueAttackWithSnap(w.ai, w.snap, combat.NewSingleTargetHit(c.Combat.DefaultTarget), 10, cb)

		return false
	}, fmt.Sprintf("perception-%v", char.Base.Key.String()))

	return w, nil
}
