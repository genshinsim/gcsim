package athameartis

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	bladeOfDaylightHoursBuff = "blade-of-the-daylight-hours-buff"
	teamBuffKey              = "athameartis-team-buff"
	buffDur                  = 3 * 60
)

type Weapon struct {
	Index    int
	refine   int
	c        *core.Core
	char     *character.CharWrapper
	teamBuff []float64
	src      int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	r := w.refine
	c := w.c
	char := w.char

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
	if c.Player.GetHexereiCount() >= 2 {
		selfVal[attributes.ATKP] = selfATK * 1.75
	} else {
		selfVal[attributes.ATKP] = selfATK
	}

	selfStatMod := character.StatMod{
		Base:         modifier.NewBase(bladeOfDaylightHoursBuff, buffDur),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return selfVal
		},
	}

	w.teamBuff = make([]float64, attributes.EndStatType)
	if c.Player.GetHexereiCount() >= 2 {
		w.teamBuff[attributes.ATKP] = teamATK * 1.75
	} else {
		w.teamBuff[attributes.ATKP] = teamATK
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

		if !char.StatusIsActive(bladeOfDaylightHoursBuff) {
			// only start a new task if the previous task should've expired
			char.AddStatMod(selfStatMod)
			w.src = c.F
			w.refreshTeamBuff(c.F)()
		} else {
			char.AddStatMod(selfStatMod)
		}
	}, fmt.Sprintf("athame-%v", char.Base.Key.String()))
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		refine: p.Refine,
		c:      c,
		char:   char,
	}

	return w, nil
}

func (w *Weapon) refreshTeamBuff(src int) func() {
	return func() {
		if w.src != src {
			return
		}

		if !w.char.StatusIsActive(bladeOfDaylightHoursBuff) {
			return
		}

		w.char.QueueCharTask(w.refreshTeamBuff(src), 60)

		active := w.c.Player.ActiveChar()
		if active.Index() == w.char.Index() {
			return
		}

		duration := buffDur

		if w.c.Player.GetHexereiCount() >= 2 {
			duration = 2 * 60
		}

		if active.StatModIsActive(teamBuffKey) {
			// if the buff already exists on the active character from a different athame artis
			// the buff is refreshed using the first weapon's refine instead of the new one
			dur := active.StatusDuration(teamBuffKey)
			active.ExtendStatus(teamBuffKey, duration-dur)
		} else {
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(teamBuffKey, duration),
				AffectedStat: attributes.ATKP,
				Amount: func() []float64 {
					return w.teamBuff
				},
			})
		}
	}
}
