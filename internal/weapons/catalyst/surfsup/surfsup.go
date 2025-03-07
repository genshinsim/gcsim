package surfsup

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	icdKey       = "surfs-up-icd"
	buffKey      = "surfs-up-buff"
	loseStackIcd = "surfs-up-stack-loss-icd"
	gainStackIcd = "surfs-up-stack-gain-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.SurfsUp, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Max HP increased by 40%.
// Once every 15s, for the 14s after using an Elemental Skill:
// Gain 4 Scorching Summer stacks.
// Each stack increases Normal Attack DMG by 24%.
// For the duration of the effect,
// once every 1.5s, lose 1 stack after a Normal Attack hits an opponent;
// once every 1.5s, gain 1 stack after triggering a Vaporize reaction on an opponent.
// Max 4 Scorching Summer stacks.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmgPerStack := 0.09 + float64(r)*0.03

	mHP := make([]float64, attributes.EndStatType)
	mHP[attributes.HPP] = 0.15 + float64(r)*0.05
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("surfs-up-hp%", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return mHP, true
		},
	})

	mNA := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}

		char.AddStatus(icdKey, 15*60, true)
		w.stacks = 4

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 14*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag == attacks.AttackTagNormal {
					mNA[attributes.DmgP] = dmgPerStack * float64(min(w.stacks, 4))
					return mNA, true
				}
				return nil, false
			},
		})

		return false
	}, fmt.Sprintf("surfs-up-skill-%v", char.Base.Key.String()))

	// Gain stack on vape
	c.Events.Subscribe(event.OnVaporize, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			return false
		}
		if char.StatusIsActive(gainStackIcd) {
			return false
		}

		if w.stacks < 5 { // limit to 5 so it can carry over when NA hits
			w.stacks++
		}
		if w.stacks == 5 {
			char.QueueCharTask(func() {
				w.stacks = 4
			}, .5*60)
		}

		c.Log.NewEvent("Surf's Up gained stack", glog.LogWeaponEvent, char.Index)
		char.AddStatus(gainStackIcd, 1.5*60, true)

		return false
	}, fmt.Sprintf("surfs-up-vape-%v", char.Base.Key.String()))

	// Lose stack on NA
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		if c.Player.Active() != char.Index {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			return false
		}
		if char.StatusIsActive(loseStackIcd) {
			return false
		}

		if w.stacks > 0 {
			w.stacks--
		}

		c.Log.NewEvent("Surf's Up lost stack", glog.LogWeaponEvent, char.Index)
		char.AddStatus(loseStackIcd, 1.5*60, true)

		return false
	}, fmt.Sprintf("surfs-up-dmg-%v", char.Base.Key.String()))

	return w, nil
}
