package sucrose

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Sucrose, NewChar)
}

type char struct {
	*character.Tmpl
	qInfused            core.EleType
	infuseCheckLocation core.AttackPattern
	c4Count             int
}

const eCD = 900

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
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 4
	c.InitCancelFrames()

	c.infuseCheckLocation = core.NewDefCircHit(0.1, false, core.TargettableEnemy, core.TargettablePlayer, core.TargettableObject)

	if c.Base.Cons >= 1 {
		c.SetNumCharges(core.ActionSkill, 2)
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		return 0
	}
}

func (c *char) a1() {
	m := make([]float64, core.EndStatType)
	m[core.EM] = 50

	swirlfunc := func(ele core.EleType) func(args ...interface{}) bool {
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

			expiry := c.Core.F + 60*8
			for _, char := range c.Core.Chars {
				this := char
				if this.Ele() != ele {
					continue
				}
				this.AddMod(core.CharStatMod{
					Key:    "sucrose-a1",
					Expiry: expiry,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}

			c.Core.Log.NewEvent("sucrose a1 triggered", core.LogCharacterEvent, c.Index, "reaction", "swirl-"+ele.String(), "expiry", expiry)
			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlfunc(core.Cryo), "sucrose-a1-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlfunc(core.Electro), "sucrose-a1-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlfunc(core.Hydro), "sucrose-a1-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlfunc(core.Pyro), "sucrose-a1-pyro")
}

func (c *char) a4() {
	m := make([]float64, core.EndStatType)
	m[core.EM] = c.Stat(core.EM) * .20

	dur := 60 * 8
	c.Core.Status.AddStatus("sucrosea4", dur)
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddMod(core.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: c.Core.F + dur,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	c.Core.Log.NewEvent("sucrose a4 triggered", core.LogCharacterEvent, c.Index, "em snapshot", m[core.EM], "expiry", c.Core.F+dur)
}

// Handles C4: Every 7 Normal and Charged Attacks, Sucrose will reduce the CD of Astable Anemohypostasis Creation-6308 by 1-7s
func (c *char) c4() {

	c.c4Count++
	if c.c4Count < 7 {
		return
	}
	c.c4Count = 0

	// Change can be in float. See this Terrapin video for example
	// https://youtu.be/jB3x5BTYWIA?t=20
	cdReduction := 60 * int(c.Core.Rand.Float64()*6+1)

	//we simply reduce the action cd
	c.ReduceActionCooldown(core.ActionSkill, cdReduction)

	c.Core.Log.NewEvent("sucrose c4 reducing E CD", core.LogCharacterEvent, c.Index, "cd_reduction", cdReduction)
}

func (c *char) c6() {
	m := make([]float64, core.EndStatType)
	m[core.EleToDmgP(c.qInfused)] = .20

	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "sucrose-c6",
			Expiry: c.Core.F + 60*10,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
