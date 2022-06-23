package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var chargeFramesNoWindup []int
var chargeFrames []int

const (
	//TODO: this hitmark should also change depending on windup or not?
	chargeHitmark = 70
	//TODO: stacks technically only last 15s and each stack has its own
	//timer
	conductiveTag = "lisa-conductive-stacks"
)

func init() {
	chargeFrames = frames.InitAbilSlice(94) //used to say 91??
	chargeFrames[action.ActionAttack] = 86
	chargeFrames[action.ActionCharge] = 90
	chargeFrames[action.ActionSkill] = 94
	chargeFrames[action.ActionBurst] = 93
	chargeFrames[action.ActionSwap] = 90

	//no wind up
	chargeFramesNoWindup = frames.InitAbilSlice(94) //used to say 91??
	chargeFramesNoWindup[action.ActionAttack] = 86
	chargeFramesNoWindup[action.ActionCharge] = 90
	chargeFramesNoWindup[action.ActionSkill] = 94
	chargeFramesNoWindup[action.ActionBurst] = 93
	chargeFramesNoWindup[action.ActionSwap] = 90
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	//TODO: we should be checking previous action to return frames here
	//and use chargeFramesNoWindup where it applies
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	cb := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		count := t.GetTag(conductiveTag)
		if count < 3 {
			t.SetTag(conductiveTag, count+1)
		}
	}

	c.Core.QueueAttack(ai,
		combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
		0,
		chargeHitmark, //no travel for lisa
		cb,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
