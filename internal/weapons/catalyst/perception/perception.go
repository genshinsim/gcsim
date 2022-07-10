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
	atk   *combat.AttackEvent
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
		x, y := a.Target.Shape().Pos()
		trgs := c.Combat.EnemyByDistance(x, y, a.Target.Index())
		next := -1
		for _, v := range trgs {
			trg, ok := c.Combat.Target(v).(*enemy.Enemy)
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

		atk := *w.atk
		atk.SourceFrame = c.F
		atk.Pattern = combat.NewDefSingleTarget(next, combat.TargettableEnemy)
		cb := w.chain(count+1, c, char)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.QueueAttackEvent(&atk, 10)
	}
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 2.1 * float64(r) * 0.3
	cd := (13 - r) * 60
	icd := 0

	c.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + cd

		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Eye of Preception Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		w.atk = &combat.AttackEvent{
			Info:     ai,
			Snapshot: char.Snapshot(&ai),
		}
		atk := *w.atk
		atk.SourceFrame = c.F
		atk.Pattern = combat.NewDefSingleTarget(0, combat.TargettableEnemy)
		cb := w.chain(0, c, char)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.QueueAttackEvent(&atk, 10)
		return false
	}, fmt.Sprintf("perception-%v", char.Base.Key.String()))

	return w, nil
}
