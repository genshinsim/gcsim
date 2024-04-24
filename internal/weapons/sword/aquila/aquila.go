package aquila

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

// ATK is increased by 20%. Triggers on taking DMG: the soul of the Falcon of the
// West awakens, holding the banner of resistance aloft, regenerating HP equal to
// 100% of ATK and dealing 200% of ATK as DMG to surrounding opponents. This
// effect can only occur once every 15s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = .15 + .05*float64(r)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("aquila favonia", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)
	const icdKey = "aquila-icd"

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(*info.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if di.ActorIndex != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 900, true) // 15 sec
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Aquila Favonia",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		snap := char.Snapshot(&ai)
		c.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 6), 1)

		atk := snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK]

		c.Player.Heal(info.HealInfo{
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
