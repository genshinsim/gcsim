package thoma

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const (
	burstKey     = "thoma-q"
	burstICDKey  = "thoma-q-icd"
	burstHitmark = 40
)

func init() {
	burstFrames = frames.InitAbilSlice(58)
	burstFrames[action.ActionAttack] = 57
	burstFrames[action.ActionSkill] = 56
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionSwap] = 56
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Crimson Ooyoroi",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// damage component not final
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 4),
		burstHitmark,
		burstHitmark,
	)

	d := 15
	if c.Base.Cons >= 2 {
		d = 18
	}

	c.AddStatus(burstKey, d*60, true)

	c.burstStateHook()

	// C4: restore 15 energy
	if c.Base.Cons >= 4 {
		c.Core.Tasks.Add(func() {
			c.AddEnergy("thoma-c4", 15)
		}, 8)
	}

	cd := 20
	if c.Base.Cons >= 1 {
		cd = 17 // the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 17 for sim purposes.
	}
	c.SetCD(action.ActionBurst, cd*60)
	c.ConsumeEnergy(7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill],
		State:           action.BurstState,
	}
}

func (c *char) burstStateDelayFuncGen(src int) func() {
	return func() {
		//ignore if on ICD
		if c.StatusIsActive(burstICDKey) || c.Core.Player.CurrentState() != action.NormalAttackState || c.burstHookSrc != src {
			return
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonFieryCollapse()
		c.Core.Log.NewEvent("thoma burst on state change", glog.LogCharacterEvent, c.Index).
			Write("state", action.NormalAttackState).
			Write("icd", c.StatusExpiry(burstICDKey))
		c.burstTickSrc = c.Core.F
		//use the hitlag affected queue for this
		c.QueueCharTask(c.burstTickerFunc(c.Core.F), 60) //check every 1sec
	}
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(event.OnAttack, func(args ...interface{}) bool {
		//check if buff is up
		if !c.StatusIsActive(burstKey) {
			return false
		}
		c.burstHookSrc = c.Core.F
		delay := common.Get5PercentN0Delay(c.Core.Player.ActiveChar())
		c.Core.Log.NewEvent("thoma burst delay on state change", glog.LogCharacterEvent, c.Index).
			Write("delay", delay)
		// This accounts for the delay in n0 timing needed for Thoma to trigger collapses
		c.Core.Tasks.Add(c.burstStateDelayFuncGen(c.Core.F), delay)

		return false
	}, "thoma-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if !c.StatusIsActive(burstKey) {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("thoma burst tick check ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.Player.CurrentState()

		if state != action.NormalAttackState {
			c.Core.Log.NewEvent("thoma burst tick check stopped, not in normal state", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state", state)
			return
		}
		state_start := c.Core.Player.CurrentStateStart()
		norm_counter := c.Core.Player.ActiveChar().NormalCounter
		if (norm_counter == 1) && c.Core.F-state_start < common.Get5PercentN0Delay(c.Core.Player.ActiveChar()) {
			c.Core.Log.NewEvent("thoma burst tick check stopped, not enough time since normal state start", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state_start", state_start)
			return
		}

		c.Core.Log.NewEvent("thoma burst triggered from ticker", glog.LogCharacterEvent, c.Index).
			Write("src", src).
			Write("state", state).
			Write("icd", c.StatusExpiry(burstICDKey))
		//we can trigger a collapse here b/c we're in normal state still and src is still the same
		c.summonFieryCollapse()
		//in theory this should not hit an icd?
		//use the hitlag affected queue for this
		c.QueueCharTask(c.burstTickerFunc(src), 60) //check every 1sec
	}
}

func (c *char) summonFieryCollapse() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fiery Collapse",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstproc[c.TalentLvlBurst()],
		FlatDmg:    0.022 * c.MaxHP(),
	}
	done := false
	shieldCb := func(_ combat.AttackCB) {
		if done {
			return
		}
		shieldamt := (burstshieldpp[c.TalentLvlBurst()]*c.MaxHP() + burstshieldflat[c.TalentLvlBurst()])
		c.genShield("Thoma Burst", shieldamt, true)
		done = true
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4.59), 0, 11, shieldCb)
	c.AddStatus(burstICDKey, 60, true)
}
