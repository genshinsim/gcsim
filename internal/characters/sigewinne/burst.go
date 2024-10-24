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
	chargeBurstDur = 99 - 6
	burstName      = "Super Saturated Syringing"
)

func init() {
	endLag = frames.InitAbilSlice(17) // Burst end -> walk
	endLag[action.ActionAttack] = 32
	endLag[action.ActionSkill] = 30
	endLag[action.ActionSwap] = 4
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
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Elemental Burst", c.CharWrapper.Base.Key)
	}

	if c.Base.Cons >= 2 {
		c.Core.Tasks.Add(c.addC2Shield, 1)
	}

	c.QueueCharTask(c.burstFindDroplets, chargeBurstDur)

	c.tickAnimLength = getBurstHitmark(1)
	c.chargeAi = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       burstName,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    burstDMG[c.TalentLvlAttack()] * c.MaxHP(),
	}

	c.burstStartF = c.Core.F + chargeBurstDur

	ticks, ok := p["ticks"]
	if !ok {
		ticks = -1
	} else {
		ticks = max(ticks, 2)
	}

	c.Core.Player.SwapCD = math.MaxInt16

	c.Core.Tasks.Add(c.burstTick(c.burstStartF, 1, ticks, false), chargeBurstDur+getBurstHitmark(1))

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(5)

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
			if c.Base.Cons >= 2 {
				c.removeC2Shield()
			}
			return
		}

		if last {
			c.Core.Player.SwapCD = endLag[action.ActionSwap]
			if c.Base.Cons >= 2 {
				c.removeC2Shield()
			}
			return
		}

		// tick param supplied and hit the limit -> proc wave, enable early cancel flag for next action check and stop queuing ticks
		if tick == maxTick {
			c.burstWave()
			c.burstEarlyCancelled = true
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
	c.Core.QueueAttack(c.chargeAi, ap, 0, 0)
}

func getBurstHitmark(tick int) int {
	switch tick {
	case 1:
		return 6
	default:
		return 25
	}
}
