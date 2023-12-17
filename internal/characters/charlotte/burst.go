package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// TODO frame & aoe
var (
	burstFrames     []int
	burstTickFrames = []int{0, 0, 0, 0, 0, 0, 0, 0}
)

const (
	burstStart        = 0
	burstRadius       = 0
	burstConsumeDelay = 0
	burstCD           = 1200
	burstInitialAbil  = "Still Photo: Comprehensive Confirmation"
	burstDotAbil      = "Still Photo: Kamera"
	healInitialMsg    = "Still Photo: Comprehensive Confirmation"
	healDotMsg        = "Still Photo: Kamera"
)

func init() {
	burstFrames = frames.InitAbilSlice(0)
	burstFrames[action.ActionDash] = 0
	burstFrames[action.ActionJump] = 0
	burstFrames[action.ActionSwap] = 0
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
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, burstRadius)

	snap := c.Snapshot(&ai)

	healp := snap.Stats[attributes.Heal]
	atk := snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK]
	heal := burstInitialHealFlat[c.TalentLvlBurst()] + atk*burstInitialHealPer[c.TalentLvlBurst()]
	healDot := burstDotHealFlat[c.TalentLvlBurst()] + atk*burstDotHealPer[c.TalentLvlBurst()]

	c.Core.QueueAttack(ai, ap, 0, burstStart)

	ai.Abil = burstDotAbil
	ai.Mult = burstDot[c.TalentLvlBurst()]

	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: healInitialMsg,
			Src:     heal,
			Bonus:   healp,
		})

		for _, value := range burstTickFrames {
			c.Core.Tasks.Add(func() {
				c.Core.QueueAttackWithSnap(ai, snap, ap, 0)
				if !c.Core.Combat.Player().IsWithinArea(ap) {
					return
				}
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: healDotMsg,
					Src:     healDot,
					Bonus:   healp,
				})
			}, value)
		}
	}, burstStart)

	c.ConsumeEnergy(burstConsumeDelay)
	c.SetCD(action.ActionBurst, burstCD)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // TODO earliest cancel
		State:           action.BurstState,
	}, nil
}
