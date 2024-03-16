package chiori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillPressFrames []int

const (
	skillHitmark             = 29
	skillDollForm            = 32
	skillDollAttackDelay     = 30 //TODO: shoudl be 0.6s?
	skillRockDollAttackDelay = 72 //TODO: dm thinks it's 1.2s

	// some doll stuff
	dollAttackInterval = 216 // 3.6s
	dollLife           = 17*60 + 30

	skillSecondPressDelay = 0

	particleICDKey = "chiori-particle-icd"
)

func init() {
	skillPressFrames = frames.InitAbilSlice(30) //TODO: i made this up
}

// Dashes nimbly forward with silken steps. Once this dash ends, Chiori will
// summon the automaton doll "Sode" beside her and sweep her blade upward,
// dealing AoE Geo DMG to nearby opponents based on her ATK and DEF. Holding the
// Skill will cause it to behave differently.
//
// Hold Enter Aiming Mode to adjust the dash direction.
//
// Sode
// - Will slash at nearby opponents at intervals, dealing AoE Geo DMG based on
// Chiori's ATK and DEF.
// - While active, when Geo Construct(s) are created nearby, an additional Sode
// will be summoned next to Chiori. Only 1 additional Sode can be summoned in
// this manner, and its duration is independently counted.

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// if this is second press, swap and activate a1
	if c.StatusIsActive(a1TailorMadeWindowKey) {
		return c.skillRecast()
	}
	delay := p["hold"]
	if delay < 0 {
		delay = 0
	}

	// splitting this is currently not necessary but allows for future change
	// if hold does something special
	c.handleSkill(delay)

	//TODO: frames to fix
	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionJump],
		State:           action.SkillState,
	}, nil
}

func (c *char) handleSkill(holdDelay int) {
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			Abil:       "Fluttering Hasode (Blink)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagChioriSkill,
			ICDGroup:   attacks.ICDGroupChioriSkill,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       thrustAtkScaling[c.TalentLvlSkill()],
		}
		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= thrustDefScaling[c.TalentLvlSkill()]
		//TODO: hit box size
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), 0)
	}, skillHitmark+holdDelay)

	//TODO: not sure if this is a task or should be on hit from trust
	c.Core.Tasks.Add(func() {
		//TODO: no clue if kill is on new form, or on skill press
		c.kill(c.skillDoll)
		c.kill(c.constructChecker)

		// create new doll
		doll := newTicker(c.Core, dollLife, nil) // .5s longer
		doll.cb = c.skillDollAttack
		doll.interval = dollAttackInterval // 3.6s

		// attack is delayed
		//TODO: no snapshot??
		c.Core.Tasks.Add(doll.tick, skillDollAttackDelay)

		c.skillDoll = doll

		// associated construct tracker; ticks every 0.3s
		cc := newTicker(c.Core, dollLife, nil)
		cc.cb = c.skillDollConstructCheck
		cc.interval = 18
		c.Core.Tasks.Add(cc.tick, 6) //TODO: i made this delay up; not sure how quick first check is
		c.constructChecker = cc
	}, skillDollForm+holdDelay)

	c.activateA1Window()
	c.SetCDWithDelay(action.ActionSkill, 15*60, holdDelay+1) //TODO: what's the delay here?
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
	return action.Info{
		Frames:          func(action.Action) int { return skillSecondPressDelay + c.Core.Player.Delays.Swap },
		AnimationLength: skillSecondPressDelay + c.Core.Player.Delays.Swap,
		CanQueueAfter:   skillSecondPressDelay + c.Core.Player.Delays.Swap,
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDollAttack() {
	ai := combat.AttackInfo{
		Abil:       "Fluttering Hasode (Turret)",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagChioriSkill,
		ICDGroup:   attacks.ICDGroupChioriSkill,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       turretAtkScaling[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
	ai.FlatDmg *= turretDefScaling[c.TalentLvlSkill()]

	//TODO: hit box size
	hitbox := combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2)
	if c.c1Active() {
		//TODO: c1 modify aoe size
		hitbox = combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2)
	}
	c.Core.QueueAttackWithSnap(ai, snap, hitbox, 0, c.particleCB)
}

func (c *char) rockDollAttack() {
	ai := combat.AttackInfo{
		Abil:       "Fluttering Hasode (Turret - Construct)",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagChioriSkill,
		ICDGroup:   attacks.ICDGroupChioriSkill,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       turretAtkScaling[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
	ai.FlatDmg *= turretDefScaling[c.TalentLvlSkill()]

	//TODO: hit box size
	hitbox := combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2)
	if c.c1Active() {
		//TODO: c1 modify aoe size
		hitbox = combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2)
	}
	c.Core.QueueAttackWithSnap(ai, snap, hitbox, 0, c.particleCB)
}

func (c *char) skillDollConstructCheck() {
	//TODO: not sure if this is correct; all assumptions here
	// if there is still a rock doll alive, we should kill it
	c.kill(c.rockDoll)

	// TODO: i'm assuming the c1 check happens here
	// now check for constructs; if nothing found do nothing
	if !c.c1Active() && c.Core.Constructs.Count() == 0 {
		return
	}

	// create a new rockdoll and delete this ticker (so that we don't self delete)
	rd := newTicker(c.Core, dollLife, nil)
	rd.cb = c.rockDollAttack
	rd.interval = dollAttackInterval
	c.Core.Tasks.Add(rd.tick, skillRockDollAttackDelay)
	c.rockDoll = rd

	// make sure this check doesn't happen again
	c.kill(c.constructChecker)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 3*60, false)
	count := 1.0
	if c.Core.Rand.Float64() < 0.2 {
		count = 2.0
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
}
