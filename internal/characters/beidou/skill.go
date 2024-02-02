package beidou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames       []int
	skillHitlagStages = []float64{.09, .09, .15}
	skillRadius       = []float64{6, 7, 8}
)

const (
	skillHitmark   = 23
	particleICDKey = "beidou-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(45)
	skillFrames[action.ActionAttack] = 44
	skillFrames[action.ActionDash] = 24
	skillFrames[action.ActionJump] = 24
	skillFrames[action.ActionSwap] = 44
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// 0 for base dmg, 1 for 1x bonus, 2 for max bonus
	counter := p["counter"]
	if counter >= 2 {
		counter = 2
		c.a4()
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Tidecaller (E)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           float64(100 * (counter + 1)),
		Element:            attributes.Electro,
		Durability:         50,
		Mult:               skillbase[c.TalentLvlSkill()] + skillbonus[c.TalentLvlSkill()]*float64(counter),
		HitlagFactor:       0.01,
		HitlagHaltFrames:   skillHitlagStages[counter] * 60,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillRadius[counter]),
		skillHitmark,
		skillHitmark,
		c.makeParticleCB(counter),
	)

	if counter > 0 {
		// add shield
		c.Core.Player.Shields.Add(&shield.Tmpl{
			ActorIndex: c.Index,
			Target:     c.Index,
			Src:        c.Core.F,
			ShieldType: shield.BeidouThunderShield,
			Name:       "Beidou Skill",
			HP:         shieldPer[c.TalentLvlSkill()]*c.MaxHP() + shieldBase[c.TalentLvlSkill()],
			Ele:        attributes.Electro,
			Expires:    c.Core.F + skillHitmark, // last until hitmark
		})
	}

	c.SetCDWithDelay(action.ActionSkill, 450, 4)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) makeParticleCB(counter int) combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.AddStatus(particleICDKey, 0.4*60, true)

		// 2 if no hit, 3 if 1 hit, 4 if perfect
		c.Core.QueueParticle(c.Base.Key.String(), 2+float64(counter), attributes.Electro, c.ParticleDelay)
	}
}
