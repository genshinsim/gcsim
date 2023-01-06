package eula

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillPressFrames []int
var skillHoldFrames []int
var icewhirlHitmarks = []int{79, 92}

const (
	skillPressHitmark = 20
	skillHoldHitmark  = 49
	a1Hitmark         = 108
	grimheartICD      = "eula-grimheart-icd"
	grimheartDuration = "eula-grimheart-duration"
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(48)
	skillPressFrames[action.ActionAttack] = 31
	skillPressFrames[action.ActionBurst] = 31
	skillPressFrames[action.ActionDash] = 29
	skillPressFrames[action.ActionJump] = 30
	skillPressFrames[action.ActionSwap] = 29

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(100)
	skillHoldFrames[action.ActionAttack] = 77
	skillHoldFrames[action.ActionBurst] = 77
	skillHoldFrames[action.ActionDash] = 75
	skillHoldFrames[action.ActionJump] = 75
	skillHoldFrames[action.ActionWalk] = 75
}

func (c *char) addGrimheartStack() {
	if !c.StatusIsActive(grimheartDuration) {
		c.grimheartStacks = 0
	}
	if c.grimheartStacks < 2 {
		c.grimheartStacks++
		c.Core.Log.NewEvent("eula: grimheart stack", glog.LogCharacterEvent, c.Index).
			Write("current count", c.grimheartStacks)
	}
	//refresh grimheart duration regardless
	c.AddStatus(grimheartDuration, 1080, true) //18 sec
}

func (c *char) currentGrimheartStacks() int {
	if !c.StatusIsActive(grimheartDuration) {
		c.grimheartStacks = 0
		return 0
	}
	if c.grimheartStacks > 2 {
		c.grimheartStacks = 2
	}
	return c.grimheartStacks
}

func (c *char) consumeGrimheartStacks() {
	c.grimheartStacks = 0
	c.DeleteStatus(grimheartDuration)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] != 0 {
		return c.holdSkill(p)
	}
	return c.pressSkill(p)
}

func (c *char) pressSkill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Icetide Vortex",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.particleDone = false
	//add 1 to grim heart if not capped by icd
	cb := func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.StatusIsActive(grimheartICD) {
			return
		}
		c.AddStatus(grimheartICD, 18, true)
		c.addGrimheartStack()
		if !c.particleDone {
			var count float64 = 1
			if c.Core.Rand.Float64() < .5 {
				count = 2
			}
			c.Core.QueueParticle("eula", count, attributes.Cryo, c.ParticleDelay)
			c.particleDone = true
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 3.5),
		skillPressHitmark,
		skillPressHitmark,
		cb,
	)

	c.SetCDWithDelay(action.ActionSkill, 60*4, 16)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) holdSkill(p map[string]int) action.ActionInfo {
	//hold e
	//296 to 341, but cd starts at 322
	//60 fps = 108 frames cast, cd starts 62 frames in so need to + 62 frames to cd
	lvl := c.TalentLvlSkill()
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Icetide Vortex (Hold)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillHold[lvl],
		HitlagHaltFrames:   0.12 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.particleDone = false
	energyCB := func(_ combat.AttackCB) {
		if !c.particleDone {
			var count float64 = 2
			if c.Core.Rand.Float64() < .5 {
				count = 3
			}
			c.Core.QueueParticle("eula", count, attributes.Cryo, c.ParticleDelay)
			c.particleDone = true
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 5.5),
		skillHoldHitmark,
		skillHoldHitmark,
		energyCB,
	)

	v := c.currentGrimheartStacks()

	//shred
	var shredCB combat.AttackCBFunc
	if v > 0 {
		done := false
		shredCB = func(a combat.AttackCB) {
			if done {
				return
			}
			e, ok := a.Target.(*enemy.Enemy)
			if !ok {
				return
			}
			done = true
			e.AddResistMod(enemy.ResistMod{
				Base:  modifier.NewBaseWithHitlag("eula-icewhirl-shred-cryo", 7*v*60),
				Ele:   attributes.Cryo,
				Value: -resRed[lvl],
			})
			e.AddResistMod(enemy.ResistMod{
				Base:  modifier.NewBaseWithHitlag("eula-icewhirl-shred-phys", 7*v*60),
				Ele:   attributes.Physical,
				Value: -resRed[lvl],
			})
		}
	}

	for i := 0; i < v; i++ {
		//multiple brand hits
		//TODO: need to double check if this is affected by hitlag; might be a deployable
		icewhirlAI := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Icetide Vortex (Icewhirl)",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       icewhirl[lvl],
		}
		if i == 0 {
			//per shizuka first swirl is not affected by hitlag?
			c.Core.QueueAttack(
				icewhirlAI,
				combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
				icewhirlHitmarks[i],
				icewhirlHitmarks[i],
				shredCB,
			)
		} else {
			c.QueueCharTask(func() {
				//spacing it out for stacks
				c.Core.QueueAttack(
					icewhirlAI,
					combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
					0,
					0,
					shredCB,
				)
			}, icewhirlHitmarks[i])
		}
	}

	//A1
	if v == 2 {
		// make sure this gets executed after hold e hitlag starts but before hold e is over
		// this makes it so it doesn't get affected by hitlag after Hold E is over
		aiA1 := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Icetide (Lightfall)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Physical,
			Durability: 25,
			Mult:       burstExplodeBase[c.TalentLvlBurst()] * 0.5,
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				aiA1,
				combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 6.5),
				a1Hitmark-(skillHoldHitmark+1),
				a1Hitmark-(skillHoldHitmark+1),
			)
		}, skillHoldHitmark+1)
	}

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		//TODO: check if the duration is right
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("eula-c1", (6*v+6)*60),
			AffectedStat: attributes.PhyP,
			Amount: func() ([]float64, bool) {
				return c.c1buff, true
			},
		})
	}

	c.consumeGrimheartStacks()
	cd := 10
	if c.Base.Cons >= 2 {
		cd = 4 //press and hold have same cd TODO: check if this is right
	}
	c.SetCDWithDelay(action.ActionSkill, cd*60, 46)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
