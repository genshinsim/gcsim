package escoffier

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	initialHeal = 92 // depends on ping
	hitmark     = 94
)

func init() {
	burstFrames = frames.InitAbilSlice(109)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Scoring Cuts",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), 0, 0, c.makeA4CB())
	}, hitmark)

	// initial heal
	c.QueueCharTask(func() {
		heal := burstHealFlat[c.TalentLvlBurst()] + burstHealPer[c.TalentLvlBurst()]*c.TotalAtk()
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Scoring Cuts (Healing)",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})
	}, initialHeal)

	c.QueueCharTask(c.a1, hitmark)

	c.SetCD(action.ActionBurst, int(burstCD[c.TalentLvlBurst()])*60)
	c.ConsumeEnergy(6)

	c.c1()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
