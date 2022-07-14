package heizou

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) a1() {
	swirlCB := func() func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			if c.a1icd > c.Core.F {
				return false
			}
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			//icd is triggered regardless if stacks are maxed or not
			c.a1icd = c.Core.F + 6
			c.addDecStack()
			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlCB(), "heizou-a1-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlCB(), "heizou-a1-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlCB(), "heizou-a1-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlCB(), "heizou-a1-pyro")
}

func (c *char) a4() {

	dur := 60 * 10
	c.Core.Status.AddStatus("heizoua4", dur)
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue //nothing for heizou
		}
		char.AddMod(core.CharStatMod{
			Key:    "heizou-a4",
			Expiry: c.Core.F + dur,
			Amount: func() ([]float64, bool) {
				return c.a4buff, true
			},
		})
	}

	c.Core.Log.NewEvent("heizou a4 triggered", core.LogCharacterEvent, c.Index, "em snapshot", c.a4buff[core.EM], "expiry", c.Core.F+dur)
}
