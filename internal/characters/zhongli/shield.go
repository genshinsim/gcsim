package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) addJadeShield() {
	shield := shieldBase[c.TalentLvlSkill()] + shieldPer[c.TalentLvlSkill()]*c.HPMax

	c.Core.Shields.Add(c.newShield(shield, 1200))
	c.Tags["shielded"] = 1

	//add resist mod whenever we get a shield
	res := []core.EleType{core.Pyro, core.Hydro, core.Cryo, core.Electro, core.Geo, core.Anemo, core.Physical}

	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.Core.Targets {
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

	c.Core.Log.Debugw("zhongli res shred active", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+1200, "char", c.Index)
}

func (c *char) removeJadeShield() {
	c.Tags["shielded"] = 0
	c.Tags["a2"] = 0
	//deactivate resist mods
	//add resist mod whenever we get a shield
	res := []core.EleType{core.Pyro, core.Hydro, core.Cryo, core.Electro, core.Geo, core.Anemo, core.Physical}
	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.Core.Targets {
			t.RemoveResMod(key)
		}
	}
	c.Core.Log.Debugw("zhongli res shred deactivated", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
}

func (c *char) newShield(base float64, dur int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = core.ShieldZhongliJadeShield
	n.Tmpl.Ele = core.Geo
	n.Tmpl.HP = base
	n.Tmpl.Name = "Zhongli Skill"
	n.Tmpl.Expires = c.Core.F + dur
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

func (s *shd) OnDamage(dmg float64, ele core.EleType, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	//try healing first
	if s.c.Base.Cons == 6 {
		//40% of dmg is converted into healing, but cannot exceed 8% of each char max hp
		//so we have to go through each char one at a time....

		for i, c := range s.c.Core.Chars {
			heal := 0.4 * dmg
			if heal > 0.08*c.MaxHP() {
				heal = 0.08 * c.MaxHP()
			}
			s.c.Core.Health.HealIndex(s.c.Index, i, heal)
			s.c.Core.Log.Debugw("zhongli c6 healing char", "frame", s.c.Core.F, "event", core.LogHurtEvent, "char", c.CharIndex(), "amount", heal, "final", c.HP())
		}
	}
	if !ok {
		s.c.removeJadeShield()
	}
	if s.c.Tags["a2"] < 5 {
		s.c.Tags["a2"]++
	}
	return taken, ok
}
