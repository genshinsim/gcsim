package yoimiya

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("yoimiya", NewChar)
}

type char struct {
	*character.Tmpl
	a2stack  int
	lastPart int
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3

	c.a2()
	c.onExit()
	c.burstHook()

	//add effect for burst

	return &c, nil
}

func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "yoimiya-a2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Sim.Status("yoimiyaa2") > 0 {
				val[core.Pyro] = float64(c.a2stack) * 0.02
				return val, true
			}
			c.a2stack = 0
			return nil, false
		},
	})
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if c.Sim.Status("yoimiyaskill") == 0 {
			return
		}
		if ds.AttackTag != core.AttackTagNormal {
			return
		}
		//here we can add stacks up to 10
		if c.a2stack < 10 {
			c.a2stack++
		}
		c.Sim.AddStatus("yoimiyaa2", 180)
		// c.a2expiry = c.Sim.Frame() + 180 // 3 seconds
	}, "yoimiya-a2")
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	//infusion to normal attack only
	if c.Sim.Status("yoimiyaskill") > 0 && ds.AttackTag == core.AttackTagNormal {
		ds.Element = core.Pyro
		//multiplier
		c.Log.Debugw("skill mult applied", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "prev", ds.Mult, "next", skill[c.TalentLvlSkill()]*ds.Mult, "char", c.Index)
		ds.Mult = skill[c.TalentLvlSkill()] * ds.Mult
	}
	return ds
}
