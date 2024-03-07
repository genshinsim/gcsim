package wanderer

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	skillKey           = "windfavored-state"
	particleICDKey     = "wanderer-particle-icd"
	plungeAvailableKey = "wanderer-plunge-available"
)

var skillFramesNormal []int

func init() {
	skillFramesNormal = frames.InitAbilSlice(28)
}

const skillHitmark = 2

func (c *char) skillActivate() action.Info {
	c.AddStatus(skillKey, 20*60, true)
	c.Core.Player.SwapCD = math.MaxInt16

	// Add 10 seconds worth of skydwellerPoints (1 point = 6 frames)
	c.skydwellerPoints = 100
	c.maxSkydwellerPoints = 100
	c.c6Count = 0

	c.Core.Tasks.Add(c.depleteSkydwellerPoints, 6)

	// Initial Skill Damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Hanega: Song of the Wind",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), skillHitmark, skillHitmark)

	// A1
	if c.Base.Ascension >= 1 {
		c.a1ValidBuffs = []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}
		c.absorbCheckA1()
	}

	c.c1()

	// Return ActionInfo
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFramesNormal),
		AnimationLength: skillFramesNormal[action.InvalidAction],
		CanQueueAfter:   skillFramesNormal[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillDeactivate() action.Info {
	delay := c.skillEndRoutine()

	return action.Info{
		Frames: func(next action.Action) int {
			if next == action.ActionLowPlunge {
				return 7
			}
			return delay
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
	// print("Starting skill end routine")
	c.DeleteStatus(skillKey)
	c.Core.Player.SwapCD = 26

	if c.StatusIsActive(a4Key) {
		c.DeleteStatus(a4Key)
	}
	if c.StatusIsActive(a4Prevent) {
		c.DeleteStatus(a4Prevent)
	}

	c.skydwellerPoints = 0
	c.a4Prob = 0.16
	c.SetCD(action.ActionSkill, 360)

	// Delete Ascension Buffs
	c.DeleteStatMod(a1PyroKey)
	c.DeleteStatMod(a1CryoKey)
	c.DeleteStatus(a1ElectroKey)

	// Delete c1 buff if active
	if c.StatusIsActive("wanderer-c1-atkspd") {
		c.DeleteStatus("wanderer-c1-atkspd")
	}

	// Delay due to falling
	c.Core.Log.NewEvent("adding delay due to falling", glog.LogCharacterEvent, c.Index)

	c.AddStatus(plungeAvailableKey, 26, true)

	// Shorter delay for plunging is hard coded in the plunge action
	return 26
}

func (c *char) depleteSkydwellerPoints() {
	if c.StatusIsActive(skillKey) {
		c.skydwellerPoints -= 1
		c.Core.Tasks.Add(c.depleteSkydwellerPoints, 6)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if !c.StatusIsActive(skillKey) {
		return c.skillActivate(), nil
	}
	return c.skillDeactivate(), nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if !c.StatusIsActive(skillKey) {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Anemo, c.ParticleDelay)
}
