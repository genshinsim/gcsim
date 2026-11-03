package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// TODO: frames estimated from in-game footage (60 FPS), +-2-4 frames.
var (
	burstFrames         []int
	burstBombardHitmark = []int{135, 142, 149, 156}
)

const burstRayHitmark = 186

func init() {
	burstFrames = frames.InitAbilSlice(170) // Q -> W
	burstFrames[action.ActionAttack] = 158
	burstFrames[action.ActionCharge] = 158
	burstFrames[action.ActionSkill] = 158
	burstFrames[action.ActionDash] = 158
	burstFrames[action.ActionJump] = 158
	burstFrames[action.ActionSwap] = 156
}

// Sandrone summons a large number of Prismatic Resonance Cannons for a frontal bombardment before
// firing on her foes with a Convective Inhibition Ray, dealing AoE Cryo DMG.
// Radiance: Stellar Glimmer: The Convective Inhibition Ray instead deals AoE Cryo DMG that is
// considered the corresponding Stellar Glimmer reaction DMG.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.endHovermech()

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Bombardment",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burstBombard[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6)
	for _, delay := range burstBombardHitmark {
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	c.QueueCharTask(c.burstRay, burstRayHitmark)

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, burstCD*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) burstRay() {
	lvl := c.TalentLvlBurst()
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Convective Inhibition Ray",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burstRay[lvl],
	}
	if c.toStellar(&ai, burstRayConduct[lvl], burstRaySwirl[lvl]) {
		ai.Mult *= c.a1ConsumeTactics()
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1}, 5, 14)
	c.Core.QueueAttack(ai, ap, 0, 0)
}
