package raiden

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

type char struct {
	*character.Tmpl
	eyeICD         int
	stacksConsumed float64
	stacks         float64
	restoreICD     int
	restoreCount   int
	c6Count        int
	c6ICD          int
}

func init() {
	core.RegisterCharFunc("raiden", NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 90
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	if c.Base.Cons == 6 {
		c.c6()
	}

	c.eyeOnDamage()
	c.onBurstStackCount()

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	mult := skillBurstBonus[c.TalentLvlSkill()]
	//add E hook
	for _, char := range c.Core.Chars {
		this := char
		char.AddMod(core.CharStatMod{
			Key:    "raiden-e",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Status.Duration("raidenskill") == 0 {
					return nil, false
				}
				if a != core.AttackTagElementalBurst {
					return nil, false
				}
				val := make([]float64, core.EndStatType)
				val[core.DmgP] = mult * this.MaxEnergy()
				return val, true
			},
		})
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.Status.Duration("raidenburst") == 0 {
			return 25
		}
		return 20
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	//a2 add dmg based on ER%
	excess := int(ds.Stats[core.ER] / 0.01)

	ds.Stats[core.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Core.Log.Debugw("a4 adding electro dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", excess, "final", ds.Stats[core.ElectroP])
	//
	////infusion to normal/plunge/charge
	//switch ds.AttackTag {
	//case core.AttackTagNormal:
	//case core.AttackTagExtra:
	//case core.AttackTagPlunge:
	//default:
	//	return ds
	//}
	//if c.Core.Status.Duration("raidenburst") > 0 {
	//	ds.Element = core.Electro
	//}
	return ds
}
