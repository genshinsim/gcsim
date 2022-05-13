package mona

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Mona, NewChar)
}

type char struct {
	*character.Tmpl
	c2icd int
	// c6bonus float64
}

const (
	bubbleKey = "mona-bubble"
	omenKey   = "omen-debuff"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = -1

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}

func (c *char) Init() {
	c.Tmpl.Init()

	c.burstHook()
	c.a4()

	//add damage mod for omen
	//add E hook
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = dmgBonus[c.TalentLvlBurst()]
	for _, char := range c.Core.Chars {
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "mona-omen",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				//ignore if omen or bubble not present
				if t.GetTag(bubbleKey) < c.Core.F && t.GetTag(omenKey) < c.Core.F {
					return nil, false
				}
				return val, true
			},
		})
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}
}

//Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "mona-a4",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			val[core.HydroP] = .2 * atk.Snapshot.Stats[core.ER]
			return val, true
		},
	})
}
