package goldenfrostboundoath

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	defBuffKey     = "golden-frostbound-oath-def"
	geoBuffKey     = "golden-frostbound-oath-geo"
	lcrBuffKey     = "golden-frostbound-oath-lcr"
	teamGeoBuffKey = "golden-frostbound-oath-team-geo"
	teamLcrBuffKey = "golden-frostbound-oath-team-lcr"
)

func init() {
	core.RegisterWeaponFunc(keys.GoldenFrostboundOath, NewWeapon)
}

type Weapon struct {
	Index          int
	core           *core.Core
	char           *character.CharWrapper
	teamBuffSource int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Decreases Gliding Stamina consumption by 15%. When using Aimed Shots, the DMG dealt by Charged Attacks
// increases by 6% every 0.5s. This effect can stack up to 6 times and will be removed 10s after leaving
// Aiming Mode.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core: c,
		char: char,
	}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = 0.12 + 0.04*float64(r)

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(defBuffKey, -1),
		AffectedStat: attributes.DEFP,
		Amount: func() []float64 {
			return m
		},
	})

	n := make([]float64, attributes.EndStatType)
	n[attributes.GeoP] = 0.3 + 0.1*float64(r)
	lcrBuff := 0.3 + 0.1*float64(r)

	teamM := make([]float64, attributes.EndStatType)
	teamM[attributes.GeoP] = 0.15 + 0.05*float64(r)
	teamLcrBuff := 0.15 + 0.05*float64(r)

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) { // uses OnHittingOther, current implementation doesn't buff the hit that triggers the buff
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt,
			attacks.AttackTagDirectLunarCrystallize,
			attacks.AttackTagReactionLunarCrystallize:
		default:
			return
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(geoBuffKey, 6*60),
			Amount: func() []float64 {
				return n
			},
		})

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(lcrBuffKey, 6*60),
			Amount: func(atk info.AttackInfo) float64 {
				switch atk.AttackTag {
				case attacks.AttackTagDirectLunarCrystallize,
					attacks.AttackTagReactionLunarCrystallize:
					return lcrBuff
				default:
					return 0
				}
			},
		})

		w.teamBuffSource = c.F
		w.refreshTeamBuff(c.F, teamM, teamLcrBuff)()
	}, fmt.Sprintf("goldenfrostboundoath-on-skill-or-lunar-crystallize-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) refreshTeamBuff(src int, m []float64, lcrBuff float64) func() {
	return func() {
		if w.teamBuffSource != src {
			return
		}

		if !w.char.StatusIsActive(geoBuffKey) {
			return
		}

		moondrifts, _ := w.core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)

		moondriftnearby := false

		playerPos := w.core.Combat.Player().Pos()
		for _, moondrift := range moondrifts {
			if playerPos.Distance(moondrift.Pos()) < 16 {
				moondriftnearby = true
				break
			}
		}

		if !moondriftnearby {
			w.char.QueueCharTask(w.refreshTeamBuff(src, m, lcrBuff), 1*60)
			return
		}

		for _, char := range w.core.Player.Chars() {
			if char.Index() == w.char.Index() {
				continue
			}
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag(teamGeoBuffKey, 2*60),
				Amount: func() []float64 {
					return m
				},
			})
			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(teamLcrBuffKey, 2*60),
				Amount: func(atk info.AttackInfo) float64 {
					switch atk.AttackTag {
					case attacks.AttackTagDirectLunarCrystallize,
						attacks.AttackTagReactionLunarCrystallize:
						return lcrBuff
					default:
						return 0
					}
				},
			})
		}
	}
}
