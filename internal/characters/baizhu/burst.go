package baizhu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var burstFrames []int

const burstStartFrame = 114

func init() {
	burstFrames = frames.InitAbilSlice(115)
	burstFrames[action.ActionSwap] = 114
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	for i := 0; i <= 14*60; i += 150 { //Every 2.5s, 14s duration
		c.QueueCharTask(func() {
			c.summonSeamlessShield()
			c.summonSeamlessShieldHealing()

		}, i+burstStartFrame)
	}
	if c.Base.Cons >= 4 { //TODO: Need a delay for this buff to trigger?
		c.c4()
	}

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(3) //TODO:Exact timing

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) summonSeamlessShield() {
	// add shield
	exist := c.Core.Player.Shields.Get(shield.ShieldBaizhuBurst)
	shieldamt := (burstShieldPP[c.TalentLvlBurst()]*c.MaxHP() + burstShieldFlat[c.TalentLvlBurst()])
	if exist != nil {
		c.summonSpiritvein()
	}
	c.Core.Player.Shields.Add(c.newShield(shieldamt, 151))
}

func (c *char) summonSeamlessShieldHealing() {
	//Seamless Shield Healing
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Seamless Shield Healing",
		Src:     burstHealPP[c.TalentLvlBurst()]*c.MaxHP() + burstHealFlat[c.TalentLvlBurst()],
		Bonus:   c.Stat(attributes.Heal),
	})
	c.a4()

}

func (c *char) summonSpiritvein() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spiritvein Damage",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstAtk[c.TalentLvlBurst()],
	}
	if c.Base.Cons >= 6 {
		ai.FlatDmg = c.MaxHP() * 0.08
	}

	// TODO: strike timing
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5),
		0,
		10,
	)
}
