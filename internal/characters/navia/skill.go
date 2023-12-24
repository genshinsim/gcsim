package navia

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

var (
	skillFrames     [][][]int
	skillMultiplier = []float64{
		0,                   // 0 hits
		1,                   // 1 hits
		1.05000000074505806, // 2 hits
		1.10000000149011612, // 3 hits etc
		1.15000000596046448,
		1.20000000298023224,
		1.36000001430511475,
		1.4000000059604645,
		1.6000000238418579,
		1.6660000085830688,
		1.8999999761581421,
		2,
	}
	hitscans = [][]float64{
		// width, angle, x offset, y offset
		{0.9, 0, 0, 0},
		{0.25, 7.911, 0, 0},
		{0.25, -1.826, 0, 0},
		{0.25, -4.325, 0, 0},
		{0.25, 0.773, 0, 0},
		{0.25, 6.209, 0, 0},
		{0.25, -2.752, 0, 0},
		{0.25, 7.845, 0.01, 0.01},
		{0.25, -7.933, -0.01, -0.01},
		{0.25, 2.626, 0, 0},
		{0.25, -9.993, 0, 0},
	}
)

const (
	travelDelay = 9

	skillPressCDStart = 11

	skillHoldCDStart     = 41
	skillMaxHoldDuration = 241 // add 1f to account for hold being set to 1 for activation

	bulletBoxLength = 11.5

	particleICDKey = "navia-particle-icd"
	arkheICDKey    = "navia-arkhe-icd"
)

func init() {
	skillFrames = make([][][]int, 2)

	// press
	skillFrames[0] = make([][]int, 2)

	// normal
	skillFrames[0][0] = frames.InitAbilSlice(40) // E -> E/Q
	skillFrames[0][0][action.ActionAttack] = 38
	skillFrames[0][0][action.ActionDash] = 24
	skillFrames[0][0][action.ActionJump] = 24
	skillFrames[0][0][action.ActionWalk] = 35
	skillFrames[0][0][action.ActionSwap] = 38

	// >= 3 shrapnel
	skillFrames[0][1] = frames.InitAbilSlice(41) // E -> E/Q
	skillFrames[0][1][action.ActionAttack] = 40
	skillFrames[0][1][action.ActionDash] = 26
	skillFrames[0][1][action.ActionJump] = 24
	skillFrames[0][1][action.ActionWalk] = 39
	skillFrames[0][1][action.ActionSwap] = 40

	// hold
	skillFrames[1] = make([][]int, 2)

	// normal
	skillFrames[1][0] = frames.InitAbilSlice(71) // E -> E/Q
	skillFrames[1][0][action.ActionAttack] = 70
	skillFrames[1][0][action.ActionDash] = 54
	skillFrames[1][0][action.ActionJump] = 54
	skillFrames[1][0][action.ActionWalk] = 65
	skillFrames[1][0][action.ActionSwap] = 69

	// >= 3 shrapnel
	skillFrames[1][1] = frames.InitAbilSlice(73) // E -> E/Q
	skillFrames[1][1][action.ActionDash] = 56
	skillFrames[1][1][action.ActionJump] = 56
	skillFrames[1][1][action.ActionWalk] = 70
	skillFrames[1][1][action.ActionSwap] = 71
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// frame related setup
	holdIndex := 0
	shrapnelIndex := 0
	firingTime := skillPressCDStart // assume tap as default

	// check hold
	// hold is added to minimum hold length
	hold := max(p["hold"], 0)
	if hold > 0 {
		// use hold skill frames
		holdIndex = 1
		// cap hold to max hold duration over minimum hold
		if hold > skillMaxHoldDuration {
			hold = skillMaxHoldDuration
		}
		// subtract 1 to account for needing to supply > 0 to indicate hold
		hold -= 1
		// calc firingTime
		firingTime = skillHoldCDStart + hold

		// crystal pulling related
		// At 0.2, tap is converted to hold, and the suction begins
		firingTimeF := c.Core.F + firingTime
		for i := 12; i < firingTime; i += 30 {
			c.pullCrystals(firingTimeF, i)
		}
	}
	c.SetCDWithDelay(action.ActionSkill, 9*60, firingTime)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rosula Shardshot",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
	}

	c.QueueCharTask(
		func() {
			if c.shrapnel >= 3 {
				shrapnelIndex = 1
			}
			c.Core.Log.NewEvent(fmt.Sprintf("firing %v crystal shrapnel", c.shrapnel), glog.LogCharacterEvent, c.Index)

			// Calculate buffs based on excess shrapnel
			excess := float64(max(c.shrapnel-3, 0))

			// snap and add buffs
			snap := c.Snapshot(&ai)
			c.addShrapnelBuffs(&snap, excess)

			// Looks for enemies in the path of each bullet
			// Initially trims enemies to check by scanning only the hit zone
			shots := 5 + min(c.shrapnel, 3)*2
			for _, t := range c.Core.Combat.EnemiesWithinArea(
				combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: -0.20568, Y: -0.043841}, 4.0722, 11.5461),
				nil,
			) {
				// Tallies up the hits
				hits := 0
				for i := 0; i < shots; i++ {
					if ok, _ := t.AttackWillLand(
						combat.NewBoxHitOnTarget(
							c.Core.Combat.Player(),
							geometry.Point{X: hitscans[i][2], Y: hitscans[i][3]}.Rotate(geometry.DegreesToDirection(hitscans[i][1])),
							hitscans[i][0],
							bulletBoxLength,
						)); ok {
						hits++
					}
				}
				c.Core.Log.NewEvent(fmt.Sprintf("target %v hit %v times", t.Key(), hits), glog.LogCharacterEvent, c.Index)
				// Applies damage based on the hits
				ai.Mult = skillshotgun[c.TalentLvlSkill()] * skillMultiplier[hits]
				c.Core.QueueAttackWithSnap(
					ai,
					snap,
					combat.NewSingleTargetHit(t.Key()),
					travelDelay,
					c.particleCB,
					c.c2(),
				)
			}
			c.surgingBlade(excess)

			// trigger A1 and C1 on firing
			c.a1()
			c.c1(c.shrapnel)

			// remove the shrapnel after firing
			if c.Base.Cons < 6 {
				c.shrapnel = 0
			} else {
				c.shrapnel = max(c.shrapnel-3, 0) // C6 keeps any more than the three
			}
		},
		firingTime,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[holdIndex][shrapnelIndex][next] + hold },
		AnimationLength: skillFrames[holdIndex][1][action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[holdIndex][0][action.ActionJump] + hold,
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	e := a.Target.(*enemy.Enemy)
	if e.Type() != targets.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)

	count := 3.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
}

// add the buffs by modifying snap
// needs to be done this way so that excess calculated during firing is also used for that firing's surgingBlade
func (c *char) addShrapnelBuffs(snap *combat.Snapshot, excess float64) {
	dmg := 0.15 * excess
	cr := 0.0
	cd := 0.0
	if c.Base.Cons >= 2 {
		cr = 0.12 * excess
	}
	if c.Base.Cons >= 6 {
		cd = 0.45 * excess
	}
	snap.Stats[attributes.DmgP] += dmg
	snap.Stats[attributes.CR] += cr
	snap.Stats[attributes.CD] += cd
	c.Core.Log.NewEvent("adding shrapnel buffs", glog.LogCharacterEvent, c.Index).Write("dmg%", dmg).Write("cr", cr).Write("cd", cd)
}

// shrapnelGain adds Shrapnel Stacks when a Crystallise Shield is picked up.
// Stacks should last 300s but this is way too long to bother
// When a character in the party obtains an Elemental Shard created from the Crystallize reaction,
// Navia will gain 1 Crystal Shrapnel charge. Navia can hold up to 6 charges of Crystal Shrapnel at once.
// Each time Crystal Shrapnel gain is triggered, the duration of the Shards you have already will be reset.
func (c *char) shrapnelGain() {
	c.Core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		// Check shield
		shd := args[0].(shield.Shield)
		if shd.Type() != shield.Crystallize {
			return false
		}

		if c.shrapnel < 6 {
			c.shrapnel++
			c.Core.Log.NewEvent("Crystal Shrapnel gained from Crystallise", glog.LogCharacterEvent, c.Index).Write("shrapnel", c.shrapnel)
		}
		return false
	}, "shrapnel-gain")
}

func (c *char) surgingBlade(excess float64) {
	if c.StatusIsActive(arkheICDKey) {
		return
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Surging Blade",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Geo,
		Durability: 0,
		Mult:       skillblade[c.TalentLvlSkill()],
	}

	// determine attack pos
	player := c.Core.Combat.Player()
	// shotgun area
	e := c.Core.Combat.ClosestEnemyWithinArea(combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: -0.20568, Y: -0.043841}, 4.0722, 11.5461), nil)
	// pos is at player + Y: 3.6 by default
	pos := geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 3.6}, player.Direction())
	if e != nil {
		// enemy in shotgun area: use their pos
		pos = e.Pos()
	}

	// aligned cd trigger is delayed and only once it's triggered the aligned attack task should be queued
	c.QueueCharTask(func() {
		c.AddStatus(arkheICDKey, 7*60, true)
		c.QueueCharTask(func() {
			snap := c.Snapshot(&ai)
			c.addShrapnelBuffs(&snap, excess)
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHitOnTarget(pos, nil, 3),
				0,
			)
		}, 36)
	}, 28)
}

// pull crystals to Navia. This has a range of 12m from datamine, so
// check every every 30f.
func (c *char) pullCrystals(firingTimeF, i int) {
	c.Core.Tasks.Add(func() {
		for j := 0; j < c.Core.Combat.GadgetCount(); j++ {
			cs, ok := c.Core.Combat.Gadget(j).(*reactable.CrystallizeShard)
			// skip if not a shard
			if !ok {
				continue
			}
			// If shard is out of 12m range, skip
			if !cs.IsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 12)) {
				continue
			}

			// approximate sucking in as 0.4m per frame (~8m distance took 20f to arrive at gorou)
			distance := cs.Pos().Distance(c.Core.Combat.Player().Pos())
			travel := int(math.Ceil(distance / 0.4))
			// if the crystal won't arrive before the shot is fired, skip
			if c.Core.F+travel >= firingTimeF {
				continue
			}
			// special check to account for edge case if shard just spawned and will arrive before it can be picked up
			if c.Core.F+travel < cs.EarliestPickup {
				continue
			}

			c.Core.Tasks.Add(func() {
				cs.AddShieldKillShard()
			}, travel)
		}
	}, i)
}
