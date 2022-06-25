package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.SkywardPride, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases all DMG by 8%. After using an Elemental Burst, Normal or Charged
	//Attack, on hit, creates a vacuum blade that does 80% of ATK as DMG to
	//opponents along its path. Lasts for 20s or 8 vacuum blades.
	w := &Weapon{}
	r := p.Refine

	//perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.06 + float64(r)*0.02
	char.AddStatMod("skyward pride", -1, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})

	counter := 0
	dur := 0
	dmg := 0.6 + float64(r)*0.2

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		dur = c.F + 1200
		counter = 0
		c.Log.NewEvent("Skyward Pride activated", glog.LogWeaponEvent, char.Index, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Base.Name))

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.F > dur {
			return false
		}
		if counter >= 8 {
			return false
		}

		counter++
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Skyward Pride Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		c.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 0, 1)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Base.Name))
	return w, nil
}
