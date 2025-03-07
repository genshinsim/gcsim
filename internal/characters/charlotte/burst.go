package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames     []int
	burstTickFrames = []int{95, 119, 143, 166, 179, 203, 226, 249}
)

const (
	burstStart        = 53
	burstAttackRadius = 6
	burstHealRadius   = 7
	burstOffsetX      = 0
	burstOffsetY      = 4.5
	burstConsumeDelay = 7
	burstCD           = 1200
	burstInitialAbil  = "Still Photo: Comprehensive Confirmation"
	burstDotAbil      = "Still Photo: Kamera"
	healInitialMsg    = "Still Photo: Comprehensive Confirmation"
	healDotMsg        = "Still Photo: Kamera"
)

func init() {
	burstFrames = frames.InitAbilSlice(70) // Q -> Walk
	burstFrames[action.ActionAttack] = 68
	burstFrames[action.ActionCharge] = 68
	burstFrames[action.ActionSkill] = 69
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionJump] = 58
	burstFrames[action.ActionSwap] = 68
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       burstInitialAbil,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	attackAP := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: burstOffsetX, Y: burstOffsetY}, burstAttackRadius)
	healAP := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: burstOffsetX, Y: burstOffsetY}, burstHealRadius)

	snap := c.Snapshot(&ai)

	healp := snap.Stats[attributes.Heal]
	atk := snap.Stats.TotalATK()
	heal := burstInitialHealFlat[c.TalentLvlBurst()] + atk*burstInitialHealPer[c.TalentLvlBurst()]
	healDot := burstDotHealFlat[c.TalentLvlBurst()] + atk*burstDotHealPer[c.TalentLvlBurst()]

	c.Core.QueueAttack(ai, attackAP, 0, burstStart)

	ai.Abil = burstDotAbil
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.ICDTag = attacks.ICDTagCharlotteKamera
	ai.ICDGroup = attacks.ICDGroupCharlotteKamera
	ai.Durability = 25

	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: healInitialMsg,
			Src:     heal,
			Bonus:   healp,
		})

		for _, value := range burstTickFrames {
			c.Core.Tasks.Add(func() {
				c.Core.QueueAttackWithSnap(ai, snap, attackAP, 0)
				if !c.Core.Combat.Player().IsWithinArea(healAP) {
					return
				}
				c.Core.Player.Heal(info.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: healDotMsg,
					Src:     healDot,
					Bonus:   healp,
				})
			}, value-burstStart)
		}
	}, burstStart)

	c.ConsumeEnergy(burstConsumeDelay)
	c.SetCD(action.ActionBurst, burstCD)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash],
		State:           action.BurstState,
	}, nil
}
