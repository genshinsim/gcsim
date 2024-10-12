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
	burstFrames = frames.InitAbilSlice(93) // Q -> D/J
	// burstFrames[action.ActionAttack] = 88  // Q -> N1
	// burstFrames[action.ActionSkill] = 89   // Q -> E
	// burstFrames[action.ActionSwap] = 88    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ocelotlicue Point",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burstDMG[c.TalentLvlBurst()],
		UseDef:     true,
	}

	// initial hit at 15f after burst start
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), burstStart, burstStart)

	if c.samplersConverted >= 2 {
		for i := 0; i < 8; i++ {
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
			}, burstStart+i*90+8)
		}
	} else {
		ai.Abil = "Ocelotlicue Point Ardent Rhythm"
		for i := 1; i <= 2; i++ {
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), burstStart+i*60, burstStart+i*60)
		}
	}

	c.ConsumeEnergy(2)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 2)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
