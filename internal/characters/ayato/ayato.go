package ayato

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	stacks            int
	shunsuikenCounter int
}

func init() {
	core.RegisterCharFunc(core.Ayaka, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.CharZone = core.ZoneInazuma
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5
	c.shunsuikenCounter = 3
	c.a2()
	c.a4()
	c.waveFlash()
	c.soukaiKankaHook()

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) a2() {
	c.Core.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.CharIndex() {
			return false
		}
		c.stacks = 2
		c.Core.Log.NewEvent("ayato a2 proc'd", core.LogCharacterEvent, c.Index)
		return false
	}, "ayato-a2")
}

func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.03 * c.MaxHP()
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "ayato-a4",
		Expiry: -1,
		Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
			if a.Info.AttackTag != core.AttackTagElementalBurst {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("soukaikanka") > 0 {
		switch ai.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = core.Hydro
	}
	return ds
}

func (c *char) AdvanceNormalIndex() {
	c.NormalCounter++

	if c.Core.Status.Duration("soukaikanka") > 0 {
		if c.NormalCounter == c.shunsuikenCounter {
			c.NormalCounter = 0
		}
	} else {
		if c.NormalCounter == c.NormalHitNum {
			c.NormalCounter = 0
		}

	}
}
