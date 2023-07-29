package flute

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
)

func init() {
	core.RegisterWeaponFunc(keys.TheFlute, NewWeapon)
}

// Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
// power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
// and a maximum of 1 can be gained every 0.5s.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	icdKey      = "flute-icd"
	durationKey = "flute-stack-duration"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	stacks := 0

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 30, true) //every 0.5s
		if !char.StatusIsActive(durationKey) {
			stacks = 0
		}
		stacks++
		//stacks lasts 30s
		char.AddStatus(durationKey, 1800, true)

		if stacks == 5 {
			//trigger dmg at 5 stacks
			stacks = 0
			char.DeleteStatus(durationKey)

			ai := combat.AttackInfo{
				ActorIndex: char.Index,
				Abil:       "Flute Proc",
				AttackTag:  attacks.AttackTagWeaponSkill,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Physical,
				Durability: 100,
				Mult:       0.75 + 0.25*float64(r),
			}
			trg := args[0].(combat.Target)
			c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 4), 0, 1)

		}
		return false
	}, fmt.Sprintf("flute-%v", char.Base.Key.String()))
	return w, nil
}
