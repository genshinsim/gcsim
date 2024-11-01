package kinich

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4IcdKey = "kinich-c4-icd-key"

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kinich-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			default:
				return nil, false
			}
			if atk.Info.Abil != c.scaleskiperAttackInfo.Abil {
				return nil, false
			}
			m[attributes.CD] = 1
			return m, true
		},
	})
}

func (c *char) c2ResShredCB(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Damage == 0 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("kinich-c2", 6*60),
		Ele:   attributes.Hydro,
		Value: -0.3,
	})
}

func (c *char) C2Snapshot(dmgBonus int) combat.Snapshot {
	s := c.Snapshot(&c.scaleskiperAttackInfo)
	s.Stats[attributes.DmgP] = float64(dmgBonus)
	if dmgBonus > 0 {
		c.Core.Log.NewEvent("Kinich C2 Damage Bonus", glog.LogCharacterEvent, c.Index).
			Write("bonus", dmgBonus).
			Write("final", s.Stats[attributes.DmgP])
	}
	return s
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		default:
			return false
		}

		if atk.Info.ActorIndex == c.Index {
			return false
		}

		active := c.Core.Player.ActiveChar()
		if active.Index == atk.Info.ActorIndex {
			return false
		}
		if c.StatusIsActive(c4IcdKey) {
			return false
		}
		c.AddStatus(c4IcdKey, 2.8*60, true)
		c.AddEnergy("kinich-c4", 5)

		return false
	}, "kinich-c4-energy")

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kinich-c4-dmgp", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalBurst:
			default:
				return nil, false
			}
			m[attributes.DmgP] = 0.7
			return m, true
		},
	})
}
