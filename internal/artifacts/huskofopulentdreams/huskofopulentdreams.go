package huskofopulentdreams

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
	core.RegisterSetFunc(keys.HuskOfOpulentDreams, NewSet)
}

type Set struct {
	stacks int
	// Required to check for stack gain icd
	stackGainICDKey string
	stackGainICD    int
	// Required to check for stack loss
	stackLossTimer int
	lastStackGain  int
	// Source initializes at -1
	lastSwap int
	core     *core.Core
	char     *character.CharWrapper
	Index    int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

// Initiate off-field stacking if off-field at start of the sim
func (s *Set) Init() error {
	if s.core.Player.Active() != s.char.Index {
		s.core.Tasks.Add(s.gainStackOfffield(s.lastSwap), 180)
	}
	return nil
}

/*
*
A character equipped with this Artifact set will obtain the Curiosity effect in the following conditions:
When on the field, the character gains 1 stack after hitting an opponent with a Geo attack,
triggering a maximum of once every 0.3s. When off the field, the character gains 1 stack every 3s.

Curiosity can stack up to 4 times, each providing 6% DEF and a 6% Geo DMG Bonus. When 6 seconds pass
without gaining a Curiosity stack, 1 stack is lost.
*
*/
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
		char: char,
	}
	s.lastSwap = -1
	s.stacks = param["stacks"]
	if s.stacks > 4 {
		s.stacks = 4
	}

	s.stackGainICDKey = "husk-4pc-stack-gain-icd"
	s.stackGainICD = 18 // 0.3s * 60
	s.lastStackGain = -1
	s.stackLossTimer = 360 // 6s * 60

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DEFP] = 0.30
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("husk-2pc", -1),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
			prev := args[0].(int)
			if prev != char.Index {
				return false
			}
			s.lastSwap = c.F
			c.Tasks.Add(s.gainStackOfffield(c.F), 3*60)
			return false
		}, fmt.Sprintf("husk-4pc-off-field-gain-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			// Only triggers when onfield
			if c.Player.Active() != char.Index {
				return false
			}
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			if char.StatusIsActive(s.stackGainICDKey) {
				return false
			}
			if atk.Info.Element != attributes.Geo {
				return false
			}

			if s.stacks < 4 {
				s.stacks++
			}

			c.Log.NewEvent("Husk gained on-field stack", glog.LogArtifactEvent, char.Index).
				Write("stacks", s.stacks).
				Write("last_swap", s.lastSwap).
				Write("last_stack_change", s.lastStackGain)

			char.AddStatus(s.stackGainICDKey, s.stackGainICD, true)

			s.lastStackGain = c.F
			char.QueueCharTask(s.checkStackLoss(c.F), s.stackLossTimer)

			return false
		}, fmt.Sprintf("husk-4pc-%v", char.Base.Key.String()))

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("husk-4pc", -1),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				m[attributes.DEFP] = 0.06 * float64(s.stacks)
				m[attributes.GeoP] = 0.06 * float64(s.stacks)
				return m, true
			},
		})
	}

	return &s, nil
}

// Helper function to check for stack loss
// called after every stack gain
func (s *Set) checkStackLoss(src int) func() {
	return func() {
		if s.lastStackGain != src {
			s.core.Log.NewEvent("husk stack loss check ignored, src diff", glog.LogCharacterEvent, s.char.Index).
				Write("src", src).
				Write("new src", s.lastStackGain)
			return
		}
		s.stacks--
		s.core.Log.NewEvent("Husk lost stack", glog.LogArtifactEvent, s.char.Index).
			Write("stacks", s.stacks).
			Write("last_swap", s.lastSwap).
			Write("last_stack_change", s.lastStackGain)

		// queue up again if we still have stacks
		if s.stacks > 0 {
			s.char.QueueCharTask(s.checkStackLoss(src), s.stackLossTimer)
		}
	}
}

func (s *Set) gainStackOfffield(src int) func() {
	return func() {
		s.core.Log.NewEvent("Husk check for off-field stack", glog.LogArtifactEvent, s.char.Index).
			Write("stacks", s.stacks).
			Write("last_swap", s.lastSwap).
			Write("last_stack_change", s.lastStackGain).
			Write("source", src)

		if s.core.Player.Active() == s.char.Index {
			return
		}
		// Ignore if the last source was not not the most recent swap
		if s.lastSwap != src {
			return
		}

		if s.stacks < 4 {
			s.stacks++
		}

		s.core.Log.NewEvent("Husk gained off-field stack", glog.LogArtifactEvent, s.char.Index).
			Write("stacks", s.stacks).
			Write("last_swap", s.lastSwap).
			Write("last_stack_change", s.lastStackGain)

		s.core.Tasks.Add(s.gainStackOfffield(src), 180)

		s.lastStackGain = s.core.F
		s.char.QueueCharTask(s.checkStackLoss(s.core.F), s.stackLossTimer)
	}
}
