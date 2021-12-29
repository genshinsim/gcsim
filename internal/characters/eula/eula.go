package eula

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Eula, NewChar)
}

type char struct {
	*character.Tmpl
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
	grimheartICD    int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo
	c.Energy = 80
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

	s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if c.Core.Status.Duration("eulaq") == 0 {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.burstCounterICD > c.Core.F {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		case core.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.Core.Log.Debugw("eula burst add stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.Core.Log.Debugw("eula c6 add additional stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Core.F + 6
		return false
	}, "eula-burst-counter")
	return &c, nil
}

func (c *char) a4() {
	c.Core.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.Core.Log.Debugw("eula a4 reset skill cd", "frame", c.Core.F, "event", core.LogCharacterEvent)
		c.ResetActionCooldown(core.ActionSkill)

		return false
	}, "eula-a4")
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit")
}

func (c *char) c4() {
	c.AddPreDamageMod(core.PreDamageMod{
		Expiry: -1,
		Key:    "eula-c4",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
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
