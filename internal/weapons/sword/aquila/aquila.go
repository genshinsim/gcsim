package aquila

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AquilaFavonia, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

//ATK is increased by 20%. Triggers on taking DMG: the soul of the Falcon of the
//West awakens, holding the banner of resistance aloft, regenerating HP equal to
//100% of ATK and dealing 200% of ATK as DMG to surrounding opponents. This
//effect can only occur once every 15s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	//perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = .15 + .05*float64(r)
	char.AddStatMod(character.StatMod{Base: modifier.NewBase("acquila favonia", -1), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
		return m, true
	}})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)
	last := -1

	c.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if c.F-last < 900 && last != -1 {
			return false
		}
		last = c.F
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Aquila Favonia",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		snap := char.Snapshot(&ai)
		c.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 1)

		atk := snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK]

		c.Player.Heal(player.HealInfo{
			Caller:  char.Index,
			Target:  c.Player.Active(),
			Message: "Aquila Favonia",
			Src:     atk * heal,
			Bonus:   char.Stat(attributes.Heal),
		})
		return false
	}, fmt.Sprintf("aquila-%v", char.Base.Key.String()))
	return w, nil
}
