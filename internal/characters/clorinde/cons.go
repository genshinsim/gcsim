package clorinde

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Icd              int     = 1.2 * 60
	c1AtkP             float64 = 0.3
	c1IcdKey                   = "clorinde-c1-IcdKey"
	c2A1FlatDmg        float64 = 2700
	c2A1PercentBuff    float64 = 0.3
	c6Icd              int     = 12 * 60
	c6IcdKey                   = "clorinde-c6-icd"
	c6Mitigate                 = 0.8
	c6GlimbrightIcdKey         = "glimbrightIcdKey"
	c6GlimbrightAtkP           = 2
)

var c1Hitmarks = []int{1, 1} // TODO hitmark for each c1 hit

// While Hunter's Vigil's Night Vigil state is active,
// when Electro DMG from  Clorinde's Normal Attacks hit opponents,
// they will trigger 2 coordinated attacks from a Nightvigil Shade
// summoned near the hit opponent, each dealing 30% of Clorinde's ATK as Electro DMG.
// This effect can occur once every 1.2s. DMG dealt this way is considered
// Normal Attack DMG.

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if !c.StatusIsActive(skillStateKey) {
			return false
		}
		if c.StatusIsActive(c1IcdKey) {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		if atk.Info.Element != attributes.Electro {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		c.AddStatus(c1IcdKey, c1Icd, false)
		c1AI := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             "Nightwatch Shade (C1)",
			AttackTag:        attacks.AttackTagNormal,
			ICDTag:           attacks.ICDTagClorindeCons,
			ICDGroup:         attacks.ICDGroupClorindeElementalArt,
			StrikeType:       attacks.StrikeTypeSlash,
			Element:          attributes.Electro,
			Durability:       25,
			Mult:             c1AtkP,
			HitlagHaltFrames: 0.01,
		}
		for _, hitmark := range c1Hitmarks {
			c.Core.QueueAttack(
				c1AI,
				combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -3}, 4),
				hitmark,
				hitmark,
				c.particleCB,
			)
		}
		return false
	}, "clorinde-c1")
}

// When Last Lightfall deals DMG to opponent(s),
// DMG dealt is increased based on Clorinde's Bond of Life percentage.
// Every 1% of her current Bond of Life will increase Last Lightfall DMG by 2%.
// The maximum Last Lightfall DMG increase achievable this way is 200%.

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("clorinde-c4-burst-bonus", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			m[attributes.DmgP] = min(c.currentHPDebtRatio()*100*0.02, 2)
			return m, true
		},
	})
}

// For 12s after Hunter's Vigil is used,
// Clorinde's CRIT Rate will be increased by 10%, and her CRIT DMG by 70%.
func (c *char) c6skill() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6Stacks = 6
	if !c.StatusIsActive(skillStateKey) {
		return
	}
	if c.StatusIsActive(c6IcdKey) {
		return
	}
	c.AddStatus(c6IcdKey, c6Icd, true)

	mCR := make([]float64, attributes.EndStatType)
	mCR[attributes.CR] = 0.1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("clorinde-c6-cr-bonus", c6Icd),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			return mCR, true
		},
	})

	mCD := make([]float64, attributes.EndStatType)
	mCD[attributes.CD] = 0.7
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("clorinde-c6-cd-bonus", c6Icd),
		AffectedStat: attributes.CD,
		Amount: func() ([]float64, bool) {
			return mCD, true
		},
	})
}

// Additionally, while Night Vigil is active, a Glimbright Shade
// will appear under specific circumstances, executing an attack
// that deals 200% of Clorinde's ATK as Electro DMG.
// DMG dealt this way is considered Normal Attack DMG.

func (c *char) c6() {
	if c.StatusIsActive(c6GlimbrightIcdKey) {
		return
	}

	c.AddStatus(c6GlimbrightIcdKey, 1*60, false)
	c6AI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glimbright Shade (C6)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagClorindeCons,
		ICDGroup:   attacks.ICDGroupClorindeElementalArt,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       c6GlimbrightAtkP,
	}
	c.Core.QueueAttack(
		c6AI,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 8),
		1, //TODO: c6 hitmark
		1,
	)
}
