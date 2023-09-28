package noblesse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.NoblesseOblige, NewSet)
}

type Set struct {
	core              *core.Core
	Index             int
	nob2buff          []float64
	nob4buff          []float64
	charIsSpecialCase bool
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

var specialChars = [keys.EndCharKeys]bool{}

func init() {
	specialChars[keys.AetherAnemo] = true
	specialChars[keys.LumineAnemo] = true
	specialChars[keys.Ningguang] = true
	specialChars[keys.Beidou] = true
	specialChars[keys.Sayu] = true
	specialChars[keys.Aloy] = true
	specialChars[keys.Ganyu] = true
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core: c,
	}

	if count >= 2 {
		s.nob2buff = make([]float64, attributes.EndStatType)
		s.nob2buff[attributes.DmgP] = 0.20
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("nob-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
					return nil, false
				}
				return s.nob2buff, true
			},
		})
	}
	if count >= 4 {
		const buffKey = "nob-4pc"
		buffDuration := 720 // 12s * 60
		s.nob4buff = make([]float64, attributes.EndStatType)
		s.nob4buff[attributes.ATKP] = 0.2
		s.charIsSpecialCase = specialChars[char.Base.Key]

		//TODO: this used to be post. need to check
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if c.Player.Active() != char.Index {
				return false
			}

			for _, x := range s.core.Player.Chars() {
				this := x
				// special case if applying 4 Noblesse to holder to fix this mess:
				// https://library.keqingmains.com/evidence/general-mechanics/bugs#noblesse-oblige-4pc-bonus-not-applying-to-some-bursts
				// https://docs.google.com/spreadsheets/d/1jhIP3C6B16nL1unX9DL_-LhSNaOy_wwhdr29pzikpcg/edit?usp=sharing
				// TODO: Does the char snapshot 4 Noblesse if 4 Noblesse is already up and they're refreshing the duration? (rn they would snapshot it)
				delay := 0
				if this.Base.Key == char.Base.Key && s.charIsSpecialCase {
					delay = 1
				}
				this.QueueCharTask(func() {
					this.AddStatMod(character.StatMod{
						Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
						AffectedStat: attributes.ATKP,
						Amount: func() ([]float64, bool) {
							return s.nob4buff, true
						},
					})
				}, delay)
			}
			c.Log.NewEvent("noblesse 4pc proc", glog.LogArtifactEvent, char.Index).
				Write("expiry (without hitlag)", c.F+buffDuration)
			return false
		}, fmt.Sprintf("nob-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
