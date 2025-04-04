package varesa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4Status = "diligent-refinement"

func (c *char) c2CB() func(combat.AttackCB) {
	if c.Base.Cons < 2 {
		return nil
	}

	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.AddEnergy("varesa-c2", 11.5)
	}
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 1.0
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("varesa-c4", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst && atk.Info.Abil != kablamAbil {
				return nil, false
			}
			if !c.nightsoulState.HasBlessing() && !c.StatusIsActive(apexState) {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c4Burst() {
	if c.Base.Cons < 4 {
		return
	}
	if c.nightsoulState.HasBlessing() || c.StatusIsActive(apexState) {
		return
	}
	c.AddStatus(c4Status, 15*60, true)
}

func (c *char) c4FlatBonus() float64 {
	if c.Base.Cons < 4 {
		return 0
	}
	if !c.StatusIsActive(c4Status) {
		return 0
	}
	bonus := min(c.TotalAtk()*5, 20000)
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "varesa c4 flat dmg bonus").
		Write("bonus", bonus)
	return bonus
}

func (c *char) c4CB(_ combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
	c.DeleteStatus(c4Status)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1
	m[attributes.CD] = 1.0
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("varesa-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch {
			case atk.Info.AttackTag == attacks.AttackTagElementalBurst:
			case atk.Info.AttackTag == attacks.AttackTagPlunge && atk.Info.Durability > 0: // TODO: collision?
			case atk.Info.Abil == kablamAbil:
			default:
				return nil, false
			}
			return m, true
		},
	})
}
