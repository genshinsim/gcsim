package messenger

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
	core.RegisterWeaponFunc(keys.Messenger, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Charged Attack hits on weak spots deal an additional 100/125/150/175/200% ATK DMG as CRIT DMG.
// Can only occur once every 10s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 0.75 + float64(r)*0.25
	const icdKey = "messenger-icd"

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		trg := args[0].(combat.Target)
		// don't proc if dmg not from weapon holder
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// don't proc if off-field
		if c.Player.Active() != char.Index {
			return false
		}
		// don't proc if not hitting weakspot
		if !atk.Info.HitWeakPoint {
			return false
		}
		// don't proc if on icd
		if char.StatusIsActive(icdKey) {
			return false
		}
		// set icd
		char.AddStatus(icdKey, 10*60, true) // 10s icd

		// queue single target proc
		ai := combat.AttackInfo{
			ActorIndex:   char.Index,
			Abil:         "Messenger Proc",
			AttackTag:    attacks.AttackTagNone,
			ICDTag:       attacks.ICDTagNone,
			ICDGroup:     attacks.ICDGroupDefault,
			StrikeType:   attacks.StrikeTypePierce,
			Element:      attributes.Physical,
			Durability:   100,
			Mult:         dmg,
			HitWeakPoint: true, // ensure crit by marking it as hitting weakspot
		}
		c.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)

		return false
	}, fmt.Sprintf("messenger-%v", char.Base.Key.String()))

	return w, nil
}
