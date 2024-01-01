package freminet

import (
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
	skillPressureHitmarks = []int{42, 37}
)

const (
	skillThrustHitmark   = 29
	particleICDKeyThrust = "freminet-particle-icd-thrust"
	particleICDKeyLv4    = "freminet-particle-icd-lv4"
	persTimeKey          = "freminet-pers-time"
	pressureBaseName     = "Pressurized Floe: Shattering Pressure"

	skillAlignedICDKey = "freminet-aligned-icd"
	skillAlignedICD    = 9 * 60
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

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(persTimeKey) {
		return c.detonateSkill()
	}

	c.skillStacks = 0
	c.AddStatus(persTimeKey, 10*60, true)

	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Pressurized Floe: Upward Thrust",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagElementalArt,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         75,
		Element:          attributes.Cryo,
		Durability:       25,
		Mult:             skillThrust[c.TalentLvlSkill()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.08 * 60,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 2.5),
		0,
		skillThrustHitmark,
		c.particleCBThrust,
	)
	c.skillAligned()

	cd := 600
	if c.StatusIsActive(burstKey) {
		cd = 3 * 60
	}

	c.SetCDWithDelay(action.ActionSkill, cd, 35)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillThrustFrames),
		AnimationLength: skillThrustFrames[action.InvalidAction],
		CanQueueAfter:   skillThrustFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillAligned() {
	if c.StatusIsActive(skillAlignedICDKey) {
		return
	}
	c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

	aiSpiritbreath := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Pressurized Floe: Spiritbreath Thorn",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 0,
		Mult:       skillBreath[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		aiSpiritbreath,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 2.5),
		62,
		62,
	)
}

func (c *char) detonateSkill() (action.Info, error) {
	pressureFrameIndex := 0
	if c.skillStacks == 4 {
		pressureFrameIndex = 1
	}

	if skillPressureCryo[c.skillStacks][c.TalentLvlSkill()] > 0 {
		poiseDMG := 150.0
		if c.skillStacks > 0 {
			poiseDMG = 70.0
		}
		ai := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             pressureBaseName + " (Cryo)",
			AttackTag:        attacks.AttackTagElementalArt,
			ICDTag:           attacks.ICDTagElementalArt,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeBlunt,
			PoiseDMG:         poiseDMG,
			Element:          attributes.Cryo,
			Durability:       25,
			Mult:             skillPressureCryo[c.skillStacks][c.TalentLvlSkill()],
			HitlagFactor:     0.01,
			HitlagHaltFrames: 0.10 * 60,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 2.5),
			0,
			skillPressureHitmarks[pressureFrameIndex],
		)
	}

	if skillPressurePhys[c.skillStacks][c.TalentLvlSkill()] > 0 {
		poiseDMG := 150.0
		if c.skillStacks < 4 {
			poiseDMG = 70.0
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       pressureBaseName + " (Physical)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   poiseDMG,
			Element:    attributes.Physical,
			Durability: 25,
			Mult:       skillPressurePhys[c.skillStacks][c.TalentLvlSkill()],
		}
		if c.skillStacks == 4 {
			ai.HitlagFactor = 0.01
			ai.HitlagHaltFrames = 0.09 * 60
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 2.5)
		var particleCB combat.AttackCBFunc
		if c.skillStacks == 4 {
			ap = combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: 0.5, Y: 0.5}, 3, 7)
			particleCB = c.particleCBLv4
		}

		c.Core.QueueAttack(
			ai,
			ap,
			0,
			skillPressureHitmarks[pressureFrameIndex],
			particleCB,
		)
	}

	// A1
	c.a1()

	// C2
	c.c2()

	c.DeleteStatus(persTimeKey)
	c.skillStacks = 0

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressureFrames[pressureFrameIndex]),
		AnimationLength: skillPressureFrames[pressureFrameIndex][action.InvalidAction],
		CanQueueAfter:   skillPressureFrames[pressureFrameIndex][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCBThrust(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKeyThrust) {
		return
	}
	c.AddStatus(particleICDKeyThrust, 0.3*60, true)

	particles := 2.0
	if c.StatusIsActive(burstKey) {
		particles = 1
	}
	c.Core.QueueParticle(c.Base.Key.String(), particles, attributes.Cryo, c.ParticleDelay)
}

func (c *char) particleCBLv4(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKeyLv4) {
		return
	}
	c.AddStatus(particleICDKeyLv4, 0.3*60, true)

	particles := 2.0
	if c.StatusIsActive(burstKey) {
		particles = 1
	}
	c.Core.QueueParticle(c.Base.Key.String(), particles, attributes.Cryo, c.ParticleDelay)
}
