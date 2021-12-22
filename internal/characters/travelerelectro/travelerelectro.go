package travelerelectro

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(keys.TravelerElectro, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	return &c, nil
}
