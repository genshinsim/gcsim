package archaic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ArchaicPetra, NewSet)
}

type Set struct {
	element attributes.Element
	Index   int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.GeoP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("archaic-2pc", -1),
			AffectedStat: attributes.GeoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)

		core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
			// Character that picks it up must be the petra set holder
			if core.Player.Active() != char.Index {
				return false
			}

			// Check shield
			shd := args[0].(shield.Shield)
			if shd.Type() != shield.Crystallize {
				return false
			}
			s.element = shd.Element()

			// Activate
			// TODO: cd for proc?
			core.Log.NewEvent("archaic petra proc'd", glog.LogArtifactEvent, char.Index).
				Write("ele", s.element)

			m[attributes.PyroP] = 0
			m[attributes.HydroP] = 0
			m[attributes.CryoP] = 0
			m[attributes.ElectroP] = 0
			m[attributes.AnemoP] = 0
			m[attributes.GeoP] = 0
			m[attributes.DendroP] = 0
			m[attributes.EleToDmgP(s.element)] = 0.35 // 35%

			// Apply mod to all characters
			for _, c := range core.Player.Chars() {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("archaic-4pc", 10*60),
					AffectedStat: attributes.NoStat,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}

			return false
		}, fmt.Sprintf("archaic-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
