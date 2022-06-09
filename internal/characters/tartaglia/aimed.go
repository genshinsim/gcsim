package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int

const aimedHitmark = 84

//Once fully charged, deal Hydro DMG and apply the Riptide status.
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   combat.StrikeTypePierce,
		Element:      attributes.Hydro,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		aimedHitmark,
		aimedHitmark+travel,
		//TODO: what's the ordering on these 2 callbacks?
		c.rtFlashCallback,   //call back for triggering slash
		c.aimedApplyRiptide, //call back for applying riptide
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		Post:            aimedHitmark,
		State:           action.AimState,
	}
}
