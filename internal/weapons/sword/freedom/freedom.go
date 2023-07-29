package freedom

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FreedomSworn, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// A part of the "Millennial Movement" that wanders amidst the winds.
// Increases DMG by 10%. When the character wielding this weapon triggers
// Elemental Reactions, they gain a Sigil of Rebellion. This effect can be
// triggered once every 0.5s and can be triggered even if said character is
// not on the field. When you possess 2 Sigils of Rebellion, all of them will
// be consumed and all nearby party members will obtain "Millennial Movement:
// Song of Resistance" for 12s. "Millennial Movement: Song of Resistance"
// increases Normal, Charged and Plunging Attack DMG by 16% and increases ATK
// by 20%. Once this effect is triggered, you will not gain Sigils of
// Rebellion for 20s. Of the many effects of the "Millennial Movement," buffs
// of the same type will not stack.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.075 + float64(r)*0.025
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("freedom-dmg", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	uniqueVal := make([]float64, attributes.EndStatType)
	uniqueVal[attributes.DmgP] = .12 + 0.04*float64(r)

	sharedVal := make([]float64, attributes.EndStatType)
	sharedVal[attributes.ATKP] = .15 + float64(r)*0.05

	stacks := 0
	buffDuration := 12 * 60
	const icdKey = "freedom-sworn-sigil-icd"
	icd := int(0.5 * 60)
	const cdKey = "freedom-sworn-cooldown"
	cd := 20 * 60

	stackFunc := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if char.StatusIsActive(cdKey) {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}

		char.AddStatus(icdKey, icd, true)
		stacks++
		c.Log.NewEvent("freedomsworn gained sigil", glog.LogWeaponEvent, char.Index).
			Write("sigil", stacks)
		if stacks == 2 {
			stacks = 0
			char.AddStatus(cdKey, cd, true)
			for _, char := range c.Player.Chars() {
				// Attack buff snapshots so it needs to be in a separate mod
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(common.MillennialKey, buffDuration),
					AffectedStat: attributes.ATKP,
					Amount: func() ([]float64, bool) {
						return sharedVal, true
					},
				})
				char.AddAttackMod(character.AttackMod{
					Base: modifier.NewBaseWithHitlag("freedomsworn", buffDuration),
					Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
						switch atk.Info.AttackTag {
						case attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge:
							return uniqueVal, true
						}
						return nil, false
					},
				})
			}
		}
		return false
	}

	for i := event.ReactionEventStartDelim + 1; i < event.OnShatter; i++ {
		c.Events.Subscribe(i, stackFunc, "freedom-"+char.Base.Key.String())
	}

	return w, nil
}
