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

const burstHitmark = 40

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
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	d := 15
	if c.Base.Cons >= 2 {
		d = 18
	}

	c.AddStatus("thoma-q", d*60, true)

	c.burstProc()

	// C4: restore 15 energy
	if c.Base.Cons >= 4 {
		c.Core.Tasks.Add(func() {
			c.AddEnergy("thoma-c4", 15)
		}, 8)
	}

	cd := 20
	if c.Base.Cons >= 1 {
		cd = 17 //the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 17 for sim purposes.
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
	const icdKey = "thoma-q-icd"
	icd := 60 // 1s * 60
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}

		if !c.StatusIsActive("thoma-q") {
			return false
		}
		if c.StatusIsActive(icdKey) {
			c.Core.Log.NewEvent("thoma Q (active) on icd", glog.LogCharacterEvent, c.Index).
				Write("frame", c.Core.F)
			return false
		}
		c.AddStatus(icdKey, icd, true)

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
		//trigger a chain of attacks starting at the first target
		atk := combat.AttackEvent{
			Info: ai,
		}
		atk.SourceFrame = c.Core.F
		atk.Pattern = combat.NewDefSingleTarget(t.Index(), combat.TargettableEnemy)
		cb := func(_ combat.AttackCB) {
			shieldamt := (burstshieldpp[c.TalentLvlBurst()]*c.MaxHP() + burstshieldflat[c.TalentLvlBurst()])
			c.genShield("Thoma Burst", shieldamt)
		}
		atk.Callbacks = append(atk.Callbacks, cb)
		c.Core.QueueAttackEvent(&atk, 11)

		c.Core.Log.NewEvent("thoma Q proc'd", glog.LogCharacterEvent, c.Index).
			Write("frame", c.Core.F).
			Write("char", ae.Info.ActorIndex).
			Write("attack tag", ae.Info.AttackTag)

		return false
	}, "thoma-burst")
}
