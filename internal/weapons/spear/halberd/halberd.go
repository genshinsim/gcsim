package halberd

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
	core.RegisterWeaponFunc(keys.Halberd, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Normal Attacks deal an additional 160/200/240/280/320% DMG.
// Can only occur once every 10s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "halberd-icd"
	dmg := 1.20 + float64(r)*0.40

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		trg := args[0].(combat.Target)
		// don't proc if dmg not from weapon holder
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// don't proc if not Normal Attack
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		// don't proc if on icd
		if char.StatusIsActive(icdKey) {
			return false
		}
		// set icd
		char.AddStatus(icdKey, 600, true) // 10s

		// queue single target proc
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Halberd Proc",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		c.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)

		return false
	}, fmt.Sprintf("halberd-%v", char.Base.Key.String()))

	return w, nil
}
