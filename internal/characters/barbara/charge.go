package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var chargeFrames []int

const chargeHitmark = 55

func init() {
	chargeFrames = frames.InitAbilSlice(89)
	chargeFrames[action.ActionDash] = 56
	chargeFrames[action.ActionJump] = 56
	chargeFrames[action.ActionSwap] = 55
	chargeFrames[action.ActionSkill] = 88
	chargeFrames[action.ActionBurst] = 87
	chargeFrames[action.ActionCharge] = 88
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		// check for healing
		if c.Core.Status.Duration(barbSkillKey) > 0 {
			// heal target
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Melody Loop (Charged Attack)",
				Src:     4 * (prochpp[c.TalentLvlSkill()]*c.MaxHP() + prochp[c.TalentLvlSkill()]),
				Bonus:   c.Stat(attributes.Heal),
			})
			done = true
		}
	}
	var c4CB combat.AttackCBFunc
	if c.Base.Cons >= 4 {
		energyCount := 0
		c4CB = func(a combat.AttackCB) {
			if a.Target.Type() != targets.TargettableEnemy {
				return
			}
			// check for healing
			if c.Core.Status.Duration(barbSkillKey) > 0 && energyCount < 5 {
				// regen energy
				c.AddEnergy("barbara-c4", 1)
				energyCount++
			}
		}
	}

	// skip CA windup if we're in NA/CA animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 14
	}

	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 5}, 3),
		chargeHitmark-windup,
		chargeHitmark-windup,
		cb,
		c4CB,
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}
}
