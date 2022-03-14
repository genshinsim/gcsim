package eula

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Eula, NewChar)
}

type char struct {
	*character.Tmpl
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
	grimheartICD    int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = coretype.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a4()
	c.onExitField()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	s.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if c.Core.StatusDuration("eulaq") == 0 {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.burstCounterICD > c.Core.Frame {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		case coretype.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.coretype.Log.NewEvent("eula burst add stack", coretype.LogCharacterEvent, c.Index, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.coretype.Log.NewEvent("eula c6 add additional stack", coretype.LogCharacterEvent, c.Index, "stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Core.Frame + 6
		return false
	}, "eula-burst-counter")
	return &c, nil
}

func (c *char) a4() {
	c.Core.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.coretype.Log.NewEvent("eula a4 reset skill cd", coretype.LogCharacterEvent, c.Index)
		c.ResetActionCooldown(core.ActionSkill)

		return false
	}, "eula-a4")
}

func (c *char) onExitField() {
	c.Core.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.StatusDuration("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit")
}

func (c *char) c4() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Expiry: -1,
		Key:    "eula-c4",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)

			if atk.Info.Abil != "Glacial Illumination (Lightfall)" {
				return nil, false
			}
			if !c.Core.Flags.DamageMode {
				return nil, false
			}
			if t.HP()/t.MaxHP() >= 0.5 {
				return nil, false
			}
			val[core.DmgP] += 0.25
			return val, true
		},
	})
}

func (e *char) Tick() {
	e.Tmpl.Tick()
	e.grimheartReset--
	if e.grimheartReset == 0 {
		e.Tags["grimheart"] = 0
	}
}
