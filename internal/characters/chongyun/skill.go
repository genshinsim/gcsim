package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark   = 36
	skillFieldKey  = "chongyunfield"
	particleICDKey = "chongyun-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> N1
	skillFrames[action.ActionBurst] = 51   // E -> Q
	skillFrames[action.ActionDash] = 35    // E -> D
	skillFrames[action.ActionJump] = 35    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> J
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Spirit Blade: Chonghua's Layered Frost",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           150,
		Element:            attributes.Cryo,
		Durability:         50,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.09 * 60,
		CanBeDefenseHalted: false,
	}

	// handle field expiry (a4) on field end via sac greatsword / large amount of cd reduction
	// need src to invalidate field ticks and a4 task
	src := c.Core.F
	c.fieldSrc = c.Core.F
	// if the field is still up then need to invalidate existing a4 task and do damage before resnapshotting for new a4
	// need to do this before the new skill area is determined
	if c.Core.Status.Duration(skillFieldKey) > 0 {
		c.a4(skillHitmark+45, c.Core.F, true) // ~45f from field expiring
	}

	// handle field damage / skill area
	c.skillArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.5}, 8)
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.skillArea.Shape.Pos(), nil, 2.5),
		0,
		skillHitmark,
		c.particleCB,
		c.makeC4Callback(),
	)

	// handle field creation
	c.QueueCharTask(func() {
		c.Core.Status.Add(skillFieldKey, 600)
	}, skillHitmark)

	// handle field ticks
	// TODO: delay between when frost field start ticking?
	for i := 0; i <= 600; i += 60 {
		c.Core.Tasks.Add(func() {
			if src != c.fieldSrc {
				return
			}
			if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
				return
			}
			active := c.Core.Player.ActiveChar()
			c.infuse(active)
		}, i+skillHitmark)
	}

	// handle field expiry (a4) on field end via expiry
	c.a4(655, c.Core.F, false)

	c.SetCDWithDelay(action.ActionSkill, 900, 34)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
}

func (c *char) onSwapHook() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.Status.Duration("chongyunfield") == 0 {
			return false
		}
		// add infusion on swap
		dur := int(infuseDur[c.TalentLvlSkill()] * 60)
		c.Core.Log.NewEvent("chongyun adding infusion on swap", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+dur)
		active := c.Core.Player.ActiveChar()
		c.infuse(active)
		return false
	}, "chongyun-field")
}

func (c *char) infuse(active *character.CharWrapper) {
	dur := int(infuseDur[c.TalentLvlSkill()] * 60)
	// c2 reduces CD by 15%
	if c.Base.Cons >= 2 {
		active.AddCooldownMod(character.CooldownMod{
			Base: modifier.NewBaseWithHitlag("chongyun-c2", dur),
			Amount: func(a action.Action) float64 {
				if a == action.ActionSkill || a == action.ActionBurst {
					return -0.15
				}
				return 0
			},
		})
	}

	// weapon infuse and A1
	switch active.Weapon.Class {
	case info.WeaponClassClaymore, info.WeaponClassSpear, info.WeaponClassSword:
		c.Core.Player.AddWeaponInfuse(
			active.Index,
			"chongyun-ice-weapon",
			attributes.Cryo,
			dur,
			true,
			attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
		)
		c.Core.Log.NewEvent("chongyun adding infusion", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+dur)
		// A1:
		// Sword, Claymore, or Polearm-wielding characters within the field created by
		// Spirit Blade: Chonghua's Layered Frost have their Normal ATK SPD increased by 8%.
		if c.Base.Ascension >= 1 {
			m := make([]float64, attributes.EndStatType)
			m[attributes.AtkSpd] = 0.08
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("chongyun-field", dur),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	default:
		return
	}
}
