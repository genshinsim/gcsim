package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// TODO: yoimiya based frames
var skillFrames []int

const (
	skillKey       = "wriothesleyskill"
	particleICDKey = "wriothesley-particle-icd"
	skillStart     = 11
)

func init() {
	skillFrames = frames.InitAbilSlice(34)
	skillFrames[action.ActionAttack] = 22
	skillFrames[action.ActionAim] = 22 // uses attack frames
	skillFrames[action.ActionBurst] = 23
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 31
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.Base.Ascension >= 4 {
		c.a4Stack = 0
	}
	if c.Base.Cons >= 1 {
		c.c1Proc = false
	}

	c.AddStatus(skillKey, skillStart+10*60, true) // activate for 10
	c.SetCDWithDelay(action.ActionSkill, 16*60, 11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
}

func (c *char) chillingPenalty(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Chilling Penalty",
		Amount:     0.045 * c.MaxHP(),
	})
}
