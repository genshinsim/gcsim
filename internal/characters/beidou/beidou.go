package beidou

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Beidou, NewChar)
}

type char struct {
	*character.Tmpl
	burstSnapshot core.Snapshot
	burstAtk      *coretype.AttackEvent
	burstSrc      int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 5
	c.CharZone = core.ZoneLiyue

	c.burstProc()
	c.a4()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

/**
Counterattacking with Tidecaller at the precise moment when the character is hit grants the maximum DMG Bonus.

Gain the following effects for 10s after unleashing Tidecaller with its maximum DMG Bonus:
• DMG dealt by Normal and Charged Attacks is increased by 15%. ATK SPD of Normal and Charged Attacks is increased by 15%.
• Greatly reduced delay before unleashing Charged Attacks.

c1
When Stormbreaker is used:
Creates a shield that absorbs up to 16% of Beidou's Max HP for 15s.
This shield absorbs Electro DMG 250% more effectively.

c2
Stormbreaker's arc lightning can jump to 2 additional targets.

c3
Within 10s of taking DMG, Beidou's Normal Attacks gain 20% additional Electro DMG.

c6
During the duration of Stormbreaker, the Electro RES of surrounding opponents is decreased by 15%.
**/

func (c *char) a4() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "beidou-a4",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			mod := make([]float64, core.EndStatType)
			mod[core.DmgP] = .15

			if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
				return mod, false
			}
			if c.Core.StatusDuration("beidoua4") == 0 {
				return mod, false
			}
			return mod, true
		},
	})
}

func (c *char) c4() {
	c.Core.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		c.Core.AddStatus("beidouc4", 600)
		c.coretype.Log.NewEvent("c4 triggered on damage", coretype.LogCharacterEvent, c.Index, "expiry", c.Core.Frame+600)
		return false
	}, "beidouc4")

	mod := make([]float64, core.EndStatType)
	mod[core.DmgP] = .15

	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		t := args[0].(coretype.Target)
		ae := args[1].(*coretype.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}
		if ae.Info.AttackTag != coretype.AttackTagNormal && ae.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Core.StatusDuration("beidouc4") == 0 {
			return false
		}

		c.coretype.Log.NewEvent("c4 proc'd on attack", coretype.LogCharacterEvent, c.Index, "char", c.Index)
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Beidou C4",
			AttackTag:  core.AttackTagNone,
			ICDTag:     core.ICDTagElementalBurst,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Electro,
			Durability: 25,
			Mult:       0.2,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(t.Index(), t.Type()), 0, 1)
		return false
	}, "beidou-c4")

}
