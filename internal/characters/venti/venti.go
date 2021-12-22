package venti

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type char struct {
	*character.Tmpl
	qInfuse core.EleType
}

func init() {
	core.RegisterCharFunc(keys.Venti, NewChar)
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
		c.AddMod(core.CharStatMod{
			Key:    "venti-c4",
			Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 600,
		})
		c.Core.Log.Debugw("c4 - adding anemo bonus", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)

	}
}
