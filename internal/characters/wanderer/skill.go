package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	skillKey = "windfavored-state"
)

var skillEndFrames []int

func init() {
	skillEndFrames = frames.InitAbilSlice(30)
}

const skillHitmark = 10

func (c *char) skillActivate(p map[string]int) action.ActionInfo {
	// TODO: Hitlag?
	c.AddStatus(skillKey, -1, true)

	// Add 10 seconds worth of skydwellerPoints (1 point = 6 frames)
	c.skydwellerPoints = 100
	c.maxSkydwellerPoints = 100

	c.Core.Tasks.Add(c.depleteSkydwellerPoints(), 6)

	// Initial Skill Damage
	// TODO: Does that even need to be a task?
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Hanega: Song of the Wind",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}

		// TODO: check radius
		radius := 2.92

		// TODO: Check snapshot moment
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), radius), 0, 0)
		c.SetCD(action.ActionSkill, 10*60)

	}, skillHitmark)

	// Initial A1 Absorption test
	c.absorbCheckA1(c.Core.F)

	c.c1()

	// Return ActionInfo
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillEndFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillDeactivate(p map[string]int) action.ActionInfo {

	delay := c.skillEndRoutine()

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay },
		AnimationLength: delay,
		CanQueueAfter:   delay,
		State:           action.Idle,
	}
}

func (c *char) checkForSkillEnd() int {
	if c.StatusIsActive(skillKey) && c.skydwellerPoints <= 0 {
		return c.skillEndRoutine()
	}
	return 0
}

func (c *char) skillEndRoutine() int {
	c.DeleteStatus(skillKey)
	c.skydwellerPoints = 0
	c.a4Prob = 0.16
	c.SetCD(action.ActionSkill, 360)

	// Delete Ascension Buffs
	for _, e := range c.a1Buffs {
		switch e {
		case attributes.Pyro:
			c.DeleteStatMod("wanderer-a1-pyro")
		case attributes.Cryo:
			c.DeleteStatMod("wanderer-a1-cryo")
		case attributes.Electro:
			c.Core.Events.Unsubscribe(event.OnEnemyHit, "wanderer-a1-electro")
		}
	}

	// Delete c1 buff if active
	if c.StatusIsActive("wanderer-c1-atkspd") {
		c.DeleteStatus("wanderer-c1-atkspd")
	}

	// Delay due to falling

	c.Core.Log.NewEvent("adding delay due to falling", glog.LogCharacterEvent, c.Index)

	// TODO: Insert correct frames (especially for different circumstances)
	return 14
}

func (c *char) depleteSkydwellerPoints() func() {
	return func() {
		if c.StatusIsActive(skillKey) {
			c.skydwellerPoints -= 1
			c.Core.Tasks.Add(c.depleteSkydwellerPoints(), 6)
		}
	}
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if !c.StatusIsActive(skillKey) {
		return c.skillActivate(p)
	} else {
		return c.skillDeactivate(p)
	}
}
