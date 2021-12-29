package beidou

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Beidou, NewChar)
}

type char struct {
	*character.Tmpl
	burstSnapshot core.Snapshot
	burstAtk      *core.AttackEvent
	burstSrc      int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro
	c.Energy = 80
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
	c.AddMod(core.CharStatMod{
		Key:    "beidou-a4",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			mod := make([]float64, core.EndStatType)
			mod[core.DmgP] = .15

			if a != core.AttackTagNormal && a != core.AttackTagExtra {
				return mod, false
			}
			if c.Core.Status.Duration("beidoua4") == 0 {
				return mod, false
			}
			return mod, true
		},
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		c.Core.Status.AddStatus("beidouc4", 600)
		c.Core.Log.Debugw("c4 triggered on damage", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+600)
		return false
	}, "beidouc4")

	mod := make([]float64, core.EndStatType)
	mod[core.DmgP] = .15

	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ae := args[1].(*core.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("beidouc4") == 0 {
			return false
		}

		c.Core.Log.Debugw("c4 proc'd on attack", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
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
