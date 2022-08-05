package itto

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	strStackKey = "strStack"
)

func init() {
	core.RegisterCharFunc(keys.Itto, NewChar)
}

type char struct {
	*tmpl.Character
	dasshuCount    int
	geoCharCount   int
	slashState     SlashType
	applyC4        bool
	burstCastF     int
	a1Stacks       int
	stacksConsumed int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 70
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.CharZone = character.ZoneInazuma

	c.burstCastF = -1
	c.slashState = InvalidSlash

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.onExitField()
	c.resetChargeState()
	if c.Base.Cons >= 2 {
		for _, char := range c.Core.Player.Chars() {
			if char.Base.Element == attributes.Geo {
				c.geoCharCount++
			}
		}
		if c.geoCharCount > 3 {
			c.geoCharCount = 3
		}
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.Tags[strStackKey] > 0 {
			return 0
		}
		return 20
	}
	return c.Character.ActionStam(a, p)
}

// Itto Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)
	if c.StatModIsActive(burstBuffKey) {
		//infusion to attacks only
		switch ai.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagPlunge:
		case combat.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Geo
	}
	return ds
}

func (c *char) addStrStack(inc int) {
	old := c.Tags[strStackKey]
	v := old + inc
	if v > 5 {
		v = 5
	} else if v < 0 {
		v = 5
	}
	c.Tags[strStackKey] = v

	if v != old {
		var s string
		if v > old {
			s = "gained"
		} else if v < old {
			s = "lost"
			c.stacksConsumed++
		}
		c.Core.Log.NewEvent(fmt.Sprintf("itto %v Superlative Superstrength stacks", s), glog.LogCharacterEvent, c.Index).
			Write("old_stacks", old).
			Write("cur_stacks", v)
	}
}

func (c *char) resetChargeState() {
	c.Core.Events.Subscribe(event.OnActionExec, func(args ...interface{}) bool {
		act := args[1].(action.Action)

		if act != action.ActionCharge {
			c.slashState = InvalidSlash
			c.a1Stacks = 0
			c.stacksConsumed = 0
		}

		return false
	}, "itto-na-ca-counter-rest")
}
