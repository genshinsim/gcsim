package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// TODO: heizou based frames & my assumptions
var chargeFrames []int

const (
	chargeHitmark = 24

	c6Hitmark = 34
)

func init() {
	chargeFrames = frames.InitAbilSlice(46)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionAttack] = 38
	chargeFrames[action.ActionSkill] = 38
	chargeFrames[action.ActionBurst] = 38
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.09 * 60,
		CanBeDefenseHalted: false,
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6)
	snap := c.Snapshot(&ai) // TODO: where's a snapshot?

	var particleCB combat.AttackCBFunc
	if c.StatusIsActive(skillKey) {
		ai.Abil = "Rebuke: Vaulting Fist"
		ai.HitlagFactor = 0.03
		ai.HitlagHaltFrames = 0.12 * 60
		ap = combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.8}, 4, 5)

		particleCB = c.particleCB

		if c.Base.Cons >= 6 {
			aiC6 := combat.AttackInfo{
				ActorIndex:         c.Index,
				Abil:               "Cherish the Innocent",
				AttackTag:          attacks.AttackTagExtra,
				ICDTag:             attacks.ICDTagNone,
				ICDGroup:           attacks.ICDGroupDefault,
				StrikeType:         attacks.StrikeTypeBlunt,
				Element:            attributes.Cryo,
				Durability:         25,
				Mult:               ai.Mult,
				HitlagFactor:       0.03,
				HitlagHaltFrames:   0.12 * 60,
				CanBeDefenseHalted: false,
			}
			c.Core.QueueAttackWithSnap(aiC6, snap, ap, c6Hitmark)
		}
	}

	var a1CB combat.AttackCBFunc
	if c.StatModIsActive(a1Status) {
		a1CB = c.a1Remove
	}

	c.Core.QueueAttackWithSnap(ai, snap, ap, chargeHitmark, particleCB, a1CB)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
