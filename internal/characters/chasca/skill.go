package chasca

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillHitmarks   = 6
	samplerInterval = 0.3 * 60

	skilRecastCD     = "xilonen-e-recast-cd"
	skillMaxDurKey   = "xilonen-e-limit"
	particleICDKey   = "xilonen-particle-icd"
	samplerShredKey  = "xilonen-e-shred"
	activeSamplerKey = "xilonen-samplers-activated"
)

func init() {
	skillFrames = frames.InitAbilSlice(20)
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionDash] = 15
	skillFrames[action.ActionJump] = 15
	skillFrames[action.ActionSwap] = 19
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() < 0.001 {
		c.exitNightsoul()
	}
}

func (c *char) enterNightsoul() {
	c.nightsoulState.EnterBlessing(80)
	c.nightsoulSrc = c.Core.F
	c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 6)
	c.NormalHitNum = 1
	c.NormalCounter = 0
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	c.SetCD(action.ActionSkill, 6.5*60)
	c.NormalHitNum = normalHitNum
	c.NormalCounter = 0
	c.DeleteStatus(particleICDKey)
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		c.reduceNightsoulPoints(0.5)
		// reduce 0.5 point per 6, which is 5 per second
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(src), 6)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return action.Info{}, errors.New("plunging attack from chasca skill is not implemented")
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spirit Reins, Shadow Hunt",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypePierce,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           skillResonance[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		5.5,
	)
	c.Core.QueueAttack(ai, ap, skillHitmarks, skillHitmarks)
	c.enterNightsoul()

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
	c.AddStatus(particleICDKey, 9999*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Anemo, c.ParticleDelay)
}
