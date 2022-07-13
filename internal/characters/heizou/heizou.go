package heizou

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Heizou, NewChar)
}

type char struct {
	*character.Tmpl
	decStack            int
	qInfused            core.EleType
	infuseCheckLocation core.AttackPattern
	a1icd               int
	c1icd               int
}

const (
	eCD = 600
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 5
	c.SkillCon = 3
	c.BurstCon = 5
	c.InitCancelFrames()
	c.a1icd = 0
	c.c1icd = 0

	c.infuseCheckLocation = core.NewDefCircHit(0.1, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

	if c.Base.Cons >= 1 {
		c.c1()
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()

}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		return 0
	}
}

func (c *char) a1() {
	addDecStack := func() func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			if c.a1icd > c.Core.F {
				return false
			}

			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			if c.decStack < 4 {
				c.decStack += 1
				c.a1icd = c.Core.F + 0.1*60
				c.Core.Log.NewEvent(
					"declension stack gained",
					core.LogCharacterEvent,
					c.Index,
					"stacks", c.decStack,
				)
			}
			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, addDecStack(), "heizou-a1-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, addDecStack(), "heizou-a1-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, addDecStack(), "heizou-a1-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, addDecStack(), "heizou-a1-pyro")
}

func (c *char) a4() {
	m := make([]float64, core.EndStatType)
	m[core.EM] = 80

	dur := 60 * 10
	c.Core.Status.AddStatus("heizoua4", dur)
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue //nothing for heizou
		}
		char.AddMod(core.CharStatMod{
			Key:    "heizou-a4",
			Expiry: c.Core.F + dur,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	c.Core.Log.NewEvent("heizou a4 triggered", core.LogCharacterEvent, c.Index, "em snapshot", m[core.EM], "expiry", c.Core.F+dur)
}

//For 5s after Shikanoin Heizou takes the field, his Normal Attack SPD is increased by 15%.
//He also gains 1 Declension stack for Heartstopper Strike. This effect can be triggered once every 10s.
func (c *char) c1() {
	// Add hook that monitors for crit hits. Mirrors existing favonius code
	// No log value saved as stat mod already shows up in debug view
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.c1icd > c.Core.F {
			return false
		}
		prev := args[0].(int)
		next := args[1].(int)
		if next == c.Index && prev != c.Index {
			val := make([]float64, core.EndStatType)
			val[core.AtkSpd] = 0.15
			c.AddPreDamageMod(core.PreDamageMod{
				Key:    "heizou-c1",
				Expiry: c.Core.F + 240,
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					if atk.Info.AttackTag != core.AttackTagNormal {
						return nil, false
					}
					return val, true
				},
			})
			if c.decStack < 4 {
				c.decStack++
			}
			c.c1icd = c.Core.F + 600
		}

		return false
	}, "heizou enter")

}
