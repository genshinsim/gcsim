package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillMeleeFrames []int
var skillMeleeWalkFrames []int
var skillMeleeDashFrames []int
var skillRangedFrames []int
var skillRangedWalkFrames []int
var skillRangedDashFrames []int

const skillHitmark = 16
const skillWalkHitmark = 3
const skillDashHitmark = 3

func init() {
	// skill (melee) -> x
	skillMeleeFrames = frames.InitAbilSlice(18)
	skillMeleeFrames[action.ActionAttack] = 17
	skillMeleeFrames[action.ActionBurst] = 18
	skillMeleeFrames[action.ActionDash] = 17
	skillMeleeFrames[action.ActionJump] = 17
	skillMeleeFrames[action.ActionSwap] = 16
	skillMeleeWalkFrames = frames.InitAbilSlice(24)
	skillMeleeWalkFrames[action.ActionAttack] = 5
	skillMeleeWalkFrames[action.ActionBurst] = 5
	skillMeleeWalkFrames[action.ActionDash] = 6
	skillMeleeWalkFrames[action.ActionJump] = skillWalkHitmark
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
	skillRangedWalkFrames = frames.InitAbilSlice(24)
	skillRangedWalkFrames[action.ActionAttack] = 5
	skillRangedWalkFrames[action.ActionBurst] = 4
	skillRangedWalkFrames[action.ActionDash] = 5
	skillRangedWalkFrames[action.ActionJump] = 4
	skillRangedDashFrames = frames.InitAbilSlice(24)
	skillRangedDashFrames[action.ActionAttack] = 17
	skillRangedDashFrames[action.ActionBurst] = 17
	skillRangedDashFrames[action.ActionDash] = 22
	skillRangedDashFrames[action.ActionJump] = 3
}

//Cast: AoE strong hydro damage
//Melee Stance: infuse NA/CA to hydro damage
func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.Core.Status.Duration("tartagliamelee") > 0 {
		c.onExitMeleeStance()
		c.ResetNormalCounter()
		adjustedFrames := skillMeleeFrames
		lastAction := &c.Core.Player.LastAction
		if (lastAction.Char == c.Index) {
			switch lastAction.Type {
			case action.ActionWalk:
				adjustedFrames = skillMeleeWalkFrames
			case action.ActionDash:
				adjustedFrames = skillMeleeDashFrames
			}
		}
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(adjustedFrames),
			AnimationLength: adjustedFrames[action.InvalidAction],
			CanQueueAfter:   adjustedFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}
	}

	c.eCast = c.Core.F
	c.Core.Status.Add("tartagliamelee", 30*60)
	c.Core.Log.NewEvent("Foul Legacy activated", glog.LogCharacterEvent, c.Index).
		Write("rtexpiry", c.Core.F+30*60)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Foul Legacy: Raging Tide",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	hitmark := skillHitmark
	lastAction := &c.Core.Player.LastAction
	if (lastAction.Char == c.Index && lastAction.Type == action.ActionWalk) {
		switch lastAction.Type {
		case action.ActionWalk:
			hitmark = skillWalkHitmark
		case action.ActionDash:
			hitmark = skillDashHitmark
		}
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), hitmark, hitmark)

	src := c.eCast
	c.Core.Tasks.Add(func() {
		if src == c.eCast && c.Core.Status.Duration("tartagliamelee") > 0 {
			c.onExitMeleeStance()
			c.ResetNormalCounter()
		}
	}, 30*60)
	c.SetCDWithDelay(action.ActionSkill, 60, 14)

	adjustedFrames := skillRangedFrames
	if (lastAction.Char == c.Index) {
		switch lastAction.Type {
			case action.ActionWalk:
				adjustedFrames = skillRangedWalkFrames
			case action.ActionDash:
				adjustedFrames = skillRangedDashFrames
		}
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(adjustedFrames),
		AnimationLength: adjustedFrames[action.InvalidAction],
		CanQueueAfter:   adjustedFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

// Hook to end Tartaglia's melee stance prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			//TODO: need to verify if this is correct
			//but if childe is currently in melee stance and skill is on CD that means that
			//the button has lit up yet from original skill press
			//in which case we need to reset the cooldown first
			c.ResetActionCooldown(action.ActionSkill)
			c.onExitMeleeStance()
		}
		return false
	}, "tartaglia-exit")
}

func (c *char) onExitMeleeStance() {
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
		c.SetCD(action.ActionSkill, skillCD)
	}
	c.Core.Status.Delete("tartagliamelee")
}
