package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	itzpapaInterval           = 59
	obsidianTzitzimitlHitmark = 20

	opalFireStateKey = "opal-fire-state"
	frostFallAbil    = "Frostfall Storm DMG"
)

var (
	skillFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(50) // E -> Walk
	skillFrames[action.ActionAttack] = 42
	skillFrames[action.ActionCharge] = 42
	skillFrames[action.ActionBurst] = 41
	skillFrames[action.ActionDash] = 49
	skillFrames[action.ActionJump] = 49
	skillFrames[action.ActionSwap] = 41
}

func (c *char) Skill(_ map[string]int) (action.Info, error) {
	// do initial attack
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Obsidian Tzitzimitl DMG",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
		HitlagFactor:   0.01,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6),
		obsidianTzitzimitlHitmark,
		obsidianTzitzimitlHitmark,
		c.particleCB,
	)

	// to do with delay
	c.QueueCharTask(func() {
		c.SetCD(action.ActionSkill, 16*60)
	}, 18)

	c.QueueCharTask(c.addShield, 37)

	c.QueueCharTask(func() {
		// summon Itzpapa and immediately check if Opal Fire state can be activated
		c.nightsoulState.EnterTimedBlessing(c.nightsoulState.Points()+24, 20*60, c.exitNightsoul)
		c.itzpapaSrc = c.Core.F
		c.tryEnterOpalFireState(c.itzpapaSrc)
	}, 22)

	// to do now
	if c.Base.Cons >= 1 {
		c.numStellarBlades = 10
	}

	if c.Base.Cons >= 6 {
		currentPoints := c.nightsoulState.Points()
		c.nightsoulState.ClearPoints()
		c.numC6Stacks = min(maxC6Stacks, currentPoints)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) exitNightsoul() {
	c.numC6Stacks = 0
	c.numStellarBlades = 0
	c.DeleteStatus(opalFireStateKey)
	c.nightsoulState.ExitBlessing()
}

func (c *char) generateNightsoulPoints(amount float64) {
	c.nightsoulState.GeneratePoints(amount)
	c.tryEnterOpalFireState(c.itzpapaSrc)
}

// try to activate Opal Fire each time Citlali gains NS points to avoid event subscribtion
func (c *char) tryEnterOpalFireState(src int) {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	if c.nightsoulState.Points() < 50 && c.Base.Cons < 6 {
		return
	}
	// if it's activation or REactivation (of Opal Fire state)
	if c.StatusIsActive(opalFireStateKey) {
		return
	}
	c.AddStatus(opalFireStateKey, -1, false)
	c.itzpapaHitTask(src)
	c.nightsoulPointReduceTask(src)
}

func (c *char) nightsoulPointReduceTask(src int) {
	const tickInterval = .1
	c.QueueCharTask(func() {
		if c.itzpapaSrc != src {
			return
		}
		if !c.StatusIsActive(opalFireStateKey) {
			return
		}

		// reduce 0.8 point every 6f, which is 8 per second
		c.nightsoulState.ConsumePoints(0.8)
		if c.Base.Cons >= 6 {
			c.numC6Stacks = min(maxC6Stacks, c.numC6Stacks+0.8)
		}
		if c.nightsoulState.Points() < 0.001 && c.Base.Cons < 6 {
			c.DeleteStatus(opalFireStateKey)
			return
		}

		c.nightsoulPointReduceTask(src)
	}, 60*tickInterval)
}

func (c *char) itzpapaHitTask(src int) {
	c.QueueCharTask(func() {
		if src != c.itzpapaSrc {
			return
		}
		if !c.StatusIsActive(opalFireStateKey) {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           frostFallAbil,
			AttackTag:      attacks.AttackTagElementalArt,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagCitlaliFrostfallStorm,
			ICDGroup:       attacks.ICDGroupCitlaliFrostfallStorm,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Cryo,
			Durability:     25,
			Mult:           frostfall[c.TalentLvlSkill()],
			FlatDmg:        c.a4Dmg(frostFallAbil),
			HitlagFactor:   0.01,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 6), 0, 0, c.c4SkullCB)
		c.itzpapaHitTask(src)
	}, itzpapaInterval)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Cryo, c.ParticleDelay)
}
