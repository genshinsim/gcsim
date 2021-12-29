package raiden

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
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
	core.RegisterCharFunc(keys.Raiden, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro
	c.Energy = 90
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5
	c.CharZone = core.ZoneInazuma

	if c.Base.Cons == 6 {
		c.c6()
	}

	c.eyeOnDamage()
	c.onBurstStackCount()
	c.onSwapClearBurst()

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	mult := skillBurstBonus[c.TalentLvlSkill()]
	//add E hook
	val := make([]float64, core.EndStatType)
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
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}
}

func (c *char) Snapshot(a *core.AttackInfo) core.Snapshot {
	s := c.Tmpl.Snapshot(a)

	//a2 add dmg based on ER%
	excess := int(s.Stats[core.ER] / 0.01)

	s.Stats[core.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Core.Log.Debugw("a4 adding electro dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", excess, "final", s.Stats[core.ElectroP])
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

	return s
}
