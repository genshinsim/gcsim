package venti

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	qInfuse    core.EleType
	infuseCheckLocation core.AttackPattern
	aiAbsorb   core.AttackInfo
	snapAbsorb core.Snapshot
}

func init() {
	core.RegisterCharFunc(core.Venti, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
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
	c.InitCancelFrames()

	c.infuseCheckLocation = core.NewDefCircHit(0.1, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

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
			Amount: func() ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 600,
		})
		c.Core.Log.NewEvent("c4 - adding anemo bonus", core.LogCharacterEvent, c.Index, "char", c.Index)

	}
}
