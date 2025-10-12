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
	skillHitmark          = 26
	particleICDKey        = "flins-particle-icd"
	skillKey              = "manifest-flame"
	spearStormCDKey       = "spearstorm-cd"
	thunderousSymphonyKey = "thunderous-symphony"
)

func init() {
	skillFrames = frames.InitAbilSlice(32)
	skillFrames[action.ActionSwap] = 31

	spearStormFrames = frames.InitAbilSlice(32)
	spearStormFrames[action.ActionSwap] = 31
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatModIsActive(skillKey) {
		return c.spearStorm()
	}

	c.AddStatus(skillKey, 10*60, true)
	c.SetCDWithDelay(action.ActionSkill, 16*60, 1)

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
	c.Core.QueueAttack(ai, ap, 10, 10)
	c.AddStatus(thunderousSymphonyKey, c.c1SkillCD(), true)

	c.c2OnSkill()
	return action.Info{
		Frames:          func(next action.Action) int { return spearStormFrames[next] },
		AnimationLength: spearStormFrames[action.InvalidAction],
		CanQueueAfter:   spearStormFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}
