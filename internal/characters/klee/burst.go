package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int
var waveHitmarks = []int{186, 294, 401, 503, 610, 718}

const burstStart = 146

func init() {
	burstFrames = frames.InitAbilSlice(139) // Q -> N1/CA/E
	burstFrames[action.ActionDash] = 103    // Q -> D
	burstFrames[action.ActionJump] = 104    // Q -> J
	burstFrames[action.ActionWalk] = 102    // Q -> Walk
	burstFrames[action.ActionSwap] = 101    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sparks'n'Splash",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurst,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               burst[c.TalentLvlBurst()],
		NoImpulse:          true,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	//lasts 10 seconds, starts after 2.2 seconds maybe?
	c.Core.Status.Add("kleeq", 600+burstStart)

	//every 1.8 second +on added shoots between 3 to 5, ignore the queue thing.. space it out .2 between each wave i guess

	// snapshot at end of animation?
	var snap combat.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
	}, 100)

	for _, start := range waveHitmarks {
		c.Core.Tasks.Add(func() {
			//no more if burst has ended early
			if c.Core.Status.Duration("kleeq") <= 0 {
				return
			}
			//wave 1 = 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 0)
			//wave 2 = 1 + 30% chance of 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 12)
			if c.Core.Rand.Float64() < 0.3 {
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 12)
			}
			//wave 3 = 1 + 50% chance of 1
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 24)
			if c.Core.Rand.Float64() < 0.5 {
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 24)
			}
		}, start)
	}

	//every 3 seconds add energy if c6
	if c.Base.Cons >= 6 {
		//TODO: this should eventually use hitlag affected queue and duration
		//but is not big deal right now b/c klee cant experience hitlag without getting hit
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
			x.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("klee-c6", 1500),
				AffectedStat: attributes.PyroP,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}

	c.c1(waveHitmarks[0])

	c.SetCDWithDelay(action.ActionBurst, 15*60, 9)
	c.ConsumeEnergy(12)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel frames
		State:           action.BurstState,
	}
}

// clear klee burst when she leaves the field and handle c4
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		// check if burst is active
		if c.Core.Status.Duration("kleeq") <= 0 {
			return false
		}
		c.Core.Status.Delete("kleeq")

		if c.Base.Cons >= 4 {
			//blow up
			ai := combat.AttackInfo{
				ActorIndex:         c.Index,
				Abil:               "Sparks'n'Splash C4",
				AttackTag:          attacks.AttackTagNone,
				ICDTag:             combat.ICDTagNone,
				ICDGroup:           combat.ICDGroupDefault,
				StrikeType:         combat.StrikeTypeDefault,
				Element:            attributes.Pyro,
				Durability:         50,
				Mult:               5.55,
				CanBeDefenseHalted: true,
				IsDeployable:       true,
			}
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0)
		}

		return false
	}, "klee-exit")
}
