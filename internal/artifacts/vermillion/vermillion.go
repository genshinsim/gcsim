package vermillion

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.VermillionHereafter, NewSet)
}

type Set struct {
	stacks int
	core   *core.Core
	char   *character.CharWrapper
	buff   []float64
	Index  int
}

const verm4pckey = "verm-4pc"

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }
func (s *Set) updateBuff() {
	//8% base + 10% per stack
	s.buff[attributes.ATKP] = 0.08 + float64(s.stacks)*0.1
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
		char: char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("verm-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	//4 Piece: After using an Elemental Burst, this character will gain the
	//Nascent Light effect, increasing their ATK by 8% for 16s. When the
	//character's HP decreases, their ATK will further increase by 10%. This
	//increase can occur this way a maximum of 4 times. This effect can be
	//triggered once every 0.8s. Nascent Light will be dispelled when the
	//character leaves the field. If an Elemental Burst is used again during the
	//duration of Nascent Light, the original Nascent Light will be dispelled.
	if count >= 4 {
		const icdKey = "verm-4pc-icd"
		icd := 48

		s.buff = make([]float64, attributes.EndStatType)

		//TODO: this used to be post. need to check
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}

			s.stacks = 0
			s.updateBuff()

			s.char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(verm4pckey, 16*60),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return s.buff, true
				},
			})

			return false

		}, fmt.Sprintf("verm-4pc-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
			di := args[0].(player.DrainInfo)
			if di.Amount <= 0 {
				return false
			}
			if !char.StatModIsActive(verm4pckey) {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			if s.stacks == 4 {
				return false
			}
			s.stacks++
			char.AddStatus(icdKey, icd, true)
			s.updateBuff()
			c.Log.NewEvent("Vermillion stack gained", glog.LogArtifactEvent, char.Index).Write("stacks", s.stacks)

			return false
		}, "Stack-on-hurt")

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
			char.DeleteStatMod(verm4pckey)
			s.stacks = 0 // resets stacks to 0 when the character swaps
			return false
		}, "char-exit")

	}

	return &s, nil
}
