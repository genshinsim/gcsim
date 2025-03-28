package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(72)
	burstFrames[action.ActionAttack] = 71
	burstFrames[action.ActionSkill] = 71
	burstFrames[action.ActionJump] = 70
	burstFrames[action.ActionSwap] = 69
}

const (
	burstHitmark = 34

	burstStatus     = "kinetic-energy"
	burstBuffStatus = "iansan-burst-buff"
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Three Principles of Power",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagElementalBurst,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		burstHitmark,
		burstHitmark,
	)

	c.Core.Tasks.Add(func() {
		if !c.nightsoulState.HasBlessing() {
			c.enterNightsoul(15)
		} else {
			c.nightsoulState.GeneratePoints(15)
		}
	}, burstHitmark)

	c.burstSrc = c.Core.F
	c.burstRestoreNS = 0
	c.updateATKBuff()
	c.applyBuffTask(c.burstSrc)
	c.Core.Events.Subscribe(event.OnActionExec, c.burstMovementRestore, burstBuffStatus)

	if c.Base.Cons >= 2 {
		c.a1ATK()
	}
	c.c4Stacks = 0

	duration := 12 * 60
	if c.Base.Cons >= 6 {
		duration += 3.0
	}
	c.AddStatus(burstStatus, duration, false) // TODO: hitlag affected?
	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(3)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) applyBuffTask(src int) {
	c.Core.Tasks.Add(func() {
		if c.burstSrc != src {
			return
		}
		if !c.StatusIsActive(burstStatus) {
			c.c4Stacks = 0
			return
		}

		points := float64(c.burstRestoreNS) + c.a1Points() + c.c4Points()
		c.burstRestoreNS = 0
		c.pointsOverflow = max(c.nightsoulState.Points()+points-c.nightsoulState.MaxPoints, 0.0)
		if c.pointsOverflow > 0 {
			c.c6()
		}
		if points > 0.0 {
			c.nightsoulState.GeneratePoints(points)
		}

		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(burstBuffStatus, 1*60),
			Amount: func() ([]float64, bool) {
				c.c2ATKBuff(active)
				return c.burstBuff, true
			},
		})
		c.applyBuffTask(src)
	}, 0.5*60) // TODO: refresh rate?
}

func (c *char) updateATKBuff() {
	if !c.StatusIsActive(burstStatus) {
		c.burstBuff[attributes.ATK] = 0
		return
	}

	rate := highATK[c.TalentLvlBurst()]
	if c.nightsoulState.Points() < 42 {
		rate = lowATK[c.TalentLvlBurst()] * c.nightsoulState.Points()
	}
	c.burstBuff[attributes.ATK] = min(c.TotalAtk()*rate, maxATK[c.TalentLvlBurst()])
}

func (c *char) burstMovementRestore(args ...interface{}) bool {
	if !c.StatusIsActive(burstStatus) {
		return true
	}

	param := args[2].(map[string]int)
	movement, ok := param["movement"]
	if ok {
		c.burstRestoreNS += movement
	}
	return false
}
