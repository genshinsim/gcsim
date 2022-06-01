package chongyun

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Chongyun, NewChar)
}

type char struct {
	*character.Tmpl
	fieldSrc int
	a4Snap   *core.AttackEvent
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	c.fieldSrc = -601

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.onSwapHook()

	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons == 6 && c.Core.Flags.DamageMode {
		c.c6()
	}
}

func (c *char) c4() {
	icd := 0
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.F < icd {
			return false
		}
		if !t.AuraContains(core.Cryo) {
			return false
		}

		c.AddEnergy("chongyun-c4", 2)

		c.Core.Log.NewEvent("chongyun c4 recovering 2 energy", core.LogCharacterEvent, c.Index, "final energy", c.Energy)
		icd = c.Core.F + 120

		return false
	}, "chongyun-c4")

}

func (c *char) onSwapHook() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("chongyunfield") == 0 {
			return false
		}
		//add infusion on swap
		c.Core.Log.NewEvent("chongyun adding infusion on swap", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
		active := c.Core.Chars[c.Core.ActiveChar]
		c.infuse(active)
		return false
	}, "chongyun-field")
}

func (c *char) infuse(char core.Character) {
	//c2 reduces CD by 15%
	if c.Base.Cons >= 2 {
		char.AddCDAdjustFunc(core.CDAdjust{
			Key: "chongyun-c2",
			Amount: func(a core.ActionType) float64 {
				if a == core.ActionSkill || a == core.ActionBurst {
					return -0.15
				}
				return 0
			},
			Expiry: c.Core.F + 126,
		})
	}

	// weapon infuse
	switch char.WeaponClass() {
	case core.WeaponClassClaymore:
		fallthrough
	case core.WeaponClassSpear:
		fallthrough
	case core.WeaponClassSword:
		c.Core.Log.NewEvent("chongyun adding infusion", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
		char.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "chongyun-ice-weapon",
			Ele:    core.Cryo,
			Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Core.F + infuseDur[c.TalentLvlSkill()],
		})
	default:
		return
	}

	//a1 adds 8% atkspd for 2.1 seconds
	val := make([]float64, core.EndStatType)
	val[core.AtkSpd] = 0.08
	char.AddMod(core.CharStatMod{
		Key:    "chongyun-field",
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 126,
	})
}

func (c *char) c6() {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.15
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "chongyun-c6",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagElementalBurst {
				return nil, false
			}
			if t.HP()/t.MaxHP() < c.HP()/c.MaxHP() {
				return m, true
			}
			return nil, false
		},
	})
}
