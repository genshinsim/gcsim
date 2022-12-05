package rosaria

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 24

func init() {
	skillFrames = frames.InitAbilSlice(51)
	skillFrames[action.ActionDash] = 38
	skillFrames[action.ActionJump] = 40
	skillFrames[action.ActionSwap] = 50
}

// Skill attack damage queue generator
// Includes optional argument "nobehind" for whether Rosaria appears behind her opponent or not (for her A1).
// Default behavior is to appear behind enemy - set "nobehind=1" to diasble A1 proc
func (c *char) Skill(p map[string]int) action.ActionInfo {
	// No ICD to the 2 hits
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Ravaging Confession (Hit 1)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skill[0][c.TalentLvlSkill()],
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}
	var c4cb combat.AttackCBFunc
	if c.Base.Cons >= 4 {
		c.c4completed = false
		c4cb = c.c4
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.24), skillHitmark, skillHitmark, c4cb)

	// A1 activation
	// When Rosaria strikes an opponent from behind using Ravaging Confession, Rosaria's CRIT RATE increases by 12% for 5s.
	// We always assume that it procs on hit 1 to simplify
	//TODO: does this need to change if we add player position?
	if p["nobehind"] != 1 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.12
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("rosaria-a1", 300+skillHitmark),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		c.Core.Log.NewEvent("Rosaria A1 activation", glog.LogCharacterEvent, c.Index).
			Write("ends_on", c.Core.F+300+skillHitmark)
	}

	// Rosaria E is dynamic, so requires a second snapshot
	//TODO: check snapshot timing here
	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Ravaging Confession (Hit 2)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skill[1][c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.QueueCharTask(func() {
		//second hit is 14 frames after the first (if we exclude hitlag)
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.8), 0, 0)
		// Particles are emitted after the second hit lands
		c.Core.QueueParticle("rosaria", 3, attributes.Cryo, c.ParticleDelay)
	}, skillHitmark+14)

	c.SetCDWithDelay(action.ActionSkill, 360, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
