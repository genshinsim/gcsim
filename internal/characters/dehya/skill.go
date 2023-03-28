package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int
var skillRecastFrames []int

const skillHitmark = 25
const skillRecastHitmark = 37

func init() {
	skillFrames = frames.InitAbilSlice(26) // E -> Swap (E -> E could be shorter)

	skillRecastFrames = frames.InitAbilSlice(46) // E Recast -> Swap
}

const (
	skillICDKey            = "dehya-skill-icd"
	dehyaFieldKey          = "dehya-field-status"
	sanctumPickupExtension = 24 // On recast from Burst/Skill-2 the field duration is extended by 0.4s
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.nextIsRecast {
		c.recastBefore = true
		c.nextIsRecast = false
		return c.skillRecast()
	}
	// calculate sanctum duration
	dur := 720
	if c.Base.Cons >= 2 {
		dur = 960
	}

	c.recastBefore = false
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Molten Inferno",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt, //TODO ???
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    c.c1var[1] * c.MaxHP(),
	}
	// TODO: damage frame
	c.skillSnapshot = c.Snapshot(&ai)

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := geometry.CalcOffsetPoint(c.Core.Combat.Player().Pos(), geometry.Point{Y: 0.8}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)

	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 5), skillHitmark)

	c.addField(dur)

	c.AddStatus(skillICDKey, skillHitmark+1, false)
	c.SetCD(action.ActionSkill, 1200)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}

func (c *char) skillHook() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		//atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		if !c.StatusIsActive(dehyaFieldKey) {
			return false
		}
		if c.StatusIsActive(skillICDKey) {
			return false
		}
		if dmg == 0 {
			return false
		}
		// don't proc if target hit is outside of the skill area
		if !trg.IsWithinArea(c.skillArea) {
			return false
		}

		// this ICD is most likely tied to the construct, so it's not hitlag extendable
		c.AddStatus(skillICDKey, 150, false) // proc every 2.5s

		c.Core.QueueAttackWithSnap(
			c.skillAttackInfo,
			c.skillSnapshot,
			combat.NewCircleHitOnTarget(trg, nil, 4.5),
			1,
		)

		c.Core.QueueParticle("dehya", 1, attributes.Pyro, c.ParticleDelay)

		return false
	}, "dehya-skill")
}

func (c *char) skillRecast() action.ActionInfo {
	c.recastBefore = true

	dur := c.StatusExpiry(dehyaFieldKey) + sanctumPickupExtension - c.Core.F //dur gets extended on field recast by a low margin, apparently
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ranging Flame",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt, //TODO ???
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillReposition[c.TalentLvlSkill()],
		FlatDmg:    c.c1var[1] * c.MaxHP(),
	}

	// pick up field at start
	c.Core.Log.NewEvent("sanctum removed", glog.LogCharacterEvent, c.Index).
		Write("Duration Remaining ", dur).
		Write("DoT tick CD", c.StatusDuration("dehya-skill-icd"))
	c.DeleteStatus(dehyaFieldKey)

	// save current DoT icd
	c.sanctumICD = c.StatusDuration(skillICDKey)
	c.AddStatus(skillICDKey, skillRecastHitmark+c.sanctumICD, false)

	//reposition

	// TODO: damage frame

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := geometry.CalcOffsetPoint(c.Core.Combat.Player().Pos(), geometry.Point{Y: 0.5}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)
	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 6), skillRecastHitmark)

	// place field back down
	c.QueueCharTask(func() { //place field
		c.addField(dur)
	}, skillRecastHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) addField(dur int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Molten Inferno (DoT)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt, //TODO ???
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillDotAtk[c.TalentLvlSkill()],
		FlatDmg:    (skillDotHP[c.TalentLvlSkill()] + c.c1var[1]) * c.MaxHP(),
	}
	//places field
	c.AddStatus(dehyaFieldKey, dur, false)
	c.Core.Log.NewEvent("sanctum added", glog.LogCharacterEvent, c.Index).
		Write("Duration Remaining ", dur).
		Write("New Expiry Frame", c.StatusExpiry(dehyaFieldKey)).
		Write("DoT tick CD", c.StatusDuration("dehya-skill-icd"))

	// snapshot for ticks
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)
}
