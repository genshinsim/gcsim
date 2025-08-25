package sequenceofsolitude

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const icdKey = "sequence-of-solitude-icd"

func init() {
	core.RegisterWeaponFunc(keys.SequenceOfSolitude, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	hp := 0.3 + float64(r)*0.1

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 15*60, true)
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Sequence of Solitude Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			FlatDmg:    char.MaxHP() * hp,
		}
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3), 0, 1)
		return false
	}, fmt.Sprintf("solitude-%v", char.Base.Key.String()))
	return w, nil
}
