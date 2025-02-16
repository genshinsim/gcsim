package ororon

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const a1NSBurstKey = "ororon-a1-ns-burst"
const a1ElectroHydroKey = "ororon-a1-electro-hydro"
const a1ECTriggerKey = "ororon-a1-ec"
const a1NSTriggerKey = "ororon-a1-ns"

const a1OnSkillKey = "ororon-a1"
const a1GainIcdKey = "ororon-a1-gain-icd"
const a1DamageIcdKey = "ororon-a1-dmg-icd"
const a1Abil = "Hypersense"

const a4Key = "ororon-a4"
const a4IcdKey = "ororon-a4-icd"

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.nightsoulState.GeneratePoints(40)
		return false
	}, a1NSBurstKey)

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		// ignores ororon himself
		if atk.Info.ActorIndex == c.Index {
			return false
		}

		switch atk.Info.Element {
		case attributes.Hydro:
		case attributes.Electro:
		default:
			return false
		}

		if !c.StatusIsActive(a1OnSkillKey) {
			return false
		}
		if c.StatusIsActive(a1GainIcdKey) {
			return false
		}
		c.AddStatus(a1GainIcdKey, 0.3*60, true)

		c.nightsoulState.GeneratePoints(5)
		c.SetTag(a1ElectroHydroKey, c.Tag(a1ElectroHydroKey)+1)
		if c.Tag(a1ElectroHydroKey) >= 10 {
			c.DeleteStatus(a1OnSkillKey)
		}
		return false
	}, a1ElectroHydroKey)

	c.Core.Events.Subscribe(event.OnElectroCharged, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		c.a1NightSoulAttack(atk)
		return false
	}, a1ECTriggerKey)

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if atk.Info.ActorIndex == c.Index {
			return false
		}
		if !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagNightsoul) {
			return false
		}
		c.a1NightSoulAttack(atk)

		return false
	}, a1NSTriggerKey)
}

func (c *char) a1NightSoulAttack(atk *combat.AttackEvent) {
	if c.nightsoulState.Points() < 10 {
		return
	}
	if c.StatusIsActive(a1DamageIcdKey) {
		return
	}
	c.AddStatus(a1DamageIcdKey, 1.8*60, true)
	c.a1EnterBlessing()

	c.nightsoulState.ConsumePoints(10)
	c.hypersense(1.6, a1Abil, atk.Pattern.Shape.Pos())
}

func (c *char) hypersense(mult float64, abil string, initialTargetPos geometry.Point) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               abil,
		AttackTag:          attacks.AttackTagNone,
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               mult,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(initialTargetPos, nil, 15), nil)
	for i := 0; len(enemies) < 4 && i < len(enemies); i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(
				enemies[i].Pos(),
				nil,
				0.2,
				0.2,
			),
			12,
			12,
		)
	}
	c.c6onHypersense()
}

// When Ororon has the jump blessing, do nothing. Blessing will exit when jump is done.
func (c *char) a1ExitBlessing() {
	c.inA1Blessing = false
	if !c.inTransmissionBlessing {
		c.nightsoulState.ExitBlessing()
	}
}

func (c *char) a1EnterBlessing() {
	c.nightsoulState.EnterBlessing(c.nightsoulState.Points())
	c.inA1Blessing = true
	c.a1Src = c.Core.F
	src := c.a1Src
	c.QueueCharTask(func() {
		if src != c.a1Src {
			return
		}
		c.a1ExitBlessing()
	}, 6*60)
}

func (c *char) a1OnSkill() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(a1OnSkillKey, 15*60, true)
	c.SetTag(a1OnSkillKey, 0)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		active := c.Core.Player.ActiveChar()
		if atk.Info.ActorIndex != active.Index {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		if !c.StatusIsActive(a4Key) {
			return false
		}
		if c.StatusIsActive(a4IcdKey) {
			return false
		}
		c.AddStatus(a4IcdKey, 60, true)

		active.AddEnergy(a4Key, 3)
		if active.Index != c.Index {
			c.AddEnergy(a4Key, 3)
		}

		c.SetTag(a4Key, c.Tag(a4Key)+1)
		if c.Tag(a4Key) >= 3 {
			c.DeleteStatus(a4Key)
		}
		return false
	}, a4Key)
}

func (c *char) makeA4cb() func(combat.AttackCB) {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		c.AddStatus(a4Key, 15*60, true)
		c.SetTag(a4Key, 0)
	}
}
