package thoma

import (
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
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4),
		burstHitmark,
		burstHitmark,
	)

	d := 15
	if c.Base.Cons >= 2 {
		d = 18
	}

	c.AddStatus(burstKey, d*60, true)

	c.burstProc()

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

func (c *char) burstProc() {
	// does not deactivate on death
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}
		next := args[1].(action.AnimationState)
		if next != action.NormalAttackState {
			return false
		}
		if c.StatusIsActive(burstICDKey) {
			c.Core.Log.NewEvent("thoma Q (active) on icd", glog.LogCharacterEvent, c.Index).
				Write("frame", c.Core.F)
			return false
		}
		c.summonFieryCollapse()
		c.Core.Log.NewEvent("thoma burst on state change", glog.LogCharacterEvent, c.Index).
			Write("frame", c.Core.F).
			Write("char", c.Core.Player.Active()).
			Write("icd", c.StatusExpiry(burstICDKey))
		c.burstTickSrc = c.Core.F
		c.QueueCharTask(c.burstTickFunc(c.Core.F), 60)
		return false
	}, "thoma-burst-animation-check")
}

func (c *char) burstTickFunc(src int) func() {
	return func() {
		if !c.StatusIsActive(burstKey) {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("thoma burst tick stopped, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.burstTickSrc)
			return
		}
		state := c.Core.Player.CurrentState()
		if state != action.NormalAttackState {
			c.Core.Log.NewEvent("thoma burst tick stopped, not normal state", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state", state)
			return
		}
		c.Core.Log.NewEvent("thoma burst triggered from tick", glog.LogCharacterEvent, c.Index).
			Write("src", src).
			Write("state", state).
			Write("icd", c.StatusExpiry(burstICDKey))
		c.summonFieryCollapse()
		c.QueueCharTask(c.burstTickFunc(src), 60)
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
	// TODO: moving hitbox
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4.5, 8),
		0,
		11,
		shieldCb,
	)
	c.AddStatus(burstICDKey, 60, true)
}
