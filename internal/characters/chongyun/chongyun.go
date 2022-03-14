package chongyun

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Chongyun, NewChar)
}

type char struct {
	*character.Tmpl
	fieldSrc int
	a4Snap   *coretype.AttackEvent
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

	c.onSwapHook()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons == 6 && c.Core.Flags.DamageMode {
		c.c6()
	}

	return &c, nil
}

func (c *char) c4() {
	icd := 0
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		t := args[0].(coretype.Target)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Frame < icd {
			return false
		}
		if !t.AuraContains(coretype.Cryo) {
			return false
		}

		c.AddEnergy("chongyun-c4", 2)

		c.coretype.Log.NewEvent("chongyun c4 recovering 2 energy", coretype.LogCharacterEvent, c.Index, "final energy", c.Energy)
		icd = c.Core.Frame + 120

		return false
	}, "chongyun-c4")

}

func (c *char) onSwapHook() {
	c.Core.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.StatusDuration("chongyunfield") == 0 {
			return false
		}
		//add infusion on swap
		c.coretype.Log.NewEvent("chongyun adding infusion on swap", coretype.LogCharacterEvent, c.Index, "expiry", c.Core.Frame+infuseDur[c.TalentLvlSkill()])
		active := c.Core.Chars[c.Core.ActiveChar]
		c.infuse(active)
		return false
	}, "chongyun-field")
}

func (c *char) infuse(char coretype.Character) {
	switch char.WeaponClass() {
	case core.WeaponClassClaymore:
		fallthrough
	case core.WeaponClassSpear:
		fallthrough
	case core.WeaponClassSword:
		c.coretype.Log.NewEvent("chongyun adding infusion", coretype.LogCharacterEvent, c.Index, "expiry", c.Core.Frame+infuseDur[c.TalentLvlSkill()])
		char.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "chongyun-ice-weapon",
			Ele:    coretype.Cryo,
			Tags:   []core.AttackTag{coretype.AttackTagNormal, coretype.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Core.Frame + infuseDur[c.TalentLvlSkill()],
		})
	default:
		return
	}

	//a2 adds 8% atkspd for 2.1 seconds
	val := make([]float64, core.EndStatType)
	val[core.AtkSpd] = 0.08
	char.AddMod(coretype.CharStatMod{
		Key:    "chongyun-field",
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.Frame + 126,
	})
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
			Expiry: c.Core.Frame + 126,
		})
	}
}

func (c *char) c6() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "chongyun-c6",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {

			val := make([]float64, core.EndStatType)
			if atk.Info.Abil != "Spirit Blade: Cloud-Parting Star" {
				return nil, false
			}
			if t.HP()/t.MaxHP() < c.HPCurrent/c.HPMax {
				val[core.DmgP] += 0.15
				return val, true
			}
			return nil, false
		},
	})
}
