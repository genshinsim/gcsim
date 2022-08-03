package noblesse

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
	core.RegisterSetFunc(keys.NoblesseOblige, NewSet)
}

type Set struct {
	core  *core.Core
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.20
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("nob-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagElementalBurst {
					return nil, false
				}
				return m, true
			},
		})
	}
	if count >= 4 {
		const buffKey = "nob-4pc"
		buffDuration := 720 // 12s * 60

		charsToCheck := [6]keys.Char{keys.TravelerAnemo, keys.Ningguang, keys.Beidou, keys.Sayu, keys.Aloy, keys.Ganyu}

		//TODO: this used to be post. need to check
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if c.Player.Active() != char.Index {
				return false
			}

			for _, this := range s.core.Player.Chars() {
				m := make([]float64, attributes.EndStatType)
				m[attributes.ATKP] = 0.2
				smod := character.StatMod{
					Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
					AffectedStat: attributes.ATKP,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				}
				if this.Base.Key != char.Base.Key {
					this.AddStatMod(smod)
				} else {
					// special case if applying 4 Noblesse to holder to fix this mess:
					// https://library.keqingmains.com/evidence/general-mechanics/bugs#noblesse-oblige-4pc-bonus-not-applying-to-some-bursts
					// https://docs.google.com/spreadsheets/d/1jhIP3C6B16nL1unX9DL_-LhSNaOy_wwhdr29pzikpcg/edit?usp=sharing
					// TODO: Does the char snapshot 4 Noblesse if 4 Noblesse is already up and they're refreshing the duration? (rn they would snapshot it)
					for i := range charsToCheck {
						if this.Base.Key == charsToCheck[i] {
							this.QueueCharTask(func() {
								this.AddStatMod(smod)
							}, 1)
						}
					}
				}
			}
			c.Log.NewEvent("noblesse 4pc proc", glog.LogArtifactEvent, char.Index).
				Write("expiry (without hitlag)", c.F+buffDuration)
			return false
		}, fmt.Sprintf("nob-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
