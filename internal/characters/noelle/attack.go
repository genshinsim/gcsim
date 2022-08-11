package noelle

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var attackFrames [][]int
var attackHitmarks = []int{28, 25, 20, 42}
var attackHitlagHaltFrame = []float64{0.10, 0.10, 0.09, 0.15}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 38)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 107)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	r := 0.3
	if c.StatModIsActive(burstBuffKey) {
		r = 2
		if c.NormalCounter == 2 {
			//q-n3 has different hit lag
			ai.HitlagHaltFrames = 0.1 * 60
		}
	}
	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		//check for healing
		if c.Core.Player.Shields.Get(shield.ShieldNoelleSkill) != nil {
			var prob float64
			if c.Base.Cons >= 1 && c.StatModIsActive(burstBuffKey) {
				prob = 1
			} else {
				prob = healChance[c.TalentLvlSkill()]
			}
			if c.Core.Rand.Float64() < prob {
				//heal target
				x := a.AttackEvent.Snapshot.BaseDef*(1+a.AttackEvent.Snapshot.Stats[attributes.DEFP]) + a.AttackEvent.Snapshot.Stats[attributes.DEF]
				heal := shieldHeal[c.TalentLvlSkill()]*x + shieldHealFlat[c.TalentLvlSkill()]
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  -1,
					Message: "Breastplate (Attack)",
					Src:     heal,
					Bonus:   a.AttackEvent.Snapshot.Stats[attributes.Heal],
				})
				done = true
			}
		}

	}
	// need char queue because of potential hitlag from C4
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), r, false, combat.TargettableEnemy),
			0,
			0,
			cb,
		)
	}, attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	c.a4Counter++
	if c.a4Counter == 4 {
		c.a4Counter = 0
		if c.Cooldown(action.ActionSkill) > 0 {
			c.ReduceActionCooldown(action.ActionSkill, 60)
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
