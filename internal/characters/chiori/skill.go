package chiori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames            [][]int
	skillHitmarks          = []int{21, 37}
	skillCDStarts          = []int{19, 34} // same as doll spawn
	skillA1WindowStarts    = []int{26, 42}
	skillA1WindowDurations = []int{78, 77}
)

const (
	skillDollDuration               = int(17.5 * 60)
	skillDollStartDelay             = int(0.6 * 60)
	skillDollAttackInterval         = int(3.6 * 60)
	skillDollConstructCheckInterval = int(0.5 * 60)
	skillDollAttackDelay            = 5 // should be 0.08s
	skillDollXOffset                = 1.2
	skillDollYOffset                = -0.3
	skillDollAoE                    = 1.2

	skillRockDollStartDelay = int(1.2 * 60)

	skillCD = 16 * 60

	particleICDKey = "chiori-particle-icd"
)

func init() {
	skillFrames = make([][]int, 2)

	// Tap E
	skillFrames[0] = frames.InitAbilSlice(51) // Tap E -> Walk
	skillFrames[0][action.ActionAttack] = 42
	skillFrames[0][action.ActionSkill] = 30
	skillFrames[0][action.ActionBurst] = 43
	skillFrames[0][action.ActionDash] = 42
	skillFrames[0][action.ActionJump] = 42
	skillFrames[0][action.ActionSwap] = 49

	// Hold E
	skillFrames[1] = frames.InitAbilSlice(88) // Hold E -> N1/Q/D/J
	skillFrames[1][action.ActionLowPlunge] = 52
	skillFrames[1][action.ActionSkill] = 44
	skillFrames[1][action.ActionWalk] = 86
	skillFrames[1][action.ActionSwap] = 87
}

// Dashes nimbly forward with silken steps. Once this dash ends, Chiori will
// summon the automaton doll "Tamoto" beside her and sweep her blade upward,
// dealing AoE Geo DMG to nearby opponents based on her ATK and DEF. Holding the
// Skill will cause it to behave differently.
//
// Hold Enter Aiming Mode to adjust the dash direction.
//
// Tamoto
// - Will slash at nearby opponents at intervals, dealing AoE Geo DMG based on
// Chiori's ATK and DEF.
// - While active, when Geo Construct(s) are created nearby, an additional Tamoto
// will be summoned next to Chiori. Only 1 additional Tamoto can be summoned in
// this manner, and its duration is independently counted.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	if hold < 0 {
		hold = 0
	}
	if hold > 1 {
		hold = 1
	}

	// if this is second press, swap and activate a1
	if c.StatusIsActive(a1WindowKey) {
		return c.skillRecast()
	}

	// splitting this is currently not necessary but allows for future change
	// if hold does something special
	c.handleSkill(hold)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[hold]),
		AnimationLength: skillFrames[hold][action.InvalidAction],
		CanQueueAfter:   skillFrames[hold][action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillRecast() (action.Info, error) {
	c.a1Tapestry()
	// find next char
	next := c.Index + 1
	if next >= len(c.Core.Player.Chars()) {
		next = 0
	}
	k := c.Core.Player.ByIndex(next).Base.Key
	c.Core.Tasks.Add(func() {
		c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "forcing swap to ", k.String())
		c.Core.Player.Exec(action.ActionSwap, k, nil)
	}, 1)
	// TODO: doesn't seem like this duration actually matters because of the forced swap, it just needs to cover the 1f until the swap is executed
	return action.Info{
		Frames:          func(action.Action) int { return c.Core.Player.Delays.Swap },
		AnimationLength: c.Core.Player.Delays.Swap,
		CanQueueAfter:   c.Core.Player.Delays.Swap,
		State:           action.SkillState,
	}, nil
}

func (c *char) handleSkill(hold int) {
	// handle upward sweep
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			Abil:       "Fluttering Hasode (Upward Sweep)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagChioriSkill,
			ICDGroup:   attacks.ICDGroupChioriSkill,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   50,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       thrustAtkScaling[c.TalentLvlSkill()],
		}

		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= thrustDefScaling[c.TalentLvlSkill()]

		c.Core.QueueAttackWithSnap(ai, snap, combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 5.5, 3.5), 0)
	}, skillHitmarks[hold])

	// trigger cd and activate a1 window
	c.SetCDWithDelay(action.ActionSkill, skillCD, skillCDStarts[hold])
	c.activateA1Window(skillA1WindowStarts[hold], skillA1WindowDurations[hold])

	// handle doll spawn
	c.Core.Tasks.Add(func() {
		// create 1st doll
		c.createDoll()

		// create construct checker that should spawn the rock doll if not c1
		if !c.c1Active {
			c.createDollConstructChecker()
			return
		}

		// create rock doll from c1 and trigger a4
		c.Core.Log.NewEvent("c1 spawning rock doll", glog.LogCharacterEvent, c.Index)
		c.createRockDoll()
		c.applyA4Buff()
	}, skillCDStarts[hold])
}

func (c *char) createDoll() {
	// kill existing doll
	c.kill(c.skillDoll)

	// determine doll pos
	player := c.Core.Combat.Player()
	dollPos := geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{X: skillDollXOffset, Y: skillDollYOffset},
		player.Direction(),
	)

	c.Core.Log.NewEvent("spawning doll", glog.LogCharacterEvent, c.Index)

	// spawn new doll
	doll := newTicker(c.Core, skillDollDuration, nil)
	doll.cb = c.skillDollAttack(c.Core.F, "Fluttering Hasode (Tamato)", dollPos)
	doll.interval = skillDollAttackInterval
	c.Core.Tasks.Add(doll.tick, skillDollStartDelay)
	c.skillDoll = doll
}

func (c *char) createDollConstructChecker() {
	// kill existing construct checker
	c.kill(c.constructChecker)

	// spawn associated construct checker to spawn the rock doll
	cc := newTicker(c.Core, skillDollDuration, nil)
	cc.cb = c.skillDollConstructCheck
	cc.interval = skillDollConstructCheckInterval
	cc.tick() // start ticking at t = 0s
	c.constructChecker = cc
}

func (c *char) skillDollAttack(src int, abil string, pos geometry.Point) func() {
	return func() {
		c.Core.Tasks.Add(func() {
			ai := combat.AttackInfo{
				Abil:       abil,
				ActorIndex: c.Index,
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagChioriSkill,
				ICDGroup:   attacks.ICDGroupChioriSkill,
				StrikeType: attacks.StrikeTypeBlunt,
				PoiseDMG:   0,
				Element:    attributes.Geo,
				Durability: 25,
				Mult:       turretAtkScaling[c.TalentLvlSkill()],
			}

			snap := c.Snapshot(&ai)
			ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
			ai.FlatDmg *= turretDefScaling[c.TalentLvlSkill()]

			// if the player has an attack target it will always choose this enemy
			// so just need to make sure that it is within the search AoE
			t := c.Core.Combat.PrimaryTarget()
			if !t.IsWithinArea(combat.NewCircleHitOnTarget(pos, nil, c.skillSearchAoE)) {
				return
			}

			c.Core.Log.NewEvent("doll attacking", glog.LogCharacterEvent, c.Index).Write("src", src)

			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(t, nil, skillDollAoE), 0, c.particleCB)
		}, skillDollAttackDelay)
	}
}

func (c *char) skillDollConstructCheck() {
	// cannot spawn rock doll if construct exists
	if c.rockDoll != nil && c.rockDoll.alive {
		return
	}
	// TODO: technically should check for constructs within 30m radius of the skill doll
	// doll pos is already passed to the attack func so can reuse that here
	// not too important and construct handler doesn't directly expose that so not doing that for now
	if c.Core.Constructs.Count() == 0 {
		return
	}

	c.Core.Log.NewEvent("construct spawning rock doll", glog.LogCharacterEvent, c.Index)
	c.createRockDoll()

	// make sure this check doesn't happen again
	c.kill(c.constructChecker)
}

func (c *char) createRockDoll() {
	// kill existing
	c.kill(c.rockDoll)

	// determine doll pos
	player := c.Core.Combat.Player()
	dollPos := geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{X: skillDollXOffset, Y: skillDollYOffset},
		player.Direction(),
	)

	// spawn new rock doll
	rd := newTicker(c.Core, skillDollDuration, nil)
	rd.cb = c.skillDollAttack(c.Core.F, "Fluttering Hasode (Tamato - Construct)", dollPos)
	rd.interval = skillDollAttackInterval
	c.Core.Tasks.Add(rd.tick, skillRockDollStartDelay)
	c.rockDoll = rd
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 3*60, true)

	count := 1.0
	if c.Core.Rand.Float64() < 0.2 {
		count = 2.0
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
}
