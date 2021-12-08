package perception

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("eye of perception", weapon)
	core.RegisterWeaponFunc("eyeofperception", weapon)
}

//Normal and Charged Attacks have a 50% chance to fire a Bolt of Perception,
//dealing 240/270/300/330/360% ATK as DMG. This bolt can bounce between enemies a maximum of 4 times.
//This effect can occur once every 12/11/10/9/8s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	dmg := 2.1 * float64(r) * 0.3
	cd := (13 - r) * 60
	icd := 0
	var w weap

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + cd
		w.snap = char.Snapshot(
			"Eye of Perception Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		w.snap.OnHitCallback = w.chainQ(0, c.F, 1, c, char)
		char.QueueDmg(&w.snap, 1)
		return false
	}, fmt.Sprintf("perception-%v", char.Name()))

	//bounce...
	//d.OnHitCallback = char.chainQ(t.Index(), char.Sim.Frame(), 1)

}

type weap struct {
	snap core.Snapshot
}

func (w *weap) chainQ(index int, src int, count int, c *core.Core, char core.Character) func(t core.Target) {
	if count == 4 {
		return nil
	}
	//check number of targets, if target < 2 then no bouncing
	//figure out the next target
	l := len(c.Targets)
	if l < 2 {
		return nil
	}
	index++
	if index >= l {
		index = 0
	}
	//trigger dmg based on a clone of d
	return func(next core.Target) {
		// c.Log.Printf("hit target %v, frame %v, done proc %v, queuing next index: %v\n", next.Index(), char.Sim.Frame(), count, index)
		d := w.snap.Clone()
		d.Targets = index
		d.SourceFrame = c.F
		d.OnHitCallback = w.chainQ(index, src, count+1, c, char)
		char.QueueDmg(&d, 1)
	}
}
