package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var burstDiceHitmarks = []int{25, 30, 36, 41} //c2 hitmark not framecounted

// initial hit
const burstHitmark = 76

func init() {
	burstFrames = frames.InitAbilSlice(93)
	burstFrames[action.ActionAttack] = 92
	burstFrames[action.ActionAim] = 92
	burstFrames[action.ActionJump] = 92
	burstFrames[action.ActionSwap] = 91
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Depth-Clarion Dice",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       0,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	//apply hydro every 3rd hit
	//triggered on normal attack or yelan's skill

	//Initial hit
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	//TODO: check if we need to add f to this
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add(burstStatus, 15*60)
		c.a4() //TODO: does this call need to be delayed?
	}, burstHitmark)

	if c.Base.Cons >= 6 { //C6 passive, lasts 20 seconds
		c.Core.Status.Add(c6Status, 20*60)
		c.c6count = 0
	}
	c.Core.Log.NewEvent("burst activated", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+15*60)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) exquisiteThrowSkillProc() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Exquisite Throw",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagYelanBurst,
		ICDGroup:   combat.ICDGroupYelanBurst,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    burstDice[c.TalentLvlBurst()] * c.MaxHP(),
	}
	for i := 0; i < 3; i++ {
		//TODO: probably snapshots before hitmark
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstDiceHitmarks[i], burstDiceHitmarks[i])
	}
}

func (c *char) summonExquisiteThrow() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Exquisite Throw",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagYelanBurst,
		ICDGroup:   combat.ICDGroupYelanBurst,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    burstDice[c.TalentLvlBurst()] * c.MaxHP(),
	}
	for i := 0; i < 3; i++ {
		//TODO: probably snapshots before hitmark
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstDiceHitmarks[i], burstDiceHitmarks[i])
	}
	if c.Base.Cons >= 2 && c.c2icd <= c.Core.F {
		ai.Abil = "Yelan C2 Proc"
		ai.FlatDmg = 14.0 / 100 * c.MaxHP()
		c.c2icd = c.Core.F + 1.8*60
		//TODO: frames timing on this?
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstDiceHitmarks[3], burstDiceHitmarks[3])
	}

	c.burstDiceICD = c.Core.F + 60
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration(burstStatus) <= 0 {
			return false
		}
		next := args[1].(action.AnimationState)
		//ignore if not normal
		if next != action.NormalAttackState {
			return false
		}
		//ignore if on ICD
		if c.burstDiceICD > c.Core.F {
			return false
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonExquisiteThrow()
		c.Core.Log.NewEvent("yelan burst on state change", glog.LogCharacterEvent, c.Index, "state", next, "icd", c.burstDiceICD)
		c.burstTickSrc = c.Core.F
		c.Core.Tasks.Add(c.burstTickerFunc(c.Core.F), 60) //check every 1sec

		return false
	}, "yelan-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if c.Core.Status.Duration(burstStatus) <= 0 {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("yelan burst tick check ignored, src diff", glog.LogCharacterEvent, c.Index, "src", src, "new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.Player.CurrentState()
		if state != action.NormalAttackState {
			c.Core.Log.NewEvent("yelan burst tick check stopped, not normal state", glog.LogCharacterEvent, c.Index, "src", src, "state", state)
			return
		}
		c.Core.Log.NewEvent("yelan burst triggered from ticker", glog.LogCharacterEvent, c.Index, "src", src, "state", state, "icd", c.burstDiceICD)
		//we can trigger a wave here b/c we're in normal state still and src is still the same
		c.summonExquisiteThrow()
		//in theory this should not hit an icd?
		c.Core.Tasks.Add(c.burstTickerFunc(src), 60) //check every 1sec
	}
}
