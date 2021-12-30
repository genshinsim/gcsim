package hutao

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Hutao, NewChar)
}

type char struct {
	*character.Tmpl
	paraParticleICD int
	// chargeICDCounter   int
	// chargeCounterReset int
	ppBonus    float64
	tickActive bool
	c6icd      int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Pyro
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 6
	c.CharZone = core.ZoneLiyue

	c.ppHook()
	c.onExitField()
	c.a4()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 1 {
			return 0
		}
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

func (c *char) a4() {
	c.AddMod(core.CharStatMod{
		Key:    "hutao-a4",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			val[core.PyroP] = 0.33
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			if c.HPCurrent/c.HPMax <= 0.5 {
				return val, true
			}
			return nil, false
		},
	})
}

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		c.checkc6()
		return false
	}, "hutao-c6")
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6icd && c.c6icd != 0 {
		return
	}
	//check if hp less than 25%
	if c.HPCurrent/c.HPMax > .25 {
		return
	}
	//if dead, revive back to 1 hp
	if c.HPCurrent == -1 {
		c.HPCurrent = 1
	}
	//increase crit rate to 100%
	val := make([]float64, core.EndStatType)
	val[core.CR] = 1
	c.AddMod(core.CharStatMod{
		Key:    "hutao-c6",
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 600,
	})

	c.c6icd = c.Core.F + 3600
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("paramita") > 0 {
		switch ai.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = core.Pyro
	}
	return ds
}
