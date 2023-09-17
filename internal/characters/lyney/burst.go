package lyney

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const (
	burstKey           = "lyney-q"
	burstMarkKey       = "lyney-q-mark"
	burstStart         = 100
	burstFirstInterval = 12
	burstInterval      = 0.15 * 60
	burstDuration      = 182
	burstCD            = 15 * 60
	burstEnergyDelay   = 8
)

func init() {
	burstFrames = frames.InitAbilSlice(321) // Q -> Walk
	burstFrames[action.ActionAttack] = 297
	burstFrames[action.ActionAim] = 297
	burstFrames[action.ActionSkill] = 101
	burstFrames[action.ActionDash] = 287
	burstFrames[action.ActionJump] = 285
	burstFrames[action.ActionSwap] = 283
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) action.Info {
	c.QueueCharTask(func() {
		c.AddStatus(burstKey, burstDuration, true)
		c.QueueCharTask(c.burstTick, burstFirstInterval)
		c.QueueCharTask(c.explosiveFirework, burstDuration)
	}, burstStart)

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(burstEnergyDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill],
		State:           action.BurstState,
	}
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

	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), func(e combat.Enemy) bool {
		return !e.StatusIsActive(burstMarkKey)
	})
	for _, enemy := range enemies {
		c.Core.QueueAttack(tickAI, combat.NewSingleTargetHit(enemy.Key()), 0, 0)
		enemy.AddStatus(burstMarkKey, c.StatusDuration(burstKey), true)
	}

	c.QueueCharTask(c.burstTick, burstInterval)
}

func (c *char) explosiveFirework() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.DeleteStatus(burstKey)
	for _, v := range c.Core.Combat.Enemies() {
		e, ok := v.(*enemy.Enemy)
		if !ok {
			continue
		}
		if !e.StatusIsActive(burstMarkKey) {
			continue
		}
		e.DeleteStatus(burstMarkKey)
	}

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
	c.Core.QueueAttack(explodeAI, combat.NewCircleHitOnTarget(qPos, nil, 6), 6, 6)

	// kill existing hat if reached limit
	if len(c.hats) == c.maxHatCount {
		c.hats[0].Kill()
	}
	g := c.newGrinMalkinHat(qPos, false, grinMalkinHatBurstDuration)
	c.hats = append(c.hats, g)
	c.Core.Combat.AddGadget(g)

	c.increasePropSurplusStacks(1)
}
