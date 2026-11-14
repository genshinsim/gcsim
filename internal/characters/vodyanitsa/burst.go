package vodyanitsa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHitmark = 95
)

func init() {
	burstFrames = frames.InitAbilSlice(102)
	burstFrames[action.ActionAttack] = 101
	burstFrames[action.ActionCharge] = 114 - burstWindup
	burstFrames[action.ActionDash] = 101
	burstFrames[action.ActionJump] = 101
	burstFrames[action.ActionWalk] = 100
	burstFrames[action.ActionSwap] = 99
}

// Leads a song of the Water Imps which lures its listeners to their doom. This deals AoE Hydro DMG
// to the opponent.
// When in the Song of Ages Past state, DMG dealt is further increased.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	mult := burst[c.TalentLvlBurst()]
	if c.StatusIsActive(skillKey) {
		mult += burstBonus[c.TalentLvlBurst()]
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sink With Thee",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       mult,
		UseHP:      true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
	)

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
