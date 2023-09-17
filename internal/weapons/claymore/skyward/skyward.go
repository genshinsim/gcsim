package skyward

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
	core.RegisterWeaponFunc(keys.SkywardPride, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Increases all DMG by 8%. After using an Elemental Burst, Normal or Charged
	// Attack, on hit, creates a vacuum blade that does 80% of ATK as DMG to
	// opponents along its path. Lasts for 20s or 8 vacuum blades.
	w := &Weapon{}
	r := p.Refine

	// perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.06 + float64(r)*0.02
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("skyward pride", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	const durKey = "skyward-pride-active"
	counter := 0
	dmg := 0.6 + float64(r)*0.2

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatus(durKey, 1200, true)
		counter = 0
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if !char.StatusIsActive(durKey) {
			return false
		}
		if counter >= 8 {
			return false
		}

		counter++
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Skyward Pride Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewBoxHitOnTarget(trg, nil, 0.1, 0.1), 0, 1)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Base.Key.String()))
	return w, nil
}
