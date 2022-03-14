package perception

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("eye of perception", weapon)
	core.RegisterWeaponFunc("eyeofperception", weapon)
}

//Normal and Charged Attacks have a 50% chance to fire a Bolt of Perception,
//dealing 240/270/300/330/360% ATK as DMG. This bolt can bounce between enemies a maximum of 4 times.
//This effect can occur once every 12/11/10/9/8s.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dmg := 2.1 * float64(r) * 0.3
	cd := (13 - r) * 60
	icd := 0
	var w weap

	c.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*coretype.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return false
		}
		if icd > c.Frame {
			return false
		}
		icd = c.Frame + cd

		ai := core.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Eye of Preception Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		w.atk = &coretype.AttackEvent{
			Info:     ai,
			Snapshot: char.Snapshot(&ai),
		}
		atk := *w.atk
		atk.SourceFrame = c.Frame
		atk.Pattern = core.NewDefSingleTarget(0, coretype.TargettableEnemy)
		cb := w.chain(0, c, char)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Combat.QueueAttackEvent(&atk, 10) //TODO: no idea actual travel time
		return false
	}, fmt.Sprintf("perception-%v", char.Name()))

	//bounce...
	//d.OnHitCallback = char.chainQ(t.Index(), char.Sim.Frame(), 1)

	//on hit find next target not marked. marks lasts 60 seconds
	return "eyeofperception"
}

type weap struct {
	atk *coretype.AttackEvent
}

const bounceKey = "eye-of-preception-bounce"

func (w *weap) chain(count int, c *core.Core, char coretype.Character) func(a core.AttackCB) {
	if count == 4 {
		return nil
	}
	return func(a core.AttackCB) {
		//mark the current target, then grab nearest target not marked
		//and trigger another attack while count < 4
		a.Target.SetTag(bounceKey, c.Frame+36) //lock out for 0.6s
		x, y := a.Target.Shape().Pos()
		trgs := c.EnemyByDistance(x, y, a.Target.Index())
		next := -1
		for _, v := range trgs {
			trg := c.Targets[v]
			if trg.GetTag(bounceKey) < c.Frame {
				next = v
				break
			}
		}
		//do nothing if no targets found
		if next == -1 {
			return
		}
		//we have a target so trigger an atk
		atk := *w.atk
		atk.SourceFrame = c.Frame
		atk.Pattern = core.NewDefSingleTarget(next, coretype.TargettableEnemy)
		cb := w.chain(count+1, c, char)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Combat.QueueAttackEvent(&atk, 10) //TODO: no idea actual travel time
	}
}
