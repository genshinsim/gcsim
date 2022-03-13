package player

import "github.com/genshinsim/gcsim/pkg/coretype"

type stamMod struct {
	f   func(a coretype.ActionType) (float64, bool)
	key string
}

func (p *Player) AddStamMod(f func(a coretype.ActionType) (float64, bool), key string) {
	ind := -1
	for i, v := range p.stamModifier {
		if key == v.key {
			ind = i
		}
	}
	if ind > -1 {
		p.core.NewEvent("char stam mod replaced", coretype.LogCharacterEvent, -1, "overwrite", true, "key", key)
		// c.Log.Debugw("char stam mod replaced", "frame", c.F, "event", LogCharacterEvent, "overwrite", true, "key", key)
		p.stamModifier[ind].f = f
		p.stamModifier[ind].key = key
	} else {
		p.core.NewEvent("char stam mod added", coretype.LogCharacterEvent, -1, "overwrite", false, "key", key)
		// c.Log.Debugw("char stam mod added", "frame", c.F, "event", LogCharacterEvent, "overwrite", false, "key", key)
		p.stamModifier = append(p.stamModifier, stamMod{
			f:   f,
			key: key,
		})
	}
}

func (p *Player) StamPercentMod(a coretype.ActionType) float64 {
	var m float64 = 1
	n := 0
	for _, mod := range p.stamModifier {
		v, done := mod.f(a)
		if !done {
			p.stamModifier[n] = mod
			n++
		}
		m += v
	}
	p.stamModifier = p.stamModifier[:n]
	return m
}

func (p *Player) RestoreStam(v float64) {
	p.Stam += v
	if p.Stam > MaxStam {
		p.Stam = MaxStam
	}
}
