package chongyun

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("chongyun", NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5

	c.onSwapHook()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons == 6 && s.Flags().DamageMode {
		c.c6()
	}

	return &c, nil
}

func (c *char) c4() {
	icd := 0
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if c.Sim.Frame() < icd {
			return
		}
		if !t.AuraContains(core.Cryo) {
			return
		}

		c.AddEnergy(2)

		c.Log.Debugw("chongyun c4 recovering 2 energy", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "final energy", c.Energy)
		icd = c.Sim.Frame() + 120
	}, "chongyun-c4")
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 24 //frames from keqing lib
		case 1:
			f = 62 - 24
		case 2:
			f = 124 - 62
		case 3:
			f = 204 - 124
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 30 //frames from keqing lib
	case core.ActionSkill:
		return 57
	case core.ActionBurst:
		return 135 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-1)
	if c.NormalCounter == 3 && c.Base.Cons >= 1 {
		d := c.Snapshot(
			"Chongyun C1",
			core.AttackTagNormal,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			25,
			.5,
		)
		//3 blades
		for i := 0; i < 3; i++ {
			x := d.Clone()
			c.QueueDmg(&x, f+i*5) //TODO: frames
		}
	}
	c.AdvanceNormalIndex()

	return f
}

func (c *char) Skill(p map[string]int) int {

	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Spirit Blade: Chonghua's Layered Frost",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-1)

	//TODO: energy count; lib says 3:4?
	c.QueueParticle("Chongyun", 4, core.Cryo, 100)

	//a4 delayed damage + cryo resist shred
	c.AddTask(func() {

		d := c.Snapshot(
			"Spirit Blade: Chonghua's Layered Frost (Ar)",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Cryo,
			25,
			skill[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		c.Sim.ApplyDamage(&d)
		//add res mod after dmg
		d.OnHitCallback = func(t core.Target) {
			t.AddResMod("Chongyun A4", core.ResistMod{
				Duration: 480, //10 seconds
				Ele:      core.Cryo,
				Value:    -0.10,
			})
		}

	}, "Chongyun-Skill", f+600)

	c.Sim.AddStatus("chongyunfield", 600)

	//TODO: delay between when frost field start ticking?
	for i := 60; i <= 600; i += 60 {
		c.AddTask(func() {
			active, _ := c.Sim.CharByPos(c.Sim.ActiveCharIndex())
			c.infuse(active)
		}, "chongyun-field", i)
	}

	c.SetCD(core.ActionSkill, 900)
	return f
}

func (c *char) onSwapHook() {
	c.Sim.AddEventHook(func(s core.Sim) bool {
		if s.Status("chongyunfield") == 0 {
			return false
		}
		//add infusion on swap
		c.Log.Debugw("chongyun adding infusion on swap", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "expiry", c.Sim.Frame()+infuseDur[c.TalentLvlSkill()])
		active, _ := c.Sim.CharByPos(c.Sim.ActiveCharIndex())
		c.infuse(active)
		return false
	}, "chongyun-field", core.PostSwapHook)
}

func (c *char) infuse(char core.Character) {
	switch char.WeaponClass() {
	case core.WeaponClassClaymore:
		fallthrough
	case core.WeaponClassSpear:
		fallthrough
	case core.WeaponClassSword:
		c.Log.Debugw("chongyun adding infusion", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "expiry", c.Sim.Frame()+infuseDur[c.TalentLvlSkill()])
		char.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "chongyun-ice-weapon",
			Ele:    core.Cryo,
			Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Sim.Frame() + infuseDur[c.TalentLvlSkill()],
		})
	default:
		return
	}

	//a2 adds 8% atkspd for 2.1 seconds
	val := make([]float64, core.EndStatType)
	val[core.AtkSpd] = 0.08
	char.AddMod(core.CharStatMod{
		Key:    "chongyun-field",
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Sim.Frame() + 126,
	})
	//c2 reduces CD by 15%
	if c.Base.Cons >= 2 {
		char.AddCDAdjustFunc(core.CDAdjust{
			Key: "chongyun-c2",
			Amount: func(a core.ActionType) float64 {
				if a == core.ActionSkill || a == core.ActionBurst {
					return -0.15
				}
				return 0
			},
			Expiry: c.Sim.Frame() + 126,
		})
	}
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)

	d := c.Snapshot(
		"Spirit Blade: Cloud-Parting Star",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		25,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	count := 3
	if c.Base.Cons == 6 {
		count = 4

	}

	for i := 0; i < count; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f+10*i)
	}

	c.SetCD(core.ActionBurst, 720)
	c.Energy = 0
	return f //TODO: frames
}

func (c *char) c6() {
	c.Sim.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Spirit Blade: Cloud-Parting Star" {
			return
		}
		if t.HP()/t.MaxHP() < c.HPCurrent/c.HPMax {
			ds.Stats[core.DmgP] += 0.15
			c.Log.Debugw("c6 add bonus dmg", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "final", ds.Stats[core.DmgP])
		}
	}, "chongyun-c6")
}
