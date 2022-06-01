package yanfei

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Yanfei, NewChar)
}

type char struct {
	*character.Tmpl
	maxTags           int
	sealStamReduction float64
	sealExpiry        int
	// burstBuffExpiry   int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Pyro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 3

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

func (c *char) Init() {
	c.Tmpl.Init()

	c.onExitField()
	c.a4()

	if c.Base.Cons >= 2 {
		c.c2()
	}
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
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "yanfei-c2",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if atk.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			if t.HP()/t.MaxHP() >= .5 {
				return nil, false
			}
			m[core.CR] = 0.20
			return m, true
		},
	})
}

// A4 Hook
// When Yan Fei's Charged Attacks deal CRIT Hits, she will deal an additional instance of AoE Pyo DMG equal to 80% of her ATK. This DMG counts as Charged Attack DMG.
func (c *char) a4() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil == "Blazing Eye (A4)" {
			return false
		}
		if !((atk.Info.AttackTag == core.AttackTagExtra) && crit) {
			return false
		}

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blazing Eye (A4)",
			AttackTag:  core.AttackTagExtra,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       0.8,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 1, 1)

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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
