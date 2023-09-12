package itto

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	strStackKey = "strStack"
)

func init() {
	core.RegisterCharFunc(keys.Itto, NewChar)
}

type char struct {
	*tmpl.Character
	savedNormalCounter int
	slashState         SlashType
	a1Stacks           int
	stacksConsumed     int
	burstCastF         int
	c2GeoMemberCount   int
	applyC4            bool
	c6Proc             bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	c.burstCastF = -1
	c.slashState = InvalidSlash

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()               // A1 setup
	c.onExitField()      // burst-exit hook
	c.resetChargeState() // post-charge hook

	if c.Base.Cons >= 2 {
		for _, char := range c.Core.Player.Chars() {
			if char.Base.Element == attributes.Geo {
				c.c2GeoMemberCount++
			}
		}
		if c.c2GeoMemberCount > 3 {
			c.c2GeoMemberCount = 3
		}
	}
	if c.Base.Cons >= 6 {
		c.c6()
		c.c6Proc = c.Base.Cons >= 6 && c.Core.Rand.Float64() < 0.5 // setup c6 proc
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		// CA in Q state don't consume stamina
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
		// apply infusion to attacks only
		switch ai.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagPlunge:
		case attacks.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Geo
	}
	return ds
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
	}, "itto-ca-counter-reset")
}

// used to increment/decrement the amount of Superlative Strength stacks
func (c *char) addStrStack(src string, inc int) {
	old := c.Tags[strStackKey]
	v := old + inc
	if v > 5 {
		v = 5
	} else if v < 0 {
		v = 0
	}

	s := "current"
	if v > old {
		s = "gained"
	} else if v < old {
		s = "lost"
		c.stacksConsumed++
	}
	c.Tags[strStackKey] = v

	c.Core.Log.NewEvent(fmt.Sprintf("itto %v SSS stacks from %v", s, src), glog.LogCharacterEvent, c.Index).
		Write("old_stacks", old).
		Write("inc", inc).
		Write("cur_stacks", v)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "will-c6-proc":
		return c.c6Proc, nil
	case "slash-type":
		if len(fields) < 2 {
			break
		}
		switch fields[1] {
		case "idle":
			return int(InvalidSlash), nil
		case "saichi":
			return int(SaichiSlash), nil
		case "left":
			return int(LeftSlash), nil
		case "right":
			return int(RightSlash), nil
		case "final":
			return int(FinalSlash), nil
		}
	case "slash":
		if len(fields) < 2 {
			break
		}
		switch fields[1] {
		case "current":
			return int(c.slashState), nil
		case "next":
			return int(c.slashState.Next(c.Tags[strStackKey], c.c6Proc)), nil
		}
	}
	return c.Character.Condition(fields)
}
