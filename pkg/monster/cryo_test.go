package monster

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/def"
)

func TestCryoAura(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	sim := dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		target = New(0, s, logger, 0, def.EnemyProfile{
			Level:  88,
			Resist: defaultResMap(),
		})

		s.OnDamage = func(ds *def.Snapshot) {
			// log.Println(ds)
			dmgCount++
			target.Attack(ds)
		}

		s.OnShielded = func(shd def.Shield) {
			// log.Println(shd.CurrentHP())
			shdCount++
		}

	})

	//TEST ATTACH

	fmt.Println("----cryo testing----")

	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	expect("initial durability", 20, target.aura.Durability())
	if target.aura.Durability() != 20 {
		t.Error("intial attach: invalid durability")
		t.FailNow()
	}

	tickCount := 0
	//TEST DECAY
	for i := 0; i < 285; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
		tickCount++
	}

	expect("decay durability after 4.75 seconds (tolerance 0.01)", 10, target.aura.Durability())
	if !durApproxEqual(10, target.aura.Durability(), 0.01) {
		t.Error("decay test: invalid durability")
		t.FailNow()
	}

	//TEST REFRESH
	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	expect("refresh 50 units on 10 existing (tolerance 0.01)", 60, target.aura.Durability())
	if !durApproxEqual(60, target.aura.Durability(), 0.01) {
		t.Error("refresh test: invalid durability")
		t.FailNow()
	}

}
