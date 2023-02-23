package tartaglia

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	skillMeleeFrames      []int
	skillMeleeWalkFrames  []int
	skillMeleeDashFrames  []int
	skillRangedFrames     []int
	skillRangedWalkFrames []int
	skillRangedDashFrames []int
)

const (
	skillHitmark     = 16
	skillWalkHitmark = 3
	skillDashHitmark = 3
)

func init() {
	// skill (melee) -> x
	skillMeleeFrames = frames.InitAbilSlice(18)
	skillMeleeFrames[action.ActionAttack] = 17
	skillMeleeFrames[action.ActionBurst] = 18
	skillMeleeFrames[action.ActionDash] = 17
	skillMeleeFrames[action.ActionJump] = 17
	skillMeleeFrames[action.ActionSwap] = 16

	// skill (melee, walk) -> x
	skillMeleeWalkFrames = frames.InitAbilSlice(24)
	skillMeleeWalkFrames[action.ActionAttack] = 5
	skillMeleeWalkFrames[action.ActionBurst] = 5
	skillMeleeWalkFrames[action.ActionDash] = 6
	skillMeleeWalkFrames[action.ActionJump] = skillWalkHitmark

	// skill (melee, dash) -> x
	skillMeleeDashFrames = frames.InitAbilSlice(23)
	skillMeleeDashFrames[action.ActionAttack] = 13
	skillMeleeDashFrames[action.ActionBurst] = 16
	skillMeleeDashFrames[action.ActionDash] = 22
	skillMeleeDashFrames[action.ActionJump] = skillDashHitmark

	// skill (ranged) -> x
	skillRangedFrames = frames.InitAbilSlice(39)
	skillRangedFrames[action.ActionAttack] = 19
	skillRangedFrames[action.ActionBurst] = 19
	skillRangedFrames[action.ActionDash] = 19
	skillRangedFrames[action.ActionJump] = 21

	// skill (ranged, walk) -> x
	skillRangedWalkFrames = frames.InitAbilSlice(24)
	skillRangedWalkFrames[action.ActionAttack] = 5
	skillRangedWalkFrames[action.ActionBurst] = 4
	skillRangedWalkFrames[action.ActionDash] = 5
	skillRangedWalkFrames[action.ActionJump] = 4

	// skill (ranged, dash) -> x
	skillRangedDashFrames = frames.InitAbilSlice(24)
	skillRangedDashFrames[action.ActionAttack] = 17
	skillRangedDashFrames[action.ActionBurst] = 17
	skillRangedDashFrames[action.ActionDash] = 22
	skillRangedDashFrames[action.ActionJump] = 3
}

// Cast: AoE strong hydro damage
// Melee Stance: infuse NA/CA to hydro damage
func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(MeleeKey) {
		cdDelay := 11
		switch c.Core.Player.CurrentState() {
		case action.WalkState,
			action.DashState:
			cdDelay = 0
		}
		c.onExitMeleeStance(cdDelay)
		c.ResetNormalCounter()
		adjustedFrames := skillMeleeFrames
		switch c.Core.Player.CurrentState() {
		case action.WalkState:
			adjustedFrames = skillMeleeWalkFrames
		case action.DashState:
			adjustedFrames = skillMeleeDashFrames
		}
		canQueueAfter := math.MaxInt
		for _, f := range adjustedFrames {
			if f < canQueueAfter {
				canQueueAfter = f
			}
		}
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(adjustedFrames),
			AnimationLength: adjustedFrames[action.InvalidAction],
			CanQueueAfter:   canQueueAfter,
			State:           action.SkillState,
		}
	}

	c.eCast = c.Core.F
	c.AddStatus(MeleeKey, 30*60, true)
	c.Core.Log.NewEvent("Foul Legacy activated", glog.LogCharacterEvent, c.Index).
		Write("rtexpiry", c.Core.F+30*60)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Foul Legacy: Raging Tide",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	cdDelay := 14
	hitmark := skillHitmark
	switch c.Core.Player.CurrentState() {
	case action.WalkState:
		hitmark = skillWalkHitmark
		cdDelay = 0
	case action.DashState:
		hitmark = skillDashHitmark
		cdDelay = 0
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), hitmark, hitmark)

	src := c.eCast
	c.QueueCharTask(func() {
		if src == c.eCast && c.StatusIsActive(MeleeKey) {
			c.onExitMeleeStance(0)
			c.ResetNormalCounter()
		}
	}, 30*60)
	c.SetCDWithDelay(action.ActionSkill, 60, cdDelay)

	adjustedFrames := skillRangedFrames
	switch c.Core.Player.CurrentState() {
	case action.WalkState:
		adjustedFrames = skillRangedWalkFrames
	case action.DashState:
		adjustedFrames = skillRangedDashFrames
	}
	canQueueAfter := math.MaxInt
	for _, f := range adjustedFrames {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(adjustedFrames),
		AnimationLength: adjustedFrames[action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 3*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
}

// Hook to end Tartaglia's melee stance prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatusIsActive(MeleeKey) {
			// TODO: need to verify if this is correct
			// but if childe is currently in melee stance and skill is on CD that means that
			// the button has lit up yet from original skill press
			// in which case we need to reset the cooldown first
			c.ResetActionCooldown(action.ActionSkill)
			c.onExitMeleeStance(0)
		}
		return false
	}, "tartaglia-exit")
}

func (c *char) onExitMeleeStance(delay int) {
	// Precise skill CD from Risuke:
	// Aligns with separate table on wiki except the 4 second duration one
	// https://discord.com/channels/763583452762734592/851428030094114847/899416824117084210
	// https://media.discordapp.net/attachments/778615842916663357/781978094495727646/unknown-20.png

	skillCD := 0

	switch timeInMeleeStance := c.Core.F - c.eCast; {
	case timeInMeleeStance < 2*60:
		skillCD = 7 * 60
	case 2*60 <= timeInMeleeStance && timeInMeleeStance < 4*60:
		skillCD = 8 * 60
	case 4*60 <= timeInMeleeStance && timeInMeleeStance < 5*60:
		skillCD = 9 * 60
	case 5*60 <= timeInMeleeStance && timeInMeleeStance < 8*60:
		skillCD = 5*60 + timeInMeleeStance
	case 8*60 <= timeInMeleeStance && timeInMeleeStance < 30*60:
		skillCD = 6*60 + timeInMeleeStance
	case timeInMeleeStance >= 30*60:
		skillCD = 45 * 60
	}

	if c.Base.Cons >= 1 {
		skillCD = int(float64(skillCD) * 0.8)
	}

	if c.mlBurstUsed {
		c.ResetActionCooldown(action.ActionSkill)
		c.mlBurstUsed = false
	} else {
		c.SetCDWithDelay(action.ActionSkill, skillCD, delay)
	}
	c.DeleteStatus(MeleeKey)
}
