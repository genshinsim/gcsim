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

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillHitmark           = 20
	skillRecastHitmark     = 40
	skillMitigationAbil    = "Fiery Sanctum Mitigation"
	skillSelfDoTAbil       = "Redmane's Blood"
	skillSelfDoTStatus     = "dehya-redmanes-blood"
	skillSelfDoTStart      = 0.1 * 60 // looks like initial dot tick happens 0.1s after mitigating
	skillSelfDoTDuration   = 10 * 60  // total of 10 ticks at 0.1s, 1.1s, ..., 9.1s
	skillSelfDoTRatio      = 0.1
	skillSelfDoTInterval   = 1 * 60
	skillICDKey            = "dehya-skill-icd"
	dehyaFieldKey          = "dehya-field-status"
	sanctumPickupExtension = 24 // On recast from Burst/Skill-2 the field duration is extended by 0.4s
)

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

func (c *char) Skill(p map[string]int) (action.Info, error) {
	burstAction := c.UseBurstAction()
	if burstAction != nil {
		return *burstAction, nil
	}
	if c.StatusIsActive(dehyaFieldKey) {
		// If recast has been used, sanctum needs to be placed anew
		if c.hasRecastSkill {
			c.pickUpField()
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
		FlatDmg:            c.c1FlatDmgRatioE * c.MaxHP(),
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

func (c *char) skillDmgHook() {
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
		FlatDmg:            c.c1FlatDmgRatioE * c.MaxHP(),
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}

	// pick up field at start
	c.pickUpField()

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

// pick up field and save current ICD and duration with implicit extension
func (c *char) pickUpField() {
	c.a1Reduction()
	c.sanctumICD = c.StatusDuration(skillICDKey)
	c.sanctumSavedDur = c.StatusExpiry(dehyaFieldKey) + sanctumPickupExtension - c.Core.F // dur gets extended on field recast by a low margin, apparently
	c.Core.Log.NewEvent("sanctum picked up", glog.LogCharacterEvent, c.Index).
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
		FlatDmg:            (c.c1FlatDmgRatioE + skillDotHP[c.TalentLvlSkill()]) * c.MaxHP(),
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

// Active characters within this field have their resistance to interruption increased, (not implemented)
// and when such characters take DMG, a portion of that damage will be mitigated and flow into Redmane's Blood.
// Dehya will then take this DMG over 10s. When the mitigated DMG stored by Redmane's Blood reaches
// or goes over a certain percentage of Dehya's Max HP, she will stop mitigating DMG in this way.
func (c *char) skillHurtHook() {
	// mitigates true dmg
	// should not mitigate corrosion (probably will never be added to sim...)
	c.Core.Events.Subscribe(event.OnPlayerPreHPDrain, func(args ...interface{}) bool {
		di := args[0].(*player.DrainInfo)
		// only mitigate external damage
		if !di.External {
			return false
		}
		// no need to mitigate if 0 dmg
		if di.Amount <= 0 {
			return false
		}
		// field needs to be active for mitigation
		if !c.StatusIsActive(dehyaFieldKey) {
			return false
		}
		// player needs to be in field for mitigation
		if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
			return false
		}
		// ignore self dot
		if di.Abil == skillSelfDoTAbil {
			return false
		}
		// stop mitigating dmg if reached threshold
		if c.skillRedmanesBlood >= 2*c.MaxHP() {
			return false
		}
		beforeAmount := di.Amount
		// calc mitigation based on skill level
		mitigation := di.Amount * skillMitigation[c.TalentLvlSkill()]
		// adjust redmane's blood
		c.skillRedmanesBlood += mitigation
		// modify hp drain
		di.Amount = max(di.Amount-mitigation, 0)
		// log mitigation
		c.Core.Log.NewEvent("dehya mitigating dmg", glog.LogCharacterEvent, c.Index).
			Write("hurt_before", beforeAmount).
			Write("mitigation", mitigation).
			Write("hurt", di.Amount)
		// add self dot status
		c.AddStatus(skillSelfDoTStatus, skillSelfDoTDuration, true)
		// queue up DoT if not already queued
		// -> retrigger should not reset interval (unsure)
		// -> has to be like this otherwise if you keep mitigating between DoT ticks then Dehya will never get damaged
		if c.skillSelfDoTQueued {
			return false
		}
		c.skillSelfDoTQueued = true
		c.QueueCharTask(c.skillSelfDoT, skillSelfDoTStart)
		return false
	}, "dehya-field-dmgtaken")
}

func (c *char) skillSelfDoT() {
	if !c.StatusIsActive(skillSelfDoTStatus) {
		c.skillSelfDoTQueued = false
		return
	}

	// queue next tick
	c.QueueCharTask(c.skillSelfDoT, skillSelfDoTInterval)

	// do not do self DoT if in burst iframes
	if c.Core.Player.Active() == c.Index && c.Core.Player.CurrentState() == action.BurstState {
		return
	}

	// recalculate the dmg on every tick
	dmg := c.skillRedmanesBlood * skillSelfDoTRatio

	// reduce redmane's blood (before considering shield mitigation/a1!)
	c.skillRedmanesBlood = max(c.skillRedmanesBlood-dmg, 0)

	// modify the dmg if a1 is active (redmane's blood is reduced by full amount before this is checked)
	if c.StatusIsActive(a1ReductionKey) {
		dmgBefore := dmg
		dmg *= 1 - a1ReductionMult
		c.Core.Log.NewEvent("dehya a1 reducing redmane's blood dmg", glog.LogCharacterEvent, c.Index).
			Write("dmg_before", dmgBefore).
			Write("dmg", dmg)
	}

	// do self DoT
	// TODO: hack because system is not designed to hit a character directly which is off-field
	// this is true physical dmg so dmg formula/element resist does not matter
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       skillSelfDoTAbil,
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 0,
		FlatDmg:    dmg,
	}
	ap := combat.NewSingleTargetHit(c.Core.Combat.Player().Key())
	snap := c.Snapshot(&ai)
	ae := &combat.AttackEvent{
		Info:        ai,
		Pattern:     ap,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}

	c.Core.Combat.Events.Emit(event.OnPlayerHit, c.Index, ae)
	dmgLeft := c.Core.Player.Shields.OnDamage(c.Index, c.Core.Player.Active(), dmg, ae.Info.Element)
	if dmgLeft > 0 {
		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: c.Index,
			Abil:       ae.Info.Abil,
			Amount:     dmgLeft,
			External:   true,
		})
	}
}
