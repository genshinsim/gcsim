package travelerelectro

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-1, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Blade",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	hits, ok := p["hits"]
	if !ok {
		hits = 1
	}

	maxAmulets := 2
	if c.Base.Cons > 0 {
		maxAmulets = 3
	}

	// clear existing amulets
	c.Amulets = make([]abundanceAmulet, maxAmulets)

	// accept param input to disable dashing to amulets on swap
	ignoreAmulets, ok := p["ignore_amulets"]
	if !ok {
		ignoreAmulets = 0
	}

	// should emc wait and collect one of the amulets?
	forceMcToCollectAmulet, ok := p["force_collect"]
	if !ok {
		forceMcToCollectAmulet = 0
	}

	// Counting from the frame E is pressed, it takes an average of 1.79 seconds for a character to be able to pick one up
	// https://library.keqingmains.com/evidence/characters/electro/traveler-electro#amulets-delay
	amuletDelay := 107 // ~1.79s

	c.QueueParticle(c.Name(), 1, core.Electro, f+100)

	for i := 0; i < hits; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f)

		c.AddTask(func() {
			// generate amulet if generated amulets < limit
			if i >= maxAmulets {
				return
			}

			// log amulet generated event
			a := abundanceAmulet{
				Collected:    false,
				CollectableF: c.Core.F + amuletDelay,
			}

			c.Amulets = append(c.Amulets, a)
		}, fmt.Sprintf("emc-amulet-generate-%d", i), amuletDelay)
	}

	c.SetCD(core.ActionSkill, 810+21) //13.5s, starts 21 frames in

	if forceMcToCollectAmulet == 1 && len(c.Amulets) > 0 {
		// wait amuletDelay frames
		f += amuletDelay

		// then try to collect an amulet
		for i := 0; i < len(c.Amulets); i++ {
			ok := c.Amulets[i].tryToCollect(c.Core.F, c, c)
			if ok {
				break
			}
		}
	}

	if ignoreAmulets == 0 && len(c.Amulets) > 0 {
		// Assume next swap(s) will dash to the amulet if they can
		c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
			for i := 0; i < len(c.Amulets); i++ {
				next := args[2].(*core.Character)
				ok := c.Amulets[i].tryToCollect(c.Core.F, *next, c)
				if ok {
					break
				}
			}

			return false
		}, "check-pickup-abundance-amulet")
	}

	return f, a
}

/**
[12:01 PM] pai: never tried to measure it but emc burst looks like it has roughly 1~1.5 abyss tile of range, skill goes a bit further i think
[12:01 PM] pai: the 3 hits from the skill also like split out and kind of auto target if that's useful information
**/
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Bellowing Thunder",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), 0, f)

	//1573 start, 1610 cd starts, 1612 energy drained, 1633 first swapable
	c.ConsumeEnergy(42)
	c.SetCD(core.ActionBurst, 1200+37)
	return f, a
}
