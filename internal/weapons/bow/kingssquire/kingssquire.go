package kingssquire

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.KingsSquire, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Obtain the Teachings of the Forest effect when unleashing Elemental Skills and Bursts, increasing Elemental Mastery by 60/80/100/120/140 for 12s.
// This effect will be removed when switching characters.
// When the Teachings of the Forest effect ends or is removed, it will deal 100/120/140/160/180% of ATK as DMG to 1 nearby opponent.
// The Teachings of the Forest effect can be triggered once every 20s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const buffKey = "kingssquire"
	const icdKey = "kingssquire-icd"

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 40 + float64(r)*20

	triggerAttack := func() {
		if !char.StatModIsActive(buffKey) {
			return
		}
		char.DeleteStatMod(buffKey)

		// determine attack pos
		player := c.Combat.Player()
		enemy := c.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(player, nil, 15), nil)
		var pos combat.Point
		if enemy == nil {
			pos = player.Pos()
		} else {
			pos = enemy.Pos()
		}

		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "King's Squire Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Mult:       0.8 + float64(r)*0.2,
		}
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(pos, nil, 1.6), 0, 1)
	}

	f := func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 20*60, true)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 12*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		char.QueueCharTask(triggerAttack, 12*60)
		return false
	}

	c.Events.Subscribe(event.OnSkill, f, fmt.Sprintf("kingssquire-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnBurst, f, fmt.Sprintf("kingssquire-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev == char.Index {
			triggerAttack()
		}
		return false
	}, fmt.Sprintf("kingssquire-%v", char.Base.Key.String()))

	return w, nil
}
