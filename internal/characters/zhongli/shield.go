package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) addJadeShield() {
	shield := shieldBase[c.TalentLvlSkill()] + shieldPer[c.TalentLvlSkill()]*c.MaxHP()

	c.Core.Player.Shields.Add(c.newShield(shield, 1200))
	c.Tags["shielded"] = 1

	//add resist mod whenever we get a shield
	res := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro, attributes.Geo, attributes.Anemo, attributes.Physical}

	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.Core.Combat.Targets() {
			e, ok := t.(core.Enemy)
			if !ok {
				continue
			}
			e.AddResistMod(key, 1200, v, -0.2)
		}
	}

	c.Core.Log.NewEvent("zhongli res shred active", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+1200, "char", c.Index)
}

func (c *char) removeJadeShield() {
	c.Tags["shielded"] = 0
	c.Tags["a1"] = 0
	//deactivate resist mods
	//add resist mod whenever we get a shield
	res := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro, attributes.Geo, attributes.Anemo, attributes.Physical}
	for _, v := range res {
		key := fmt.Sprintf("zhongli-%v", v.String())
		for _, t := range c.Core.Combat.Targets() {
			e, ok := t.(core.Enemy)
			if !ok {
				continue
			}
			e.DeleteResistMod(key)
		}
	}
	c.Core.Log.NewEvent("zhongli res shred deactivated", glog.LogCharacterEvent, c.Index, "char", c.Index)
}

func (c *char) newShield(base float64, dur int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = shield.ShieldZhongliJadeShield
	n.Tmpl.Ele = attributes.Geo
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

func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	//try healing first
	if s.c.Base.Cons >= 6 {
		//40% of dmg is converted into healing, but cannot exceed 8% of each char max hp
		//so we have to go through each char one at a time....

		active := s.c.Core.Player.ActiveChar()
		heal := 0.4 * dmg
		maxhp := s.c.MaxHP()
		if heal > 0.08*maxhp {
			heal = 0.08 * maxhp
		}
		s.c.Core.Player.Heal(player.HealInfo{
			Caller:  s.c.Index,
			Target:  active.Index,
			Message: "Chrysos, Bounty of Dominator",
			Src:     heal,
			Bonus:   s.c.Stat(attributes.Heal),
		})
	}
	if !ok {
		s.c.removeJadeShield()
	}
	if s.c.Tags["a1"] < 5 {
		s.c.Tags["a1"]++
	}
	return taken, ok
}
