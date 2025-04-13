package mizuki

const (
	c1Key      = "mizuki-c1-key"
	c1Interval = 3.5 * 60
	c1Duration = 3 * 60
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	//swirlfunc := func(args ...interface{}) bool {
	//	if _, ok := args[0].(*enemy.Enemy); !ok {
	//		return false
	//	}
	//
	//	// Only when dream drifter is active
	//	if !c.StatusIsActive(dreamDrifterStateKey) {
	//		return false
	//	}
	//
	//	atk := args[1].(*combat.AttackEvent)
	//
	//	if !atk.Reacted {
	//		return false
	//	}
	//
	//	e := args[0].(*enemy.Enemy)
	//	if !e.StatusIsActive(c1Key) {
	//		return false
	//	}
	//
	//	e.DeleteStatus(c1Key)
	//	cbs := args[2].(*[]combat.AttackCBFunc)
	//
	//	*cbs = append(*cbs, func(cb combat.AttackCB) {
	//		cb.AttackEvent.
	//	})
	//
	//	return false
	//}
	//
	//c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc, "mizuki-c1-pyro-swirl")
	//c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc, "mizuki-c1-hydro-swirl")
	//c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc, "mizuki-c1-electro-swirl")
	//c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc, "mizuki-c1-cryo-swirl")
}
