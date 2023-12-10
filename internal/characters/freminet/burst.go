package freminet

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

const (
	burstKey     = "freminet-stalking"
	burstHitmark = 44
)

func init() {
	burstFrames = frames.InitAbilSlice(65)
	burstFrames[action.ActionAttack] = 52
	burstFrames[action.ActionSkill] = 52
	burstFrames[action.ActionDash] = 53
	burstFrames[action.ActionJump] = 52
	burstFrames[action.ActionSwap] = 51
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.AddStatus(burstKey, 600, true)

	c.ResetActionCooldown(action.ActionSkill)

	// TODO: Freminet; Update Info
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Shadowhunter's Ambush",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Cryo,
		Durability:         50,
		Mult:               burst[c.TalentLvlBurst()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.09 * 60,
		CanBeDefenseHalted: false,
	}

	// TODO: Freminet; Insert Hitbox
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.5}, 8)

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(skillArea.Shape.Pos(), nil, 2.5),
		0,
		burstHitmark,
	)

	c.SetCD(action.ActionBurst, 60*15)
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
