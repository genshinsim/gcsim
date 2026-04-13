package linnea

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key                    = "linnea-a1"
	a4Key                    = "linnea-a4"
	lunarcrystallizeBonusKey = "linnea-lcr-bonus"
)

func (c *char) a1OnLumi(src int) {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Ticker(src)
}

func (c *char) a1Ticker(src int) {
	if c.skillSrc != src {
		return
	}
	if !c.StatusIsActive(skillStandardPower) && !c.StatusIsActive(skillSuperPower) {
		return
	}
	shred := 0.15
	if c.Core.Player.GetMoonsignLevel() >= 2 {
		shred += 0.15
	}
	for _, e := range c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.LumiPos(), nil, 5), nil) {
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(a1Key, 1.5*60),
			Ele:   attributes.Geo,
			Value: -shred,
		})
	}
	c.Core.Tasks.Add(func() { c.a1Ticker(src) }, 60)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)

	// only give buff to other moonsign chars
	for _, char := range c.Core.Player.Chars() {
		if char.Moonsign == 0 {
			continue
		}
		if char.Index() == c.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key, -1),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				if c.Core.Player.Active() != char.Index() {
					return nil
				}

				m[attributes.EM] = c.TotalDef(true) * 0.05
				return m
			},
		})
	}

	// give buff to self
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a4Key, -1),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			if c.Core.Player.ActiveChar().Moonsign > 0 && c.Core.Player.Active() != c.Index() {
				return nil
			}

			m[attributes.EM] = c.TotalDef(true) * 0.05
			return m
		},
	})
}

func (c *char) moonsignInit() {
	c.Core.Flags.Custom[reactable.LunarCrystallizeEnableKey] = 1
	f := func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCrystallize:
		case attacks.AttackTagReactionLunarCrystallize:
		default:
			return
		}

		bonus := min(c.TotalDef(true)/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("linnea adding lunar crystallize base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, f, lunarcrystallizeBonusKey)
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, f, lunarcrystallizeBonusKey)
}
