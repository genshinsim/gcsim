package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var paramitaChargeFrames []int

const chargeHitmark = 19
const paramitaChargeHitmark = 6

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(62)
	chargeFrames[action.ActionAttack] = 57
	chargeFrames[action.ActionSkill] = 57
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// charge (paramita) -> x
	paramitaChargeFrames = frames.InitAbilSlice(44)
	paramitaChargeFrames[action.ActionBurst] = 35
	paramitaChargeFrames[action.ActionDash] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionJump] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionSwap] = 42
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	var hitmark int
	var act action.ActionInfo
	var bbcb combat.AttackCBFunc

	if c.Core.Status.Duration("paramita") > 0 {
		//[3:56 PM] Isu: My theory is that since E changes attack animations, it was coded
		//to not expire during any attack animation to simply avoid the case of potentially
		//trying to change animations mid-attack, but not sure how to fully test that
		//[4:41 PM] jstern25| ₼WHO_SUPREMACY: this mostly checks out
		//her e can't expire during q as well
		if paramitaChargeHitmark > c.Core.Status.Duration("paramita") {
			c.Core.Status.Add("paramita", paramitaChargeHitmark)
			// c.S.Status["paramita"] = c.Core.F + f //extend this to barely cover the burst
		}
		bbcb = c.applyBB
		//charge land 182, tick 432, charge 632, tick 675
		//charge land 250, tick 501, charge 712, tick 748

		//e cast at 123, animation ended 136 should end at 664 if from cast or 676 if from animation end, tick at 748 still buffed?

		// adjust frames in paramita state
		hitmark = chargeHitmark
		act = action.ActionInfo{
			Frames:          frames.NewAbilFunc(paramitaChargeFrames),
			AnimationLength: paramitaChargeFrames[action.InvalidAction],
			CanQueueAfter:   hitmark,
			State:           action.ChargeAttackState,
		}
	} else {
		hitmark = paramitaChargeHitmark
		act = action.ActionInfo{
			Frames:          frames.NewAbilFunc(chargeFrames),
			AnimationLength: chargeFrames[action.InvalidAction],
			CanQueueAfter:   hitmark,
			State:           action.ChargeAttackState,
		}
	}

	//check for particles
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupPole,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), hitmark, hitmark, c.ppParticles, bbcb)

	return act
}
