package arlecchino

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	spikeHitmark      = 17
	finalHitmark      = 38
	particleICDKey    = "arlecchino-particle-icd"
	directiveLimitKey = "directive-limit"
	directiveKey      = "directive"
	directiveSrcKey   = "directive-src"
)

func init() {
	skillFrames = frames.InitAbilSlice(77)
	skillFrames[action.ActionAttack] = 70
	skillFrames[action.ActionCharge] = 65
	skillFrames[action.ActionBurst] = 70
	skillFrames[action.ActionDash] = 72
	skillFrames[action.ActionJump] = 73
	skillFrames[action.ActionSwap] = 60
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
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 0.5)
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

	skillCleaveArea := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 4.5}, 6)
	c.Core.QueueAttack(ai, skillCleaveArea, finalHitmark, finalHitmark, c.particleCB, c.bloodDebtDirective)
	c.QueueCharTask(c.debtLimit, finalHitmark+1)

	c.SetCDWithDelay(action.ActionSkill, 30*60, 16)
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

func (c *char) absorbDirectives() {
	area := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 3.0}, 6.5)
	enemies := c.Core.Combat.EnemiesWithinArea(area, nil)
	for _, e := range enemies {
		if !e.StatusIsActive(directiveKey) {
			continue
		}

		level := e.GetTag(directiveKey)

		newDebt := directiveScaling[level] * c.MaxHP()
		if c.StatusIsActive(directiveLimitKey) {
			newDebt = min(c.skillDebtMax-c.skillDebt, newDebt)
		}

		if newDebt > 0 {
			c.skillDebt += newDebt
			c.ModifyHPDebtByAmount(newDebt)
		}
		e.RemoveTag(directiveKey)
		e.RemoveTag(directiveSrcKey)
		e.DeleteStatus(directiveKey)

		c.c4OnAbsorb()
		if level >= 2 {
			c.c2OnAbsorbDue()
		}
	}
}
