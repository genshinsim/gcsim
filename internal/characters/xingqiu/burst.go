package xingqiu

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) summonSwordWave() {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guhua Sword: Raincutter",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	//only if c.nextRegen is true and first sword
	var c2cb, c6cb func(a core.AttackCB)
	if c.nextRegen {
		c6cb = func(a core.AttackCB) {
			c.AddEnergy("xingqiu-c6", 3)
		}
	}
	if c.Base.Cons >= 2 {
		icd := -1
		c2cb = func(a core.AttackCB) {
			if c.Core.F < icd {
				return
			}
			icd = c.Core.F + 1

			c.AddTask(func() {
				a.Target.AddResMod("xingqiu-c2", core.ResistMod{
					Ele:      core.Hydro,
					Value:    -0.15,
					Duration: 4 * 60,
				})
			}, "xq-sword-debuff", 1)
		}
	}

	for i := 0; i < c.numSwords; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 20, 20, c2cb, c6cb)
		c6cb = nil
		c.burstCounter++
	}

	//figure out next wave # of swords
	switch c.numSwords {
	case 2:
		c.numSwords = 3
		c.nextRegen = false
	case 3:
		if c.Base.Cons == 6 {
			c.numSwords = 5
			c.nextRegen = true
		} else {
			c.numSwords = 2
			c.nextRegen = false
		}
	case 5:
		c.numSwords = 2
		c.nextRegen = false
	}

	c.burstSwordICD = c.Core.F + 60
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(core.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return false
		}
		next := args[1].(core.AnimationState)
		//ignore if not normal
		if next != core.NormalAttackState {
			return false
		}
		//ignore if on ICD
		if c.burstSwordICD > c.Core.F {
			return false
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonSwordWave()
		c.Core.Log.NewEvent("xq burst on state change", core.LogCharacterEvent, c.Index, "state", next, "icd", c.burstSwordICD)
		c.burstTickSrc = c.Core.F
		c.AddTask(c.burstTickerFunc(c.Core.F), "xq-ticker", 60) //check every 1sec

		return false
	}, "xq-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("xq burst tick check ignored, src diff", core.LogCharacterEvent, c.Index, "src", src, "new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.State()
		if state != core.NormalAttackState {
			c.Core.Log.NewEvent("xq burst tick check stopped, not normal state", core.LogCharacterEvent, c.Index, "src", src, "state", state)
			return
		}
		c.Core.Log.NewEvent("xq burst triggered from ticker", core.LogCharacterEvent, c.Index, "src", src, "state", state, "icd", c.burstSwordICD)
		//we can trigger a wave here b/c we're in normal state still and src is still the same
		c.summonSwordWave()
		//in theory this should not hit an icd?
		c.AddTask(c.burstTickerFunc(src), "xq-ticker", 60) //check every 1sec
	}
}

/**
func (c *char) burstHook() {
	c.Core.Events.Subscribe(core.PostAttack, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return false
		}
		delay := 0 //wait 5 frames into attack animation
		//check if off ICD
		if c.burstSwordICD > c.Core.F {
			f := args[0].(int)
			//if burst icd is between current frame and f (animation end frame)
			//then we should queue up a sword anyways at when the icd comes up
			if c.burstSwordICD <= c.Core.F+f+10 {
				delay = c.burstSwordICD - c.Core.F
				c.Core.Log.NewEvent("Xingqiu Q on animation delay", core.LogCharacterEvent, c.Index,  "icd", c.burstSwordICD, "f", f, "check", c.Core.F+f+10)
			} else {
				return false
			}
		}

		//trigger swords, only first sword applies hydro
		for i := 0; i < c.numSwords; i++ {

			wave := i

			if delay > 0 {

				c.AddTask(func() {

					d := c.Snapshot(
						"Guhua Sword: Raincutter",
						core.AttackTagElementalBurst,
						core.ICDTagElementalBurst,
						core.ICDGroupDefault,
						core.StrikeTypePierce,
						core.Hydro,
						25,
						burst[c.TalentLvlBurst()],
					)
					d.Targets = 0 //only hit main target
					d.OnHitCallback = func(t core.Target) {
						//check energy
						if c.nextRegen && wave == 0 {
							c.AddEnergy(3)
						}
						//check c2
						if c.Base.Cons >= 2 {
							c.AddTask(func() {
								t.AddResMod("xingqiu-c2", core.ResistMod{
									Ele:      core.Hydro,
									Value:    -0.15,
									Duration: 4 * 60,
								})
							}, "xq-sword-debuff", 1)

						}
					}

					c.QueueDmg(&d, 20)

				}, "sword-wave", delay)
			} else {
				d := c.Snapshot(
					"Guhua Sword: Raincutter",
					core.AttackTagElementalBurst,
					core.ICDTagElementalBurst,
					core.ICDGroupDefault,
					core.StrikeTypePierce,
					core.Hydro,
					25,
					burst[c.TalentLvlBurst()],
				)
				d.Targets = 0 //only hit main target
				d.OnHitCallback = func(t core.Target) {
					//check energy
					if c.nextRegen && wave == 0 {
						c.AddEnergy(3)
					}
					//check c2
					if c.Base.Cons >= 2 {
						c.AddTask(func() {
							t.AddResMod("xingqiu-c2", core.ResistMod{
								Ele:      core.Hydro,
								Value:    -0.15,
								Duration: 4 * 60,
							})
						}, "xq-sword-debuff", 1)

					}
				}

				c.QueueDmg(&d, 20)

			}

			c.burstCounter++
		}

		//figure out next wave # of swords
		switch c.numSwords {
		case 2:
			c.numSwords = 3
			c.nextRegen = false
		case 3:
			if c.Base.Cons == 6 {
				c.numSwords = 5
				c.nextRegen = true
			} else {
				c.numSwords = 2
				c.nextRegen = false
			}
		case 5:
			c.numSwords = 2
			c.nextRegen = false
		}

		//estimated 1 second ICD
		c.burstSwordICD = c.Core.F + 60 + delay

		return false
	}, "xq-burst")
}**/
