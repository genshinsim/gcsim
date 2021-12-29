package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Sucrose, NewChar)
}

type char struct {
	*character.Tmpl
	a4EM []float64
	// a4snap   core.Snapshot
	qInfused core.EleType
	//charges
	eChargeMax int
	eCharges   int

	c4Count int
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
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 4
	c.a4EM = make([]float64, core.EndStatType)

	c.eChargeMax = 1
	if c.Base.Cons >= 1 {
		c.eChargeMax = 2
	}
	c.eCharges = c.eChargeMax

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	c.a2()
	c.a4()

	if c.Base.Cons == 6 {
		c.c6()
	}

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

func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	val[core.EM] = 50
	for _, char := range c.Core.Chars {
		this := char
		if this.Ele() == core.Anemo || this.Ele() == core.Geo {
			continue //nothing for geo/anemo char
		}
		this.AddMod(core.CharStatMod{
			Key:    "sucrose-a2",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				var f int
				var ok bool

				// c.Core.Log.Debugw("sucrose a2 check", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", this.CharIndex(), "ele", this.Ele())
				switch this.Ele() {
				case core.Pyro:
					f, ok = c.Tags["a2-pyro"]
				case core.Cryo:
					f, ok = c.Tags["a2-cryo"]
				case core.Hydro:
					f, ok = c.Tags["a2-hydro"]
				case core.Electro:
					f, ok = c.Tags["a2-electro"]
				default:
					return nil, false
				}
				// c.Core.Log.Debugw("sucrose a2 adding", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", this.CharIndex(), "ele", this.Ele(), "expiry", f, "ok", ok)
				return val, f > c.Core.F && ok
			},
		})
	}

	swirlfunc := func(tag string) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			//TODO: not sure if sucrose needs to be active
			c.Tags[tag] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "reaction", tag, "expiry", c.Core.F+480)
			return false
		}
	}

	c.Core.Events.Subscribe(core.OnSwirlCryo, swirlfunc("a2-cryo"), "a2-cryo")
	c.Core.Events.Subscribe(core.OnSwirlElectro, swirlfunc("a2-electro"), "a2-electro")
	c.Core.Events.Subscribe(core.OnSwirlHydro, swirlfunc("a2-hydro"), "a2-hydro")
	c.Core.Events.Subscribe(core.OnSwirlPyro, swirlfunc("a2-pyro"), "a2-pyro")
}

func (c *char) a4() {
	c.a4EM = make([]float64, core.EndStatType)

	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddMod(core.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Status.Duration("sucrosea4") == 0 {
					return nil, false
				}
				return c.a4EM, true
			},
		})
	}
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

	c.Core.Log.Debugw("sucrose c4 reducing E CD", "frame", c.Core.F, "event", core.LogCharacterEvent, "cd_reduction", cdReduction)
}

func (c *char) c6() {
	c.AddMod(core.CharStatMod{
		Key: "sucrose-c6",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.Status.Duration("sucrosec6") == 0 {
				return nil, false
			}
			p := core.EleToDmgP(c.qInfused)
			val := make([]float64, core.EndStatType)
			val[p] = 0.2
			return val, true
		},
		Expiry: -1,
	})
}
