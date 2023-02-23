package qiqi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const (
	skillHitmark = 32
	skillBuffKey = "qiqi-e"
)

func init() {
	skillFrames = frames.InitAbilSlice(57) // E -> N1/Swap
	skillFrames[action.ActionBurst] = 58   // E -> Q
	skillFrames[action.ActionDash] = 6     // E -> D
	skillFrames[action.ActionJump] = 5     // E -> J
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// +1 to avoid end duration issues
	// Qiqi E is a deployable after Initial Hit, so it shouldn't be hitlag extendable
	c.AddStatus(skillBuffKey, 15*60+1, false)
	c.skillLastUsed = c.Core.F
	src := c.Core.F

	// Initial damage
	// Both healing and damage are snapshot
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Herald of Frost: Initial Damage",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             combat.ICDTagElementalArt,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Durability:         25,
			Mult:               skillInitialDmg[c.TalentLvlSkill()],
			HitlagFactor:       0.05,
			HitlagHaltFrames:   0.05 * 60,
			CanBeDefenseHalted: true,
		}
		snap := c.Snapshot(&ai)

		// One healing proc happens immediately on cast
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Herald of Frost (Tick)",
			Src:     c.healSnapshot(&snap, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()),
			Bonus:   snap.Stats[attributes.Heal],
		})

		// Healing and damage instances are snapshot
		// Separately cloned snapshots are fed into each function to ensure nothing interferes with each other

		// Queue up continuous healing instances
		// No exact frame data on when the healing ticks happen. Just roughly guessing here
		// Healing ticks happen 3 additional times during the skill - assume ticks are roughly 4.5s apart
		// so in sec (0 = skill cast), 1, 5.5, 10, 14.5
		c.skillHealSnapshot = snap
		c.Core.Tasks.Add(c.skillHealTickTask(src), 4.5*60)

		// Queue up damage swipe instances.
		// No exact frame data on when the damage ticks happen. Just roughly guessing here
		// Occurs 9 times over the course of the skill
		// Once shortly after initial cast, then 8 additional procs over the rest of the duration
		// Each proc occurs in "pairs" of two swipes each spaced around 2.25s apart
		// The time between each swipe in a pair is about 1s
		// No exact frame data available plus the skill duration is affected by hitlag
		// Damage procs occur (in sec 0 = skill cast): 1.5, 3.75, 4.75, 7, 8, 10.25, 11.25, 13.5, 14.5

		aiTick := ai
		aiTick.Abil = "Herald of Frost: Skill Damage"
		aiTick.Mult = skillDmgCont[c.TalentLvlSkill()]
		aiTick.IsDeployable = true // ticks still apply hitlag but is a deployable so doesnt affect qiqi

		snapTick := c.Snapshot(&aiTick)
		tickAE := &combat.AttackEvent{
			Info:        aiTick,
			Snapshot:    snapTick,
			SourceFrame: c.Core.F,
		}

		// assumes ruin guard hitlag as extra delay (E Tick 1 gets delayed by E Initial hitlag)
		// can't use char queue for this because hitlag from other sources shouldn't count
		c.Core.Tasks.Add(c.skillDmgTickTask(src, tickAE, 60), 57+7)

		// Apply damage needs to take place after above takes place to ensure stats are handled correctly
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), 0)
	}, skillHitmark)

	c.SetCDWithDelay(action.ActionSkill, 1800, 3) // 30s * 60

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}
}

// Handles skill damage swipe instances
// Also handles C1:
// When the Herald of Frost hits an opponent marked by a Fortune-Preserving Talisman, Qiqi regenerates 2 Energy.
func (c *char) skillDmgTickTask(src int, ae *combat.AttackEvent, lastTickDuration int) func() {
	return func() {
		if !c.StatusIsActive(skillBuffKey) {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		// Clones initial snapshot
		tick := *ae //deference the pointer here
		// pattern shouldn't snapshot on attack event creation because the skill follows the player
		tick.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5)

		if c.Base.Cons >= 1 {
			tick.Callbacks = append(tick.Callbacks, c.c1)
		}

		c.Core.QueueAttackEvent(&tick, 0)

		nextTick := 60
		if lastTickDuration == 60 {
			nextTick = 135
		}
		c.Core.Tasks.Add(c.skillDmgTickTask(src, ae, nextTick), nextTick)
	}
}

// Handles skill auto healing ticks
func (c *char) skillHealTickTask(src int) func() {
	return func() {
		if !c.StatusIsActive(skillBuffKey) {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Herald of Frost (Tick)",
			Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()),
			Bonus:   c.skillHealSnapshot.Stats[attributes.Heal],
		})

		// Queue next instance
		c.Core.Tasks.Add(c.skillHealTickTask(src), 4.5*60)
	}
}
