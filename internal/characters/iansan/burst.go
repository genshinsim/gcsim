package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(45) // Q -> Walk
	burstFrames[action.ActionAttack] = 43
	burstFrames[action.ActionSkill] = 42
	burstFrames[action.ActionDash] = 44
	burstFrames[action.ActionJump] = 43
	burstFrames[action.ActionSwap] = 41
}

const (
	burstHitmark = 38

	burstStatus     = "iansan-burst"
	burstBuffStatus = "iansan-burst-buff"

	restoreNSCap   = 99999
	restoreNSRatio = 0.67
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "The Three Principles of Power",
		AttackTag:        attacks.AttackTagElementalBurst,
		AdditionalTags:   []attacks.AttackTag{attacks.AttackTagNightsoul},
		PoiseDMG:         180,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		Element:          attributes.Electro,
		Durability:       25,
		Mult:             burst[c.TalentLvlBurst()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.05,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
	)

	c.Core.Tasks.Add(func() {
		if !c.nightsoulState.HasBlessing() {
			c.enterNightsoul(15)
		} else {
			c.nightsoulState.GeneratePoints(15)
		}
	}, 40)

	c.c2OnBurst()
	c.c4OnBurst()

	duration := 12*60 + c.c6BurstBonusDur()

	c.AddStatus(burstStatus, duration, true)
	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(6)

	c.burstSrc = c.Core.F
	c.burstRestoreNS = 0
	c.updateATKBuff(c.burstSrc)()
	c.restorePointsTask(c.burstSrc)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) burstInit() {
	c.burstBuff = make([]float64, attributes.EndStatType)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:  modifier.NewBaseWithHitlag(burstBuffStatus, -1),
			Extra: true,
			Amount: func() []float64 {
				if !c.StatusIsActive(burstStatus) {
					return nil
				}
				return c.burstBuff
			},
		})
	}

	c.Core.Events.Subscribe(event.OnMovement, c.burstMovementRestore, burstBuffStatus)
}

func (c *char) restorePointsTask(src int) {
	c.Core.Tasks.Add(func() {
		if c.burstSrc != src {
			return
		}
		if !c.StatusIsActive(burstStatus) {
			c.c4Stacks = 0
			return
		}

		points := float64(c.burstRestoreNS)*restoreNSRatio + c.a1Points() + c.c4Points()
		c.burstRestoreNS = 0
		c.pointsOverflow = max(c.nightsoulState.Points()+points-c.nightsoulState.MaxPoints, 0.0)
		if c.pointsOverflow > 0 {
			c.c6OnOverflow()
		}
		if points > 0.0 {
			c.nightsoulState.GeneratePoints(points)
		}
		if points >= 1.0 {
			c.a4Heal()
		}

		c.restorePointsTask(src)
	}, 1*60)
}

func (c *char) updateATKBuff(src int) func() {
	return func() {
		if c.burstSrc != src {
			return
		}

		if !c.StatusIsActive(burstStatus) {
			c.burstBuff[attributes.ATK] = 0
			return
		}

		rate := highATK
		if c.nightsoulState.Points() < 42 {
			rate = lowATK * c.nightsoulState.Points()
		}

		nonExtraAtk := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK).TotalATK()
		c.burstBuff[attributes.ATK] = min(nonExtraAtk*rate, maxATK[c.TalentLvlBurst()])

		c.QueueCharTask(c.updateATKBuff(src), 0.3*60)
	}
}

func (c *char) burstMovementRestore(args ...any) {
	if !c.StatusIsActive(burstStatus) {
		return
	}

	movement := args[0].(float64)
	c.burstRestoreNS = min(c.burstRestoreNS+movement, restoreNSCap)
}
