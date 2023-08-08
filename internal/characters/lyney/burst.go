package lyney

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// 3 * 60 duration
// burst, none, default, default, no hitlag
// every 0.15*60 perform true ST attack with burst
// at end of duration: explosiveFirework with 6m aoe around player
// burst ends either via duration or swap
// 15s cooldown
// energy drain

var burstFrames []int

const (
	// TODO: proper frames
	burstKey         = "lyney-q"
	burstInterval    = 0.15 * 60
	burstDuration    = 3*60 + 1 // + 1 for final tick
	burstCD          = 15 * 60
	burstCDStart     = 1
	burstEnergyDelay = 7
)

func init() {
	// TODO: proper min frames, currently using thoma
	burstFrames = frames.InitAbilSlice(58)
	burstFrames[action.ActionAttack] = 57
	burstFrames[action.ActionSkill] = 56
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionSwap] = 56
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.AddStatus(burstKey, burstDuration, true)
	c.QueueCharTask(c.burstTick, burstInterval)
	c.QueueCharTask(c.explosiveFirework, burstDuration)

	c.SetCDWithDelay(action.ActionBurst, burstCD, burstCDStart)
	c.ConsumeEnergy(burstEnergyDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill], // TODO: proper frames, should be earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		c.explosiveFirework()
		return false
	}, "lyney-exit")
}

func (c *char) burstTick() {
	if !c.StatusIsActive(burstKey) {
		return
	}

	tickAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wondrous Trick: Miracle Parade",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), nil)
	for i := 0; i < len(enemies); i++ {
		c.Core.QueueAttack(tickAI, combat.NewSingleTargetHit(enemies[i].Key()), 0, 0)
	}

	c.QueueCharTask(c.burstTick, burstInterval)
}

func (c *char) explosiveFirework() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.DeleteStatus(burstKey)

	explodeAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wondrous Trick: Miracle Parade (Explosive Firework)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       explosiveFirework[c.TalentLvlBurst()],
	}
	qPos := c.Core.Combat.Player().Pos()
	// TODO: proper frames
	c.Core.QueueAttack(explodeAI, combat.NewCircleHitOnTarget(qPos, nil, 6), 5, 5)

	if len(c.hats) < c.maxHatCount {
		g := c.newGrinMalkinHat(qPos, false) // TODO: double check hp drain on this
		c.Core.Combat.AddGadget(g)
		c.hats = append(c.hats, g)
	}

	c.increasePropSurplusStacks(1)
}
