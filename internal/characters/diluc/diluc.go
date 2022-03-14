package diluc

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Diluc, NewChar)
}

type char struct {
	*character.Tmpl
	eStarted    bool
	eStartFrame int
	eLastUse    int
	eCounter    int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Pyro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4

	if c.Base.Cons >= 1 && c.Core.Flags.DamageMode {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

func (c *char) c1() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "diluc-c1",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if t.HP()/t.MaxHP() > 0.5 {
				val[core.DmgP] = 0.15
				return val, true
			}
			return nil, false
		},
	})
}

func (c *char) c2() {
	stack := 0
	last := 0
	c.Core.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if last != 0 && c.Core.Frame-last < 90 {
			return false
		}
		//last time is more than 10 seconds ago, reset stacks back to 0
		if c.Core.Frame-last > 600 {
			stack = 0
		}
		stack++
		if stack > 3 {
			stack = 3
		}
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.1 * float64(stack)
		val[core.AtkSpd] = 0.05 * float64(stack)
		c.AddMod(coretype.CharStatMod{
			Key:    "diluc-c2",
			Amount: func() ([]float64, bool) { return val, true },
			Expiry: c.Core.Frame + 600,
		})
		return false
	}, "diluc-c2")

}

func (c *char) c4() {
	c.AddMod(coretype.CharStatMod{
		Key:    "diluc-c4",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if c.Core.StatusDuration("dilucc4") > 0 {
				val[core.DmgP] = 0.4
				return val, true
			}
			return nil, false
		},
	})
}
func (c *char) Tick() {
	c.Tmpl.Tick()

	if c.eStarted {
		//check if 4 second has passed since last use
		if c.Core.Frame-c.eLastUse >= 240 {
			//if so, set ability to be on cd equal to 10s less started
			cd := 600 - (c.Core.Frame - c.eStartFrame)
			c.coretype.Log.NewEvent("diluc skill going on cd", coretype.LogCharacterEvent, c.Index, "duration", cd, "last", c.eLastUse)
			c.SetCD(core.ActionSkill, cd)
			//reset
			c.eStarted = false
			c.eStartFrame = -1
			c.eLastUse = -1
			c.eCounter = 0
		}
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		// With A1
		return 20
	default:
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}
