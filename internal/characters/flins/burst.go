package flins

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

var (
	burstFrames    []int
	symphonyFrames []int
)

var burstHitmarkMid = []int{17, 17 + 9, 17 + 9 + 9, 17 + 9 + 9 + 9}

const (
	burstHitmarkInitial  = 93
	burstHitmarkFinal    = 17 + 9 + 9 + 9 + 27
	symphonyHitmark      = 44
	symphonyExtraHitmark = 44 + 18
)

func init() {
	burstFrames = frames.InitAbilSlice(130)
	burstFrames[action.ActionAttack] = 102
	burstFrames[action.ActionSkill] = 103
	burstFrames[action.ActionDash] = 103
	burstFrames[action.ActionJump] = 103
	burstFrames[action.ActionWalk] = 130
	burstFrames[action.ActionSwap] = 91

	symphonyFrames = frames.InitAbilSlice(81)
	symphonyFrames[action.ActionAttack] = 53
	symphonyFrames[action.ActionSkill] = 53
	symphonyFrames[action.ActionDash] = 51
	symphonyFrames[action.ActionJump] = 50
	symphonyFrames[action.ActionWalk] = 58
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(thunderousSymphonyKey) {
		return c.thunderousSymphony()
	}

	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Cometh the Night (Initial)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       burstInit[c.TalentLvlBurst()],
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 7)
		c.Core.QueueAttack(ai, ap, 0, 0)

		ai = info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Cometh the Night (Mid)",
			AttackTag:        attacks.AttackTagDirectLunarCharged,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDirectLunarCharged,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Electro,
			IgnoreDefPercent: 1,
			Mult:             burstMid[c.TalentLvlBurst()],
		}

		midHits := 2
		if c.getMoonsignLevel() == 2 {
			midHits += 2
		}

		for i := range midHits {
			c.Core.QueueAttack(ai, ap, 0, burstHitmarkMid[i])
		}

		ai.Abil = "Cometh the Night (Final)"
		ai.Mult = burstFinal[c.TalentLvlBurst()]

		c.Core.QueueAttack(ai, ap, 0, burstHitmarkFinal)
	}, burstHitmarkInitial)

	c.SetCD(action.ActionBurst, 15*60)

	c.ConsumeEnergy(4)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) thunderousSymphony() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Thunderous Symphony",
		AttackTag:        attacks.AttackTagDirectLunarCharged,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDirectLunarCharged,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
		Mult:             symphony[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 7)
	c.Core.QueueAttack(ai, ap, symphonyHitmark, symphonyHitmark)

	if c.getMoonsignLevel() == 2 && c.Core.Status.Duration(reactable.LcKey) > 0 {
		ai.Mult = symphonyExtra[c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, ap, symphonyHitmark, symphonyExtraHitmark)
	}

	c.ConsumeEnergyPartial(3, 30)
	c.DeleteStatus(thunderousSymphonyKey)

	return action.Info{
		Frames:          frames.NewAbilFunc(symphonyFrames),
		AnimationLength: symphonyFrames[action.InvalidAction],
		CanQueueAfter:   symphonyFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
