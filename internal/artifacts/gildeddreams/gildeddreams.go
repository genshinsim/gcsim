package gildeddreams

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.GildedDreams, NewSet)
}

type Set struct {
	buff  []float64
	c     *core.Core
	char  *character.CharWrapper
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	emCount := 0
	atkCount := 0

	for _, this := range s.c.Player.Chars() {
		if s.char.Index == this.Index {
			continue
		}
		if this.Base.Element != s.char.Base.Element {
			emCount++
		} else {
			atkCount++
		}
	}

	if emCount > 3 {
		emCount = 3
	}
	if atkCount > 3 {
		atkCount = 3
	}

	s.buff = make([]float64, attributes.EndStatType)
	s.buff[attributes.EM] = 50 * float64(emCount)
	s.buff[attributes.ATKP] = 0.14 * float64(atkCount)

	return nil
}

// 2-Piece Bonus: Elemental Mastery +80.
// 4-Piece Bonus: Within 8s of triggering an Elemental Reaction, the character equipping this will obtain buffs based on the Elemental
// Type of the other party members. ATK is increased by 14% for each party member whose Elemental Type is the same as the equipping
// character, and Elemental Mastery is increased by 50 for every party member with a different Elemental Type. Each of the aforementioned
// buffs will count up to 3 characters. This effect can be triggered once every 8s. The character who equips this can still trigger its
// effects when not on the field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		c:    c,
		char: char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("gd-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		const icdKey = "gd-4pc-icd"
		add := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, 8*60, true)

			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("gd-4pc", 8*60),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return s.buff, true
				},
			})
			c.Log.NewEvent("gilded dreams proc'd", glog.LogArtifactEvent, char.Index).
				Write("em", s.buff[attributes.EM]).
				Write("atk", s.buff[attributes.ATKP])
			return false
		}

		for i := event.ReactionEventStartDelim + 1; i < event.OnShatter; i++ {
			c.Events.Subscribe(i, add, fmt.Sprintf("gd-4pc-%v", char.Base.Key.String()))
		}
	}

	return &s, nil
}
