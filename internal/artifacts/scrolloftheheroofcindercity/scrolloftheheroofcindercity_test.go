package scrolloftheheroofcindercity

import (
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func init() {
	testhelper.RegisterTestCharacter()
	testhelper.RegisterTestWeapon()
}

// A Stellar Swirl triggered by an Anemo wearer gives the party the Anemo and Cryo DMG Bonus, like
// a Cryo Swirl (Relic_NyxTime_4 lists SwirlIce and StarSwirlIce together).
func TestStellarSwirlTriggersBuff(t *testing.T) {
	c, _ := core.New(core.Opt{Seed: time.Now().Unix()})
	c.Combat.SetPlayer(avatar.New(c, info.Point{}, 1))
	trg := enemy.New(c, info.EnemyProfile{
		Level:  100,
		Resist: make(map[attributes.Element]float64),
		Pos:    info.Coord{R: 1},
	})
	c.Combat.AddEnemy(trg)

	prof := testhelper.DefaultProfile(testhelper.TestCharKey, testhelper.TestWeaponKey)
	prof.Base.Element = attributes.Anemo
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("error adding char: %v", err)
	}
	c.Player.SetActive(idx)
	char := c.Player.ByIndex(idx)
	if _, err := NewSet(c, char, 4, nil); err != nil {
		t.Fatalf("error adding set: %v", err)
	}
	if err := c.Init(); err != nil {
		t.Fatalf("error initializing core: %v", err)
	}

	c.Events.Emit(event.OnStellarSwirl, trg, &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}})

	for _, ele := range []attributes.Element{attributes.Anemo, attributes.Cryo} {
		if got := char.Stat(attributes.EleToDmgP(ele)); got != 0.12 {
			t.Errorf("%v DMG bonus after Stellar Swirl: got %v, want 0.12", ele, got)
		}
	}
}
