package gorou

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Gorou, NewChar)
}

type char struct {
	*character.Tmpl
	eFieldSrc     int
	eFieldHealSrc int
	qFieldSrc     int
	gorouBuff     []float64
	geoCharCount  int
	c2Extension   int
	c6buff        []float64
}

const (
	defenseBuffKey           = "goroubuff"
	generalWarBannerKey      = "generalwarbanner"
	generalGloryKey          = "generalglory"
	generalWarBannerDuration = 600    //10s
	generalGloryDuration     = 9 * 60 //9 s
	heedlessKey              = "headlessbuff"
	c6key                    = "gorou-c6"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 4
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneInazuma

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) Init() {
	c.Tmpl.Init()

	/**
	Provides up to 3 buffs to active characters within the skill's AoE based on the number of Geo characters in the party at the time of casting:
	• 1 Geo character: Adds "Standing Firm" - DEF Bonus.
	• 2 Geo characters: Adds "Impregnable" - Increased resistance to interruption.
	• 3 Geo characters: Adds "Crunch" - Geo DMG Bonus.
	**/

	for _, char := range c.Core.Chars {
		if char.Ele() == core.Geo {
			c.geoCharCount++
		}
	}

	c.gorouBuff = make([]float64, core.EndStatType)
	c.gorouBuff[core.DEF] = skillDefBonus[c.TalentLvlSkill()]
	if c.geoCharCount > 2 {
		c.gorouBuff[core.GeoP] = 0.15 // 15% geo damage
	}
	/**
	For 12s after using Inuzaka All-Round Defense or Juuga: Forward Unto Victory, increases the CRIT DMG of all nearby party members' Geo DMG based on the buff level of the skill's field at the time of use:
	• "Standing Firm": +10%
	• "Impregnable": +20%
	• "Crunch": +40%
	This effect cannot stack and will take reference from the last instance of the effect that is triggered.
	**/
	c.c6buff = make([]float64, core.EndStatType)
	switch c.geoCharCount {
	case 1:
		c.c6buff[core.CD] = 0.1
	case 2:
		c.c6buff[core.CD] = 0.2
	default:
		//can't be less than 1 so this is 3 or 4
		c.c6buff[core.CD] = 0.4
	}

	if c.Base.Cons > 0 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
}
