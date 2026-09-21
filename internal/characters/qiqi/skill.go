package qiqi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	skillHitmark        = 32
	skillBuffKey        = "qiqi-e"
	skillICDKey         = "qiqi-e-icd"
	skillParticleICDKey = "qiqi-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(58) // E -> Q
	skillFrames[action.ActionAttack] = 57  // E -> N1
	skillFrames[action.ActionDash] = 6     // E -> D
	skillFrames[action.ActionJump] = 5     // E -> J
	skillFrames[action.ActionSwap] = 57    // E -> Swap
}

// Using the Icevein Talisman, Qiqi brings forth the Herald of Frost, dealing Cryo DMG to
// surrounding opponents.
//
// Herald of Frost
// · On hit, Qiqi's Normal and Charged Attacks regenerate HP for your own party members and nearby
// teammates, with healing scaled based on Qiqi's ATK.
// · Periodically regenerates your active character's HP.
// · Follows the character around, dealing Cryo DMG to opponents in the character's path.
// If your current active character hits an opponent with an attack, the Herald of Frost will
// perform a coordinated attack on the opponent, dealing Cryo DMG.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	// +1 to avoid end duration issues
	// Qiqi E is a deployable after Initial Hit, so it shouldn't be hitlag extendable
	c.AddStatus(skillBuffKey, 15*60+1, false)
	c.skillLastUsed = c.Core.F
	src := c.Core.F

	// Initial damage
	// Both healing and damage are snapshot
	c.Core.Tasks.Add(func() {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Herald of Frost: Initial Damage",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagQiqiElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
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
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
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

		// assumes ruin guard hitlag as extra delay (E Tick 1 gets delayed by E Initial hitlag)
		// can't use char queue for this because hitlag from other sources shouldn't count
		c.Core.Tasks.Add(c.skillDmgTickTask(src, &aiTick, &snapTick, 60), 57+7)

		// Apply damage needs to take place after above takes place to ensure stats are handled correctly
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), 0)
	}, skillHitmark)

	c.QueueCharTask(func() {
		c.SetCD(action.ActionSkill, 1800) // 30s * 60
		c.revelationSkillCDReduction()
	}, 3)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDmgTickTask(src int, ai *info.AttackInfo, snap *info.Snapshot, lastTickDuration int) func() {
	return func() {
		if !c.StatusIsActive(skillBuffKey) {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		// pattern shouldn't snapshot on attack event creation because the skill follows the player
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5)

		c.Core.QueueAttackWithSnap(*ai, *snap, ap, 0, c.particleCB, c.c1CB, c.c1RevelationCB)

		nextTick := 60
		if lastTickDuration == 60 {
			nextTick = 135
		}
		c.Core.Tasks.Add(c.skillDmgTickTask(src, ai, snap, nextTick), nextTick)
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

		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  c.Core.Player.Active(),
			Message: "Herald of Frost (Tick)",
			Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()),
			Bonus:   c.skillHealSnapshot.Stats[attributes.Heal],
		})

		// Queue next instance
		c.Core.Tasks.Add(c.skillHealTickTask(src), 4.5*60)
	}
}

func (c *char) skillInit() {
	if !c.revelation {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		// only be triggered by on field
		if c.Core.Player.Active() != atk.Info.ActorIndex {
			return
		}

		// ignore EC, hydro swirl, and burning damage
		// this clause is here since these damage types are sourced to the target rather than character
		switch atk.Info.AttackTag {
		case attacks.AttackTagECDamage, attacks.AttackTagBurningDamage, attacks.AttackTagSwirlHydro:
			return
		}

		if !c.StatusIsActive(skillBuffKey) {
			return
		}

		// Coordinated attack ICD is per character, not global
		char := c.Core.Player.ActiveChar()
		if char.StatusIsActive(skillICDKey) {
			return
		}

		char.AddStatus(skillICDKey, 2.2*60, true)

		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Herald of Frost: Coordinated Attack",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagQiqiElementalArt,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Durability:         25,
			Mult:               skillCoord[c.TalentLvlSkill()],
			HitlagFactor:       0.05,
			HitlagHaltFrames:   0.05 * 60,
			CanBeDefenseHalted: true,
			IsDeployable:       true,
		}

		ap := combat.NewCircleHitOnTarget(e, nil, 2)

		// based on video it's the same frame. But to prevent issues due to delayed aura attachment,
		// it's set to 1 frame after
		c.Core.QueueAttack(ai, ap, 1, 1, c.particleCB, c.c1CB, c.c1RevelationCB)
	}, "qiqi-e-hook")
}

func (c *char) skillHealCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	// Qiqi NA/CA healing proc in skill duration
	if !c.StatusIsActive(skillBuffKey) {
		return
	}
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  -1,
		Message: "Herald of Frost (Attack)",
		Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealOnHitPer, skillHealOnHitFlat, c.TalentLvlSkill()),
		Bonus:   c.skillHealSnapshot.Stats[attributes.Heal],
	})
}

func (c *char) particleCB(a info.AttackCB) {
	if !c.revelation {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(skillParticleICDKey) {
		return
	}
	c.AddStatus(skillParticleICDKey, 6*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Cryo, c.ParticleDelay)
}
