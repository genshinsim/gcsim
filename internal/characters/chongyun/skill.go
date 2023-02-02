package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 36

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> N1
	skillFrames[action.ActionBurst] = 51   // E -> Q
	skillFrames[action.ActionDash] = 35    // E -> D
	skillFrames[action.ActionJump] = 35    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> J
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	//if fieldSrc is < duration then this is prob a sac proc
	//we need to stop the old field from ticking (by changing fieldSrc)
	//and also trigger a4 delayed damage
	src := c.Core.F

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Spirit Blade: Chonghua's Layered Frost",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagElementalArt,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Cryo,
		Durability:         50,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.09 * 60,
		CanBeDefenseHalted: false,
	}
	c.skillArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.5}, 8)
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.skillArea.Shape.Pos(), nil, 2.5),
		0,
		skillHitmark,
	)

	c.Core.QueueParticle("chongyun", 4, attributes.Cryo, skillHitmark+c.ParticleDelay)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Chonghua's Layered Frost (A4)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	cb := func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("chongyun-a4", 480),
			Ele:   attributes.Cryo,
			Value: -0.10,
		})
	}
	snap := c.Snapshot(&ai)

	// A4:
	// When the field created by Spirit Blade: Chonghua's Layered Frost disappears,
	// another spirit blade will be summoned to strike nearby opponents,
	// dealing 100% of Chonghua's Layered Frost's Skill DMG as AoE Cryo DMG.
	// Opponents hit by this blade will have their Cryo RES decreased by 10% for 8s.
	hasA4 := c.Base.Ascension >= 4

	// if field is overwriting last
	// TODO: should really just make this a struct, keep a reference, and compare the reference instead
	// of playing around with this int field
	if src-c.fieldSrc < 600 && hasA4 {
		// we're overriding previous field so trigger A4 here
		atk := c.a4Snap
		c.Core.QueueAttackEvent(atk, 1)
	}

	c.fieldSrc = src

	if hasA4 {
		// override previous snap
		c.a4Snap = &combat.AttackEvent{
			Info:     ai,
			Snapshot: snap,
		}
		c.a4Snap.Callbacks = append(c.a4Snap.Callbacks, cb)

		// A4 delayed damage + cryo resist shred
		// TODO: assuming this is NOT affected by hitlag since it should be tied to deployable?
		c.Core.Tasks.Add(func() {
			// if src changed then that means the field changed already
			if src != c.fieldSrc {
				return
			}
			enemy := c.Core.Combat.ClosestEnemyWithinArea(c.skillArea, nil)
			if enemy != nil {
				c.a4Snap.Pattern = combat.NewCircleHitOnTarget(enemy, nil, 3.5)
			} else {
				c.a4Snap.Pattern = combat.NewCircleHitOnTarget(c.skillArea.Shape.Pos(), nil, 3.5)
			}
			// TODO: this needs to be fixed still for sac gs
			c.Core.QueueAttackEvent(c.a4Snap, 0)
		}, 665)
	}

	c.Core.Status.Add("chongyunfield", 600)

	//TODO: delay between when frost field start ticking?
	for i := skillHitmark - 1; i <= 600; i += 60 {
		c.Core.Tasks.Add(func() {
			if src != c.fieldSrc {
				return
			}
			if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
				return
			}
			active := c.Core.Player.ActiveChar()
			c.infuse(active)
		}, i)
	}

	c.SetCDWithDelay(action.ActionSkill, 900, 34)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) onSwapHook() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.Status.Duration("chongyunfield") == 0 {
			return false
		}
		//add infusion on swap
		c.Core.Log.NewEvent("chongyun adding infusion on swap", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
		active := c.Core.Player.ActiveChar()
		c.infuse(active)
		return false
	}, "chongyun-field")
}

func (c *char) infuse(active *character.CharWrapper) {
	//c2 reduces CD by 15%
	if c.Base.Cons >= 2 {
		active.AddCooldownMod(character.CooldownMod{
			Base: modifier.NewBase("chongyun-c2", 126),
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
	case weapon.WeaponClassClaymore, weapon.WeaponClassSpear, weapon.WeaponClassSword:
		c.Core.Player.AddWeaponInfuse(
			active.Index,
			"chongyun-ice-weapon",
			attributes.Cryo,
			infuseDur[c.TalentLvlSkill()],
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		)
		c.Core.Log.NewEvent("chongyun adding infusion", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
		// A1:
		// Sword, Claymore, or Polearm-wielding characters within the field created by
		// Spirit Blade: Chonghua's Layered Frost have their Normal ATK SPD increased by 8%.
		if c.Base.Ascension >= 1 {
			m := make([]float64, attributes.EndStatType)
			m[attributes.AtkSpd] = 0.08
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("chongyun-field", 126),
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
