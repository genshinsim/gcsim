package chiori

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
	// can hold to aim
	if p["hold"] == 1 {
		//TODO: handle skill hold
		return action.Info{}, errors.New("chiori skill hold not implemented yet")
	}

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
	}, skillHitmark)

	//TODO: not sure if this is a task or should be on hit from trust
	c.Core.Tasks.Add(func() {
		//TODO: no clue if kill is on new form, or on skill press
		c.kill(c.skillDoll)
		c.kill(c.constructChecker)

		// create new doll
		doll := newTicker(c.Core, dollLife) // .5s longer
		doll.cb = c.skillDollAttack
		doll.interval = dollAttackInterval // 3.6s

		// attack is delayed
		//TODO: no snapshot??
		c.Core.Tasks.Add(doll.tick, skillDollAttackDelay)

		c.skillDoll = doll

		// associated construct tracker; ticks every 0.3s
		cc := newTicker(c.Core, dollLife)
		cc.cb = c.skillDollConstructCheck
		cc.interval = 18
		c.Core.Tasks.Add(cc.tick, 6) //TODO: i made this delay up; not sure how quick first check is

		c.constructChecker = cc
	}, skillDollForm)

	c.a1Window()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 1) //TODO: what's the delay here?

	//TODO: frames to fix
	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionJump],
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
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), 0)
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
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), 0)
}

func (c *char) skillDollConstructCheck() {
	//TODO: not sure if this is correct; all assumptions here
	// if there is still a rock doll alive, we should kill it
	c.kill(c.rockDoll)

	// now check for constructs; if nothing found do nothing
	if c.Core.Constructs.Count() == 0 {
		return
	}

	// create a new rockdoll and delete this ticker (so that we don't self delete)
	rd := newTicker(c.Core, dollLife)
	rd.cb = c.rockDollAttack
	rd.interval = dollAttackInterval
	c.Core.Tasks.Add(rd.tick, skillRockDollAttackDelay)
	c.rockDoll = rd

	c.a4()

	// make sure this check doesn't happen again
	c.kill(c.constructChecker)
}
