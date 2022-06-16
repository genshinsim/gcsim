package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const burstStart = 101

func init() {
	burstFrames = frames.InitAbilSlice(101)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sparks'n'Splash",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
		NoImpulse:  true,
	}
	//lasts 10 seconds, starts after 2.2 seconds maybe?
	c.Core.Status.Add("kleeq", 600+132)

	//every 1.8 second +on added shoots between 3 to 5, ignore the queue thing.. space it out .2 between each wave i guess

	// snapshot at end of animation?
	var snap combat.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
	}, 100)

	for i := 132; i < 732; i += 108 {
		c.Core.Tasks.Add(func() {
			//no more if burst has ended early
			if c.Core.Status.Duration("kleeq") <= 0 {
				return
			}
			//wave 1 = 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 0)
			//wave 2 = 1 + 30% chance of 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 12)
			if c.Core.Rand.Float64() < 0.3 {
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 12)
			}
			//wave 3 = 1 + 50% chance of 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 24)
			if c.Core.Rand.Float64() < 0.5 {
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 24)
			}
		}, i)
	}

	//every 3 seconds add energy if c6
	if c.Base.Cons >= 6 {
		for i := burstStart + 180; i < burstStart+600; i += 180 {
			c.Core.Tasks.Add(func() {
				//no more if burst has ended early
				if c.Core.Status.Duration("kleeq") <= 0 {
					return
				}

				for i, x := range c.Core.Player.Chars() {
					if i == c.Index {
						continue
					}
					x.AddEnergy("klee-c6", 3)
				}
			}, i)
		}

		// add 10% pyro for 25s
		m := make([]float64, attributes.EndStatType)
		m[attributes.PyroP] = .1
		for _, x := range c.Core.Player.Chars() {
			x.AddStatMod("klee-c6", 1500, attributes.PyroP, func() ([]float64, bool) {
				return m, true
			})
		}
	}

	c.c1(132)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 15)
	c.ConsumeEnergy(15)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstStart,
		State:           action.BurstState,
	}
}

// clear klee burst when she leaves the field and handle c4
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// check if burst is active
		if c.Core.Status.Duration("kleeq") <= 0 {
			return false
		}
		c.Core.Status.Delete("kleeq")

		if c.Base.Cons >= 4 {
			//blow up
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Sparks'n'Splash C4",
				AttackTag:  combat.AttackTagNone,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    attributes.Pyro,
				Durability: 50,
				Mult:       5.55,
			}
			c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0, 0)
		}

		return false
	}, "klee-exit")
}
