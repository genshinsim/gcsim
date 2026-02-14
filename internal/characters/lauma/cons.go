package lauma

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key            = "lauma-c1"
	c1IcdKey         = "lauma-c1-icd"
	c1HitMark        = 5
	c4IcdKey         = "lauma-c4-icd"
	c6ElevationBonus = 0.25
	c6SkillHitName   = "Frostgrove Sanctuary C6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	// on lb proc heal
	c.Core.Events.Subscribe(event.OnLunarBloom, func(args ...any) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		if !c.StatusIsActive(c1Key) {
			return false
		}

		if c.StatusIsActive(c1IcdKey) {
			return false
		}

		c.AddStatus(c1IcdKey, 1.9*60, true)

		healAmt := 5.0 * c.Stat(attributes.EM)

		// heal active character
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(info.HealInfo{
				Type:    info.HealTypeAbsolute,
				Message: "Lauma C1 (Heal)",
				Src:     healAmt,
			})
		}, c1HitMark)

		return true
	}, "lauma-c1")
}

func (c *char) c1OnBurst() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddStatus(c1Key, 20*60, true)
}

func (c *char) c1DeerStamMod() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 0.6
}

func (c *char) c1DeerDurMod() int {
	if c.Base.Cons < 1 {
		return 0
	}
	return 5 * 60
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	if !c.ascendantGleam {
		return
	}

	for _, x := range c.Core.Player.Chars() {
		x.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("lauma-c2-lunarbloom-buff", -1),
			Amount: func(atk info.AttackInfo) (float64, bool) {
				if atk.AttackTag != attacks.AttackTagDirectLunarBloom {
					return 0, false
				}
				return 0.4, false
			},
		})
	}
}

func (c *char) c2PaleHymnScalingBloom() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	return 5
}

func (c *char) c2PaleHymnScalingLunarBloom() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	return 4
}

func (c *char) c4RefundCB(a info.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4IcdKey) {
		return
	}
	c.AddEnergy("lauma-c4", 5)
	c.AddStatus(c4IcdKey, 5*60, true)
}

func (c *char) addC6PaleHymnCB(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}

	c.addC6PaleHymn(2)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	if !c.ascendantGleam {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return false
		}

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma-c6-elevation", glog.LogCharacterEvent, c.Index()).Write("bonus", c6ElevationBonus)
		}

		atk.Info.Elevation += c6ElevationBonus
		return false
	}, lunarbloomBonusKey+"-c6")
}

func (c *char) c6OnSkill() {
	if c.Base.Cons < 6 {
		return
	}
	// TODO: Does clearing on skill use have a delay?
	c.DeleteStatus(paleHymnC6Key)
	c.paleHymn[paleHymnC6] = 0
	c.paleHymnSrc[paleHymnC6] = 0
	c.c6Count = 0
}

func (c *char) c6OnFrostgroveTick() {
	if c.Base.Cons < 6 {
		return
	}

	if c.c6Count >= 8 {
		return
	}

	c.c6Count++

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             c6SkillHitName,
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Durability:       0,
		UseEM:            true,
		Mult:             1.85,
		IgnoreDefPercent: 1,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player().Pos(),
			nil,
			6,
		),
		16, // 0.26s delay from DM
		16,
		c.addC6PaleHymnCB,
	)
}
