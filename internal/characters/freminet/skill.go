package freminet

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillThrustFrames     []int
	skillPressureFrames   [][]int
	skillPressureHitmarks []int = []int{42, 37}
)

const (
	skillThrustHitmark = 29
	particleICDKey     = "freminet-particle-icd"
	persTimeKey        = "freminet-pers-time"
	pressureBaseName   = "Pressurized Floe: Shattering Pressure"
)

func init() {
	skillThrustFrames = frames.InitAbilSlice(46)
	skillThrustFrames[action.ActionAttack] = 42
	skillThrustFrames[action.ActionSkill] = 31
	skillThrustFrames[action.ActionBurst] = 42
	skillThrustFrames[action.ActionDash] = 32
	skillThrustFrames[action.ActionJump] = 31

	skillPressureFrames = make([][]int, 2)
	// < Lv.4
	skillPressureFrames[0] = frames.InitAbilSlice(55)
	skillPressureFrames[0][action.ActionAttack] = 53
	skillPressureFrames[0][action.ActionSkill] = 47
	skillPressureFrames[0][action.ActionBurst] = 47
	skillPressureFrames[0][action.ActionDash] = 47
	skillPressureFrames[0][action.ActionJump] = 47
	skillPressureFrames[0][action.ActionSwap] = 51
	// == Lv.4
	skillPressureFrames[1] = frames.InitAbilSlice(59)
	skillPressureFrames[1][action.ActionAttack] = 53
	skillPressureFrames[1][action.ActionSkill] = 42
	skillPressureFrames[1][action.ActionBurst] = 42
	skillPressureFrames[1][action.ActionDash] = 43
	skillPressureFrames[1][action.ActionJump] = 41
	skillPressureFrames[1][action.ActionSwap] = 51
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if !c.StatusIsActive(persTimeKey) {
		c.skillStacks = 0

		c.AddStatus(persTimeKey, 600, true)
		c.persID = c.Core.F

		// TODO: Freminet; Update Info
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Pressurized Floe: Upward Thrust",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Cryo,
			Durability:         50,
			Mult:               skillThrust[c.TalentLvlSkill()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: false,
		}

		// TODO: Freminet; Insert Hitbox
		skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.5}, 8)

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(skillArea.Shape.Pos(), nil, 2.5),
			0,
			skillThrustHitmark,
			c.particleCB,
		)

		// TODO: Freminet; Update Info
		aiSpiritbreath := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Pressurized Floe: Spiritbreath Thorn",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Cryo,
			Durability:         50,
			Mult:               skillBreath[c.TalentLvlSkill()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: false,
		}

		currentID := c.Core.F

		c.Core.Tasks.Add(func() {
			if c.StatusIsActive(persTimeKey) && currentID == c.persID {
				c.Core.QueueAttack(aiSpiritbreath, combat.NewCircleHitOnTarget(skillArea.Shape.Pos(), nil, 2.5), 0, 0)
			}
		}, 62)

		cd := 600

		if c.StatusIsActive(burstKey) {
			cd = 3 * 60
		}

		c.SetCDWithDelay(action.ActionSkill, cd, 35)

		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(skillThrustFrames),
			AnimationLength: skillThrustFrames[action.InvalidAction],
			CanQueueAfter:   skillThrustFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}
	} else {
		// Manual Cancel
		actionInfo := c.detonateSkill()

		return actionInfo
	}
}

func (c *char) detonateSkill() action.ActionInfo {

	// TODO: Freminet; Insert Hitbox
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.5}, 8)

	pressureFrameIndex := 0
	if c.skillStacks == 4 {
		pressureFrameIndex = 1
	}

	if skillPressureCryo[c.skillStacks][c.TalentLvlSkill()] > 0 {
		// TODO: Freminet; Update Info
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               pressureBaseName + " (Cryo)",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Cryo,
			Durability:         50,
			Mult:               skillPressureCryo[c.skillStacks][c.TalentLvlSkill()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: false,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(skillArea.Shape.Pos(), nil, 2.5),
			0,
			skillPressureHitmarks[pressureFrameIndex],
			c.particleCB,
		)
	}

	if skillPressurePhys[c.skillStacks][c.TalentLvlSkill()] > 0 {
		// TODO: Freminet; Update Info
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               pressureBaseName + " (Physical)",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Physical,
			Durability:         50,
			Mult:               skillPressurePhys[c.skillStacks][c.TalentLvlSkill()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: false,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(skillArea.Shape.Pos(), nil, 2.5),
			0,
			skillPressureHitmarks[pressureFrameIndex],
			c.particleCB,
		)
	}

	// A1
	if c.Base.Ascension >= 1 && c.skillStacks < 4 {
		c.ReduceActionCooldown(action.ActionSkill, 60)
	}

	// C2
	if c.Base.Cons >= 2 {
		if c.skillStacks < 4 {
			c.AddEnergy(c1Key, 2)
		} else {
			c.AddEnergy(c1Key, 3)
		}
	}

	c.DeleteStatus(persTimeKey)
	c.skillStacks = 0

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressureFrames[pressureFrameIndex]),
		AnimationLength: skillPressureFrames[pressureFrameIndex][action.InvalidAction],
		CanQueueAfter:   skillPressureFrames[pressureFrameIndex][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	// TODO: Freminet; Check Particle amount and ICD
	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
}

func (c *char) persTimeCB() func(combat.AttackCB) {
	done := false
	return func(a combat.AttackCB) {
		if done {
			return
		}

		frostMod := skillAddNA[c.TalentLvlSkill()]

		if c.StatusIsActive(burstKey) {
			frostMod = frostMod * 2
		}

		// TODO: Freminet; Update Info
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Pressurized Floe: Pers Time Frost",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       frostMod,
		}

		// TODO: Freminet; Update Hitbox
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)

		c.Core.QueueAttack(ai, ap, 1, 1)

		if c.skillStacks < 4 {
			if c.StatusIsActive(burstKey) {
				c.skillStacks = int(math.Min(float64(c.skillStacks+2), 4))
			} else {
				c.skillStacks++
			}
		}

		done = true
	}

}
