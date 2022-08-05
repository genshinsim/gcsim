package itto

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Itto, NewChar)
}

type char struct {
	*tmpl.Character
	savedNormalCounter int
	chargedCount       int
	burstBuffKey       string
	burstBuffDuration  int
	stackKey           string
	a1Stacks           int
	stacksConsumed     int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	// boilerplate
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 70
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneInazuma

	// needed for NA reset mechanic (Dasshu)
	c.savedNormalCounter = c.NormalCounter

	// needed to determine A1 buff
	c.a1Stacks = 0

	// needed to keep track of Superlative Strength stacks for CAs
	c.stackKey = "strStack"
	c.Tags[c.stackKey] = 0

	// used to determine what CA was used
	// 0 = "Saichimonji Slash"
	// 1 = "Arataki Kesagiri Combo Slash Left"
	// 2 = "Arataki Kesagiri Combo Slash Right"
	// 3 = "Arataki Kesagiri Final Slash"
	c.chargedCount = -1
	c.stacksConsumed = 1

	// used for burst stuff
	c.burstBuffKey = burstBuffKey
	c.burstBuffDuration = 660 + 90 + 45 // barely cover basic combo

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.onExitField()
	if c.Base.Cons >= 6 {
		c.c6ChargedCritDMG()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		// CA in Q state don't consume stamina
		if c.Tags[c.stackKey] > 0 {
			return 0
		}
		return 20
	}
	return c.Character.ActionStam(a, p)
}

// Itto Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.StatModIsActive(c.burstBuffKey) {
		// apply infusion to attacks only
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

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.StatModIsActive(c.burstBuffKey) {
			c.DeleteStatMod(c.burstBuffKey)
			c.DeleteStatMod(c.burstBuffKey + "-atkspd")
		}
		c.a1Stacks = 0
		c.savedNormalCounter = 0
		c.chargedCount = -1
		c.stacksConsumed = 1
		return false
	}, "itto-exit")
}

// used to increment/decrement the amount of Superlative Strength stacks
func (c *char) changeStacks(amount int) {
	c.Tags[c.stackKey] += amount
	if c.Tags[c.stackKey] > 5 {
		c.Tags[c.stackKey] = 5
	}
	if c.Tags[c.stackKey] < 0 {
		c.Tags[c.stackKey] = 0
	}
}
