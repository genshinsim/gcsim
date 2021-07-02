package combat

import "github.com/genshinsim/gsim/pkg/def"

func (s *Sim) HealActive(hp float64) {
	s.chars[s.active].ModifyHP(s.healBonusMult() * hp)
	s.log.Debugw("healing", "frame", s.f, "event", def.LogHealEvent, "frame", s.f, "char", s.active, "amount", hp, "bonus", s.healBonusMult(), "final", s.chars[s.active].HP())
}

func (s *Sim) HealAll(hp float64) {
	for i, c := range s.chars {
		c.ModifyHP(s.healBonusMult() * hp)
		s.log.Debugw("healing (all)", "frame", s.f, "event", def.LogHealEvent, "frame", s.f, "char", i, "amount", hp, "bonus", s.healBonusMult(), "final", s.chars[s.active].HP())
	}
}

func (s *Sim) HealAllPercent(percent float64) {
	for i, c := range s.chars {
		hp := c.MaxHP() * percent
		c.ModifyHP(s.healBonusMult() * hp)
		s.log.Debugw("healing (all)", "frame", s.f, "event", def.LogHealEvent, "frame", s.f, "char", i, "amount", hp, "bonus", s.healBonusMult(), "final", s.chars[s.active].HP())
	}
}

func (s *Sim) HealIndex(index int, hp float64) {
	s.chars[index].ModifyHP(s.healBonusMult() * hp)
	s.log.Debugw("healing", "frame", s.f, "event", def.LogHealEvent, "frame", s.f, "char", index, "amount", hp, "bonus", s.healBonusMult(), "final", s.chars[s.active].HP())
}

func (s *Sim) healBonusMult() float64 {
	var sum float64 = 1
	for _, f := range s.IncHealBonusFunc {
		sum += f()
	}
	return sum
}

func (s *Sim) AddIncHealBonus(f func() float64) {
	s.IncHealBonusFunc = append(s.IncHealBonusFunc, f)
}
