package arlecchino

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const burstHitmarks = 110
const balemoonRisingHealAbil = "Balemoon Rising (Heal)"

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(146)
	burstFrames[action.ActionAttack] = 113
	burstFrames[action.ActionCharge] = 124
	burstFrames[action.ActionDash] = 111
	burstFrames[action.ActionJump] = 113
	burstFrames[action.ActionSwap] = 145
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Balemoon Rising",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)
	c.QueueCharTask(c.absorbDirectives, 22)

	c.QueueCharTask(func() { c.ResetActionCooldown(action.ActionSkill) }, 107)
	c.Core.QueueAttack(ai, skillArea, burstHitmarks, burstHitmarks)

	// video seems to have a lot of delay due to ping
	// Probably should be burst hitmark +1
	c.QueueCharTask(c.balemoonRisingHeal, 123)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(12)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) balemoonRisingHeal() {
	amt := 1.5*c.CurrentHPDebt() + 1.5*c.getTotalAtk()
	// call the template healing method directly to bypass Heal override
	c.Character.Heal(&info.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: balemoonRisingHealAbil,
		Src:     amt,
		Bonus:   c.Stat(attributes.Heal),
	})
}
