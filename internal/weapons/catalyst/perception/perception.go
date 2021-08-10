package perception

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("eye of perception", weapon)
}

//Normal and Charged Attacks have a 50% chance to fire a Bolt of Perception,
//dealing 240/270/300/330/360% ATK as DMG. This bolt can bounce between enemies a maximum of 4 times.
//This effect can occur once every 12/11/10/9/8s.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	dmg := 2.1 * float64(r) * 0.3
	cd := (13 - r) * 60
	icd := 0
	var w weap

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + cd
		w.snap = c.Snapshot(
			"Eye of Perception Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		w.snap.OnHitCallback = w.chainQ(0, s.Frame(), 1, s, c)
		c.QueueDmg(&w.snap, 1)

	}, fmt.Sprintf("perception-%v", c.Name()))

	//bounce...
	//d.OnHitCallback = c.chainQ(t.Index(), c.Sim.Frame(), 1)

}

type weap struct {
	snap core.Snapshot
}

func (w *weap) chainQ(index int, src int, count int, s core.Sim, c core.Character) func(t core.Target) {
	if count == 4 {
		return nil
	}
	//check number of targets, if target < 2 then no bouncing
	//figure out the next target
	l := len(s.Targets())
	if l < 2 {
		return nil
	}
	index++
	if index >= l {
		index = 0
	}
	//trigger dmg based on a clone of d
	return func(next core.Target) {
		// log.Printf("hit target %v, frame %v, done proc %v, queuing next index: %v\n", next.Index(), c.Sim.Frame(), count, index)
		d := w.snap.Clone()
		d.Targets = index
		d.SourceFrame = s.Frame()
		d.OnHitCallback = w.chainQ(index, src, count+1, s, c)
		c.QueueDmg(&d, 1)
	}
}
