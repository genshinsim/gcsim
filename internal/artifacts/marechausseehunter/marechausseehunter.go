package marechausseehunter

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.MarechausseeHunter, NewSet)
}

type Set struct {
	stacks int
	core   *core.Core
	char   *character.CharWrapper
	buff   []float64
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func (s *Set) onChangeHP() {
	const buffKey = "mh-4pc"

	if !s.char.StatModIsActive(buffKey) {
		s.stacks = 0
	}
	if s.stacks < 3 {
		s.stacks++
	}

	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(buffKey, 5*60),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			s.buff[attributes.CR] = 0.12 * float64(s.stacks)
			return s.buff, true
		},
	})
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core: c,
		char: char,
	}

	// Normal and Charged Attack DMG +15%
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("mh-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
		})
	}

	// When current HP increases or decreases, CRIT Rate will be increased by 12% for 5s. Max 3 stacks.
	if count < 4 {
		return &s, nil
	}

	s.buff = make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(info.DrainInfo)
		if di.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}

		s.onChangeHP()
		return false
	}, fmt.Sprintf("mh-4pc-drain-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if c.Player.Active() != char.Index {
			return false
		}
		if index != char.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		// do not trigger if at max hp already
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		s.onChangeHP()
		return false
	}, fmt.Sprintf("mh-4pc-heal-%v", char.Base.Key.String()))

	// TODO: OnCharacterHurt?

	return &s, nil
}
