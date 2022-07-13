package heizou

import "github.com/genshinsim/gcsim/pkg/core"

//The first Windmuster Iris explosion in each Windmuster Kick will regenerate 9 Elemental Energy for Shikanoin Heizou.
//Every subsequent explosion in that Windmuster Kick will each regenerate an additional 1.5 Energy for Heizou.
//One Windmuster Kick can regenerate a total of 13.5 Energy for Heizou in this manner.
func (c *char) c4(i int) {
	energy := 0.0
	switch i {
	case 1:
		energy += 9.5
	case 2, 3:
		energy += 1.5
	case 4:
		energy += 1.0
	}
	c.AddEnergy("heizou c4", energy)
}

//Each Declension stack will increase the CRIT Rate of the Heartstopper Strike unleashed by 4%.
//When Heizou possesses Conviction, this Heartstopper Strike's CRIT DMG is increased by 32%.

func (c *char) c6(snap *core.Snapshot) {
	snap.Stats[core.CR] += float64(c.decStack) * 0.04
	if c.decStack == 4 {
		snap.Stats[core.CD] += 0.32
	}
}
