package mailedflower

import (
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
	core.RegisterWeaponFunc(keys.MailedFlower, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Within 8s after the character's Elemental Skill hits an opponent or the character triggers an Elemental Reaction,
// their ATK and Elemental Mastery will be increased by 12%/15%/18%/21%/24% and 48/60/72/84/96 respectively.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := 0.09 + float64(r)*0.03
	em := 36 + float64(r)*12
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = atk
	m[attributes.EM] = em

	f := func() {
		if c.Player.Active() != char.Index {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("mailedflower", 8*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	fDamage := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		f()
		return false
	}
	fReact := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		f()
		return false
	}

	c.Events.Subscribe(event.OnEnemyDamage, fDamage, "mailedflower-skill-"+char.Base.Key.String())
	// considers shatter as an elemental reaction
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, fReact, "mailedflower-"+char.Base.Key.String())
	}

	return w, nil
}
