package mika

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Mika, NewChar)
}

type char struct {
	*tmpl.Character
	maxDetectorStacks int
	healIcd           int
	a4Stack           bool
	c4Count           int

	skillbuff []float64
	c6buff    []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.maxDetectorStacks = 3
	c.healIcd = 2.5 * 60

	c.onBurstHeal()
	if c.Base.Ascension >= 1 {
		c.a1()
	}
	if c.Base.Ascension >= 4 {
		c.a4()
		c.maxDetectorStacks++
	}

	// The Soulwind state of Starfrost Swirl can decrease the healing interval between instances caused by Skyfeather Song's Eagleplume state.
	// This decrease percentage is equal to the ATK SPD increase provided by Soulwind.
	if c.Base.Cons >= 1 {
		c.healIcd = int(float64(c.healIcd) * (1.0 - atkSpdBuff[c.TalentLvlSkill()]))
	}

	c.skillbuff = make([]float64, attributes.EndStatType)
	c.skillbuff[attributes.AtkSpd] = atkSpdBuff[c.TalentLvlSkill()]
	if c.Base.Cons >= 6 {
		c.c6buff = make([]float64, attributes.EndStatType)
		c.c6buff[attributes.CD] = 0.6
		c.maxDetectorStacks++
	}

	return nil
}
