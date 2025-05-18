package raiden

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

/**
The Raiden Shogun unveils a shard of her Euthymia, dealing Electro DMG to nearby opponents, and granting nearby party members the Eye of Stormy Judgment.
Eye of Stormy Judgment
**/

var skillFrames []int

const (
	skillHitmark   = 51
	skillKey       = "raiden-e"
	particleICDKey = "raiden-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 17
	skillFrames[action.ActionSwap] = 36
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Eye of Stormy Judgement",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		skillHitmark,
		skillHitmark,
	)

	// Add pre-damage mod
	mult := skillBurstBonus[c.TalentLvlSkill()]
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		this := char
		// starts 1s after cd delay
		c.Core.Tasks.Add(func() {
			this.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(skillKey, 1500),
				Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
						return nil, false
					}

					m[attributes.DmgP] = mult * this.EnergyMax
					return m, true
				},
			})
		}, 6+60)
	}

	c.SetCDWithDelay(action.ActionSkill, 600, 6)

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
	c.AddStatus(particleICDKey, 0.8*60, true)
	if c.Core.Rand.Float64() < 0.5 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
	}
}

/*
*
When characters with this buff attack and hit opponents, the Eye will unleash a coordinated attack, dealing AoE Electro DMG at the opponent's position.
The Eye can initiate one coordinated attack every 0.9s per party.
*
*/
func (c *char) eyeOnDamage() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		// ignore if eye on icd
		if c.eyeICD > c.Core.F {
			return false
		}
		// ignore if eye status not active on char that's doing dmg
		if !c.Core.Player.ByIndex(ae.Info.ActorIndex).StatusIsActive(skillKey) {
			return false
		}
		// ignore EC, hydro swirl, and burning damage
		// this clause is here since these damage types are sourced to the target rather than character
		if ae.Info.AttackTag == attacks.AttackTagECDamage || ae.Info.AttackTag == attacks.AttackTagBurningDamage ||
			ae.Info.AttackTag == attacks.AttackTagSwirlHydro {
			return false
		}
		// ignore self dmg
		if ae.Info.ActorIndex == c.Index &&
			ae.Info.AttackTag == attacks.AttackTagElementalArt &&
			ae.Info.StrikeType == attacks.StrikeTypeSlash {
			return false
		}
		// ignore 0 damage
		if dmg == 0 {
			return false
		}

		// hit mark 857, eye land 862
		// electro appears to be applied right away
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Eye of Stormy Judgement (Strike)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillTick[c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 2 && c.StatusIsActive(BurstKey) {
			ai.IgnoreDefPercent = 0.6
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 4), 5, 5, c.particleCB)

		c.eyeICD = c.Core.F + 54 // 0.9 sec icd
		return false
	}, "raiden-eye")
}
