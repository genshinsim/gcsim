package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) addJadeShield() {
	shield := shieldBase[c.TalentLvlSkill()] + shieldPer[c.TalentLvlSkill()]*c.HPMax

	c.Core.Shields.Add(c.newShield(shield, 1200))
	c.Tags["shielded"] = 1

	//add resist mod whenever we get a shield
	res := []coretype.EleType{core.Pyro, core.Hydro, coretype.Cryo, core.Electro, core.Geo, core.Anemo, core.Physical}

	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.coretype.Targets {
			t.AddResMod(
				key,
				core.ResistMod{
					Ele:      v,
					Value:    -0.2,
					Duration: 1200,
				},
			)
		}
	}

	c.coretype.Log.NewEvent("zhongli res shred active", coretype.LogCharacterEvent, c.Index, "expiry", c.Core.Frame+1200, "char", c.Index)
}

func (c *char) removeJadeShield() {
	c.Tags["shielded"] = 0
	c.Tags["a2"] = 0
	//deactivate resist mods
	//add resist mod whenever we get a shield
	res := []coretype.EleType{core.Pyro, core.Hydro, coretype.Cryo, core.Electro, core.Geo, core.Anemo, core.Physical}
	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.coretype.Targets {
			t.RemoveResMod(key)
		}
	}
	c.coretype.Log.NewEvent("zhongli res shred deactivated", coretype.LogCharacterEvent, c.Index, "char", c.Index)
}

func (c *char) newShield(base float64, dur int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.Frame
	n.Tmpl.ShieldType = core.ShieldZhongliJadeShield
	n.Tmpl.Ele = core.Geo
	n.Tmpl.HP = base
	n.Tmpl.Name = "Zhongli Skill"
	n.Tmpl.Expires = c.Core.Frame + dur
	n.c = c
	return n
}

type shd struct {
	*shield.Tmpl
	c *char
}

func (s *shd) OnExpire() {
	s.c.removeJadeShield()
}

func (s *shd) OnDamage(dmg float64, ele coretype.EleType, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	//try healing first
	if s.c.Base.Cons == 6 {
		//40% of dmg is converted into healing, but cannot exceed 8% of each char max hp
		//so we have to go through each char one at a time....

		c := s.c.Core.Chars[s.c.Core.ActiveChar]
		heal := 0.4 * dmg
		if heal > 0.08*c.MaxHP() {
			heal = 0.08 * c.MaxHP()
		}
		s.c.Core.Health.Heal(core.HealInfo{
			Caller:  s.c.Index,
			Target:  s.c.Core.ActiveChar,
			Message: "Chrysos, Bounty of Dominator",
			Src:     heal,
			Bonus:   c.Stat(core.Heal),
		})
	}
	if !ok {
		s.c.removeJadeShield()
	}
	if s.c.Tags["a2"] < 5 {
		s.c.Tags["a2"]++
	}
	return taken, ok
}
