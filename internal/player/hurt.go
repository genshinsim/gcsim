package player

import "github.com/genshinsim/gcsim/pkg/coretype"

func (p *Player) HurtChar(dmg float64, ele coretype.EleType) {
	//reduce damage by damage reduction first, do so via a hook
	var dr float64
	n := 0
	for _, f := range p.damageReduction {
		amt, done := f()
		dr += amt
		if !done {
			p.damageReduction[n] = f
			n++
		}
	}
	p.damageReduction = p.damageReduction[:n]
	dmg = dmg * (1 - dr)

	//apply damage to all shields
	post := p.OnDamage(dmg, ele)

	//reduce character's hp by damage
	c := p.Chars[p.ActiveChar]
	c.ModifyHP(-post)

	p.core.NewEvent("damage taken", coretype.LogHurtEvent, p.ActiveChar, "dmg", dmg, "taken", post, "shielded", dmg-post, "char_hp", c.HP(), "shield_count", p.ShieldCount())

	if post > 0 {
		p.core.Emit(coretype.OnCharacterHurt, post)
	}
}
