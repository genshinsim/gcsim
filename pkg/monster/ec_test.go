package monster

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/def"
)

func TestElectroOnHydro(t *testing.T) {

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

	//TEST SWIRL
	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Hydro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	// expect("apply 25 hydro to 25 electro (tolerance 0.01)", 20, target.aura.Durability())
	// if !durApproxEqual(20, target.aura.Durability(), 0.01) {
	// 	t.Error("ec test: invalid durability")
	// 	t.FailNow()
	// }
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()

	fmt.Println("----testing applying 25 hydro on 25 electro (no delay)----")

	expect("checking ec dealt 1 dmg tick after initial app", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//check durability
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}

	expect("check electro durability after tick, t=0s (frame 1)", 20.0*(1-1.0/570.0), ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-1.0/570.0), ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//tick for 6 more frames, wane should have happened now
	for i := 0; i < 6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}

	expect("check electro durability after 0.1s, t=0s (frame 1)", (20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//expecting another tick at frame 61
	for i := 0; i < 60-6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("check electro durability at at t=1s (frame 61)", 20.0*(1-61.0/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-61.0/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	expect("checking ec dealt 1 dmg tick at t=1s (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//tick for 6 more frames, ec should be gone now
	for i := 0; i < 6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("expecting ec to be gone now at t=1.1s (frame 67), nothing left", nil, target.aura)
	if target.aura != nil {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}
}

func TestHydroOnElectro(t *testing.T) {

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
	fmt.Println("----testing applying 25 electro on 25 hydro (no delay)----")

	//TEST SWIRL
	target.Attack(&def.Snapshot{
		Durability: 25,
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
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	// expect("apply 25 hydro to 25 electro (tolerance 0.01)", 20, target.aura.Durability())
	// if !durApproxEqual(20, target.aura.Durability(), 0.01) {
	// 	t.Error("ec test: invalid durability")
	// 	t.FailNow()
	// }
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking ec dealt 1 dmg tick after initial app", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//check durability
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, t=0s (frame 1)", 20.0*(1-1.0/570.0), ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-1.0/570.0), ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//tick for 6 more frames, wane should have happened now
	for i := 0; i < 6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("check electro durability after 0.1s, t=0s (frame 1)", (20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//expecting another tick at frame 61
	for i := 0; i < 60-6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("check electro durability at at t=1s (frame 61)", 20.0*(1-61.0/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-61.0/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	expect("checking ec dealt 1 dmg tick at t=1s (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//tick for 6 more frames, ec should be gone now
	for i := 0; i < 6; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("expecting ec to be gone now at t=1.1s (frame 67), nothing left", nil, target.aura)
	if target.aura != nil {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}
}

func TestECChain(t *testing.T) {

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

	fmt.Println("----testing 25 hydro + 25 electro, wait 1 sec, + 25 electro----")
	target.Attack(&def.Snapshot{
		Durability: 25,
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
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	tickCount := 61

	//tick to frame 61, then refresh
	for i := 0; i < tickCount; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	//after 61 frames should have 1 tick wane + decay; 2nd wane not in yet
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, t=1s (frame 61)", (20.0*(1-def.Durability(tickCount)/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-def.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should have 2 dmg counts
	expect("checking ec dealt 2 dmg tick at t=1s (frame 61)", 2, dmgCount)
	if dmgCount != 2 {
		t.Errorf("ec test: expecting 2 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//add 10 electro, should trigger 1 dmg immediately, + wane in 6 frames
	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("check electro durability after + 25 electro, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10+25, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after + 25 electro, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should tick once next frame
	sim.F++
	target.AuraTick()
	target.Tick()
	tickCount++
	expect("checking ec dealt 1 dmg tick after reapply electro (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//wait 5 more frames, we should get 1 wane from the initial tick
	tickCount += 5
	for i := 0; i < 5; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("check electro durability after wane #2 (frame 67)", 20.0*(1-def.Durability(tickCount)/570.0)-20+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-20+25, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	//at this point hydro should be gone already since we had 2 decay
	expect("expecting ec to be gone now at t=1.1s (frame 67), left with electro", ec.electro, target.aura)
	if target.aura != ec.electro {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}

	//tick again, should have another wane
	tickCount++
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("check electro durability after wane #3; even though hydro gone already (frame 68)", 20.0*(1-def.Durability(tickCount)/570.0)-30+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-30+25, ec.electro.CurrentDurability, 0.01) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	//tick a bunch more times to make sure no more damage and wane
	for i := 0; i < 120; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("expecting no more damage ticks", 0, dmgCount)
	if dmgCount != 0 {
		t.Errorf("ec test: expecting 0 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

}

func TestECHydroChain(t *testing.T) {

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

	fmt.Println("----testing 25 hydro + 25 electro, wait 1 sec, + 25 hydro----")
	target.Attack(&def.Snapshot{
		Durability: 25,
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
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	tickCount := 61

	//tick to frame 61, then refresh
	for i := 0; i < tickCount; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	//after 61 frames should have 1 tick wane + decay; 2nd wane not in yet
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, t=1s (frame 61)", (20.0*(1-def.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability)
	if !durApproxEqual((20.0*(1-def.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should have 2 dmg counts
	expect("checking ec dealt 2 dmg tick at t=1s (frame 61)", 2, dmgCount)
	if dmgCount != 2 {
		t.Errorf("ec test: expecting 2 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//add 10 electro, should trigger 1 dmg immediately, + wane in 6 frames
	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Hydro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("check electro durability after + 25 hydro, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after + 25 hydro, t=1s (frame 61)", 20.0*(1-def.Durability(tickCount)/570.0)-10+25, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10+25, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should tick once next frame
	sim.F++
	target.AuraTick()
	target.Tick()
	tickCount++
	expect("checking ec dealt 1 dmg tick after reapply electro (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//wait 5 more frames, we should get 1 wane from the initial tick
	tickCount += 5
	for i := 0; i < 5; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
	}
	expect("check hydro durability after wane #2 (frame 67)", 20.0*(1-def.Durability(tickCount)/570.0)-20+25, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-20+25, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//at this point electro should be gone already since we had 2 decay
	expect("expecting ec to be gone now at t=1.1s (frame 67), left with hydro", ec.hydro, target.aura)
	if target.aura != ec.hydro {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}

	//tick again, should have another wane
	tickCount++
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("check hydro durability after wane #3; even though hydro gone already (frame 68)", 20.0*(1-def.Durability(tickCount)/570.0)-30+25, target.aura.Durability())
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-30+25, target.aura.Durability(), 0.01) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}

}

func TestECSwirl(t *testing.T) {

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
			a, _ := target.Attack(ds)
			if a > 0 {
				dmgCount++
			}
		}

	})

	fmt.Println("----testing 25/25 ec +  25 anemo swirl----")

	target.Attack(&def.Snapshot{
		Durability: 25,
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
		Element:    def.Electro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	tickCount := 0 //enough for 1 tick at f = 1, and wane at f = 1 + 6
	for i := 0; i < 7; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
		tickCount++
	}

	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, f = 7", 20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-def.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, f = 7", (20.0*(1-def.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability)
	if !durApproxEqual((20.0*(1-def.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should have 2 dmg counts
	expect("checking ec dealt 1 dmg tick at f = 7", 2, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	//apply anemo
	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Anemo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	//expecting only hydro left
	hleft := 25 - 2*((20.0*(1-def.Durability(tickCount)/570.0))-10) //amount anemo left after reducing electro
	hleft = ((20.0 * (1 - def.Durability(tickCount)/570.0)) - 10) - 0.5*hleft
	expect("check hydro durability after swirl, f = 7", hleft, ec.hydro.CurrentDurability)
	if !durApproxEqual(hleft, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}

	//tick once, expecting 2 swirls
	for i := 0; i < 1; i++ {
		sim.F++
		target.AuraTick()
		target.Tick()
		tickCount++
	}
	expect("checking swirl dealt 2 dmg tick at f = 8", 2, dmgCount)
	if dmgCount != 2 {
		t.Errorf("ec test: expecting 2 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	log.Println(target.aura.Type())
	//expecting hydro?
	expect("check if ending aura after swirl dmg is hydro", def.Hydro, target.aura.Type())
	if target.aura.Type() != def.Hydro {
		t.Error("ec test: expecting residual hydro aura")
		t.FailNow()
	}
}
