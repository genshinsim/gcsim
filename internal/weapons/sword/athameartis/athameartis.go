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
	bladeOfDaylightHoursStatus = "blade-of-the-daylight-hours-status"
	bladeOfDaylightHoursBuff   = "blade-of-the-daylight-hours-buff"
	teamBuffKey                = "athameartis-team-buff"
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

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(bladeOfDaylightHoursBuff, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			if !char.StatusIsActive(bladeOfDaylightHoursStatus) {
				return nil
			}

			if c.Player.GetHexereiCount() >= 2 {
				selfVal[attributes.ATKP] = selfATK * 1.75
			} else {
				selfVal[attributes.ATKP] = selfATK
			}

			return selfVal
		},
	})

	teamVal := make([]float64, attributes.EndStatType)

	// TODO: When there are Athame Artis of different refines on the team,
	// the team ATK% buff should use the higher refine value.

	applyTeamBuff := func() {
		active := c.Player.ActiveChar()

		duration := 3 * 60

		if c.Player.GetHexereiCount() >= 2 {
			duration = 2 * 60
			teamVal[attributes.ATKP] = teamATK * 1.75
		} else {
			teamVal[attributes.ATKP] = teamATK
		}

		// TODO: When there are Athame Artis of different refines on the team,
		// the team atk% buff should prioritize higher buff value

		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(teamBuffKey, duration),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return teamVal
			},
		})
	}
	var refreshTeamBuff func()

	refreshTeamBuff = func() {
		if !char.StatusIsActive(bladeOfDaylightHoursStatus) {
			return
		}

		applyTeamBuff()

		c.Tasks.Add(refreshTeamBuff, 60)
	}

	// TODO: ATK% buff should affect hit that triggered it.
	// Currently gcsim snapshots before OnEnemyHit event so it isn't buffing the triggering hit.

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
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

		char.AddStatus(bladeOfDaylightHoursStatus, 3*60, true)

		applyTeamBuff()
		c.Tasks.Add(refreshTeamBuff, 60)
	}, fmt.Sprintf("athame-%v", char.Base.Key.String()))

	return w, nil
}
