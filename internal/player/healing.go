package player

import "github.com/genshinsim/gcsim/pkg/coretype"

func (p *Player) HealIndex(hi *coretype.HealInfo, index int) {
	c := p.Chars[index]

	bonus := p.healBonusMult(index) + hi.Bonus
	hp := .0
	switch hi.Type {
	case coretype.HealTypeAbsolute:
		hp = hi.Src
	case coretype.HealTypePercent:
		hp = c.MaxHP() * hi.Src
	}
	heal := hp * bonus

	c.ModifyHP(heal)
	p.core.Emit(coretype.OnHeal, hi.Caller, index, heal)
	p.core.NewEvent(hi.Message, coretype.LogHealEvent, index, "amount", hp, "bonus", bonus, "final", c.HP())
}

func (p *Player) Heal(hi coretype.HealInfo) {
	if hi.Target == -1 { // all
		for i := range p.Chars {
			p.HealIndex(&hi, i)
		}
	} else {
		p.HealIndex(&hi, hi.Target)
	}
}

func (p *Player) healBonusMult(healedCharIndex int) float64 {
	var sum float64 = 1
	for _, f := range p.healingBonus {
		sum += f(healedCharIndex)
	}
	return sum
}

func (p *Player) AddIncHealBonus(f func(healedCharIndex int) float64) {
	p.healingBonus = append(p.healingBonus, f)
}

func (p *Player) AddDamageReduction(f func() (float64, bool)) {
	p.damageReduction = append(p.damageReduction, f)
}
