package kaeya

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Kaeya, NewChar)
}

type char struct {
	*character.Tmpl
	c4icd     int
	icicleICD []int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5

	c.icicleICD = make([]int, 4)

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	if c.Base.Cons > 0 {
		c.c1()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(.3, false, core.TargettableEnemy), f-1, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge 1",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[0][c.TalentLvlAttack()],
	}
	//TODO: damage frame
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f-15, f-15)
	ai.Abil = "Charge 2"
	ai.Mult = charge[1][c.TalentLvlAttack()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f-5, f-5)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Frostgnaw",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	a4count := 0
	cb := func(a core.AttackCB) {
		heal := .15 * (a.AttackEvent.Snapshot.BaseAtk*(1+a.AttackEvent.Snapshot.Stats[core.ATKP]) + a.AttackEvent.Snapshot.Stats[core.ATK])
		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Cold-Blooded Strike",
			Src:     heal,
			Bonus:   c.Stat(core.Heal),
		})
		//if target is frozen after hit then drop additional energy;
		if a4count == 2 {
			return
		}
		if a.Target.AuraContains(core.Frozen) {
			a4count++
			c.QueueParticle("kaeya", 1, core.Cryo, 100)
			c.Core.Log.NewEvent("kaeya a4 proc", core.LogCharacterEvent, c.Index)
		}
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, 28, cb)

	//2 or 3 1:1 ratio
	count := 2
	if c.Core.Rand.Float64() < 0.67 {
		count = 3
	}
	c.QueueParticle("kaeya", count, core.Cryo, f+100)

	c.SetCD(core.ActionSkill, 360+28) //+28 since cd starts 28 frames in
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Waltz",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	//duration starts counting 49 frames in per kqm lib
	//hits around 13 times

	//each icicle takes 120frames to complete a rotation and has a internal cooldown of 0.5
	count := 3
	if c.Base.Cons == 6 {
		count++
	}
	offset := 120 / count

	for i := 0; i < count; i++ {

		//each icicle will start at i * offset (i.e. 0, 40, 80 OR 0, 30, 60, 90)
		//assume each icicle will last for 8 seconds
		//assume damage dealt every 120 (since only hitting at the front)
		//on icicle collision, it'll trigger an aoe dmg with radius 2
		//in effect, every target gets hit every time icicles rotate around
		for j := f + offset*i; j < f+480; j += 120 {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), j)
		}
	}

	c.ConsumeEnergy(55)
	if c.Base.Cons == 6 {
		c.AddTask(func() { c.AddEnergy("kaeya-c6", 15) }, "kaeya-c6", 56)
	}

	c.SetCDWithDelay(core.ActionBurst, 900, 55)
	return f, a
}

// func (c *char) burstICD() {
// 	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
// 		atk := args[1].(*core.AttackEvent)
// 		if atk.Info.ActorIndex != c.Index {
// 			return false
// 		}
// 		if ds.Abil != "Glacial Waltz" {
// 			return false
// 		}
// 		//check icd
// 		if c.icicleICD[ds.ExtraIndex] > c.Core.F {
// 			ds.Cancelled = true
// 			return false
// 		}
// 		c.icicleICD[ds.ExtraIndex] = c.Core.F + 30
// 		return false
// 	}, "kaeya-burst-icd")
// }
