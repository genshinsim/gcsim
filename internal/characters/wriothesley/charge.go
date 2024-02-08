package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var chargeFrames []int

const (
	chargeHitmark = 19
)

func init() {
	chargeFrames = frames.InitAbilSlice(52) // CA -> N1/E/Q
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionWalk] = 51
	chargeFrames[action.ActionSwap] = 49
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.LastAction.Type != action.ActionAttack {
		return action.Info{}, player.ErrInvalidChargeAction
	}

	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Charge Attack",
		AttackTag:        attacks.AttackTagExtra,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         110,
		Element:          attributes.Cryo,
		Durability:       25,
		Mult:             charge[c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.09 * 60,
	}

	// TODO: snapshot timing
	snap := c.Snapshot(&ai)
	var ap combat.AttackPattern
	var rebukeCB combat.AttackCBFunc
	var particleCB combat.AttackCBFunc
	var c6Attack bool
	if c.Base.Ascension >= 1 {
		if c.Base.Cons >= 1 {
			rebukeCB, c6Attack = c.c1(&ai, &snap)
		} else {
			rebukeCB = c.a1(&ai, &snap)
		}

		if rebukeCB != nil {
			particleCB = c.particleCB
			ap = combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.8}, 4, 5)
		} else {
			ap = combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6)
		}
	}

	c.Core.QueueAttackWithSnap(ai, snap, ap, chargeHitmark, rebukeCB, particleCB)
	// When released, it will also unleash an icicle that deals 100% of Rebuke: Vaulting Fist's Base
	// DMG. DMG dealt this way is regarded as Charged Attack DMG.
	// You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
	if c6Attack {
		ai.Abil += " (C6)"
		ai.StrikeType = attacks.StrikeTypeDefault
		ai.PoiseDMG = 50
		c.Core.QueueAttackWithSnap(ai, snap, ap, chargeHitmark, rebukeCB, particleCB)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}
