package monster

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/def"
)

func TestFrozenDuration(t *testing.T) {
	var target *Target

	dmgCount := 0

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

	})

	fmt.Println("----testing applying 25 cryo on 50 hydro (no delay)----")

	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Hydro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	dur := 0

	for i := 0; i < 220; i++ {
		sim.F++

		if target.aura != nil {
			if target.aura.Type() == def.Frozen {
				dur++
			}
			target.AuraTick()
			target.Tick()

		}

	}
	expect("checking freeze duration in frames (25 cryo on 50 hydro no delay)", 209, dur)
	if dur != 209 {
		t.Errorf("freeze test: expecting 209 frames in duration, got %v", dur)
		t.FailNow()
	}

	//extending should add to existing durability? capped at?

}
