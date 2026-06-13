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
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	bladeOfDaylightHours = "blade-of-the-daylight-hours"
	teamBuffKey          = "athameartis-team-buff"
)

func init() {
	core.RegisterWeaponFunc(keys.AthameArtis, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
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

	selfVal := make([]float64, attributes.EndStatType)
	selfVal[attributes.ATKP] = selfATK

	selfValHex := make([]float64, attributes.EndStatType)
	selfValHex[attributes.ATKP] = selfATK * 1.75

	teamVal := make([]float64, attributes.EndStatType)
	teamVal[attributes.ATKP] = teamATK

	teamValHex := make([]float64, attributes.EndStatType)
	teamValHex[attributes.ATKP] = teamATK * 1.75

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(bladeOfDaylightHours, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			if !char.StatusIsActive(bladeOfDaylightHours) {
				return nil
			}

			if c.Player.GetHexereiCount() >= 2 {
				return selfValHex
			}

			return selfVal
		},
	})

	applyTeamBuff := func() {
		active := c.Player.ActiveChar()

		duration := 3 * 60
		buff := teamVal

		if c.Player.GetHexereiCount() >= 2 {
			duration = 2 * 60
			buff = teamValHex
		}

		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(teamBuffKey, duration),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return buff
			},
		})
	}
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return
		}

		char.AddStatus(bladeOfDaylightHours, 3*60, true)
		applyTeamBuff()
	}, fmt.Sprintf("athame-%v", char.Base.Key.String()))

	return w, nil
}
