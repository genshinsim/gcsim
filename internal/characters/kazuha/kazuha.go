package kazuha

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Kazuha, NewChar)
}

type char struct {
	*character.Tmpl
	a4Expiry            int
	a1Ele               core.EleType
	qInfuse             core.EleType
	c6Active            int
	infuseCheckLocation core.AttackPattern
}

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
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5
	c.CharZone = core.ZoneInazuma

	c.infuseCheckLocation = core.NewDefCircHit(1.5, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
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
			atk := args[1].(*core.AttackEvent)
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
				char.AddMod(core.CharStatMod{
					Key:    "kazuha-a4-" + key,
					Expiry: c.Core.F + 60*8,
					Amount: func() ([]float64, bool) {

						m[core.CryoP] = 0
						m[core.ElectroP] = 0
						m[core.HydroP] = 0
						m[core.PyroP] = 0

						m[ele] = dmg
						return m, true
					},
				})
			}

			c.Core.Log.NewEvent("kazuha a4 proc", core.LogCharacterEvent, c.Index, "reaction", ele.String(), "char", c.CharIndex())

			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlfunc(core.CryoP, "cryo"), "kazuha-a4-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlfunc(core.ElectroP, "electro"), "kazuha-a4-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlfunc(core.HydroP, "hydro"), "kazuha-a4-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlfunc(core.PyroP, "pyro"), "kazuha-a4-pyro")
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Base.Cons < 6 {
		return ds
	}

	if c.c6Active <= c.Core.F {
		return ds
	}

	//add 0.2% dmg for every EM
	ds.Stats[core.DmgP] += 0.002 * ds.Stats[core.EM]

	c.Core.Log.NewEvent("c6 adding dmg", core.LogCharacterEvent, c.Index, "em", ds.Stats[core.EM], "final", ds.Stats[core.DmgP])

	return ds
}
