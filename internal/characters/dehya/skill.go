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
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int
var skillRecastFrames []int

const skillHitmark = 20
const skillRecastHitmark = 40

func init() {
	skillFrames = frames.InitAbilSlice(25) // E -> Swap/Dash/Walk
	skillFrames[action.ActionAttack] = 39  // E -> N1
	skillFrames[action.ActionSkill] = 30   // E -> E
	skillFrames[action.ActionJump] = 28    // E -> J
	skillFrames[action.ActionBurst] = 29   // E -> Q

	skillRecastFrames = frames.InitAbilSlice(44) // E -> Swap/Walk
	skillRecastFrames[action.ActionAttack] = 74  // E -> N1
	skillRecastFrames[action.ActionSkill] = 45   // E -> E
	skillRecastFrames[action.ActionDash] = 45    // E -> D
	skillRecastFrames[action.ActionJump] = 50    // E -> J
	skillRecastFrames[action.ActionBurst] = 45   // E -> Q
}

const (
	skillICDKey            = "dehya-skill-icd"
	dehyaFieldKey          = "dehya-field-status"
	sanctumPickupExtension = 24 // On recast from Burst/Skill-2 the field duration is extended by 0.4s
)

func (c *char) Skill(p map[string]int) (action.Info, error) {
	burstAction := c.UseBurstAction()
	if burstAction != nil {
		return *burstAction, nil
	}
	if c.StatusIsActive(dehyaFieldKey) {
		// If recast has been used, sanctum needs to be placed anew
		if c.hasRecastSkill {
			c.removeField()
		} else {
			c.hasRecastSkill = true
			return c.skillRecast()
		}
	}

	c.hasRecastSkill = false
	c.hasC2DamageBuff = false
	// Initial cast duration is always 12s
	dur := 720

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Molten Inferno",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           50,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		FlatDmg:            c.c1var[1] * c.MaxHP(),
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}
	// TODO: damage frame
	c.skillSnapshot = c.Snapshot(&ai)

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := geometry.CalcOffsetPoint(c.Core.Combat.Player().Pos(), geometry.Point{Y: 0.8}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)

	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 5), skillHitmark)

	c.Core.Tasks.Add(func() { // place field
		c.addField(dur)
	}, skillHitmark+1)

	c.AddStatus(skillICDKey, skillHitmark+1, false)
	c.SetCDWithDelay(action.ActionSkill, 20*60, 18)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHook() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		// atk := args[1].(*combat.AttackEvent)
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
			2,
		)

		// Set buff flag to false with 2f delay to line up with activation delay
		if c.hasC2DamageBuff {
			c.Core.Tasks.Add(func() {
				c.hasC2DamageBuff = false
			}, 2)
		}

		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Pyro, c.ParticleDelay)

		return false
	}, "dehya-skill")
}

func (c *char) skillRecast() (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Ranging Flame",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           50,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skillReposition[c.TalentLvlSkill()],
		FlatDmg:            c.c1var[1] * c.MaxHP(),
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}

	// pick up field at start
	c.removeField()

	// Add icd extension
	c.AddStatus(skillICDKey, skillRecastHitmark+c.sanctumICD, false)

	// reposition

	// TODO: damage frame

	player := c.Core.Combat.Player()
	// assuming tap e for hitbox offset
	skillPos := geometry.CalcOffsetPoint(c.Core.Combat.Player().Pos(), geometry.Point{Y: 0.5}, player.Direction())
	c.skillArea = combat.NewCircleHitOnTarget(skillPos, nil, 10)
	c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewCircleHitOnTarget(skillPos, nil, 6), skillRecastHitmark)

	// place field back down
	c.Core.Tasks.Add(func() { // place field
		// if C2, duration will be extended by 6s on recreation
		if c.Base.Cons >= 2 {
			c.sanctumSavedDur += 360
		}
		c.addField(c.sanctumSavedDur)
	}, skillRecastHitmark+1)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

// Remove field and save current ICD and duration with implicit extension
func (c *char) removeField() {
	c.sanctumICD = c.StatusDuration(skillICDKey)
	c.sanctumSavedDur = c.StatusExpiry(dehyaFieldKey) + sanctumPickupExtension - c.Core.F // dur gets extended on field recast by a low margin, apparently
	c.Core.Log.NewEvent("sanctum removed", glog.LogCharacterEvent, c.Index).
		Write("Duration Remaining ", c.sanctumSavedDur).
		Write("DoT tick CD", c.StatusDuration(skillICDKey))
	c.Core.Tasks.Add(func() {
		c.DeleteStatus(dehyaFieldKey)
	}, 1)
}

func (c *char) addField(dur int) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Molten Inferno (DoT)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skillDotAtk[c.TalentLvlSkill()],
		FlatDmg:            (skillDotHP[c.TalentLvlSkill()] + c.c1var[1]) * c.MaxHP(),
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}
	// places field
	c.AddStatus(dehyaFieldKey, dur, false)
	c.Core.Log.NewEvent("sanctum added", glog.LogCharacterEvent, c.Index).
		Write("Duration Remaining ", dur).
		Write("New Expiry Frame", c.StatusExpiry(dehyaFieldKey)).
		Write("DoT tick CD", c.StatusDuration(skillICDKey))

	// snapshot for ticks
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)
}

func (c *char) skillHurtHook() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if !c.StatusIsActive(dehyaFieldKey) {
			return false
		}
		if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
			c.Core.Log.NewEvent("dehya-c2-check failed", glog.LogCharacterEvent, c.Index).
				Write("Target not within field", c.Core.Combat.Player().Pos())
			return false
		}
		// Damage mitigation in the future if target is not Dehya
		if c.Base.Cons >= 2 {
			c.Core.Log.NewEvent("dehya-sanctum-c2-damage activated", glog.LogCharacterEvent, c.Index)
			c.hasC2DamageBuff = true
		}
		return false
	}, "dehya-field-dmgtaken")
}
