package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.YaeMiko, NewChar)
}

type char struct {
	*character.Tmpl
	kitsunes         []*kitsune
	totemParticleICD int
}

const (
	yaeTotemStatus = "yae_oldest_totem_expiry"
	yaeTotemCount  = "totems"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 90
	}
	c.Energy = float64(e)
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3
	c.BurstCon = 5
	c.SkillCon = 3

	c.SetNumCharges(core.ActionSkill, 3)

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
		return 50
	default:
		return 0
	}
}

func (c *char) a4() {
	m := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "yaemiko-a1",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != core.AttackTagElementalArt {
				return nil, false
			}
			m[core.DmgP] = c.Stat(core.EM) * 0.0015
			return m, true
		},
	})
}

// When Sesshou Sakura lightning hits opponents, the Electro DMG Bonus of all nearby party members is increased by 20% for 5s.
func (c *char) c4() {
	m := make([]float64, core.EndStatType)
	m[core.ElectroP] = .20

	// TODO: does this trigger for yaemiko too? assuming it does
	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "yaemiko-c4",
			Expiry: c.Core.F + 5*60,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
