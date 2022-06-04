package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

var skillFrames []int

const skillHitmark = 56

func (c *char) Skill(p map[string]int) action.ActionInfo {

	//if fieldSrc is < duration then this is prob a sac proc
	//we need to stop the old field from ticking (by changing fieldSrc)
	//and also trigger a4 delayed damage
	src := c.Core.F

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Chonghua's Layered Frost",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(3, false, combat.TargettableEnemy), 0, skillHitmark)

	//TODO: energy count; lib says 3:4?
	c.Core.QueueParticle("chongyun", 4, attributes.Cryo, 100)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Chonghua's Layered Frost (Ar)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	cb := func(a combat.AttackCB) {
		e, ok := a.Target.(core.Enemy)
		if !ok {
			return
		}
		e.AddResistMod("chongyun-a4", 480, attributes.Cryo, -0.10)
	}
	snap := c.Snapshot(&ai)

	//if field is overwriting last
	if src-c.fieldSrc < 600 {
		//we're overriding previous field so trigger a4 here
		atk := c.a4Snap
		c.Core.QueueAttackEvent(atk, 1)
	}
	c.fieldSrc = src
	//override previous snap
	c.a4Snap = &combat.AttackEvent{
		Info:     ai,
		Snapshot: snap,
		Pattern:  combat.NewDefCircHit(3, false, combat.TargettableEnemy),
	}
	c.a4Snap.Callbacks = append(c.a4Snap.Callbacks, cb)

	//a4 delayed damage + cryo resist shred
	c.Core.Tasks.Add(func() {
		//if src changed then that means the field changed already
		if src != c.fieldSrc {
			return
		}
		//TODO: this needs to be fixed still for sac gs
		c.Core.QueueAttackEvent(c.a4Snap, 0)
	}, skillHitmark+600)

	c.Core.Status.Add("chongyunfield", 600)

	//TODO: delay between when frost field start ticking?
	for i := skillHitmark - 1; i <= 600; i += 60 {
		c.Core.Tasks.Add(func() {
			if src != c.fieldSrc {
				return
			}
			active := c.Core.Player.ActiveChar()
			c.infuse(active)
		}, i)
	}

	c.SetCD(action.ActionSkill, 900)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) onSwapHook() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("chongyunfield") == 0 {
			return false
		}
		//add infusion on swap
		c.Core.Log.NewEvent("chongyun adding infusion on swap", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
		active := c.Core.Player.ActiveChar()
		c.infuse(active)
		return false
	}, "chongyun-field")
}

func (c *char) infuse(active *character.CharWrapper) {
	//c2 reduces CD by 15%
	if c.Base.Cons >= 2 {
		active.AddCooldownMod("chongyun-c2", 126, func(a action.Action) float64 {
			if a == action.ActionSkill || a == action.ActionBurst {
				return -0.15
			}
			return 0
		})
	}

	// weapon infuse
	switch active.Weapon.Class {
	case weapon.WeaponClassClaymore:
		fallthrough
	case weapon.WeaponClassSpear:
		fallthrough
	case weapon.WeaponClassSword:
		c.Core.Player.AddWeaponInfuse(
			active.Index,
			"chongyun-ice-weapon",
			attributes.Cryo,
			infuseDur[c.TalentLvlSkill()],
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		)
		c.Core.Log.NewEvent("chongyun adding infusion", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+infuseDur[c.TalentLvlSkill()])
	default:
		return
	}

	//a1 adds 8% atkspd for 2.1 seconds
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.08
	active.AddStatMod("chongyun-field", 126, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})
}
