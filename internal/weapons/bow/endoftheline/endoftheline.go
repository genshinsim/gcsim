package endoftheline

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
)

func init() {
	core.RegisterWeaponFunc(keys.EndOfTheLine, NewWeapon)
}

type Weapon struct {
	Index     int
	procCount int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Triggers the Flowrider effect after using an Elemental Skill, dealing 80% ATK as AoE DMG upon hitting an opponent with an attack.
// Flowrider will be removed after 15s or after causing 3 instances of AoE DMG.
// Only 1 instance of AoE DMG can be caused every 2s in this way.
// Flowrider can be triggered once every 12s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	flowriderDmg := 0.60 + float64(r)*0.20
	const effectKey = "endoftheline-effect"
	const effectIcdKey = "endoftheline-effect-icd"
	const dmgIcdKey = "endoftheline-dmg-icd"

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		// do nothing if holder is off-field
		if c.Player.Active() != char.Index {
			return false
		}
		// do nothing if flowrider is on icd
		if char.StatusIsActive(effectIcdKey) {
			return false
		}
		// add flowrider status and proc flowrider icd
		char.AddStatus(effectKey, 15*60, true)
		char.AddStatus(effectIcdKey, 12*60, true)
		// reset icd
		if char.StatusIsActive(dmgIcdKey) {
			char.DeleteStatus(dmgIcdKey)
		}
		// reset proc count
		w.procCount = 0
		return false
	}, fmt.Sprintf("endoftheline-effect-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		// do nothing if holder is off-field
		if c.Player.Active() != char.Index {
			return false
		}
		// do nothing if attack not from holder
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// do nothing if flowrider is not active
		if !char.StatusIsActive(effectKey) {
			return false
		}
		// do nothing if flowrider dmg is on icd
		if char.StatusIsActive(dmgIcdKey) {
			return false
		}
		// queue up flowrider proc
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "End of the Line Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       flowriderDmg,
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 2.5), 0, 1)

		w.procCount++
		c.Log.NewEvent("endoftheline proc", glog.LogWeaponEvent, char.Index).
			Write("procCount", w.procCount)
		if w.procCount == 3 {
			char.DeleteStatus(effectKey)
		} else {
			char.AddStatus(dmgIcdKey, 2*60, true)
		}

		return false
	}, fmt.Sprintf("endoftheline-dmg-%v", char.Base.Key.String()))

	return w, nil
}
