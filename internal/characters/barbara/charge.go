package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var chargeFrames []int

const chargeHitmark = 90

func init() {
	chargeFrames = frames.InitAbilSlice(90)
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		//check for healing
		if c.Core.Status.Duration(barbSkillKey) > 0 {
			//heal target
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
			//check for healing
			if c.Core.Status.Duration(barbSkillKey) > 0 && energyCount < 5 {
				//regen energy
				c.AddEnergy("barbara-c4", 1)
				energyCount++
			}
		}
	}

	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(ai,
		combat.NewDefCircHit(2, false, combat.TargettableEnemy),
		0,
		chargeHitmark,
		cb,
		c4CB)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
