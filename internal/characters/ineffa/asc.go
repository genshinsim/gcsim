package ineffa

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
	a4Key              = "ineffa-a4"
	lunarchageBonusKey = "ineffa-lc-bonus"
)

func (c *char) a1OnDischarge() {
	if c.Base.Ascension < 1 {
		return
	}

	if c.Core.Status.Duration(reactable.LcKey) == 0 {
		return
	}

	a1Atk := func() {
		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Birgitta (A1)",
			AttackTag:        attacks.AttackTagDirectLunarCharged,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDirectLunarCharged,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Electro,
			Mult:             0.65,
			IgnoreDefPercent: 1,
		}
		c.Core.QueueAttack(ai, combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4, 4.1), 0, 0)
	}

	c.Core.Tasks.Add(a1Atk, 30)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	// TODO: Is this buff hitlag affected on Ineffa only? Or is it hitlag affected per character?
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key+"-buff", -1),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.Active() != char.Index() && c.Index() != char.Index() {
					return nil, false
				}
				if !c.StatusIsActive(a4Key) {
					return nil, false
				}

				stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
				m[attributes.EM] = stats.TotalATK() * 0.06
				return m, true
			},
		})
	}
}

func (c *char) a4OnBurst() {
	if c.Base.Ascension < 4 {
		return
	}

	// TODO when does this start and end?
	c.AddStatus(a4Key, 20*60, true)
}

func (c *char) lunarchargeInit() {
	c.Core.Flags.Custom[reactable.LunarChargeEnableKey] = 1

	// TODO: moonsign?

	// TODO: every 100 ATK that Ineffa has increasing Lunar-Charged's Base DMG by 0.7%, up to a maximum of 14%.
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCharged:
		case attacks.AttackTagReactionLunarCharge:
		default:
			return false
		}

		stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
		bonus := min(stats.TotalATK()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("ineffa adding lunarcharged base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarchageBonusKey)
}
