package kazuha

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("kazuha", NewChar)
}

type char struct {
	*character.Tmpl
	a4Expiry int
	a2Ele    core.EleType
	qInfuse  core.EleType
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

	c.a4()

	return &c, nil
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
	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Expiry: -1,
			Key:    "kazuha-a2",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.a4Expiry < c.Core.F {
					return nil, false
				}
				return val, true
			},
		})
	}
	c.Core.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		var typ core.EleType
		switch ds.ReactionType {
		case core.SwirlCryo:
			typ = core.Cryo
		case core.SwirlElectro:
			typ = core.Electro
		case core.SwirlHydro:
			typ = core.Hydro
		case core.SwirlPyro:
			typ = core.Pyro
		default:
			return false
		}
		//update expiry
		c.a4Expiry = c.Core.F + 480
		//recalc em
		em := c.Stat(core.EM)
		val[core.EleToDmgP(typ)] = 0.0004 * em
		return false
	}, "kazuha-a2")
}
