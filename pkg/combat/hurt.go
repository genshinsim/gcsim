package combat

import "github.com/genshinsim/gsim/pkg/def"

func (s *Sim) AddOnHurt(f func(s def.Sim)) {
	s.onHurt = append(s.onHurt, f)
}

func (s *Sim) DamageChar(dmg float64, ele def.EleType) {
	//reduce damage by damage reduction first, do so via a hook
	var dr float64
	for _, f := range s.DRFunc {
		dr += f()
	}
	dmg = dmg * (1 - dr)

	//apply damage to all shields
	post := s.DamageShields(dmg, ele)

	//reduce character's hp by damage
	c := s.chars[s.active]
	c.ModifyHP(-post)

	s.log.Debugw("damage taken", "frame", s.f, "event", def.LogHurtEvent, "frame", s.f, "dmg", dmg, "taken", post, "shielded", dmg-post, "char_hp", c.HP(), "shield_count", len(s.shields))

	if post > 0 {
		for _, f := range s.onHurt {
			f(s)
		}
	}
}
