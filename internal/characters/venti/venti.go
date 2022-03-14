package venti

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

type char struct {
	*character.Tmpl
	qInfuse coretype.EleType
}

func init() {
	core.RegisterCharFunc(core.Venti, NewChar)
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 6
	c.BurstCon = 3
	c.SkillCon = 5

	return &c, nil
}

func (c *char) ReceiveParticle(p core.Particle, isActive bool, partyCount int) {
	c.Tmpl.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 4 {
		//only pop this if active
		if !isActive {
			return
		}

		val := make([]float64, core.EndStatType)
		val[core.AnemoP] = 0.25
		c.AddMod(coretype.CharStatMod{
			Key:    "venti-c4",
			Amount: func() ([]float64, bool) { return val, true },
			Expiry: c.Core.Frame + 600,
		})
		c.coretype.Log.NewEvent("c4 - adding anemo bonus", coretype.LogCharacterEvent, c.Index, "char", c.Index)

	}
}
