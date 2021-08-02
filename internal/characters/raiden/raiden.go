package raiden

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

type char struct {
	*character.Tmpl
	eyeICD         int
	stacksConsumed float64
	stacks         float64
	restoreICD     int
	restoreCount   int
}

func init() {
	combat.RegisterCharFunc("raiden", NewChar)
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = def.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	c.eyeOnDamage()
	c.onBurstStackCount()

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	mult := skillBurstBonus[c.TalentLvlSkill()]
	//add E hook
	for _, char := range c.Sim.Characters() {
		char.AddMod(def.CharStatMod{
			Key:    "raiden-e",
			Expiry: -1,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				if c.Sim.Status("raidenskill") == 0 {
					return nil, false
				}
				if a != def.AttackTagElementalBurst {
					return nil, false
				}
				val := make([]float64, def.EndStatType)
				val[def.DmgP] = mult * char.MaxEnergy()
				return val, true
			},
		})
	}
}

func (c *char) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d def.Durability, mult float64) def.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	//a2 add dmg based on ER%
	excess := int(ds.Stats[def.ER] / 0.01)

	ds.Stats[def.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Log.Debugw("a4 adding electro dmg", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index, "stacks", excess, "final", ds.Stats[def.ElectroP])

	//infusion to normal/plunge/charge
	switch ds.AttackTag {
	case def.AttackTagNormal:
	case def.AttackTagExtra:
	case def.AttackTagPlunge:
	default:
		return ds
	}
	if c.Sim.Status("raidenburst") > 0 {
		ds.Element = def.Electro
	}
	return ds
}
