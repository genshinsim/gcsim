package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"math"
)

const (
	skillKey = "windfavored-state"
)

var skillFramesNormal []int

func init() {
	skillFramesNormal = frames.InitAbilSlice(28)

}

const skillHitmark = 2

func (c *char) skillActivate(p map[string]int) action.ActionInfo {
	c.AddStatus(skillKey, 20*60, true)
	c.Core.Player.SwapCD = math.MaxInt16

	// Add 10 seconds worth of skydwellerPoints (1 point = 6 frames)
	c.skydwellerPoints = 100
	c.maxSkydwellerPoints = 100
	c.c6Count = 0

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

	}, skillHitmark)

	// Initial A1 Absorption test
	c.a1ValidBuffs = []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}
	c.absorbCheckA1(c.Core.F)()

	c.c1()
	c.c6()

	// Return ActionInfo
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFramesNormal),
		AnimationLength: skillFramesNormal[action.InvalidAction],
		CanQueueAfter:   skillFramesNormal[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillDeactivate(p map[string]int) action.ActionInfo {

	delay := c.skillEndRoutine()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			if next == action.ActionLowPlunge {
				return 7
			} else {
				return delay
			}
		},
		AnimationLength: delay,
		CanQueueAfter:   7,
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
	//print("Starting skill end routine")
	c.DeleteStatus(skillKey)
	c.Core.Player.SwapCD = 26

	if c.StatusIsActive(a4Key) {
		c.DeleteStatus(a4Key)
	}

	c.skydwellerPoints = 0
	c.a4Prob = 0.16
	c.SetCD(action.ActionSkill, 360)

	// Delete Ascension Buffs
	c.DeleteStatMod("wanderer-a1-pyro")
	c.DeleteStatMod("wanderer-a1-cryo")
	c.Core.Events.Unsubscribe(event.OnEnemyHit, "wanderer-a1-electro")

	// Delete c1 buff if active
	if c.StatusIsActive("wanderer-c1-atkspd") {
		c.DeleteStatus("wanderer-c1-atkspd")
	}

	// Delete c6 buff if active
	c.Core.Events.Unsubscribe(event.OnEnemyHit, "wanderer-c6")

	// Delay due to falling
	c.Core.Log.NewEvent("adding delay due to falling", glog.LogCharacterEvent, c.Index)

	// Shorter delay for plunging is hard coded in the plunge action
	return 26
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
