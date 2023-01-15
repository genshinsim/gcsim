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
	"github.com/genshinsim/gcsim/pkg/modifier"
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

	const icdKey = "skyward-atlas-icd"
	icd := 1800 // 30s * 60

	travel, ok := p.Params["travel"]
	if !ok {
		travel = 10
	}

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			return false
		}
		c.Log.NewEvent("skywardatlas proc'd", glog.LogWeaponEvent, char.Index)

		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Skyward Atlas Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       atk,
		}
		snap := char.Snapshot(&ai)

		for i := 1; i <= 6; i++ {
			c.Tasks.Add(func() {
				enemy := c.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 15), nil)
				if enemy == nil {
					return
				}
				c.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(enemy, nil, 1.2), travel)
			}, i*147)
		}
		char.AddStatus(icdKey, icd, true)
		return false
	}, fmt.Sprintf("skyward-atlas-%v", char.Base.Key.String()))

	//permanent stat buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = dmg
	m[attributes.HydroP] = dmg
	m[attributes.CryoP] = dmg
	m[attributes.ElectroP] = dmg
	m[attributes.AnemoP] = dmg
	m[attributes.GeoP] = dmg
	m[attributes.DendroP] = dmg
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("skyward-atlas", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	return w, nil
}
