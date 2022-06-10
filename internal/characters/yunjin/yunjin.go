package yunjin

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Yunjin, NewChar)
}

type char struct {
	*tmpl.Character
	burstTriggers       [4]int
	partyElementalTypes int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	c.partyElementalTypes = 0
	for i := range c.burstTriggers {
		c.burstTriggers[i] = 30
	}

	w.Character = &c

	return nil
}

// Occurs after all characters are loaded, so getPartyElementalTypeCounts works properly
func (c *char) Init() error {
	c.getPartyElementalTypeCounts()
	c.burstProc()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

// Adds event to get the number of elemental types in the party for Yunjin A4
func (c *char) getPartyElementalTypeCounts() {
	partyElementalTypes := make(map[attributes.Element]int)
	for _, char := range c.Core.Player.Chars() {
		partyElementalTypes[char.Base.Element]++
	}
	for range partyElementalTypes {
		c.partyElementalTypes += 1
	}
	c.Core.Log.NewEvent("Yun Jin Party Elemental Types (A4)", glog.LogCharacterEvent, c.Index, "party_elements", c.partyElementalTypes)
}
