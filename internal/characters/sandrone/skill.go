package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const skillParticleICDKey = "sandrone-particle-icd"

func init() {
	skillFrames = frames.InitAbilSlice(600) // need to full exit E before actual walk occurs
	skillFrames[action.ActionAttack] = 10
	skillFrames[action.ActionCharge] = 46
	skillFrames[action.ActionSkill] = 46 // TODO: Get actual value, currently assuming same as CA which is wrong
	skillFrames[action.ActionBurst] = 6
	skillFrames[action.ActionDash] = 26
	skillFrames[action.ActionSwap] = 27
}

// When Sandrone starts hovering, if there are opponents nearby, she will also summon a Prismatic
// Resonance Cannon and fire 2 Prism Shots at any opponents who have yet to cool down, dealing Cryo
// DMG.
// Radiance: Stellar Glimmer: The second Prism Shot will instead deal Cryo DMG that is considered
// the corresponding Stellar Glimmer reaction DMG.
// When on the Tea Party Tactical Assault Hovermech, Sandrone will also elegantly and quietly repair
// Fagio, causing Decoding Power to swiftly decrease.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	bonusMult := 1 + c.a1OnSkillBonusMult()
	c.reduceDecode(100)

	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Prism Shot",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)
		c.Core.QueueAttack(ai, ap, 5, 5+travel, c.particleCB)

		c.toStellar(&ai, skillSSC[c.TalentLvlSkill()], skillSSw[c.TalentLvlSkill()])
		ai.Mult *= bonusMult
		c.Core.QueueAttack(ai, ap, 5, 5+travel, c.particleCB)
	}, 19)

	c.SetCD(action.ActionSkill, 4*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(skillParticleICDKey) {
		return
	}
	c.AddStatus(skillParticleICDKey, 2.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
}
