package eula

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstHitmark = 100
const lightfallHitmark = 35

func init() {
	burstFrames = frames.InitAbilSlice(122)
	burstFrames[action.ActionDash] = 121
	burstFrames[action.ActionJump] = 121
	burstFrames[action.ActionSwap] = 121
	burstFrames[action.ActionWalk] = 117
}

const (
	burstKey = "eula-q"
)

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.burstCounter = 0
	if c.Base.Cons >= 6 {
		c.burstCounter = 5
	}
	// lights up 9.5s from cast
	//deployable; not affected by hitlag
	c.Core.Status.Add(burstKey, 9*60+30)
	c.Core.Log.NewEvent("eula burst started", glog.LogCharacterEvent, c.Index).
		Write("stacks", c.burstCounter).
		Write("expiry", c.Core.Status.Duration(burstKey))

	//add initial damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Illumination",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burstInitial[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	// A4: When Glacial Illumination is cast, the CD of Icetide Vortex is reset and Eula gains 1 stack of Grimheart.
	if c.grimheartStacks < 2 {
		c.grimheartStacks++
	}
	c.Core.Log.NewEvent("eula: grimheart stack", glog.LogCharacterEvent, c.Index).
		Write("current count", c.grimheartStacks)
	c.ResetActionCooldown(action.ActionSkill)
	c.Core.Log.NewEvent("eula a4 reset skill cd", glog.LogCharacterEvent, c.Index)

	// lightfall hitmark is 600f from cast
	c.Core.Tasks.Add(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Core.Status.Duration(burstKey) > 0 {
			c.triggerBurst()
		}
	}, 600-lightfallHitmark)

	//energy does not deplete until after animation
	c.ConsumeEnergy(107)
	c.SetCDWithDelay(action.ActionBurst, 20*60, 97)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionWalk], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) triggerBurst() {
	if c.burstCounter > 30 {
		c.burstCounter = 30
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Illumination (Lightfall)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 50,
		Mult:       burstExplodeBase[c.TalentLvlBurst()] + burstExplodeStack[c.TalentLvlBurst()]*float64(c.burstCounter),
	}

	c.Core.Log.NewEvent("eula burst triggering", glog.LogCharacterEvent, c.Index).
		Write("stacks", c.burstCounter).
		Write("mult", ai.Mult)

	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), lightfallHitmark, lightfallHitmark)
	c.Core.Status.Delete(burstKey)
	c.burstCounter = 0
}

func (c *char) burstStacks() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Core.Status.Duration(burstKey) == 0 {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		//TODO: this looks like the icd is dependent on gadget timer. need to double check
		if c.burstCounterICD > c.Core.F {
			return false
		}
		switch atk.Info.AttackTag {
		case combat.AttackTagElementalArt:
		case combat.AttackTagElementalBurst:
		case combat.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.Core.Log.NewEvent("eula burst add stack", glog.LogCharacterEvent, c.Index).
			Write("stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.Core.Log.NewEvent("eula c6 add additional stack", glog.LogCharacterEvent, c.Index).
				Write("stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Core.F + 6
		return false
	}, "eula-burst-counter")
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration(burstKey) > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit")
}
