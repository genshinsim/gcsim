package furina

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames  []int
	chargeHitmark = 32
	chargeOffset  = 0.0
)

func init() {
	chargeFrames = frames.InitAbilSlice(253) // C -> Walk
	chargeFrames[action.ActionAttack] = 37
	chargeFrames[action.ActionSkill] = 47
	chargeFrames[action.ActionBurst] = chargeHitmark
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = chargeHitmark
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Charge %v", c.arkhe),
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		HitlagHaltFrames:   0.02 * 60,
		CanBeDefenseHalted: false,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: chargeOffset},
		2.6,
	)

	if c.Base.Cons >= 6 && c.StatusIsActive(c6Key) {
		ai.Element = attributes.Hydro
		ai.IgnoreInfusion = true
		ai.FlatDmg = c.c6BonusDMG()
		c.Core.QueueAttack(ai, ap, chargeHitmark, chargeHitmark, c.c6cb)
	} else {
		c.Core.QueueAttack(ai, ap, chargeHitmark, chargeHitmark)
	}

	if c.arkhe == ousia {
		c.QueueCharTask(func() {
			c.arkhe = pneuma
			if c.StatusIsActive(skillKey) {
				c.summonSinger(0)
			}
		}, chargeHitmark+1)
	} else {
		c.QueueCharTask(func() {
			c.arkhe = ousia
			if c.StatusIsActive(skillKey) {
				c.summonSalonMembers(0)
			}
		}, chargeHitmark+1)
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}
