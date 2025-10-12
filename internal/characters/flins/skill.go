package flins

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames      []int
	spearStormFrames []int
)

const (
	skillHitmark          = 19
	particleICDKey        = "flins-particle-icd"
	skillKey              = "manifest-flame"
	spearStormHitmark     = 23
	spearStormCDKey       = "spearstorm-cd"
	thunderousSymphonyKey = "thunderous-symphony"
)

func init() {
	skillFrames = frames.InitAbilSlice(44)
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionSkill] = 22
	skillFrames[action.ActionBurst] = 19
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 18
	skillFrames[action.ActionSwap] = 17

	spearStormFrames = frames.InitAbilSlice(42)
	spearStormFrames[action.ActionAttack] = 28
	spearStormFrames[action.ActionBurst] = 28
	spearStormFrames[action.ActionDash] = 26
	spearStormFrames[action.ActionJump] = 26
	spearStormFrames[action.ActionWalk] = 32
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatModIsActive(skillKey) {
		return c.spearStorm()
	}

	c.AddStatus(skillKey, 10*60+skillHitmark, true)
	c.SetCDWithDelay(action.ActionSkill, 16*60, skillHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) spearStorm() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Northland Spearstorm",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillDmg[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 5, 60)
	c.Core.QueueAttack(ai, ap, spearStormHitmark, spearStormHitmark)
	c.AddStatus(spearStormCDKey, c.c1SkillCD(), false)
	c.AddStatus(thunderousSymphonyKey, 6*60, true)
	c.c2OnSkill()
	return action.Info{
		Frames:          func(next action.Action) int { return spearStormFrames[next] },
		AnimationLength: spearStormFrames[action.InvalidAction],
		CanQueueAfter:   spearStormFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}
