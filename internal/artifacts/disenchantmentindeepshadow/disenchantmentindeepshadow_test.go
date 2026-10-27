package disenchantmentindeepshadow

import (
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func init() {
	testhelper.RegisterTestCharacter()
	testhelper.RegisterTestWeapon()
}

// The 4pc CRIT Rate on special reaction contributions against a Superconducted opponent is the
// wearer's own (Relic_6.6_Superconductor is a modifier on the wearer), not the party's.
func TestSpecialReactionCritOnlyForWearer(t *testing.T) {
	c, _ := core.New(core.Opt{Seed: time.Now().Unix()})
	c.Combat.SetPlayer(avatar.New(c, info.Point{}, 1))
	trg := enemy.New(c, info.EnemyProfile{
		Level:  100,
		Resist: make(map[attributes.Element]float64),
		Pos:    info.Coord{R: 1},
	})
	c.Combat.AddEnemy(trg)

	var idx [2]int
	for i := range idx {
		var err error
		idx[i], err = c.AddChar(testhelper.DefaultProfile(testhelper.TestCharKey, testhelper.TestWeaponKey))
		if err != nil {
			t.Fatalf("error adding char: %v", err)
		}
	}
	c.Player.SetActive(idx[0])
	if _, err := NewSet(c, c.Player.ByIndex(idx[0]), 4, nil); err != nil {
		t.Fatalf("error adding set: %v", err)
	}
	if err := c.Init(); err != nil {
		t.Fatalf("error initializing core: %v", err)
	}
	trg.AddStatus(reactable.SuperConductShredKey, 600, false)

	for i, want := range map[int]float64{idx[0]: 0.16, idx[1]: 0} {
		ae := &info.AttackEvent{Info: info.AttackInfo{ActorIndex: i}}
		c.Events.Emit(event.OnSpecialReactionAttack, trg, ae)
		if got := ae.Snapshot.Stats[attributes.CR]; got != want {
			t.Errorf("char %v CRIT Rate on special reaction: got %v, want %v", i, got, want)
		}
	}
}
