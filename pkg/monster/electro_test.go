package monster

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/core"
)

func TestElectroAura(t *testing.T) {

	dmgCount := 0
	shdCount := 0

	sim := dummy.NewSim(func(s *dummy.Sim) {
		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, core.EndStatType)
			c.Stats[core.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		s.OnDamage = func(ds *core.Snapshot) {
			// log.Println(ds)
			dmgCount++
		}

		s.OnShielded = func(shd core.Shield) {
			// log.Println(shd.CurrentHP())
			shdCount++
		}

	})

	target := New(0, sim, logger, 0, core.EnemyProfile{
		Level:  88,
		Resist: defaultResMap(),
	})

	//TEST ATTACH

	target.Attack(&core.Snapshot{
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	fmt.Println("----electro testing----")
	expect("initial durability", 20, target.aura.Durability())
	if target.aura.Durability() != 20 {
		t.Error("intial attach: invalid durability")
		t.FailNow()
	}

	//TEST DECAY
	for i := 0; i < 285; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("decay durability after 4.75 seconds (tolerance 0.01)", 10, target.aura.Durability())
	if !durApproxEqual(10, target.aura.Durability(), 0.01) {
		t.Error("decay test: invalid durability")
		t.FailNow()
	}

	//TEST REFRESH
	target.Attack(&core.Snapshot{
		Durability: 50,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("refresh 50 units on 10 existing (tolerance 0.01)", 60, target.aura.Durability())
	if !durApproxEqual(60, target.aura.Durability(), 0.01) {
		t.Error("refresh test: invalid durability")
		t.FailNow()
	}

}
