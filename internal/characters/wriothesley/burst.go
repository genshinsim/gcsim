package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// TODO: chongyun based frames & my assumptions
var burstFrames []int

var burstHitmarks = []int{50, 59, 67, 76, 85}

const (
	burstOusiaHitmark = 95

	burstOusiaICDKey = "wriothesley-ousia-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(79) // Q -> Swap
	burstFrames[action.ActionAttack] = 64  // Q -> N1
	burstFrames[action.ActionSkill] = 64   // Q -> E
	burstFrames[action.ActionDash] = 64    // Q -> D
	burstFrames[action.ActionJump] = 66    // Q -> J
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
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.00}, 5.00, 15.00)

	for _, hitmark := range burstHitmarks {
		c.Core.QueueAttack(ai, ap, hitmark, hitmark)
	}

	if !c.StatusIsActive(burstOusiaICDKey) {
		aiOusia := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Surging Blade",
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Mult:               burstOusia[c.TalentLvlBurst()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.03 * 60,
			CanBeDefenseHalted: false,
		}
		c.Core.QueueAttack(
			aiOusia,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.00}, 7.00, 15.00),
			burstOusiaHitmark,
			burstOusiaHitmark,
		)

		c.AddStatus(burstOusiaICDKey, 10*60, true)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(6)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
