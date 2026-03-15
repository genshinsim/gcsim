package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const skillHitmark = 24

const particleICDKey = "nefer-particle-icd"

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 26
	skillFrames[action.ActionCharge] = 29
	skillFrames[action.ActionDash] = 38
	skillFrames[action.ActionJump] = 38
	skillFrames[action.ActionSwap] = 25
	skillFrames[action.ActionBurst] = 38
	skillFrames[action.ActionWalk] = 34
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Senet Strategy: Dance of a Thousand Nights",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[0][c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * skill[1][c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: 1}, 3), skillHitmark, skillHitmark, c.skillParticleCB)
	c.AddStatus(shadowDanceKey, 9*60, true)
	c.phantasmCharges = phantasmChargesPerSkill
	c.startSeedWindow()
	c.SetCDWithDelay(action.ActionSkill, 9*60, 22)

	if c.Base.Cons >= 2 && c.ascendantGleam {
		c.addVeilStacks(2)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 0.2*60, false)

	count := 2.0
	if c.Core.Rand.Float64() < 0.66 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
}
