package clorinde

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFrames           []int
	burstSkillStateFrames []int
	burstHitmarks         = []int{97, 103, 109, 115, 121}
)

const (
	burstCD = 15 * 60
)

func init() {
	burstFrames = frames.InitAbilSlice(128) // Q - N1/E/Jump/Walk
	burstFrames[action.ActionDash] = 127
	burstFrames[action.ActionSwap] = 127

	burstSkillStateFrames = frames.InitAbilSlice(128) // Q - Jump/Walk
	burstSkillStateFrames[action.ActionAttack] = 127
	burstSkillStateFrames[action.ActionSkill] = 127
	burstSkillStateFrames[action.ActionDash] = 127
	burstSkillStateFrames[action.ActionSwap] = 127
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Burst",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagElementalBurst,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeSlash,
		Element:          attributes.Electro,
		Durability:       25,
		Mult:             burstDamage[c.TalentLvlBurst()],
		HitlagHaltFrames: 0.1,
	}
	for _, v := range burstHitmarks {
		// TODO: what's the size of this??
		ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1}, 11.2, 9)
		c.Core.QueueAttack(ai, ap, v, v)
	}

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(14)
	c.QueueCharTask(func() {
		c.ModifyHPDebtByRatio(burstBOL[c.TalentLvlBurst()])
	}, 13)

	if c.StatusIsActive(skillStateKey) {
		return action.Info{
			Frames:          frames.NewAbilFunc(burstSkillStateFrames),
			AnimationLength: burstSkillStateFrames[action.InvalidAction],
			CanQueueAfter:   burstSkillStateFrames[action.ActionSwap], // earliest cancel
			State:           action.BurstState,
		}, nil
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
