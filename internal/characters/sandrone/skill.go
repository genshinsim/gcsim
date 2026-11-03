package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

// TODO: frames and particle count estimated from in-game footage.
const (
	hovermechKey      = "sandrone-hovermech"
	hovermechDuration = 6 * 60
	hovermechStart    = 20
	prismShot1        = 30
	prismShot2        = 45
	prismRange        = 10
	particleICDKey    = "sandrone-particle-icd"
	particleCount     = 2
)

func init() {
	skillFrames = frames.InitAbilSlice(40) // E -> W
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionCharge] = 30
	skillFrames[action.ActionSkill] = 30
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 30
	skillFrames[action.ActionSwap] = 25
}

// Sandrone boards the Tea Party Tactical Assault Hovermech and hovers forward for 6s.
// When Sandrone starts hovering, if there are opponents nearby, she will also summon a Prismatic
// Resonance Cannon and fire 2 Prism Shots, dealing Cryo DMG.
// Radiance: Stellar Glimmer: The second Prism Shot will instead deal Cryo DMG that is considered
// the corresponding Stellar Glimmer reaction DMG.
// When on the Hovermech, Sandrone will also repair Fagio, causing Decoding Power to swiftly decrease.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.AddStatus(hovermechKey, hovermechDuration, true)
		c.ensureFagioTicker()

		area := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, prismRange)
		if len(c.Core.Combat.EnemiesWithinArea(area, nil)) == 0 {
			return
		}

		// A1 is decided when the cannon is summoned
		a1 := c.a1PrismActive()

		c.QueueCharTask(c.prismShot(false, false), prismShot1-hovermechStart)
		c.QueueCharTask(c.prismShot(true, a1), prismShot2-hovermechStart)
	}, hovermechStart)

	c.SetCDWithDelay(action.ActionSkill, skillCD*60, hovermechStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) prismShot(second, a1 bool) func() {
	return func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Prism Shot",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillPrism[c.TalentLvlSkill()],
		}
		if second {
			lvl := c.TalentLvlSkill()
			if c.toStellar(&ai, skillPrismConduct[lvl], skillPrismSwirl[lvl]) && a1 {
				ai.Abil += " (A1)"
				ai.Mult *= 1 + a1PrismMult
			}
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2)
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
	}
}

// Swapping out takes Sandrone off the Hovermech.
func (c *char) hovermechInit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		c.endHovermech()
	}, hovermechKey+"-swap")
}

// endHovermech is called when Sandrone performs another action that takes her off the Hovermech.
func (c *char) endHovermech() {
	if !c.StatusIsActive(hovermechKey) {
		return
	}
	c.DeleteStatus(hovermechKey)
	c.Core.Log.NewEvent("hovermech ended", glog.LogCharacterEvent, c.Index())
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Cryo, c.ParticleDelay)
}
