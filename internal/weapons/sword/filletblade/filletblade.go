package filletblade

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
)

func init() {
	core.RegisterWeaponFunc(keys.FilletBlade, NewWeapon)
}

// On hit, has 50% chance to deal 240/280/320/360/400% ATK DMG to a single enemy.
// Can only occur once every 15/14/13/12/11s.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const icdKey = "fillet-blade-icd"

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	cd := 960 - 60*r

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if c.Rand.Float64() > 0.5 {
			return false
		}
		if dmg == 0 {
			return false
		}
		// add a new action that deals % dmg immediately
		// superconduct attack
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Fillet Blade Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       2.0 + 0.4*float64(r),
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)

		// trigger cd
		char.AddStatus(icdKey, cd, true)

		return false
	}, fmt.Sprintf("fillet-blade-%v", char.Base.Key.String()))
	return w, nil
}
