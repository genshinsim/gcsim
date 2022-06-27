package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.SkywardAtlas, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {

	w := &Weapon{}
	r := p.Refine

	dmg := 0.09 + float64(r)*0.03
	atk := 1.2 + float64(r)*0.4

	icd := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Skyward Atlas Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       atk,
		}
		snap := char.Snapshot(&ai)
		for i := 0; i < 6; i++ {
			c.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), i*150)
		}
		icd = c.F + 1800
		return false
	}, fmt.Sprintf("skyward-atlast-%v", char.Base.Key.String()))

	//permanent stat buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = dmg
	m[attributes.HydroP] = dmg
	m[attributes.CryoP] = dmg
	m[attributes.ElectroP] = dmg
	m[attributes.AnemoP] = dmg
	m[attributes.GeoP] = dmg
	m[attributes.DendroP] = dmg
	char.AddStatMod("skyward-atlast", -1, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})

	return w, nil
}
