package swordofnarzissenkreuz

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
)

func init() {
	core.RegisterWeaponFunc(keys.SwordOfNarzissenkreuz, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	icdKey  = "swordofnarzissenkreuz-icd"
	icd     = 12 * 60
	hitmark = 0.1 * 60 // rough estimate
)

// When the equipping character does not have an Arkhe: When Normal Attacks, Charged Attacks, or Plunging Attacks strike,
// a Pneuma or Ousia energy blast will be unleashed, dealing 160/200/240/280/320% of ATK as DMG. This effect can be triggered once every 12s.
// The energy blast type is determined by the current type of the Sword of Narzissenkreuz.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// TODO: arkhe does not influence anything at the moment
	// 0 for Pneuma, 1 for Ousia
	arkhe, ok := p.Params["arkhe"]
	if !ok {
		arkhe = 1
	}
	if arkhe < 0 {
		arkhe = 0
	}
	if arkhe > 1 {
		arkhe = 1
	}
	c.Log.NewEvent("swordofnarzissenkreuz arkhe", glog.LogWeaponEvent, char.Index).
		Write("arkhe", arkhe)

	// no event sub if char has arkhe
	if char.HasArkhe {
		return w, nil
	}

	mult := 1.2 + float64(r)*0.4

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		trg := args[0].(combat.Target)

		// only proc on normal, charge and plunging dmg
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)

		// TODO: check actual attack info and AoE, most weapons have 3m radius so that is what's used here
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Sword of Narzissenkreuz",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       mult,
		}
		// hitmark timing is affected by hitlag
		char.QueueCharTask(func() {
			c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 0)
		}, hitmark)

		return false
	}, fmt.Sprintf("swordofnarzissenkreuz-%v", char.Base.Key.String()))

	return w, nil
}
