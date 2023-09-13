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

const (
	burstFirstShield  = 81
	burstFirstRefresh = 142
	burstRefreshRate  = 146
	burstShieldExpiry = 152
	// subject to change
	burstTickRelease = 21
	burstTickTravel  = 8
)

func init() {
	burstFrames = frames.InitAbilSlice(105) // Q -> CA/D
	burstFrames[action.ActionAttack] = 104
	burstFrames[action.ActionSkill] = 104
	burstFrames[action.ActionJump] = 104
	burstFrames[action.ActionWalk] = 104
	burstFrames[action.ActionSwap] = 102
}

func (c *char) Burst(p map[string]int) action.Info {
	// no heal on first shield
	c.Core.Tasks.Add(func() {
		c.summonSeamlessShield()
	}, burstFirstShield)

	// refresh shield 5 times
	for i := 0; i <= 4; i += 1 {
		c.Core.Tasks.Add(func() {
			c.summonSeamlessShield()
			c.summonSeamlessShieldHealing()
		}, burstFirstShield+burstFirstRefresh+burstRefreshRate*i)
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(5)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) summonSeamlessShield() {
	// add shield
	exist := c.Core.Player.Shields.Get(shield.BaizhuBurst)
	shieldamt := (burstShieldPP[c.TalentLvlBurst()]*c.MaxHP() + burstShieldFlat[c.TalentLvlBurst()])
	if exist != nil {
		c.summonSpiritvein()
	}
	c.Core.Player.Shields.Add(c.newShield(shieldamt, burstShieldExpiry))
}

func (c *char) summonSeamlessShieldHealing() {
	// Seamless Shield Healing
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

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5),
		burstTickRelease,
		burstTickRelease+burstTickTravel,
	)
}
