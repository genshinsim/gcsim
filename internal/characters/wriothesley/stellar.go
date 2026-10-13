package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarConductText = " (Stellar-Conduct)"
	a1StellarKey       = "wriothesley-a1-stellar"
	a1StellarICDKey    = "wriothesley-a1-stellar-icd"
)

func (c *char) isStellarRadiance() bool {
	return c.StatusIsActive(reactable.PolestarFieldKey)
}

func (c *char) a1StellarEnabled() bool {
	if !c.revelation {
		return false
	}
	return c.isStellarRadiance()
}

// Radiance: Stellar-Conduct: The Ascension Talent "There Shall Be a Plea for Justice" is changed to
// the following:
//
// When Wriothesley's HP is less than 60%, he will obtain a "Gracious Rebuke: Unveiling." This
// enhances the next Charged Attack of his Normal Attack: Forceful Fists of Frost into Luster:
// Vaulting Fist, which will not consume Stamina and restores 30% Max HP to Wriothesley upon hit.
// Wriothesley can gain HP in this way once every 2s.
//
// Additionally, the third and fifth strikes of Repelling Fist enhanced by Chilling Penalty will
// deal 60% and 80% of their original Cryo DMG respectively, while Luster: Vaulting Fist will deal
// 100% of its original AoE Cryo DMG. The aforementioned DMG is considered as Stellar-Conduct DMG,
// and Wriothesley deals 30% increased Stellar-Conduct DMG.
//
// Additionally, when inside a Polestar Field, Wriothesley will enter the Radiance: Stellar-Conduct
// state.
func (c *char) stellarInit() {
	if !c.revelation {
		return
	}
	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1StellarKey, -1),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.AttackTag != attacks.AttackTagDirectStellarConduct {
				return 0
			}
			return 0.3
		},
	})
}

func (c *char) a1StellarCAReady() bool {
	if c.Base.Ascension < 1 {
		return false
	}

	if !c.revelation {
		return false
	}

	if !c.isStellarRadiance() {
		return false
	}

	return c.CurrentHPRatio() < 0.6
}

func (c *char) a1StellarCA() action.Info {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Luster: Vaulting Fist" + stellarConductText,
		AttackTag:        attacks.AttackTagDirectStellarConduct,
		AdditionalTags:   []attacks.AttackTag{attacks.AttackTagWrioSSCCA},
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         110,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
		Mult:             charge[c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.12 * 60,
	}

	// TODO: snapshot timing
	snap := c.Snapshot(&ai)
	c.addC6Buff(&snap)

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -0.8}, 4, 5)

	c.QueueCharTask(func() {
		c.Core.QueueAttackWithSnap(ai, snap, ap, 0, c.a1StellarCACB, c.c1StellarCACB)
		c.c6OnStellarCA(&ai, &snap, &ap)
	}, chargeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

func (c *char) a1StellarCACB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	// do not proc if a1 not active
	if c.StatusIsActive(a1StellarICDKey) {
		return
	}

	c.AddStatus(a1StellarICDKey, 2*60, true)

	// heal
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Index(),
		Message: a1HealName,
		Src:     0.3 * c.MaxHP(),
		Bonus:   c.Stat(attributes.Heal),
	})
}
