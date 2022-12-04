package raiden

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

/**
The Raiden Shogun unveils a shard of her Euthymia, dealing Electro DMG to nearby opponents, and granting nearby party members the Eye of Stormy Judgment.
Eye of Stormy Judgment
**/

var skillFrames []int

const (
	skillHitmark = 51
	skillKey     = "raiden-e"
)

func init() {
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 17
	skillFrames[action.ActionSwap] = 36
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Eye of Stormy Judgement",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
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
		// start eye at hitmark only, eye dmg shouldn't be able to proc before the eye spawns
		c.Core.Tasks.Add(func() {
			this.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(skillKey, 1500),
				Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != combat.AttackTagElementalBurst {
						return nil, false
					}

					m[attributes.DmgP] = mult * this.EnergyMax
					return m, true
				},
			})
		}, skillHitmark)
	}

	c.SetCDWithDelay(action.ActionSkill, 600, 6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
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
		//ignore if eye on icd
		if c.eyeICD > c.Core.F {
			return false
		}
		//ignore if eye status not active on char that's doing dmg
		if !c.Core.Player.ByIndex(ae.Info.ActorIndex).StatusIsActive(skillKey) {
			return false
		}
		//ignore EC and hydro swirl damage
		if ae.Info.AttackTag == combat.AttackTagECDamage || ae.Info.AttackTag == combat.AttackTagSwirlHydro {
			return false
		}
		//ignore self dmg
		if ae.Info.Abil == "Eye of Stormy Judgement" {
			return false
		}
		//ignore 0 damage
		if dmg == 0 {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.Core.QueueParticle("raiden", 1, attributes.Electro, c.ParticleDelay)
		}

		//hit mark 857, eye land 862
		//electro appears to be applied right away
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Eye of Stormy Judgement (Strike)",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeSlash,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillTick[c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 2 && c.StatusIsActive(burstKey) {
			ai.IgnoreDefPercent = 0.6
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 4), 5, 5)

		c.eyeICD = c.Core.F + 54 //0.9 sec icd
		return false
	}, "raiden-eye")

}
