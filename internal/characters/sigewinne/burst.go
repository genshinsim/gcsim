package sigewinne

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var endLag []int

const (
	earliestCancel = 122
	chargeBurstDur = 99
)

func init() {
	endLag = frames.InitAbilSlice(271 - 241) // Burst end -> Skill
	endLag[action.ActionAttack] = 269 - 241
	endLag[action.ActionSwap] = 245 - 241
	endLag[action.ActionWalk] = 258 - 241
	endLag[action.ActionDash] = 0
	endLag[action.ActionJump] = 0
}

func (c *char) burstFindDroplets() {
	droplets := c.getSourcewaterDroplets()

	// TODO: If droplets time out before the "droplet check" it doesn't count.
	indices := c.Core.Combat.Rand.Perm(len(droplets))
	orbs := 0
	for _, ind := range indices {
		g := droplets[ind]
		c.consumeDroplet(g)
		orbs += 1
	}
	c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", orbs, " droplets"), glog.LogCharacterEvent, c.Index)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.burstEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Elemental Burst", c.Base.Key)
	}

	ticks, ok := p["ticks"]
	if !ok {
		ticks = -1
	} else {
		ticks = max(ticks, 2)
	}

	c.Core.Player.SwapCD = math.MaxInt16
	c.tickAnimLength = getBurstHitmark(1)

	c.Core.Tasks.Add(func() {
		// TODO: correct timing?
		c.burstFindDroplets()

		c.burstStartF = c.Core.F
		c.Core.Tasks.Add(c.burstTick(c.burstStartF, 1, ticks, false), getBurstHitmark(1))
	}, chargeBurstDur)

	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)
	c.ConsumeEnergy(5)

	c.addC2Shield()

	return action.Info{
		Frames: func(next action.Action) int {
			return chargeBurstDur + c.tickAnimLength + endLag[next]
		},
		AnimationLength: chargeBurstDur + c.burstMaxDuration,
		CanQueueAfter:   earliestCancel,
		State:           action.BurstState,
		OnRemoved: func(next action.AnimationState) {
			// need to calculate correct swap cd in case of early cancel
			switch next {
			case action.DashState, action.JumpState:
				c.Core.Player.SwapCD = max(player.SwapCDFrames-(c.Core.F-c.lastSwap), 0)
			}
			c.removeC2Shield()
		},
	}, nil
}

func (c *char) burstTick(src, tick, maxTick int, last bool) func() {
	return func() {
		if c.burstStartF != src {
			return
		}
		// no longer in burst anim -> no tick
		if c.Core.F > c.burstStartF+c.burstMaxDuration {
			return
		}

		if last {
			c.Core.Player.SwapCD = endLag[action.ActionSwap]
			return
		}

		// tick param supplied and hit the limit -> proc wave, enable early cancel flag for next action check and stop queuing ticks
		if tick == maxTick {
			c.burstWave()
			c.burstEarlyCancelled = true
			c.Core.Player.SwapCD = endLag[action.ActionSwap]
			return
		}

		c.burstWave()

		// next tick handling
		if maxTick == -1 || tick < maxTick {
			tickDelay := getBurstHitmark(tick + 1)
			// calc new animation length to be up until next tick happens
			nextTickAnimLength := c.Core.F - c.burstStartF + tickDelay

			c.Core.Tasks.Add(c.burstTick(src, tick+1, maxTick, false), tickDelay)

			// queue up last tick if next tick would happen after burst duration ends
			if nextTickAnimLength > c.burstMaxDuration {
				// queue up final tick to happen at end of burst duration
				c.Core.Tasks.Add(c.burstTick(src, tick+1, maxTick, true), c.burstMaxDuration-c.tickAnimLength)
				// update tickAnimLength to be equal to entire burst duration at the end
				c.tickAnimLength = c.burstMaxDuration
			} else {
				// next tick happens within burst duration -> update tickAnimLength as usual
				c.tickAnimLength = nextTickAnimLength
			}
		}
	}
}

func (c *char) burstWave() {
	// TODO: the ACTUAL hitbox???
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 10)

	// TODO: is deployable?
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Super Saturated Syringing",
		AttackTag:    attacks.AttackTagElementalBurst,
		ICDTag:       attacks.ICDTagElementalBurst,
		ICDGroup:     attacks.ICDGroupSigewinneBurst,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Hydro,
		Durability:   25,
		FlatDmg:      burstDMG[c.TalentLvlAttack()] * c.MaxHP(),
		HitlagFactor: 0.01,
	}
	c.Core.QueueAttack(ai, ap, 0, 0, c.c2CB)
}

func getBurstHitmark(tick int) int {
	switch tick {
	case 1:
		return 0
	default:
		return 25
	}
}
