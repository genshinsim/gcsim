package oceanhuedclam

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
	core.RegisterSetFunc(keys.OceanHuedClam, NewSet)
}

type Set struct {
	bubbleHealStacks     float64
	bubbleDurationExpiry int
	core                 *core.Core
	Index                int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error {
	// Shows which character currently has an active OHC proc. -1 = Non-active
	s.core.Flags.Custom["OHCActiveChar"] = -1
	return nil
}

// 2-Piece Bonus: Healing Bonus +15%
// 4-Piece Bonus: When the character equipping this artifact set heals a character in the party,
// a Sea-Dyed Foam will appear for 3 seconds, accumulating the amount of HP recovered from healing (including overflow healing).
// At the end of the duration, the Sea-Dyed Foam will explode, dealing DMG to nearby opponents based on 90% of the accumulated
// healing.

// (This DMG is calculated similarly to Reactions such as Electro-Charged, and Superconduct, but it is not affected by
// Elemental Mastery, Character Levels, or Reaction DMG Bonuses).
//
//		Only one Sea-Dyed Foam can be produced every 3.5 seconds. Each Sea-Dyed Foam can accumulate up to 30,000 HP (including
//	 overflow healing). There can be no more than one Sea-Dyed Foam active at any given time.
//		This effect can still be triggered even when the character who is using this artifact set is not on the field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core: c,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.Heal] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("ohc-2pc", -1),
			AffectedStat: attributes.Heal,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		bubbleICDExpiry := 0

		// On Heal subscription to start accumulating the healing
		c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
			src := args[0].(*info.HealInfo)
			healAmt := args[4].(float64)

			if src.Caller != char.Index {
				return false
			}

			// OHC must either be inactive or this equipped character has to have an OHC bubble active
			if c.Flags.Custom["OHCActiveChar"] != -1 && c.Flags.Custom["OHCActiveChar"] != float64(char.Index) {
				return false
			}

			s.bubbleHealStacks += healAmt
			if s.bubbleHealStacks >= 30000 {
				s.bubbleHealStacks = 30000
			}

			// Activate bubble if this character's bubble is off CD, and add the bubble pop task
			if bubbleICDExpiry < c.F {
				s.bubbleDurationExpiry = c.F + 3*60
				bubbleICDExpiry = c.F + 3.5*60

				c.Flags.Custom["OHCActiveChar"] = float64(char.Index)

				// Bubble pop task
				c.Tasks.Add(func() {
					// Bubble is physical damage. This is handled in the reaction damage function, so it is not affected by physical dmg bonus/enemy defense
					// d := c.Snapshot(
					// 	"OHC Damage",
					// 	core.AttackTagNone,
					// 	core.ICDTagNone,
					// 	core.ICDGroupDefault,
					// 	core.StrikeTypeDefault,
					// 	core.Physical,
					// 	0,
					// 	0,
					// )
					// d.Targets = core.TargetAll
					// d.IsOHCDamage = true
					// d.FlatDmg = bubbleHealStacks * .9
					// c.QueueDmg(&d, 0)

					atk := combat.AttackInfo{
						ActorIndex:       char.Index,
						Abil:             "Sea-Dyed Foam",
						AttackTag:        attacks.AttackTagNoneStat,
						ICDTag:           attacks.ICDTagNone,
						ICDGroup:         attacks.ICDGroupDefault,
						StrikeType:       attacks.StrikeTypeDefault,
						Element:          attributes.Physical,
						IgnoreDefPercent: 1,
						FlatDmg:          s.bubbleHealStacks * .9,
					}
					// snapshot -1 since we don't need stats
					c.QueueAttack(atk, combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 6), -1, 1)

					// Reset
					c.Flags.Custom["OHCActiveChar"] = -1
					s.bubbleHealStacks = 0
				}, 3*60)

				c.Log.NewEvent("ohc bubble activated", glog.LogArtifactEvent, char.Index).
					Write("bubble_pops_at", s.bubbleDurationExpiry).
					Write("ohc_icd_expiry", bubbleICDExpiry)
			}

			c.Log.NewEvent("ohc bubble accumulation", glog.LogArtifactEvent, char.Index).
				Write("bubble_pops_at", s.bubbleDurationExpiry).
				Write("bubble_total", s.bubbleHealStacks)

			return false
		}, fmt.Sprintf("ohc-4pc-heal-accumulation-%v", char.Base.Key.String()))
	}

	return &s, nil
}
