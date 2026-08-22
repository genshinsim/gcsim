package prune

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

var (
	skillFrames        []int
	skillConvertFrames []int
)

const (
	particleICDKey        = "prune-particle-icd"
	skillRecastWindowKey  = "prune-skill-recast-window"
	skillSwirlCheckICDKey = "prune-skill-swirl-check-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(44) // walk
	skillFrames[action.ActionAttack] = 56
	skillFrames[action.ActionCharge] = 59
	skillFrames[action.ActionSkill] = 41
	skillFrames[action.ActionBurst] = 28
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 29
	skillFrames[action.ActionSwap] = 27

	skillConvertFrames = frames.InitAbilSlice(82) // walk
	skillConvertFrames[action.ActionAttack] = 65
	skillConvertFrames[action.ActionCharge] = 76
	skillConvertFrames[action.ActionSkill] = 69
	skillConvertFrames[action.ActionBurst] = 67
	skillConvertFrames[action.ActionDash] = 67
	skillConvertFrames[action.ActionJump] = 66
	skillConvertFrames[action.ActionSwap] = 65
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// if this is second press and had swirled, do the convert dmg
	if c.StatusIsActive(skillRecastWindowKey) {
		return c.skillConvert()
	}

	skillFrames[action.ActionSkill] = 41
	c.skillConvertEle = attributes.NoElement

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Ring-A-Ding-Ding! Hexhunter Chime DMG",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             skill[c.TalentLvlSkill()],
		HitlagFactor:     0.05,
		HitlagHaltFrames: 0.02 * 60,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1.4}, 2.5),
		0,
		26,
		c.particleCB,
	)

	swirlfunc := func(ele attributes.Element) func(args ...any) {
		return func(args ...any) {
			// reaction target must be an enemy
			if _, ok := args[0].(*enemy.Enemy); !ok {
				return
			}

			// reaction must have been triggered by prune's attack
			atk, ok := args[1].(*info.AttackEvent)
			if !ok {
				return
			}

			if atk.Info.ActorIndex != c.Index() {
				return
			}

			// reaction must have been triggered by the initial skill hit
			if atk.Info.Abil != "Ring-A-Ding-Ding! Hexhunter Chime DMG" {
				return
			}

			// A1 has a shared 1.2s ICD across all four swirl elements
			if c.StatusIsActive(skillSwirlCheckICDKey) {
				return
			}
			c.AddStatus(skillSwirlCheckICDKey, 12, false)

			// allow skill to be recasted
			c.skillConvertEle = ele
			c.AddStatus(skillRecastWindowKey, 364, true)

			// skill recast has different cancel frame
			skillFrames[action.ActionSkill] = 28
		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.Cryo), "prune-skill-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.Electro), "prune-skill-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.Hydro), "prune-skill-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.Pyro), "prune-skill-pyro")
	// TODO: Add subscriptions for stellar-swirl when it's implemented

	c.SetCDWithDelay(action.ActionSkill, 15*60, 25)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillConvert() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Clang Clang! Witch-tribution Comes! DMG",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           75,
		Element:            c.skillConvertEle,
		Durability:         25,
		Mult:               skillConvert[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.05 * 60,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1.5}, 2.5),
		0,
		30,
		c.makeA4CB,
		c.makeC1CB,
		c.makeC2CB,
		c.c4Ricochet(c.skillConvertEle, attacks.AttackTagElementalArt),
	)

	c.DeleteStatus(skillRecastWindowKey)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillConvertFrames),
		AnimationLength: skillConvertFrames[action.InvalidAction],
		CanQueueAfter:   skillConvertFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Anemo, c.ParticleDelay)
}
