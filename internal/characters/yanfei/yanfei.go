package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("yanfei", NewChar)
}

type char struct {
	*character.Tmpl
	maxTags           int
	sealStamReduction float64
	sealExpiry        int
	burstBuffExpiry   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 3

	c.onExitField()
	c.a4()

	c.maxTags = 3
	if c.Base.Cons == 6 {
		c.maxTags = 4
	}

	c.sealStamReduction = 0.15
	if c.Base.Cons > 0 {
		c.sealStamReduction = 0.25
	}

	return &c, nil
}

// Hook that clears yanfei burst status and seals when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Tags["seal"] = 0
		c.sealExpiry = c.Core.F - 1
		c.Core.Status.DeleteStatus("yanfeiburst")
		return false
	}, "yanfei-exit")
}

// Hook for C2:
// Increases Yan Fei's Charged Attack CRIT Rate by 20% against enemies below 50% HP.
func (c *char) c2() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		target := args[0].(core.Target)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if target.HP()/target.MaxHP() >= .5 {
			return false
		}
		ds.Stats[core.CR] += 0.20

		c.Core.Log.Debugw("yanfei c2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "target", target.Index(), "target_hp_percent", target.HP()/target.MaxHP())

		return false
	}, "yanfei-c2")
}

// A4 Hook
// When Yan Fei's Charged Attacks deal CRIT Hits, she will deal an additional instance of AoE Pyo DMG equal to 80% of her ATK. This DMG counts as Charged Attack DMG.
func (c *char) a4() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		crit := args[3].(bool)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.Abil == "Blazing Eye (A4)" {
			return false
		}
		if !((ds.AttackTag == core.AttackTagExtra) && crit) {
			return false
		}

		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Blazing Eye (A4)",
				core.AttackTagExtra,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeBlunt,
				core.Pyro,
				25,
				.8,
			)
			return &d
		}, 1)

		return false
	}, "yanfei-a4")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.F > c.sealExpiry {
			c.Tags["seal"] = 0
		}
		stacks := c.Tags["seal"]
		return 50 * (1 - c.sealStamReduction*float64(stacks))
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
