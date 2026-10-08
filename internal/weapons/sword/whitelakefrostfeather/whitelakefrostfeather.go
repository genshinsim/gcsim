package whitelakefrostfeather

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	stackGainICDKey = "whitelake-frostfeather-stack-icd"
	stellarBuffKey  = "whitelake-frostfeather-stellar"
	energyICDKey    = "whitelake-frostfeather-energy-icd"
)

type Weapon struct {
	Index        int
	stackCounter int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.06 + 0.02*float64(r)

	critBuff := make([]float64, attributes.EndStatType)
	critBuff[attributes.CD] = 0.35 + 0.15*float64(r)

	energy := 3.5 + 0.5*float64(r)

	onSkillHit := func(args ...any) {
		if char.StatusIsActive(stackGainICDKey) {
			return
		}

		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		default:
			return
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(fmt.Sprintf("whitelake-frostfeather-atk-%v", w.stackCounter+1), 8*60),
			Amount: func() []float64 {
				return atkBuff
			},
		})

		char.AddStatus(stackGainICDKey, 0.1*60, true)

		prevStatModKey := fmt.Sprintf("whitelake-frostfeather-atk-%v", (w.stackCounter-1)%3+1)

		w.stackCounter++
		w.stackCounter %= 3

		if !char.StatModIsActive(prevStatModKey) {
			return
		}

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(stellarBuffKey, char.StatusDuration(prevStatModKey)),
			Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
				if atk.Info.AttackTag.IsStellarDirect() {
					return critBuff
				}

				return nil
			},
		})
	}

	refundEnergy := func() {
		if char.StatusIsActive(energyICDKey) {
			return
		}
		char.AddEnergy("whitelake-frostfeather-energy", energy)
		char.AddStatus(energyICDKey, 3.5*60, true)
	}

	onStellarReaction := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		refundEnergy()
	}

	onStellarDamage := func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if !atk.Info.AttackTag.IsStellar() {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		refundEnergy()
	}

	onStellarReactionDmg := func(args ...any) {
		if !char.StatModIsActive(stellarBuffKey) {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if !atk.Info.AttackTag.IsStellarReact() {
			return
		}

		atk.Snapshot.Stats[attributes.CD] += critBuff[attributes.CD]
	}

	c.Events.Subscribe(event.OnEnemyHit, onSkillHit, "whitelake-frostfeather-on-skill-hit-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnEnemyDamage, onStellarDamage, "whitelake-frostfeather-on-stellar-damage-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnStellarConduct, onStellarReaction, "whitelake-frostfeather-on-stellar-reaction-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnStellarSwirl, onStellarReaction, "whitelake-frostfeather-on-stellar-reaction-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnSpecialReactionAttack, onStellarReactionDmg, "whitelake-frostfeather-on-stellar-reaction-dmg-"+char.Base.Key.String())

	return w, nil
}
