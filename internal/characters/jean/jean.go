package jean

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Jean, NewChar)
}

type char struct {
	*tmpl.Character
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 6 {
		c.Core.Log.NewEvent("jean c6 not implemented", glog.LogCharacterEvent, c.Index)
	}
	return nil
}

func (c *char) ReceiveParticle(p character.Particle, isActive bool, partyCount int) {
	c.Character.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 2 {
		//only pop this if jean is active
		if !isActive {
			return
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.15
		for _, active := range c.Core.Player.Chars() {
			active.AddStatMod(character.StatMod{Base: modifier.NewBase("jean-c2", 900), AffectedStat: attributes.AtkSpd, Amount: func() ([]float64, bool) {
				return m, true
			}})
		}
	}
}
