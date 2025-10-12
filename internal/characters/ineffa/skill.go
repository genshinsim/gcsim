package ineffa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	skillHitmark   = 26
	particleICDKey = "ineffa-particle-icd"
	birgittaKey    = "birgitta"
)

func init() {
	skillFrames = frames.InitAbilSlice(32)
	skillFrames[action.ActionSwap] = 31
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Cleaning Mode: Carrier Frequency",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillInitial[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5), skillHitmark, skillHitmark)
	c.QueueCharTask(c.addShield, 4)
	c.QueueCharTask(c.summonBirgitta, skillHitmark+42)

	c.SetCDWithDelay(action.ActionSkill, 16*60, 1)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) summonBirgitta() {
	c.birgittaSrc = c.Core.F
	c.AddStatus(birgittaKey, 20*60, false)
	c.Core.Tasks.Add(c.birgittaDischarge(c.birgittaSrc), 42)
}

func (c *char) birgittaDischarge(src int) func() {
	return func() {
		// src changed, cancel these ticks
		if c.birgittaSrc != src {
			return
		}

		if !c.StatusIsActive(birgittaKey) {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Birgitta",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillDoT[c.TalentLvlSkill()],
		}

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), 0, 0, c.baseParticleCB)
		c.a1OnDischarge()
		c.Core.Tasks.Add(c.birgittaDischarge(src), 119)
	}
}

func (c *char) baseParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 0.3*60, true)

	if c.Core.Rand.Float64() < 0.66 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
	}
}
