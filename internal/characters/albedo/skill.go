package albedo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

const skillHitmark = 32

func init() {
	skillFrames = frames.InitAbilSlice(32)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Abiogenesis: Solar Isotoma",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: damage frame
	c.bloomSnapshot = c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, c.bloomSnapshot, combat.NewDefCircHit(3, false, combat.TargettableEnemy), skillHitmark)

	//snapshot for ticks
	ai.Abil = "Abiogenesis: Solar Isotoma (Tick)"
	ai.ICDTag = combat.ICDTagElementalArt
	ai.Mult = skillTick[c.TalentLvlSkill()]
	ai.UseDef = true
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)

	// Reset ICD
	c.icdSkill = c.Core.F - 1

	//create a construct
	// Construct is not fully formed until after the hit lands (exact timing unknown)
	c.Core.Tasks.Add(func() {
		c.Core.Constructs.New(c.newConstruct(1800), true)

		c.lastConstruct = c.Core.F

		c.Tags["elevator"] = 1
	}, skillHitmark)

	c.SetCD(action.ActionSkill, 240)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,

		State: action.SkillState,
	}
}

func (c *char) skillHook() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)
		if c.Tags["elevator"] == 0 {
			return false
		}
		if c.Core.F < c.icdSkill {
			return false
		}
		// Can't be triggered by itself when refreshing
		if atk.Info.Abil == "Abiogenesis: Solar Isotoma" {
			return false
		}

		c.icdSkill = c.Core.F + 120 // every 2 seconds

		snap := c.skillSnapshot

		// a1: skill tick deal 25% more dmg if enemy hp < 50%
		if c.Core.Combat.DamageMode && t.HP()/t.MaxHP() < .5 {
			snap.Stats[attributes.DmgP] += 0.25
			c.Core.Log.NewEvent("a1 proc'd, dealing extra dmg", glog.LogCharacterEvent, c.Index, "hp %", t.HP()/t.MaxHP(), "final dmg", snap.Stats[attributes.DmgP])
		}

		c.Core.QueueAttackWithSnap(c.skillAttackInfo, snap, combat.NewDefCircHit(3, false, combat.TargettableEnemy), 1)

		//67% chance to generate 1 geo orb
		if c.Core.Rand.Float64() < 0.67 {
			c.Core.QueueParticle("albedo", 1, attributes.Geo, 100)
		}

		// c1: skill tick regen 1.2 energy
		if c.Base.Cons >= 1 {
			c.AddEnergy("albedo-c1", 1.2)
			c.Core.Log.NewEvent("c1 restoring energy", glog.LogCharacterEvent, c.Index)
		}

		// c2: skill tick grant stacks, lasts 30s; each stack increase burst dmg by 30% of def, stack up to 4 times
		if c.Base.Cons >= 2 {
			if c.Core.Status.Duration("albedoc2") == 0 {
				c.Tags["c2"] = 0
			}
			c.Core.Status.Add("albedoc2", 1800) //lasts 30 seconds
			c.Tags["c2"]++
			if c.Tags["c2"] > 4 {
				c.Tags["c2"] = 4
			}
		}

		return false
	}, "albedo-skill")
}
