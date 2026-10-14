package breezebornerefrain

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	erBufffKey      = "breezeborne-refrain-er"
	viperKey        = "thus-lied-the-viper"
	hymnStackICDKey = "hymn-of-the-pure-stack-icd"
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

	erBuff := make([]float64, attributes.EndStatType)
	erBuff[attributes.ER] = 0.15 + float64(r)*0.05

	stellarBuff := 0.18 + 0.06*float64(r)

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase(erBufffKey, -1),
		Amount: func() []float64 {
			return erBuff
		},
	})

	// TODO: how does this work with multiple holders that are different refine levels?
	for _, ch := range c.Player.Chars() {
		ch.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(viperKey, -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !char.StatusIsActive(viperKey) {
					return 0
				}
				switch ai.AttackTag {
				case
					attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
					return stellarBuff
				default:
					return 0
				}
			},
		})
	}

	onSkillOrBurstHit := func(args ...any) {
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
		case attacks.AttackTagElementalBurst:
		default:
			return
		}

		if char.StatusIsActive(viperKey) {
			return
		}

		if char.StatusIsActive(hymnStackICDKey) {
			return
		}

		w.stackCounter++
		char.AddStatus(hymnStackICDKey, 2, true)

		if w.stackCounter < 3 {
			return
		}

		w.stackCounter = 0
		char.AddStatus(viperKey, 12*60, true)
	}

	c.Events.Subscribe(event.OnEnemyHit, onSkillOrBurstHit, "breezeborne-refrain-on-hit-"+char.Base.Key.String())

	return w, nil
}
