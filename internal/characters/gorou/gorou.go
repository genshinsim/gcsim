package gorou

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	defenseBuffKey           = "goroubuff"
	generalWarBannerKey      = "generalwarbanner"
	generalGloryKey          = "generalglory"
	generalWarBannerDuration = 600    //10s
	generalGloryDuration     = 9 * 60 //9 s
	heedlessKey              = "heedlessbuff"
	c6key                    = "gorou-c6"
)

func init() {
	core.RegisterCharFunc(keys.Gorou, NewChar)
}

type char struct {
	*tmpl.Character
	eFieldSrc      int
	eFieldHealSrc  int
	qFieldSrc      int
	gorouBuff      []float64
	geoCharCount   int
	c2Extension    int
	c6buff         []float64
	a2buff         []float64
	healFieldStats [attributes.EndStatType]float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = character.ZoneInazuma

	c.c6buff = make([]float64, attributes.EndStatType)
	c.gorouBuff = make([]float64, attributes.EndStatType)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a2buff = make([]float64, attributes.EndStatType)
	c.a2buff[attributes.DEFP] = .25

	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == attributes.Geo {
			c.geoCharCount++
		}
	}

	/**
	Provides up to 3 buffs to active characters within the skill's AoE based on the number of Geo characters in the party at the time of casting:
	• 1 Geo character: Adds "Standing Firm" - DEF Bonus.
	• 2 Geo characters: Adds "Impregnable" - Increased resistance to interruption.
	• 3 Geo characters: Adds "Crunch" - Geo DMG Bonus.
	**/
	c.gorouBuff[attributes.DEF] = skillDefBonus[c.TalentLvlSkill()]
	if c.geoCharCount > 2 {
		c.gorouBuff[attributes.GeoP] = 0.15 // 15% geo damage
	}

	/**
	For 12s after using Inuzaka All-Round Defense or Juuga: Forward Unto Victory, increases the CRIT DMG of all nearby party members' Geo DMG based on the buff level of the skill's field at the time of use:
	• "Standing Firm": +10%
	• "Impregnable": +20%
	• "Crunch": +40%
	This effect cannot stack and will take reference from the last instance of the effect that is triggered.
	**/
	switch c.geoCharCount {
	case 1:
		c.c6buff[attributes.CD] = 0.1
	case 2:
		c.c6buff[attributes.CD] = 0.2
	default:
		//can't be less than 1 so this is 3 or 4
		c.c6buff[attributes.CD] = 0.4
	}

	if c.Base.Cons > 0 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}
