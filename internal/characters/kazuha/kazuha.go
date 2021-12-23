package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Kazuha, NewChar)
}

type char struct {
	*character.Tmpl
	a4Expiry int
	a2Ele    core.EleType
	qInfuse  core.EleType
	c6Active int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5
	c.CharZone = core.ZoneInazuma

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	c.a4()

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
	val := make([]float64, core.EndStatType)
	expiry := make([]int, core.EndStatType)
	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Expiry: -1,
			Key:    "kazuha-a4",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				m := make([]float64, core.EndStatType)
				ok := false
				for i, exp := range expiry {
					if exp > c.Core.F {
						m[i] = val[i]
						ok = true
					}
				}
				if !ok {
					return nil, false
				}
				return m, true
			},
		})
	}

	swirlfunc := func(ele core.StatType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			//update expiry
			expiry[ele] = c.Core.F + 480
			// c.a4Expiry = c.Core.F + 480
			//recalc em
			em := c.Stat(core.EM)
			val[ele] = 0.0004 * em
			c.Core.Log.Debugw("kazuah a4 proc", "frame", c.Core.F, "event", core.LogCharacterEvent, "reaction", ele.String(), "char", c.CharIndex())

			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlfunc(core.CryoP), "kazuha-a4-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlfunc(core.ElectroP), "kazuha-a4-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlfunc(core.HydroP), "kazuha-a4-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlfunc(core.PyroP), "kazuha-a4-pyro")
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Base.Cons < 6 {
		return ds
	}

	if c.c6Active <= c.Core.F {
		return ds
	}

	//infusion to normal/plunge/charge
	switch ai.AttackTag {
	case core.AttackTagNormal:
	case core.AttackTagExtra:
	case core.AttackTagPlunge:
	default:
		return ds
	}
	ai.Element = core.Anemo

	//add 0.2% dmg for every EM
	ds.Stats[core.DmgP] += 0.002 * ds.Stats[core.EM]

	c.Core.Log.Debugw("c6 adding dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "em", ds.Stats[core.EM], "final", ds.Stats[core.DmgP])

	return ds
}
