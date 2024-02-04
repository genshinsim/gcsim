package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

var burstHitmarks = []int{99, 104, 109, 114, 119}

const (
	burstOusiaHitmark = 160
	burstOusiaICDKey  = "wriothesley-ousia-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(133)
	burstFrames[action.ActionSwap] = 94
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Darkgold Wolfbite",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 8, 16)

	// TODO: snapshot timing
	snap := c.Snapshot(&ai)
	c.c2(&snap)
	for _, hitmark := range burstHitmarks {
		c.Core.QueueAttackWithSnap(ai, snap, ap, hitmark)
	}

	c.QueueCharTask(func() {
		if c.StatusIsActive(burstOusiaICDKey) {
			return
		}
		c.AddStatus(burstOusiaICDKey, 10*60, true)

		aiOusia := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             "Surging Blade",
			AttackTag:        attacks.AttackTagElementalBurst,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Cryo,
			Durability:       0,
			Mult:             burstOusia[c.TalentLvlBurst()],
			HitlagFactor:     0.01,
			HitlagHaltFrames: 0.03 * 60,
		}
		c.Core.QueueAttackWithSnap(
			aiOusia,
			snap,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 10, 16),
			0,
		)
	}, burstOusiaHitmark)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
