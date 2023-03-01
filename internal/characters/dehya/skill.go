package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
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
	skillICDKey = "dehya-skill-icd"
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
	c1var := 0.0
	if c.Base.Cons >= 1 {
		c1var = 0.036
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Molten Inferno",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt, //TODO ???
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    c1var * c.MaxHP(),
	}
	// TODO: damage frame
	c.skillSnapshot = c.Snapshot(&ai)

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := combat.CalcOffsetPoint(c.Core.Combat.Player().Pos(), combat.Point{Y: 3}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)

	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 5), skillHitmark)
	c.sanctumActive = true
	c.sanctumSource = c.Core.F
	c.sanctumExpiry = c.sanctumSource + dur + skillHitmark
	c.Core.Tasks.Add(c.removeSanctum(c.sanctumExpiry), c.sanctumExpiry-c.Core.F)

	// snapshot for ticks
	ai.Abil = "Molten Inferno (DoT)"
	ai.ICDTag = combat.ICDTagElementalArt
	ai.Mult = skillDotAtk[c.TalentLvlSkill()]
	ai.FlatDmg = skillDotHP[c.TalentLvlSkill()] * c.MaxHP()
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)

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
		if !c.sanctumActive {
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
			combat.NewCircleHitOnTarget(trg, nil, 3.4),
			1,
		)

		c.Core.QueueParticle("dehya", 1, attributes.Pyro, c.ParticleDelay)

		return false
	}, "dehya-skill")
}

func (c *char) skillRecast() action.ActionInfo {
	c.recastBefore = true
	c1var := 0.0
	if c.Base.Cons >= 1 {
		c1var = 0.036
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ranging Flame",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt, //TODO ???
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillReposition[c.TalentLvlSkill()],
		FlatDmg:    c1var * c.MaxHP(),
	}

	// pick up field at start
	c.Core.Log.NewEvent("Sanctum Expiration Info ", glog.LogCharacterEvent, c.Index).
		Write("Duration Remaining ", c.sanctumExpiry+c.sanctumPickupExtension-c.Core.F).
		Write("New Expiry Frame", c.sanctumExpiry+c.sanctumPickupExtension+skillRecastHitmark).
		Write("Field Source", c.sanctumSource).
		Write("DoT tick CD", c.StatusDuration("dehya-skill-icd"))
	c.sanctumActive = false
	c.sanctumExpiry += c.sanctumPickupExtension + skillRecastHitmark

	// save current DoT icd
	c.sanctumICD = c.StatusDuration(skillICDKey)
	c.AddStatus(skillICDKey, skillRecastHitmark+c.sanctumICD, false)

	//reposition

	// TODO: damage frame

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := combat.CalcOffsetPoint(c.Core.Combat.Player().Pos(), combat.Point{Y: 3}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)
	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 5), skillRecastHitmark)

	// place field back down
	c.Core.Tasks.Add(func() {
		c.sanctumActive = true
		c.Core.Tasks.Add(c.removeSanctum(c.sanctumExpiry), c.sanctumExpiry-c.Core.F)
		// snapshot for ticks
		ai.Abil = "Molten Inferno (DoT)"
		ai.ICDTag = combat.ICDTagElementalArt
		ai.Mult = skillDotAtk[c.TalentLvlSkill()]
		ai.FlatDmg = skillDotHP[c.TalentLvlSkill()] * c.MaxHP()
		c.skillAttackInfo = ai
		c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)
	}, skillRecastHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) removeSanctum(src int) func() {
	return func() {
		// if expiration has changed, then this is no longer the same sanctum, do nothing
		if c.sanctumExpiry != src {
			c.Core.Log.NewEvent("sanctum not removed, src changed", glog.LogCharacterEvent, c.Index).
				Write("src", src)
			return
		}
		c.Core.Log.NewEvent("sanctum removed", glog.LogCharacterEvent, c.Index).
			Write("src", src)
		c.sanctumActive = false
	}
}
