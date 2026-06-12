package nicole

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
	c1ICDKey = "nicole-c1-icd"
	c2Key    = "nicole-c2"
	c4Key    = "nicole-c4"
	c4ICDKey = "nicole-c4-icd"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if c.StatusIsActive(c1ICDKey) {
			return
		}

		c.AddStatus(c1ICDKey, 6*60, true)

		char := c.Core.Player.Chars()[ae.Info.ActorIndex]

		ai := info.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Arcane Projection: Unity",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    char.Base.Element,
			Mult:       6,
			FlatDmg:    c.hexereiOnProjection(char),
		}

		ap := combat.NewCircleHitOnTarget(t, nil, 3)

		c.Core.QueueAttack(ai, ap, projectionHitmark, projectionHitmark)
	}, "nicole-c1")
}

func (c *char) c2OnSkillRemoveBuff() {
	if c.Base.Cons < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.DeleteStatMod(c2Key)
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Buff = make([]float64, attributes.EndStatType)
	c.c2Buff[attributes.ATK] = 300
}

func (c *char) c2OnSkillAddBuff() {
	if c.Base.Cons < 2 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c2Key, skillDur),
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				return c.c2Buff
			},
		})
	}
}

func (c *char) c2OnUpgrade(char *character.CharWrapper) {
	if c.Base.Cons < 2 {
		return
	}

	if !char.StatusIsActive(a1Key) {
		return
	}

	src := c.Core.F
	char.Tags[c2Key] = src
	c.c2Ticker(char, src)
}

func (c *char) c2Ticker(char *character.CharWrapper, src int) {
	if !char.StatusIsActive(a1Key) {
		return
	}

	if char.Tags[c2Key] != src {
		return
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
		e, ok := e.(*enemy.Enemy)
		if !ok {
			continue
		}
		elem := char.Base.Element
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("nicole-c2-"+elem.String(), 1*60),
			Ele:   elem,
			Value: -0.25,
		})
	}
	c.Core.Tasks.Add(func() { c.c2Ticker(char, src) }, 30)
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if !char.StatusIsActive(c4Key) {
			return
		}

		if char.Tags[c4Key] > 0 {
			amt := 0.7 * c.TotalAtk()
			char.Tags[c4Key]--

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Nicole c4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("effect_ends_at", char.StatusExpiry(c4Key)).
					Write("stacks_left", char.Tags[c4Key])
			}

			atk.Info.FlatDmg += amt
		}
	}, "nicole-c4-hook")
}

func (c *char) c4OnUpgrade(char *character.CharWrapper) {
	if c.Base.Cons < 4 {
		return
	}

	if char.StatusIsActive(c4ICDKey) {
		return
	}
	char.AddStatus(c4ICDKey, 16*60, true)
	char.AddStatus(c4Key, 20*60, true)
	char.Tags[c4Key] = 8
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !char.StatModIsActive(a1Key) {
			return
		}
		atk.Info.IgnoreDefPercent = min(atk.Info.IgnoreDefPercent+0.4, 0.9)
	}, "nicole-c6-hook")
}

func (c *char) c6SelfBuffDur() int {
	if c.Base.Cons < 6 {
		return 8 * 60
	}
	return -1
}

func (c *char) c6OnUpgrade(char *character.CharWrapper) {
	if c.Base.Cons < 6 {
		return
	}

	if char.Index() == c.Index() {
		for i, otherChar := range c.Core.Player.Chars() {
			if i == c.Index() {
				continue
			}
			c.a1UpgradeBuff(otherChar, -1)
		}
	}
}

func (c *char) c6DeleteA1OnSwap() bool {
	return c.Base.Cons < 6
}
