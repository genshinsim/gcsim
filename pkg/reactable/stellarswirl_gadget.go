package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

var (
	sswStackMult   = []float64{0, 2, 2, 3, 3, 3, 3}
	sswStackRadius = []float64{0, 4, 4, 6, 6, 6, 6}
)

const (
	sswDuration      = 181
	sswThinkInterval = 181
)

type StellarVortex struct {
	*gadget.Gadget
	r         *Reactable
	fieldArea info.AttackPattern
}

func (r *Reactable) newStellarVortex() *StellarVortex {
	p := &StellarVortex{r: r}

	p.Gadget = gadget.New(r.core, r.core.Combat.Player().Pos(), 1, info.GadgetTypPolestarField)
	p.ThinkInterval = sswThinkInterval
	p.Duration = sswDuration
	p.OnKill = p.explode
	r.core.Combat.AddGadget(p)
	return p
}

func (p *StellarVortex) explode() {
	owner := p.r.sswOwner()
	contribMap := p.r.sswContributors()
	p.r.removeSSwContributors()
	ai := info.AttackInfo{
		ActorIndex:       owner,
		DamageSrc:        p.Key(),
		Abil:             "Stellar Swirl Detonation",
		AttackTag:        attacks.AttackTagReactionStellarSwirl,
		ICDTag:           attacks.ICDTagNone,      // TODO: use stellar swirl ICD tag
		ICDGroup:         attacks.ICDGroupDefault, // TODO: use stellar swirl ICD group
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Cryo,
		Durability:       25,
		IgnoreDefPercent: 1,
	}

	stacks := p.r.sswStacks()
	ap := combat.NewCircleHitOnTarget(p, nil, sswStackRadius[stacks])

	for _, e := range p.r.core.Combat.Enemies() {
		if willHit, _ := e.AttackWillLand(ap); !willHit {
			continue
		}
		ai, snap := p.r.calcStellarSwirlDmg(e, ai, ap, contribMap, sswStackMult[stacks])
		ai.ActorIndex = owner
		p.r.core.QueueAttackWithSnap(ai, snap, combat.NewSingleTargetHit(e.Key()), 0)
	}

	p.r.core.Flags.Custom[sswStackKey] = 0

	for _, c := range p.r.core.Player.Chars() {
		// TODO: find exact duration of jump buff
		// TODO: if chars are currently jumping, need to make the jump boosted height
		c.AddStatus(player.StellarSwirlAirborneBuff, 6*60, true) // buff needs to be removed on jump?
	}
}

func (p *StellarVortex) HandleAttack(atk *info.AttackEvent) float64 { return 0 }

func (p *StellarVortex) Tick() {
	p.Gadget.Tick()
}
