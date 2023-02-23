package fischl

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillOzSpawn     = 32
	skillRecastCD    = 92 // 2f CD delay
	skillRecastCDKey = "fischl-skill-recast-cd"
	particleICDKey   = "fischl-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 16
	skillFrames[action.ActionSwap] = 42

	skillRecastFrames = frames.InitAbilSlice(37)
	skillRecastFrames[action.ActionAttack] = 36
	skillRecastFrames[action.ActionBurst] = 35
	skillRecastFrames[action.ActionDash] = 4
	skillRecastFrames[action.ActionJump] = 5
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["recast"] != 0 && c.ozActive && !c.StatusIsActive(skillRecastCDKey) {
		return c.skillRecast()
	}
	// always trigger electro no ICD on initial summon
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Oz (Summon)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupFischl,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       birdSum[c.TalentLvlSkill()],
	}

	radius := 2.0
	if c.Base.Cons >= 2 {
		ai.Mult += 2
		radius = 3
	}
	// hitmark is 5 frames after oz spawns
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), combat.Point{Y: 1.5}, radius),
		skillOzSpawn,
		skillOzSpawn+5,
	)

	// CD Delay is 18 frames, but things break if Delay > CanQueueAfter
	// so we add 18 to the duration instead. this probably mess up CDR stuff
	c.SetCD(action.ActionSkill, 25*60+18) // 18 frames until CD starts

	c.Core.Tasks.Add(func() {
		c.AddStatus(skillRecastCDKey, skillRecastCD, false)
	}, 18)

	// set oz to active at the start of the action
	c.ozActive = true
	// set on field oz to be this one
	c.queueOz("Skill", skillOzSpawn)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.1*60, true)
	if c.Core.Rand.Float64() < .67 {
		// TODO: this delay used to be 120
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
	}
}

func (c *char) skillRecast() action.ActionInfo {
	c.AddStatus(skillRecastCDKey, skillRecastCD, false)
	c.Core.Tasks.Add(func() {
		c.ozTickSrc = c.Core.F // reset attack timer
		c.Core.Tasks.Add(c.ozTick(c.ozTickSrc), 60)
		c.ozSnapshot.Snapshot = c.Snapshot(&c.ozSnapshot.Info)
		c.Core.Log.NewEvent("Recasting oz", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+60)
	}, 2) // 2f delay
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) queueOz(src string, ozSpawn int) {
	// calculate oz duration
	dur := 600
	if c.Base.Cons == 6 {
		dur += 120
	}
	spawnFn := func() {
		// setup variables for tracking oz
		c.ozSource = c.Core.F
		c.ozTickSrc = c.Core.F
		c.ozActiveUntil = c.Core.F + dur
		// queue up oz removal at the end of the duration for gcsl conditional
		c.Core.Tasks.Add(c.removeOz(c.Core.F), dur)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Oz (%v)", src),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupFischl,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       birdAtk[c.TalentLvlSkill()],
		}
		player := c.Core.Combat.Player()
		c.ozPos = combat.CalcOffsetPoint(player.Pos(), combat.Point{Y: 1.5}, player.Direction())

		snap := c.Snapshot(&ai)
		c.ozSnapshot = combat.AttackEvent{
			Info:        ai,
			Snapshot:    snap,
			SourceFrame: c.Core.F,
		}
		c.ozSnapshot.Callbacks = append(c.ozSnapshot.Callbacks, c.particleCB)

		c.Core.Tasks.Add(c.ozTick(c.Core.F), 60)
		c.Core.Log.NewEvent("Oz activated", glog.LogCharacterEvent, c.Index).
			Write("source", src).
			Write("expected end", c.ozActiveUntil).
			Write("next expected tick", c.Core.F+60)
	}
	if ozSpawn > 0 {
		c.Core.Tasks.Add(spawnFn, ozSpawn)
	} else {
		spawnFn()
	}
}

func (c *char) ozTick(src int) func() {
	return func() {
		// if src != ozSource then this is no longer the same oz, do nothing
		if src != c.ozTickSrc {
			return
		}
		c.Core.Log.NewEvent("Oz ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+60).
			Write("active", c.ozActiveUntil).
			Write("src", src)
		// trigger damage
		ae := c.ozSnapshot
		ae.Pattern = combat.NewBoxHit(
			c.ozPos,
			c.Core.Combat.PrimaryTarget(),
			combat.Point{Y: -0.5},
			0.1,
			1,
		)
		c.Core.QueueAttackEvent(&ae, c.ozTravel)

		// queue up next hit only if next hit oz is still active
		if c.Core.F+60 <= c.ozActiveUntil {
			c.Core.Tasks.Add(c.ozTick(src), 60)
		}
	}
}

func (c *char) removeOz(src int) func() {
	return func() {
		// if src != ozSource then this is no longer the same oz, do nothing
		if c.ozSource != src {
			c.Core.Log.NewEvent("Oz not removed, src changed", glog.LogCharacterEvent, c.Index).
				Write("src", src)
			return
		}
		c.Core.Log.NewEvent("Oz removed", glog.LogCharacterEvent, c.Index).
			Write("src", src)
		c.ozActive = false
	}
}
