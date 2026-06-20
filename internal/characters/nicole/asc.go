package nicole

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Key = "guidance-of-theosis"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Buff = make([]float64, attributes.EndStatType)
	c.a1Buff[attributes.ATK] = 300

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		if prev != c.Index() && c.c6DeleteA1OnSwap() {
			char := c.Core.Player.Chars()[prev]
			char.DeleteStatMod(a1Key)
		}

		next := args[1].(int)
		char := c.Core.Player.Chars()[next]
		src := c.Core.F
		c.a1Src = src

		delay := 3 * 60
		if char.IsHexerei {
			delay = 0
		}

		c.QueueCharTask(func() {
			if c.a1Src != src {
				return
			}
			c.a1UpgradeBuff(char, -1)
		}, delay)
	}, "nicole-a1")
}

func (c *char) a1OnSkillRemoveBuff() {
	if c.Base.Ascension < 1 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.DeleteStatMod(a1Key)
	}
}

func (c *char) a4Init() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		switch ae.Info.Element {
		case attributes.Pyro:
		case attributes.Hydro:
		case attributes.Electro:
		case attributes.Cryo:
		case attributes.Anemo:
		case attributes.Geo:
		case attributes.Dendro:
		default:
			return
		}

		dur := c.c6SelfBuffDur()
		c.a1UpgradeBuff(c.CharWrapper, dur)
	}, "nicole-a4")
}

// Upgrades the skill buff for A1. If the dur is less than 0, the duration of the upgrade is equal to the duration of the skill buff.
func (c *char) a1UpgradeBuff(char *character.CharWrapper, dur int) {
	if dur < 0 {
		dur = char.StatusDuration(skillBuffKey)
	} else {
		dur = min(dur, char.StatusDuration(skillBuffKey))
	}
	if dur == 0 {
		return
	}
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(a1Key, dur),
		AffectedStat: attributes.ATK,
		Amount: func() []float64 {
			return c.a1Buff
		},
	})
	c.c2OnUpgrade(char)
	c.c4OnUpgrade(char)
	c.c6OnUpgrade(char)
}

func (c *char) hexereiOnProjection(char *character.CharWrapper) float64 {
	if !c.IsHexerei {
		return 0
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return 0
	}

	if !char.IsHexerei {
		return 0
	}

	return c.TotalAtk() * 3
}
