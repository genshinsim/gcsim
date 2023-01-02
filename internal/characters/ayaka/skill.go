package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 33

func init() {
	skillFrames = frames.InitAbilSlice(49)
	skillFrames[action.ActionBurst] = 48
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 48
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Hyouka",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	//a1 increase normal + ca dmg by 30% for 6s
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.3
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("ayaka-a1", 360),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == combat.AttackTagNormal || atk.Info.AttackTag == combat.AttackTagExtra
		},
	})

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5), 0, skillHitmark)

	// 4 or 5, 1:1 ratio
	var count float64 = 4
	if c.Core.Rand.Float64() < 0.5 {
		count = 5
	}
	c.Core.QueueParticle("ayaka", count, attributes.Cryo, skillHitmark+c.ParticleDelay)

	c.SetCD(action.ActionSkill, 600)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
