package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames [][]int

var aimedHitmarks = []int{15 - 12, 15, 86}

func init() {
	aimedFrames = make([][]int, 3)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(23 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(23)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot
	aimedFrames[2] = frames.InitAbilSlice(94)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]
}

// Once fully charged, deal Hydro DMG and apply the Riptide status.
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(meleeKey) {
		c.Core.Log.NewEvent("aim called when not in ranged stance", glog.LogActionEvent, c.Index).
			Write("action", action.ActionAim)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 2
	}
	if hold > 2 {
		hold = 2
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Hydro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60, // deployable hitlag, but only on weakspot
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	if hold < 2 {
		ai.Abil = "Aimed Shot"
		if hold == 0 {
			ai.Abil += " (ARCC)"
		}
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
			aimedHitmarks[hold],
			aimedHitmarks[hold]+travel,
		)
	} else {
		c.Core.QueueAttack(
			ai,
			combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
			aimedHitmarks[hold],
			aimedHitmarks[hold]+travel,
			// TODO: what's the ordering on these 2 callbacks?
			c.rtFlashCallback,   // call back for triggering slash
			c.aimedApplyRiptide, // call back for applying riptide
		)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}
}
