package debateclub

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
	core.RegisterWeaponFunc(keys.DebateClub, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After using an Elemental Skill, Normal or Charged Attacks, on hit, deal an additional 60/75/90/105/120% ATK DMG in a small area.
// Effect lasts 15s. DMG can only occur once every 3s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const effectKey = "debate-club-effect"
	const icdKey = "debate-club-icd"
	dmg := 0.45 + float64(r)*0.15

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		// don't activate if skill not from weapon holder
		if c.Player.Active() != char.Index {
			return false
		}
		// reset icd
		char.DeleteStatus(icdKey)
		// add debate club effect
		char.AddStatus(effectKey, 900, true) // 15s
		return false
	}, fmt.Sprintf("debate-club-activation-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		trg := args[0].(combat.Target)
		// don't proc if dmg not from weapon holder
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// don't proc if off-field
		if c.Player.Active() != char.Index {
			return false
		}
		// don't proc if debate club effect isn't active
		if !char.StatusIsActive(effectKey) {
			return false
		}
		// don't proc if not Normal/Charged Attack
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		// don't proc if on icd
		if char.StatusIsActive(icdKey) {
			return false
		}
		// set icd
		char.AddStatus(icdKey, 180, true) // 3s

		// queue proc
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Debate Club Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 1)

		return false
	}, fmt.Sprintf("debate-club-proc-%v", char.Base.Key.String()))

	return w, nil
}
