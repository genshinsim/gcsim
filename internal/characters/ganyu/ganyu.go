package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("ganyu", NewChar)
}

type char struct {
	*character.Tmpl
	a2expiry int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	//add a2
	var val [core.EndStatType]float64
	val[core.CR] = 0.2
	c.AddMod(core.CharStatMod{
		Key: "ganyu-a2",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			return val, c.a2expiry > c.Core.F && a == core.AttackTagExtra
		},
		Expiry: -1,
	})

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.Tags["last"] = -1
	}

	return &c, nil
}
