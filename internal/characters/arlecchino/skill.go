package arlecchino

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	spikeHitmark      = 14
	finalHitmark      = 37
	particleICDKey    = "arlecchino-particle-icd"
	directiveLimitKey = "directive-limit"
	directiveKey      = "directive"
	directiveSrcKey   = "directive-src"
)

func init() {
	skillFrames = frames.InitAbilSlice(37)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "All is Ash (Spike)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupArlecchinoElementalArt,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillSpike[c.TalentLvlSkill()],
	}
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3)
	c.Core.QueueAttack(ai, skillArea, spikeHitmark, spikeHitmark)

	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "All is Ash (Cleave)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		CanBeDefenseHalted: true,
		Mult:               skillFinal[c.TalentLvlSkill()],
	}

	skillArea = combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3)
	c.Core.QueueAttack(ai, skillArea, finalHitmark, finalHitmark, c.particleCB, c.bloodDebtDirective)
	c.QueueCharTask(c.debtLimit, finalHitmark+1)

	c.SetCDWithDelay(action.ActionSkill, 30*60, 0)
	c.QueueCharTask(c.c6skill, finalHitmark)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}

func (c *char) bloodDebtDirective(a combat.AttackCB) {
	// TODO: is this a redundant check?
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}

	trg, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	trg.AddStatus(directiveKey, 30*60, true)
	trg.SetTag(directiveSrcKey, c.Core.F)
	trg.SetTag(directiveKey, c.initialDirectiveLevel)
	trg.QueueEnemyTask(c.directiveTickFunc(c.Core.F, 2, trg), 5*60)
	c.a1Upgrade(trg, c.Core.F)
}

func (c *char) directiveTickFunc(src, count int, trg *enemy.Enemy) func() {
	return func() {
		// do nothing if source changed
		if trg.Tags[directiveSrcKey] != src {
			return
		}
		if !trg.StatusIsActive(directiveKey) {
			return
		}
		c.Core.Log.NewEvent("Blood Debt Directive checking for tick", glog.LogCharacterEvent, c.Index).
			Write("src", src)

		// queue up one damage instance
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blood Debt Directive",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupArlecchinoElementalArt,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillSigil[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 0)

		if count-1 > 0 {
			// queue up next instance
			trg.QueueEnemyTask(c.directiveTickFunc(src, count-1, trg), 5*60)
		}
	}
}

func (c *char) debtLimit() {
	c.AddStatus(directiveLimitKey, 35*60, true)
	c.skillDebtMax = 1.45 * c.MaxHP()
	c.skillDebt = 0
}
