package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	lumidouceArkheCD = "lumidouce-arkhe-cd"
	particleICDKey   = "skill-particle-icd"

	lumidouceSummonHitmark = 37
	lumidouceSpawn         = 16
	lumidouceArkheHitmark  = 62

	particleICD = 2.5 * 60
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(37) // E -> Walk
	skillFrames[action.ActionAttack] = 28
	skillFrames[action.ActionBurst] = 24
	skillFrames[action.ActionDash] = 23
	skillFrames[action.ActionJump] = 24
	skillFrames[action.ActionSwap] = 34
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	var ok bool
	c.caseTravel, ok = p["travel"]
	if !ok {
		c.caseTravel = lumidouceAttackTravel
	}

	player := c.Core.Combat.Player()

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lumidouce Case (Summon)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(player, geometry.Point{Y: 2.6}, 4.5),
		lumidouceSummonHitmark,
		lumidouceSummonHitmark,
		c.c2,
	)

	if c.Tag(lumidouceLevel) != 3 { // spawn if no burst
		c.QueueCharTask(func() {
			c.spawnLumidouceCase(1, geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 2.6}, player.Direction()))
			c.c6()
		}, lumidouceSpawn)
	}
	c.arkheAttack()
	c.SetCD(action.ActionSkill, int(skillCD[c.TalentLvlSkill()]*60))

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) arkheAttack() {
	if c.StatusIsActive(lumidouceArkheCD) {
		return
	}
	c.AddStatus(lumidouceArkheCD, int(skillArkeCD[c.TalentLvlSkill()]*60), true)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spiritbreath Thorn",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillArkeCD[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 2.6}, 4.5),
		lumidouceArkheHitmark,
		lumidouceArkheHitmark,
		c.c2,
	)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Dendro, c.ParticleDelay)
}
