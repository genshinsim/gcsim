package player

import "github.com/genshinsim/gcsim/pkg/coretype"

func (p *Player) ShieldCount() int {
	return len(p.shields)
}

func (p *Player) IsCharShielded(char int) bool {
	return len(p.shields) > 0 && char == p.ActiveChar
}

func (p *Player) GetShield(t coretype.ShieldType) coretype.Shield {
	for _, v := range p.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

func (p *Player) AddShieldBonus(f func() float64) {
	p.shieldBonusFunc = append(p.shieldBonusFunc, f)
}

func (p *Player) Add(shd coretype.Shield) {
	//we always assume over write of the same type
	ind := -1
	for i, v := range p.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind > -1 {
		p.core.NewEvent("shield overridden", coretype.LogShieldEvent, -1, "overwrite", true, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
		p.shields[ind].OnOverwrite()
		p.shields[ind] = shd
	} else {
		p.shields = append(p.shields, shd)
		p.core.NewEvent("shield added", coretype.LogShieldEvent, -1, "overwrite", false, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
	}
	p.core.Emit(coretype.OnShielded, shd)
}

func (p *Player) OnDamage(dmg float64, ele coretype.EleType) float64 {
	var bonus float64
	//find shield bonuses
	for _, f := range p.shieldBonusFunc {
		bonus += f()
	}
	min := dmg //min of damage taken
	n := 0
	for _, v := range p.shields {
		taken, ok := v.OnDamage(dmg, ele, bonus)
		if taken < min {
			min = taken
		}
		if ok {
			p.shields[n] = v
			n++
		}
	}
	p.shields = p.shields[:n]
	return min
}
