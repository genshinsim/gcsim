package xilonen

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames []int

	burstDamageHitmarks = []int{96, 128, 164}
	burstHealHitmarks   = []int{185, 274, 362, 450, 539, 627, 716, 805}
)

func init() {
	burstFrames = frames.InitAbilSlice(101) // Q -> W
	burstFrames[action.ActionAttack] = 95   // Q -> N1
	burstFrames[action.ActionSkill] = 93    // Q -> E
	burstFrames[action.ActionDash] = 95     // Q -> D
	burstFrames[action.ActionJump] = 94     // Q -> J
	burstFrames[action.ActionSwap] = 92     // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Ocelotlicue Point!",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalBurst,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       75.0,
		Element:        attributes.Geo,
		Durability:     25,
		Mult:           burstDMG[c.TalentLvlBurst()],
		UseDef:         true,
		HitlagFactor:   0.05,
	}

	if c.samplersConverted >= 2 {
		c.burstDamage(ai)
	} else {
		c.burstHeal(ai)
	}

	c.ConsumeEnergy(16)
	c.SetCD(action.ActionBurst, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstDamage(ai combat.AttackInfo) {
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
	for i, hitmark := range burstDamageHitmarks {
		if i > 0 {
			ai.Abil = "Follow-Up Beat"
			ai.Mult = burstFollowDMG[c.TalentLvlBurst()]
		}
		c.Core.QueueAttack(ai, ap, hitmark, hitmark)
	}
}

func (c *char) burstHeal(ai combat.AttackInfo) {
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), burstDamageHitmarks[0], burstDamageHitmarks[0])

	hi := info.HealInfo{
		Caller:  c.Index,
		Message: "Ebullient Rhythm",
	}
	for _, hitmark := range burstHealHitmarks {
		c.Core.Tasks.Add(func() {
			hi.Target = c.Core.Player.Active()
			hi.Src = burstHealBase[c.TalentLvlBurst()] + c.TotalDef()*burstHealPer[c.TalentLvlBurst()]
			hi.Bonus = c.Stat(attributes.Heal)
			c.Core.Player.Heal(hi)
		}, hitmark)
	}
}
