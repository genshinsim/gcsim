package xilonen

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const burstStart = 96

func init() {
	burstFrames = frames.InitAbilSlice(101) // Q -> W
	burstFrames[action.ActionAttack] = 94   // Q -> N1
	burstFrames[action.ActionSkill] = 94    // Q -> E
	burstFrames[action.ActionDash] = 95     // Q -> D
	burstFrames[action.ActionJump] = 94     // Q -> J
	burstFrames[action.ActionSwap] = 92     // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Ocelotlicue Point",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalBurst,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Geo,
		Durability:     25,
		Mult:           burstDMG[c.TalentLvlBurst()],
		UseDef:         true,
	}

	// initial hit at 15f after burst start
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), burstStart, burstStart)

	if c.samplersConverted >= 2 {
		for i := 1; i <= 8; i++ {
			hpplus := c.Stat(attributes.Heal)
			heal := burstHealBase[c.TalentLvlBurst()] + c.TotalDef()*burstHealPer[c.TalentLvlBurst()]
			c.Core.Tasks.Add(func() {
				c.Core.Player.Heal(info.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: "Ocelotlicue Point Ebullient Rhythm",
					Src:     heal,
					Bonus:   hpplus,
				})
			}, burstStart+int(88.5*float64(i))) // alternate between 88 and 89 frames
		}
	} else {
		ai.Abil = "Ocelotlicue Point Ardent Rhythm"
		for i := 1; i <= 2; i++ {
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), burstStart+i*33, burstStart+i*33)
		}
	}

	c.ConsumeEnergy(14)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
