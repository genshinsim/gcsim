package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	PolestarFieldStacksKey = PolestarFieldKey + "-stacks"
	recordKey              = PolestarFieldKey + "-recorded"
	fieldDur               = 6 * 60
	stackICDKey            = PolestarFieldKey + "-icd"
	stackICDDur            = 0.1 * 60
	maxStacks              = 12
	thinkInterval          = 4 * 60
	buffDur                = thinkInterval + 1 // one frame higher so it isn't constantly expiring and then reapplying on the same frame
)

var stellarConductBuff = []float64{
	0.2,
	0.29,
	0.3,
	0.31,
	0.32,
	0.33,
	0.34,
	0.35,
	0.36,
	0.37,
	0.38,
	0.39,
	0.4,
}

type PolestarField struct {
	*gadget.Gadget
	mCryo     []float64
	mElectro  []float64
	fieldArea info.AttackPattern
}

func (r *Reactable) newPolestarField() *PolestarField {
	p := &PolestarField{}

	p.Gadget = gadget.New(r.core, r.self.Pos(), 1, info.GadgetTypPolestarField)
	p.ThinkInterval = thinkInterval
	p.Duration = fieldDur
	p.mCryo = make([]float64, attributes.EndStatType)
	p.mElectro = make([]float64, attributes.EndStatType)
	p.OnThinkInterval = p.refreshBuffs
	p.OnKill = p.unsub
	p.fieldArea = combat.NewCircleHitOnTarget(p, nil, 15) // R15 shape for the field
	p.Core.Events.Subscribe(event.OnElementApplied, func(args ...any) {
		target := args[0].(info.Target)
		if target.Type() != info.TargettableEnemy {
			return
		}

		if !target.IsWithinArea(p.fieldArea) {
			return
		}

		ai := args[1].(*info.AttackInfo)

		switch ai.Element {
		case attributes.Electro:
		case attributes.Cryo:
		// TODO: does adding frozen aura count?
		default:
			return
		}

		if ai.Durability <= info.ZeroDur {
			return
		}

		char := p.Core.Player.Chars()[ai.ActorIndex]
		if char.StatusIsActive(stackICDKey) {
			return
		}

		char.AddStatus(stackICDKey, stackICDDur, true)

		p.Core.Flags.Custom[recordKey] = min(p.Core.Flags.Custom[recordKey]+1, maxStacks)
		p.Core.Log.NewEvent("Adding polestar field stored stacks", glog.LogElementEvent, -1).Write("new_stacks", int(p.Core.Flags.Custom[recordKey]))
	}, p.subscriptionKey())

	p.refreshBuffs()

	r.core.Combat.AddGadget(p)
	return p
}

func (p *PolestarField) HandleAttack(atk *info.AttackEvent) float64 { return 0 }

func (p *PolestarField) resetDuration() {
	p.Duration = fieldDur
}

func (p *PolestarField) subscriptionKey() string {
	return "stellarconduct-hook-" + p.Pos().String()
}

func (p *PolestarField) unsub() {
	p.Core.Events.Unsubscribe(event.OnElementApplied, p.subscriptionKey())
}

func (p *PolestarField) refreshBuffs() {
	oldStacks := int(p.Core.Flags.Custom[PolestarFieldStacksKey])
	newStacks := int(p.Core.Flags.Custom[recordKey])
	p.Core.Flags.Custom[PolestarFieldStacksKey] = float64(newStacks)
	p.Core.Flags.Custom[recordKey] = 0
	p.Core.Log.NewEvent("Updating polestar field buff stacks", glog.LogElementEvent, -1).Write("old_stacks", oldStacks).Write("new_stacks", newStacks)

	for _, e := range p.Core.Combat.EnemiesWithinArea(p.fieldArea, nil) {
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(StellarConductShredKey, buffDur),
			Ele:   attributes.Physical,
			Value: -0.40,
		})
	}

	if !p.Core.Combat.Player().IsWithinArea(p.fieldArea) {
		return
	}
	for _, char := range p.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(PolestarFieldKey+"-cryo", buffDur),
			AffectedStat: attributes.CryoP,
			Amount: func() []float64 {
				if !char.StatusIsActive(PolestarFieldKey) {
					return nil
				}

				buff := stellarConductBuff[newStacks]
				p.mCryo[attributes.CryoP] = buff
				return p.mCryo
			},
		})

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(PolestarFieldKey+"-electro", buffDur),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				if !char.StatusIsActive(PolestarFieldKey) {
					return nil
				}

				buff := stellarConductBuff[newStacks]
				p.mElectro[attributes.ElectroP] = buff
				return p.mElectro
			},
		})

		char.AddStatus(PolestarFieldKey, buffDur, true)
		char.SetTag(PolestarFieldStacksKey, newStacks)
	}
}

func (p *PolestarField) Tick() {
	p.Gadget.Tick()
}
