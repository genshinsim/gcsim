package raiden

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
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
	QSnapshot      core.Snapshot
}

func init() {
	combat.RegisterCharFunc("raiden", NewChar)
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
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
	for _, char := range c.Sim.Characters() {
		char.AddMod(core.CharStatMod{
			Key:    "raiden-e",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Sim.Status("raidenskill") == 0 {
					return nil, false
				}
				if a != core.AttackTagElementalBurst {
					return nil, false
				}
				val := make([]float64, core.EndStatType)
				val[core.DmgP] = mult * char.MaxEnergy()
				return val, true
			},
		})
	}
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	//a2 add dmg based on ER%
	excess := int(ds.Stats[core.ER] / 0.01)

	ds.Stats[core.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Log.Debugw("a4 adding electro dmg", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "char", c.Index, "stacks", excess, "final", ds.Stats[core.ElectroP])
	//
	////infusion to normal/plunge/charge
	//switch ds.AttackTag {
	//case core.AttackTagNormal:
	//case core.AttackTagExtra:
	//case core.AttackTagPlunge:
	//default:
	//	return ds
	//}
	//if c.Sim.Status("raidenburst") > 0 {
	//	ds.Element = core.Electro
	//}
	return ds
}
