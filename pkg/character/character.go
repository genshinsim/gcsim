package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Character struct {
	Core  *core.Core
	Index int
	character.Character

	SkillCon int
	BurstCon int

	//normal attack counter
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int
}

func (c *Character) ConsumeEnergy(delay int) {
	if delay == 0 {
		c.Core.Player.ByIndex(c.Index).Energy = 0
		return
	}
	c.Core.Tasks.Add(func() {
		c.Core.Player.ByIndex(c.Index).Energy = 0
	}, delay)
}

func (c *Character) AddEnergy(src string, e float64) {
	char := c.Core.Player.ByIndex(c.Index)
	preEnergy := char.Energy
	char.Energy += e
	if char.Energy > char.EnergyMax {
		char.Energy = char.EnergyMax
	}
	if char.Energy < 0 {
		char.Energy = 0
	}

	c.Core.Events.Emit(event.OnEnergyChange, c, preEnergy, e, src)
	c.Core.Log.NewEvent("adding energy", glog.LogEnergyEvent, c.Index,
		"rec'd", e,
		"post_recovery", char.Energy,
		"source", src,
		"max_energy", char.EnergyMax,
	)
}

func (c *Character) ReceiveParticle(p character.Particle, isActive bool, partyCount int) {
	char := c.Core.Player.ByIndex(c.Index)
	var amt, er, r float64
	r = 1.0
	if !isActive {
		r = 1.0 - 0.1*float64(partyCount)
	}
	//recharge amount - particles: same = 3, non-ele = 2, diff = 1
	//recharge amount - orbs: same = 9, non-ele = 6, diff = 3 (3x particles)
	switch {
	case p.Ele == char.Base.Element:
		amt = 3
	case p.Ele == attributes.NoElement:
		amt = 2
	default:
		amt = 1
	}
	amt = amt * r //apply off field reduction
	//apply energy regen stat

	er = c.Core.Player.ByIndex(c.Index).Stat(attributes.ER)

	amt = amt * (1 + er) * float64(p.Num)

	pre := char.Energy

	char.Energy += amt
	if char.Energy > char.EnergyMax {
		char.Energy = char.EnergyMax
	}

	c.Core.Events.Emit(event.OnEnergyChange, c, pre, amt, p.Source)
	c.Core.Log.NewEvent(
		"particle",
		glog.LogEnergyEvent,
		c.Index,

		"source", p.Source,
		"count", p.Num,
		"ele", p.Ele,
		"ER", er,
		"is_active", isActive,
		"party_count", partyCount,
		"pre_recovery", pre,
		"amt", amt,
		"post_recovery", char.Energy,
		"max_energy", char.EnergyMax,
	)
}

func (c *Character) ResetNormalCounter() {
	c.NormalCounter = 0
}

func (c *Character) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *Character) NextNormalCounter() int {
	return c.NormalCounter + 1
}
