package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Bennett, NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.Base.Element = core.Pyro

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	if c.Base.Cons >= 2 {
		c.c2()
	}
}

func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.ER] = .3

	c.AddMod(core.CharStatMod{
		Key:          "bennett-c2",
		Expiry:       -1,
		AffectedStat: core.ER, // to avoid infinite loop when calling MaxHP
		Amount: func() ([]float64, bool) {
			return val, c.HP()/c.MaxHP() < 0.7
		},
	})
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	var cd int
	var cdDelay int

	switch p["hold"] {
	case 1:
		c.skillHoldShort(false)
		cd = 450 - 90
		cdDelay = 43
	case 2:
		c.skillHoldLong()
		cd = 600 - 120
		cdDelay = 110
	default:
		c.skillPress()
		cd = 300 - 60
		cdDelay = 14
	}

	//A4
	if c.ModIsActive("bennett-field") {
		cd = cd / 2
	}

	c.SetCDWithDelay(core.ActionSkill, cd, cdDelay)

	return f, a

}

func (c *char) skillPress() {

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Passion Overload (Press)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 16, 16)

	//25 % chance of 3 orbs
	count := 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.QueueParticle("bennett", count, core.Pyro, 120)
}

func (c *char) skillHoldShort() {

	delay := []int{45, 57}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Passion Overload (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
	}

	for i, v := range skill1 {
		ai.Mult = v[c.TalentLvlSkill()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), delay[i], delay[i])
	}

	//25 % chance of 3 orbs
	count := 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.QueueParticle("bennett", count, core.Pyro, 215)
}

func (c *char) skillHoldLong() {

	delay := []int{112, 121}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Passion Overload (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
	}

	for i, v := range skill2 {
		ai.Mult = v[c.TalentLvlSkill()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), delay[i], delay[i])
	}

	ai.Mult = explosion[c.TalentLvlSkill()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 166, 166)

	//Bennett Hold E is guaranteed 3 orbs
	c.QueueParticle("bennett", 3, core.Pyro, 298)

}

const burstStartFrame = 34

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//add field effect timer
	c.Core.Status.AddStatus("btburst", 720+burstStartFrame)
	//hook for buffs; active right away after cast

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fantastic Voyage",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	//TODO: review bennett AOE size
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 37, 37)
	stats, _ := c.SnapshotStats()

	//apply right away
	c.applyBennettField(stats)()

	//add 12 ticks starting at t = 1 to t= 12
	// Buff appears to start ticking right before hit
	// https://discord.com/channels/845087716541595668/869210750596554772/936507730779308032
	for i := burstStartFrame; i <= 720+burstStartFrame; i += 60 {
		c.AddTask(c.applyBennettField(stats), "bennett-field", i)
	}

	c.ConsumeEnergy(36)
	c.SetCDWithDelay(core.ActionBurst, 900, 34)
	return f, a
}

const bennettSelfInfusionDurationInFrames = 126

func (c *char) applyBennettField(stats [core.EndStatType]float64) func() {
	hpplus := stats[core.Heal]
	heal := bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP()
	pc := burstatk[c.TalentLvlBurst()]
	if c.Base.Cons >= 1 {
		pc += 0.2
	}
	atk := pc * float64(c.Base.Atk+c.Weapon.Atk)
	return func() {
		c.Core.Log.NewEvent("bennett field ticking", core.LogCharacterEvent, -1)

		//self infuse
		player, ok := c.Core.Targets[0].(*player.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		player.ApplySelfInfusion(core.Pyro, 25, bennettSelfInfusionDurationInFrames)

		active := c.Core.Chars[c.Core.ActiveChar]
		//heal if under 70%
		if active.HP()/active.MaxHP() < .7 {
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.ActiveChar,
				Message: "Inspiration Field",
				Src:     heal,
				Bonus:   hpplus,
			})
		}

		//add attack if over 70%
		threshold := .7
		if c.Base.Cons >= 1 {
			threshold = 0
		}
		// Activate attack buff
		if active.HP()/active.MaxHP() > threshold {
			//add 2.1s = 126 frames
			val := make([]float64, core.EndStatType)
			val[core.ATK] = atk

			// 15% Pyro damage percent bonus applies to all characters in the field, regardless of weapon type
			if c.Base.Cons == 6 {
				val[core.PyroP] = 0.15
			}

			active.AddMod(core.CharStatMod{
				Key: "bennett-field",
				Amount: func() ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.F + 126,
			})
			c.Core.Log.NewEvent("bennett field - adding attack", core.LogCharacterEvent, c.Index, "threshold", threshold)
			//if c6 add weapon infusion and 15% pyro
			if c.Base.Cons == 6 {
				switch active.WeaponClass() {
				case core.WeaponClassClaymore:
					fallthrough
				case core.WeaponClassSpear:
					fallthrough
				case core.WeaponClassSword:
					active.AddWeaponInfuse(core.WeaponInfusion{
						Key:    "bennett-fire-weapon",
						Ele:    core.Pyro,
						Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
						Expiry: c.Core.F + 126,
					})
				}

			}
		}
	}
}
