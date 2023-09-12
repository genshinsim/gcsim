package bennett

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstStartFrame   = 34
	burstBuffDuration = 126
	burstKey          = "bennettburst"
	burstFieldKey     = "bennett-field"
)

func init() {
	burstFrames = frames.InitAbilSlice(53)
	burstFrames[action.ActionDash] = 49
	burstFrames[action.ActionJump] = 50
	burstFrames[action.ActionSwap] = 51
}

func (c *char) Burst(p map[string]int) action.Info {
	// add field effect timer
	// deployable thus not hitlag
	c.Core.Status.Add(burstKey, 720+burstStartFrame)
	// hook for buffs; active right away after cast

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fantastic Voyage",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	const radius = 6.0
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, radius)
	c.Core.QueueAttack(ai, burstArea, 37, 37)

	// add 13 ticks starting from t=0s to t=12s
	// buff appears to start ticking right before hit (t=0s)
	// https://discord.com/channels/845087716541595668/869210750596554772/936507730779308032
	stats, _ := c.Stats()
	for i := 0; i <= 12*60; i += 60 {
		c.Core.Tasks.Add(func() {
			if c.Core.Combat.Player().IsWithinArea(burstArea) {
				c.applyBennettField(stats)()
			}
		}, i+burstStartFrame)
	}

	c.ConsumeEnergy(36)
	c.SetCDWithDelay(action.ActionBurst, 900, 34)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) applyBennettField(stats [attributes.EndStatType]float64) func() {
	hpplus := stats[attributes.Heal]
	heal := bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP()
	pc := burstatk[c.TalentLvlBurst()]
	if c.Base.Cons >= 1 {
		pc += 0.2
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATK] = pc * float64(c.Base.Atk+c.Weapon.BaseAtk)
	if c.Base.Cons >= 6 {
		m[attributes.PyroP] = 0.15
	}

	return func() {
		c.Core.Log.NewEvent("bennett field ticking", glog.LogCharacterEvent, -1)

		// self infuse
		p, ok := c.Core.Combat.Player().(*avatar.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		p.ApplySelfInfusion(attributes.Pyro, 25, burstBuffDuration)

		active := c.Core.Player.ActiveChar()
		// heal if under 70%
		if active.CurrentHPRatio() < 0.7 {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  active.Index,
				Message: "Inspiration Field",
				Src:     heal,
				Bonus:   hpplus,
			})
		}

		// add attack if over 70%
		threshold := .7
		if c.Base.Cons >= 1 {
			threshold = 0
		}
		// Activate attack buff
		if active.CurrentHPRatio() > threshold {
			// add weapon infusion
			if c.Base.Cons >= 6 {
				switch active.Weapon.Class {
				case info.WeaponClassClaymore:
					fallthrough
				case info.WeaponClassSpear:
					fallthrough
				case info.WeaponClassSword:
					c.Core.Player.AddWeaponInfuse(
						active.Index,
						"bennett-fire-weapon",
						attributes.Pyro,
						burstBuffDuration,
						true,
						attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
					)
				}
			}

			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(burstFieldKey, burstBuffDuration),
				AffectedStat: attributes.NoStat,
				Extra:        true,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})

			c.Core.Log.NewEvent("bennett field - adding attack", glog.LogCharacterEvent, c.Index).
				Write("threshold", threshold)
		}
	}
}
