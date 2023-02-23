package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	chargeFrames   []int
	chargeHitmarks = []int{14, 27}
)

func init() {
	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionSkill] = 29
	chargeFrames[action.ActionBurst] = 29
	chargeFrames[action.ActionDash] = 14
	chargeFrames[action.ActionJump] = 15
	chargeFrames[action.ActionSwap] = 52
}

// since E is aoe, so this should be considered aoe too
// hitWeakPoint: tartaglia can proc Prototype Cresent's Passive on Geovishap's weakspots.
// Evidence: https://youtu.be/oOfeu5pW0oE
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	if !c.StatusIsActive(MeleeKey) {
		c.Core.Log.NewEvent("charge called when not in melee stance", glog.LogActionEvent, c.Index).
			Write("action", action.ActionCharge)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}

	hitWeakPoint, ok := p["hitWeakPoint"]
	if !ok {
		hitWeakPoint = 0
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Charged Attack",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       combat.ICDTagExtraAttack,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeSlash,
		Element:      attributes.Hydro,
		Durability:   25,
		HitWeakPoint: hitWeakPoint != 0,
	}

	for i, mult := range eCharge {
		ai.Mult = mult[c.TalentLvlSkill()]
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.2),
			chargeHitmarks[i],
			chargeHitmarks[i],
			c.makeA4CB(),      // callback for applying riptide
			c.rtSlashCallback, // callback for triggering slash
		)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash], // earliest cancel
		State:           action.ChargeAttackState,
	}
}
