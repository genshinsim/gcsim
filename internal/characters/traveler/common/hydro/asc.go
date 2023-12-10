package hydro

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const a1ICDKey = "sourcewater-droplet-icd"

// After the Dewdrop fired by the Hold Mode of the Aquacrest Saber hits an opponent, a Sourcewater Droplet will be
// generated near to the Traveler. If the Traveler picks it up, they will restore 7% HP.
// 1 Droplet can be created this way every second, and each use of Aquacrest Saber can create 4 Droplets at most.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	count := 0
	return func(a combat.AttackCB) {
		if count >= 4 {
			return
		}
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(a1ICDKey) {
			return
		}

		count++
		droplet := c.newDroplet()
		c.Core.Combat.AddGadget(droplet)
		c.AddStatus(a1ICDKey, 60, true)
	}
}

func (c *char) a1PickUp(count int) {
	for _, g := range c.Core.Combat.Gadgets() {
		if count == 0 {
			return
		}

		droplet, ok := g.(*common.SourcewaterDroplet)
		if !ok {
			continue
		}
		droplet.Kill()
		count--

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Spotless Waters",
			Src:     c.MaxHP() * 0.07,
			Bonus:   c.Stat(attributes.Heal),
		})

		// Picking up a Sourcewater Droplet will restore 2 Energy to the Traveler.
		// Requires the Passive Talent "Spotless Waters."
		if c.Base.Cons >= 1 {
			c.AddEnergy("travelerhydro-c1", 2)
		}

		if c.Base.Cons >= 6 {
			c.c6()
		}
	}
}

func (c *char) newDroplet() *common.SourcewaterDroplet {
	player := c.Core.Combat.Player()
	pos := geometry.CalcRandomPointFromCenter(
		geometry.CalcOffsetPoint(
			player.Pos(),
			geometry.Point{Y: 3.5},
			player.Direction(),
		),
		0.3,
		3,
		c.Core.Rand,
	)
	droplet := common.NewSourcewaterDroplet(c.Core, pos, combat.GadgetTypSourcewaterDropletHydroTrav)
	return droplet
}
