package athameartis

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AthameArtis, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(
	c *core.Core,
	char *character.CharWrapper,
	p info.WeaponProfile,
) (info.Weapon, error) {
	w := &Weapon{}

	r := p.Refine

	burstCD := 0.12 + 0.04*float64(r)
	selfATK := 0.15 + 0.05*float64(r)
	teamATK := 0.12 + 0.04*float64(r)

	burstCDMod := make([]float64, attributes.EndStatType)
	burstCDMod[attributes.CD] = burstCD

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("athame-burst-cd", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil
			}
			return burstCDMod
		},
	})

	// TODO:
	// Passive specifies "when an Elemental Burst hits an opponent".
	// Currently implemented using OnEnemyDamage as a conservative
	// approximation until hit-vs-damage behavior is verified.

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return
		}

		// Can trigger while off-field.

		procSelfATK := selfATK
		procTeamATK := teamATK

		// TODO verify exact Hexerei condition
		if c.Player.GetHexereiCount() >= 2 {
			procSelfATK *= 1.75
			procTeamATK *= 1.75
		}

		selfVal := make([]float64, attributes.EndStatType)
		selfVal[attributes.ATKP] = procSelfATK

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("athame-self-atk", 3*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return selfVal
			},
		})

		teamVal := make([]float64, attributes.EndStatType)
		teamVal[attributes.ATKP] = procTeamATK

		for _, p := range c.Player.Chars() {
			if p.Index() == char.Index() {
				continue
			}

			p.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("athame-team-atk", 3*60),
				AffectedStat: attributes.ATKP,
				Amount: func() []float64 {
					return teamVal
				},
			})
		}
	}, fmt.Sprintf("athame-%v", char.Base.Key.String()))

	return w, nil
}
