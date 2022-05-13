package xiangling

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Xiangling, NewChar)
}

type char struct {
	*character.Tmpl
	guoba *panda
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
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue
	c.Base.Element = core.Pyro

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	//add in a guoba
	c.guoba = newGuoba(c.Core)
	c.Core.AddTarget(c.guoba)
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

func (c *char) c6(dur int) {
	m := make([]float64, core.EndStatType)
	m[core.PyroP] = 0.15

	c.Core.Status.AddStatus("xlc6", dur)

	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "xiangling-c6",
			Expiry: c.Core.F + dur,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

var hitmarks = [][]int{{12}, {8}, {11, 18}, {5, 15, 24, 29}, {21}}

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
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(0.1, false, core.TargettableEnemy),
			hitmarks[c.NormalCounter][i],
			hitmarks[c.NormalCounter][i],
		)
	}

	//if n = 5, add explosion for c2
	if c.Base.Cons >= 2 && c.NormalCounter == 4 {
		// According to TCL, does not snapshot and has no ability type scaling tags
		// TODO: Does not mention ICD or pyro aura strength?
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Oil Meets Fire (C2)",
			AttackTag:  core.AttackTagNone,
			ICDTag:     core.ICDTagNormalAttack,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       .75,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 120, 120) //todo: explosion frames
	}
	//add a 75 frame attackcounter reset
	c.AdvanceNormalIndex()
	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       nc[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	//return animation cd
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guoba",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}

	var cb core.AttackCBFunc
	if c.Base.Cons >= 1 {
		cb = func(a core.AttackCB) {
			a.Target.AddResMod("xiangling-c1", core.ResistMod{
				Ele:      core.Pyro,
				Value:    -0.15,
				Duration: 6 * 60,
			})
		}

	}

	delay := 126 //first tick at 126
	c.Core.Status.AddStatus("xianglingguoba", 500+f)

	//lasts 7.3 seconds, shoots every 100 frames
	snap := c.Snapshot(&ai)
	for i := 0; i < 4; i++ {
		c.AddTask(func() {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 10, cb)
			c.guoba.pyroWindowStart = c.Core.F
			c.guoba.pyroWindowEnd = c.Core.F + 20
		}, "guoba-shoot", delay+i*100-10) //10 frame window to swirl
		//TODO: check guoba particle generation
		c.QueueParticle("xiangling", 1, core.Pyro, delay+i*100+150)
	}

	c.SetCDWithDelay(core.ActionSkill, 12*60, 13)
	//return animation cd
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	lvl := c.TalentLvlBurst()

	delay := []int{18, 33, 56}
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}
	for i := 0; i < len(pyronadoInitial); i++ {
		ai.Abil = fmt.Sprintf("Pyronado Hit %v", i+1)
		ai.Mult = pyronadoInitial[i][lvl]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), delay[i], delay[i])
	}

	//approx 73 frames per cycle
	//max is either 10s or 14s, plus animation
	max := 10*60 + a
	if c.Base.Cons >= 4 {
		max = 14*60 + a
	}

	ai = core.AttackInfo{
		Abil:       "Pyronado",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       pyronadoSpin[lvl],
	}

	c.Core.Status.AddStatus("xianglingburst", max)

	for delay := 56; delay <= max; delay += 73 { //first hit on same frame as 3rd initial hit
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2.5, false, core.TargettableEnemy), 54, delay)
	}

	//add an effect starting at frame 55 to end of duration to increase pyro dmg by 15% if c6
	if c.Base.Cons >= 6 {
		//wait 55 frames, add effect.
		c.AddTask(func() { c.c6(max) }, "xiangling-c6", 55)
	}

	//add cooldown to sim
	c.SetCDWithDelay(core.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(24)

	//return animation cd
	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}

func (c *char) Jump(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionJump, p)
	return f, a
}
