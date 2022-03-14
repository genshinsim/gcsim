package kazuha

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Kazuha, NewChar)
}

type char struct {
	*character.Tmpl
	a4Expiry int
	a2Ele    coretype.EleType
	qInfuse  coretype.EleType
	c6Active int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5
	c.CharZone = core.ZoneInazuma

	c.InitCancelFrames()

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.a4()

}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

//Upon triggering a Swirl reaction, Kaedehara Kazuha will grant all party members a 0.04%
//Elemental DMG Bonus to the element absorbed by Swirl for every point of Elemental Mastery
//he has for 8s. Bonuses for different elements obtained through this method can co-exist.
//
//this ignores any EM he gets from Sucrose A4, which is: When Astable Anemohypostasis Creation
//- 6308 or Forbidden Creation - Isomer 75 / Type II hits an opponent, increases all party
//members' (excluding Sucrose) Elemental Mastery by an amount equal to 20% of Sucrose's
//Elemental Mastery for 8s.
//
//he still benefits from sucrose em but just cannot share it

func (c *char) a4() {
	m := make([]float64, core.EndStatType)

	swirlfunc := func(ele core.StatType, key string) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			atk := args[1].(*coretype.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			// do not overwrite mod if same frame
			if c.Core.F < icd {
				return false
			}
			icd = c.Core.F + 1

			//recalc em
			dmg := 0.0004 * c.Stat(core.EM)

			for _, char := range c.Core.Chars {
				char.AddMod(coretype.CharStatMod{
					Key:    "kazuha-a4-" + key,
					Expiry: c.Core.Frame + 60*8,
					Amount: func() ([]float64, bool) {

						m[coretype.CryoP] = 0
						m[core.ElectroP] = 0
						m[core.HydroP] = 0
						m[core.PyroP] = 0

						m[ele] = dmg
						return m, true
					},
				})
			}

			c.coretype.Log.NewEvent("kazuha a4 proc", coretype.LogCharacterEvent, c.Index, "reaction", ele.String(), "char", c.Index())

			return false
		}
	}

	c.Core.Subscribe(coretype.OnSwirlCryo, swirlfunc(coretype.CryoP, "cryo"), "kazuha-a4-cryo")
	c.Core.Subscribe(coretype.OnSwirlElectro, swirlfunc(core.ElectroP, "electro"), "kazuha-a4-electro")
	c.Core.Subscribe(coretype.OnSwirlHydro, swirlfunc(core.HydroP, "hydro"), "kazuha-a4-hydro")
	c.Core.Subscribe(coretype.OnSwirlPyro, swirlfunc(core.PyroP, "pyro"), "kazuha-a4-pyro")
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Base.Cons < 6 {
		return ds
	}

	if c.c6Active <= c.Core.Frame {
		return ds
	}

	//infusion to normal/plunge/charge
	switch ai.AttackTag {
	case coretype.AttackTagNormal:
	case coretype.AttackTagExtra:
	case core.AttackTagPlunge:
	default:
		return ds
	}
	ai.Element = core.Anemo

	//add 0.2% dmg for every EM
	ds.Stats[core.DmgP] += 0.002 * ds.Stats[core.EM]

	c.coretype.Log.NewEvent("c6 adding dmg", coretype.LogCharacterEvent, c.Index, "em", ds.Stats[core.EM], "final", ds.Stats[core.DmgP])

	return ds
}
