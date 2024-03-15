package gaming

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// TODO: Kazuha Plunge Frames
var plungeFrames []int
var ePlungeKey = "Charmed CloudStrider"

const (
	plungePressHitmark = 36
	particleICDKey     = "gaming-particle-icd"
)

// TODO: missing plunge -> skill
// TODO: Taken from Kazuha
func init() {
	// skill (press) -> high plunge -> x
	plungeFrames = frames.InitAbilSlice(55) // max
	plungeFrames[action.ActionDash] = 43
	plungeFrames[action.ActionJump] = 50
	plungeFrames[action.ActionSwap] = 50

	// from video https://streamable.com/arcm3e is 73f from E -> P -> N1
	plungeFrames[action.ActionAttack] = 50
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	// last action must be skill without glide cancel
	if c.Core.Player.LastAction.Type != action.ActionSkill ||
		c.Core.Player.LastAction.Param["glide_cancel"] != 0 {
		c.Core.Log.NewEvent("only plunge after skill without glide cancel", glog.LogActionEvent, c.Index).
			Write("action", action.ActionLowPlunge)
		return action.Info{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}, nil
	}

	act := action.Info{
		State: action.PlungeAttackState,
	}

	//TODO: is this accurate?? these should be the hitmarks
	var hitmark int
	hitmark = plungePressHitmark
	act.Frames = frames.NewAbilFunc(plungeFrames)
	act.AnimationLength = plungeFrames[action.InvalidAction]
	act.CanQueueAfter = plungeFrames[action.ActionDash] // earliest cancel

	// aoe dmg
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           ePlungeKey,
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}

	// only drain HP when above 10% HP
	if c.CurrentHPRatio() > 10 {
		currentHP := c.CurrentHP()
		maxHP := c.MaxHP()
		hpdrain := 0.15 * currentHP
		// The HP consumption from using this skill can only bring him to 10% HP.
		if (currentHP-hpdrain)/maxHP <= 0.1 {
			hpdrain = currentHP - 0.1*maxHP
		}
		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: c.Index,
			Abil:       "Charmed Cloudstrider",
			Amount:     hpdrain,
		})
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 4.5),
		hitmark,
		hitmark,
		c.particleCB,
		c.makeA1CB(),
		c.onHit,
	)

	return act, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, gamingICD*60, false)

	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Pyro, c.ParticleDelay)
}

func (c *char) onHit(a combat.AttackCB) {
	c.spawnManchai(a)
	c.c4()
}
