package cyno

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarConductText = " (Stellar-Conduct)"
)

func (c *char) isStellarRadiance() bool {
	if !c.revelation {
		return false
	}
	return c.StatusIsActive(reactable.PolestarFieldKey)
}

// If Cyno unleashes the Elemental Skill Secret Rite: Chasmic Soulfarer while
// not in the Pactsworn Pathclearer state, he will enter Pactsworn Pathclearer
// for 6s.
func (c *char) revelationOnSkill() {
	if !c.revelation {
		return
	}

	if c.StatusIsActive(burstKey) {
		return
	}

	dur := 6 * 60
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstKey, dur), // TODO: extra duration?
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return c.burstEM
		},
	})
	c.burstSrc = c.Core.F

	// cannot extend the duration of Pactsworn when the state is granted by the
	// revelation passive
	c.burstExtension = 2
	c.a1OnSkillRevelation()

	c.c1RevelationOnPactsworn(dur)
	c.c2RevelationOnPactsworn()
}
