package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1ICDkey = "yaoyao-c1-stam-icd"
const c2ICDkey = "yaoyao-c2-icd"
const c6MegaRadishRad = 4.0

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DendroP] = 0.15
	active := c.Core.Player.ActiveChar()
	active.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("yaoyao-c1", 8.8*60),
		AffectedStat: attributes.DendroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	if c.StatusIsActive(c1ICDkey) {
		return
	}
	c.Core.Player.RestoreStam(15)
	c.AddStatus(c1ICDkey, 5*60, false)
}

func (c *char) makeC2CB() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	if !c.StatusIsActive(burstKey) {
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c2ICDkey) {
			return
		}
		c.AddEnergy("yaoyao-c2", 3)
		c.AddStatus(c2ICDkey, 0.8*60, false)
	}
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = c.MaxHP() * 0.003
	if m[attributes.EM] > 120 {
		m[attributes.EM] = 120
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("yaoyao-c4", 8*60),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (yg *yuegui) c6(target geometry.Point, ai combat.AttackInfo, hi player.HealInfo) {
	ai.Abil = "Mega Radish"
	ai.Mult = 0.75
	ai.AttackTag = attacks.AttackTagNone
	ai.ICDTag = attacks.ICDTagNone
	ai.ICDGroup = attacks.ICDGroupDefault

	hi.Message = "Radish C6"
	hi.Src = yg.c.MaxHP() * 0.075

	c6MegaRadishAoE := combat.NewCircleHitOnTarget(target, nil, c6MegaRadishRad)
	yg.Core.QueueAttackWithSnap(ai, yg.snap, c6MegaRadishAoE, c6TravelDelay)
	if yg.Core.Combat.Player().IsWithinArea(c6MegaRadishAoE) {
		yg.c.Core.Player.Heal(hi)
	}
}
