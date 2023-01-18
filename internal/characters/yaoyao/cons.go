package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1ICDkey = "yaoyao-c1-stam-icd"
const c2ICDkey = "yaoyao-c2-icd"
const c6megaRadishRad = 3.0

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DendroP] = 0.15
	active := c.Core.Player.ActiveChar()
	active.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("yaoyao-c1", 8),
		AffectedStat: attributes.DendroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	if c.StatusIsActive(c1ICDkey) {
		return
	}
	c.Core.Player.RestoreStam(15)
	c.AddStatus(c1ICDkey, 5*60-1, false)
}

func (c *char) c2() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	if c.StatusIsActive(c2ICDkey) {
		return
	}
	c.AddEnergy("yaoyao-c2", 3)
	c.AddStatus(c2ICDkey, 0.8*60, false)
}
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = c.MaxHP() * 0.003
	if m[attributes.EM] > 120 {
		m[attributes.EM] = 120
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("yaoyao-c4", 8),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (yg *yuegui) c6(ai combat.AttackInfo, hi player.HealInfo, radishRad float64) (combat.AttackInfo, player.HealInfo, float64) {
	if yg.GadgetTyp() == combat.GadgetTypYueguiThrowing && (yg.throwCounter == 2 || yg.throwCounter == 5) {
		ai.Abil = "Mega Radish"
		ai.Mult = 0.75

		hi.Src = yg.c.MaxHP() * 0.075

		radishRad = c6megaRadishRad
	}
	return ai, hi, radishRad
}
