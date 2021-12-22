package venti

import "github.com/genshinsim/gcsim/pkg/core"

func c2cb(a core.AttackCB) {
	a.Target.AddResMod("venti-c2-phys", core.ResistMod{
		Ele:      core.Physical,
		Value:    -0.12,
		Duration: 600,
	})
	a.Target.AddResMod("venti-c2-anemo", core.ResistMod{
		Ele:      core.Physical,
		Value:    -0.12,
		Duration: 600,
	})
}

func c6cb(ele core.EleType) func(a core.AttackCB) {
	return func(a core.AttackCB) {
		a.Target.AddResMod("venti-c6-anemo", core.ResistMod{
			Ele:      ele,
			Value:    -0.20,
			Duration: 600,
		})
	}
}

// func (c *char) applyC2(ds *core.Snapshot) {
// 	ds.OnHitCallback = func(t core.Target) {
// 		t.AddResMod("venti-c2-phys", core.ResistMod{
// 			Ele:      core.Physical,
// 			Value:    -0.12,
// 			Duration: 600,
// 		})
// 		t.AddResMod("venti-c2-anemo", core.ResistMod{
// 			Ele:      core.Physical,
// 			Value:    -0.12,
// 			Duration: 600,
// 		})
// 	}
// }

// func (c *char) applyC6(ds *core.Snapshot, ele core.EleType) {
// 	ds.OnHitCallback = func(t core.Target) {
// 		t.AddResMod("venti-c6-anemo", core.ResistMod{
// 			Ele:      ele,
// 			Value:    -0.20,
// 			Duration: 600,
// 		})
// 	}
// }
