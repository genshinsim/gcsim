package ganyu

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Ganyu, NewChar)
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
	c.Base.Element = core.Cryo

	e, ok := p.Params["energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	//add a2
	val := make([]float64, core.EndStatType)
	val[core.CR] = 0.2
	c.AddPreDamageMod(core.PreDamageMod{
		Key: "ganyu-a2",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			return val, c.a2expiry > c.Core.F && atk.Info.AttackTag == core.AttackTagExtra
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
