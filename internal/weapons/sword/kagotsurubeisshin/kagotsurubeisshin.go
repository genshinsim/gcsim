package kagotsurubeisshin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.KagotsurubeIsshin, NewWeapon)
}

// When a Normal, Charged, or Plunging Attack hits an opponent, it will whip up a
// Hewing Gale, dealing AoE DMG equal to 180% of ATK and increasing ATK by 15% for
// 8s. This effect can be triggered once every 8s.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const icdKey = "kagotsurube-isshin-icd"

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}

	duration := 480
	cd := 480

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra && atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		val := make([]float64, attributes.EndStatType)
		val[attributes.ATKP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("kagotsurube-isshin", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		// add a new action that deals % dmg immediately
		// superconduct attack
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Kagotsurube Isshin Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       1.8,
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 1)

		// trigger cd
		char.AddStatus(icdKey, cd, true)

		return false
	}, fmt.Sprintf("kagotsurube-isshin-%v", char.Base.Key.String()))
	return w, nil
}
