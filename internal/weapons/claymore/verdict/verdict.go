package verdict

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
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Verdict, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey           = "verdict-skill-dmg"
	buffDuration      = 15 * 60
	dmgWindowKey      = "verdict-dmg-window"
	dmgWindowDuration = 0.2 * 60
)

// Increases ATK by 20/25/30/35/40%. When party members obtain Elemental Shards from Crystallize reactions,
// the equipping character will gain 1 Seal, increasing Elemental Skill DMG by 18/22.5/27/31.5/36%.
// The Seal lasts for 15s, and the equipper may have up to 2 Seals at once.
// All of the equipper's Seals will disappear 0.2s after their Elemental Skill deals DMG.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(r)*0.05
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("verdict-atk", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	// seal gain on crystallize shard pickup
	c.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		// Check shield
		shd := args[0].(shield.Shield)
		if shd.Type() != shield.Crystallize {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 2 {
			w.stacks++
		}
		c.Log.NewEvent("verdict adding stack", glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks)
		char.AddStatus(buffKey, buffDuration, true)
		return false
	}, fmt.Sprintf("verdict-seal-%v", char.Base.Key.String()))

	// skill dmg increase while seals active
	skillDmg := 0.135 + float64(r)*0.045
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		// don't do anything if not in buff
		if !char.StatusIsActive(buffKey) {
			return false
		}
		// otherwise if this is first time proccing
		// - set duration for dmg window
		// - reset buff once window is over
		if !char.StatusIsActive(dmgWindowKey) {
			char.AddStatus(dmgWindowKey, dmgWindowDuration, true)
			char.QueueCharTask(func() {
				char.DeleteStatus(buffKey)
				w.stacks = 0
			}, dmgWindowDuration)
		}
		skillDmgAdd := skillDmg * float64(w.stacks)
		atk.Snapshot.Stats[attributes.DmgP] += skillDmgAdd

		c.Log.NewEvent("verdict adding skill dmg", glog.LogPreDamageMod, char.Index).
			Write("skill_dmg_added", skillDmgAdd)
		return false
	}, fmt.Sprintf("verdict-onhit-%v", char.Base.Key.String()))

	return w, nil
}
