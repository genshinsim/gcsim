package monster

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/def"
)

func TestMelt(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		target = New(0, s, logger, def.EnemyProfile{
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

	fmt.Println("----melt testing----")

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	ds := &def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	}
	target.Attack(ds)
	expect("apply 25 cryo to 50 pyro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("melt test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	expect("checking melt is set", def.Melt, ds.ReactionType)
	if ds.ReactionType != def.Melt {
		t.Errorf("melt test: expecting melt flag set, got %v", ds.ReactionType)
	}
	expect("checking melt multiplier", 1.5, ds.ReactMult)
	if !floatApproxEqual(1.5, ds.ReactMult, 0.0000001) {
		t.Errorf("melt test: expecting 1.5 multiplier, got %v", ds.ReactMult)
	}

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 100,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	ds = &def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	}
	target.Attack(ds)

	expect("apply 25 pyro to 100 cryo (tolerance 0.01)", 30, target.aura.Durability())
	if !durApproxEqual(30, target.aura.Durability(), 0.01) {
		t.Error("melt test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	expect("checking melt is set", def.Melt, ds.ReactionType)
	if ds.ReactionType != def.Melt {
		t.Errorf("melt test: expecting melt flag set, got %v", ds.ReactionType)
	}
	expect("checking melt multiplier", 2, ds.ReactMult)
	if !floatApproxEqual(2, ds.ReactMult, 0.0000001) {
		t.Errorf("melt test: expecting 1.5 multiplier, got %v", ds.ReactMult)
	}

}

func TestSuperconduct(t *testing.T) {

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

		target = New(0, s, logger, def.EnemyProfile{
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

	fmt.Println("----superconduct testing----")

	//TEST SUPERCONDUCT
	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Cryo,
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
	expect("apply 25 electro to 25 cryo (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("superconduct test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking superconduct dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("superconduct test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
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
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 cryo to 25 electro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("superconduct test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking superconduct dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("superconduct test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

}

func TestOverload(t *testing.T) {

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

		target = New(0, s, logger, def.EnemyProfile{
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

	fmt.Println("----overload testing----")

	target.aura = nil
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
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 pyro to 25 electro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("overload test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking overload dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Pyro,
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
	expect("apply 25 electro to 25 pyro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("overload test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking overload dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

}

func TestVaporize(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		target = New(0, s, logger, def.EnemyProfile{
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

	fmt.Println("----vaporize testing----")

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 100,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	ds := &def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Hydro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	}
	target.Attack(ds)
	expect("apply 25 hydro to 100 pyro (tolerance 0.01)", 30, target.aura.Durability())
	if !durApproxEqual(30, target.aura.Durability(), 0.01) {
		t.Error("vaporize test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	expect("checking vape is set", def.Vaporize, ds.ReactionType)
	if ds.ReactionType != def.Vaporize {

		t.Errorf("vaporize test: expecting vaporize flag set, got %v", ds.ReactionType)
	}
	expect("checking vape multiplier", 2, ds.ReactMult)
	if !floatApproxEqual(2, ds.ReactMult, 0.0000001) {
		t.Errorf("vaporize test: expecting 2.0 multiplier, got %v", ds.ReactMult)
	}

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Hydro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	ds = &def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	}
	target.Attack(ds)
	expect("apply 25 pyro to 50 hydro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("vaporize test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	expect("checking vape is set", def.Vaporize, ds.ReactionType)
	if ds.ReactionType != def.Vaporize {
		t.Errorf("vaporize test: expecting vaporize flag set, got %v", ds.ReactionType)
	}
	expect("checking vape multiplier", 1.5, ds.ReactMult)
	if !floatApproxEqual(1.5, ds.ReactMult, 0.0000001) {
		t.Errorf("vaporize test: expecting 2.0 multiplier, got %v", ds.ReactMult)
	}

}

func TestCrystallize(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var shdHP float64
	var target *Target
	var shdEle def.EleType

	dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		target = New(0, s, logger, def.EnemyProfile{
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
			shdHP = shd.CurrentHP()
			shdEle = shd.Element()
		}

	})

	fmt.Println("----crystallize testing----")

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Geo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 geo to 50 pyro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("crystallize test: invalid durability")
		t.FailNow()
	}
	expect("checking crystallize added shield", 1, shdCount)
	if shdCount != 1 {
		t.Errorf("crystallize test: expecting 1 shield, got %v", dmgCount)
	}
	shdCount = 0
	expect("checking crystallize shield base hp (lvl 90)", 1851.06030273438, shdHP)
	if !floatApproxEqual(1851.06030273438, shdHP, 0.0000001) {
		t.Errorf("crystallize test: expecting shield hp = 1851.06030273438, got %v", shdHP)
	}
	expect("checking crystallize shield ele type", def.Pyro, shdEle)
	if shdEle != def.Pyro {
		t.Errorf("crystallize test: expecting pyro shield, got %v", shdEle)
	}

	target.aura = nil
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
		Element:    def.Geo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 geo to 50 hydro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("crystallize test: invalid durability")
		t.FailNow()
	}
	expect("checking crystallize added shield", 1, shdCount)
	if shdCount != 1 {
		t.Errorf("crystallize test: expecting 1 shield, got %v", dmgCount)
	}
	shdCount = 0
	expect("checking crystallize shield base hp (lvl 90)", 1851.06030273438, shdHP)
	if !floatApproxEqual(1851.06030273438, shdHP, 0.0000001) {
		t.Errorf("crystallize test: expecting shield hp = 1851.06030273438, got %v", shdHP)
	}
	expect("checking crystallize shield ele type", def.Hydro, shdEle)
	if shdEle != def.Hydro {
		t.Errorf("crystallize test: expecting hydro shield, got %v", shdEle)
	}

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})

	target.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Geo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 geo to 50 cryo (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("crystallize test: invalid durability")
		t.FailNow()
	}
	expect("checking crystallize added shield", 1, shdCount)
	if shdCount != 1 {
		t.Errorf("crystallize test: expecting 1 shield, got %v", dmgCount)
	}
	shdCount = 0
	expect("checking crystallize shield base hp (lvl 90)", 1851.06030273438, shdHP)
	if !floatApproxEqual(1851.06030273438, shdHP, 0.0000001) {
		t.Errorf("crystallize test: expecting shield hp = 1851.06030273438, got %v", shdHP)
	}
	expect("checking crystallize shield ele type", def.Cryo, shdEle)
	if shdEle != def.Cryo {
		t.Errorf("crystallize test: expecting cryo shield, got %v", shdEle)
	}

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 50,
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
		Element:    def.Geo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 geo to 50 electro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("crystallize test: invalid durability")
		t.FailNow()
	}
	expect("checking crystallize added shield", 1, shdCount)
	if shdCount != 1 {
		t.Errorf("crystallize test: expecting 1 shield, got %v", dmgCount)
	}
	shdCount = 0
	expect("checking crystallize shield base hp (lvl 90)", 1851.06030273438, shdHP)
	if !floatApproxEqual(1851.06030273438, shdHP, 0.0000001) {
		t.Errorf("crystallize test: expecting shield hp = 1851.06030273438, got %v", shdHP)
	}
	expect("checking crystallize shield ele type", def.Electro, shdEle)
	if shdEle != def.Electro {
		t.Errorf("crystallize test: expecting electro shield, got %v", shdEle)
	}

}

func TestSwirl(t *testing.T) {

	dmgCount := 0
	attackCount := 0
	shdCount := 0
	var target *Target

	sim := dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		target = New(0, s, logger, def.EnemyProfile{
			Level:  88,
			Resist: defaultResMap(),
		})

		s.OnDamage = func(ds *def.Snapshot) {
			// log.Println(ds)
			// log.Println(target.attackWillLand(ds))
			// log.Println(ds.Durability)
			attackCount++
			dmg := target.Attack(ds)
			if dmg > 0 {
				dmgCount++
			}
		}

		s.OnShielded = func(shd def.Shield) {
			// log.Println(shd.CurrentHP())
			shdCount++
		}

	})

	fmt.Println("----swirl testing----")

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
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
	expect("apply 25 anemo to 25 cryo (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking swirl triggered 2 attacks", 2, attackCount)
	if attackCount != 2 {
		t.Errorf("swirl test: expecting 2 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
	target.Attack(&def.Snapshot{
		Durability: 25,
		Element:    def.Pyro,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
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
	expect("apply 25 anemo to 25 pyro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking swirl triggered 2 attacks", 2, attackCount)
	if attackCount != 2 {
		t.Errorf("swirl test: expecting 2 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
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
		Element:    def.Anemo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 25 hydro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking swirl triggered 2 attacks", 2, attackCount)
	if attackCount != 2 {
		t.Errorf("swirl test: expecting 2 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
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
		Element:    def.Anemo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 25 electro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	target.AuraTick()
	target.Tick()
	expect("checking swirl triggered 2 attacks", 2, attackCount)
	if attackCount != 2 {
		t.Errorf("swirl test: expecting 2 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
}

func TestSwirlMultiTarget(t *testing.T) {

	var dmgACount, dmgBCount, attackCount int
	shdCount := 0
	var targetA *Target
	var targetB *Target

	sim := dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		char := dummy.NewChar(func(c *dummy.Char) {
			c.Stats = make([]float64, def.EndStatType)
			c.Stats[def.EM] = 100
		})

		s.Chars = append(s.Chars, char)

		targetA = New(0, s, logger, def.EnemyProfile{
			Level:  88,
			Resist: defaultResMap(),
		})
		targetB = New(1, s, logger, def.EnemyProfile{
			Level:  88,
			Resist: defaultResMap(),
		})

		s.OnDamage = func(ds *def.Snapshot) {
			// log.Println(ds)
			// log.Println(target.attackWillLand(ds))
			// log.Println(ds.Durability)
			attackCount++
			if targetA.Attack(ds) > 0 {
				dmgACount++
			}
			if targetB.Attack(ds) > 0 {
				dmgBCount++
			}
		}

		s.OnShielded = func(shd def.Shield) {
			// log.Println(shd.CurrentHP())
			shdCount++
		}

	})

	fmt.Println("----multitarget swirl testing----")

	targetA.aura = nil
	targetB.aura = nil
	targetA.Attack(&def.Snapshot{
		Durability: 50,
		Element:    def.Cryo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	targetA.Attack(&def.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    def.Anemo,
		ICDTag:     def.ICDTagNone,
		ICDGroup:   def.ICDGroupDefault,
		Stats:      make([]float64, def.EndStatType),
		Targets:    def.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 50 cryo on target 1 (tolerance 0.001)", 27.5, targetA.aura.Durability())
	if !durApproxEqual(27.5, targetA.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	sim.F++
	targetA.AuraTick()
	targetB.AuraTick()
	targetA.Tick()
	targetB.Tick()
	// targetB.Tick()
	expect("checking swirl triggered 2 attacks", 2, attackCount)
	if attackCount != 2 {
		t.Errorf("swirl test: expecting 2 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick to target A", 1, dmgACount)
	if dmgACount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgACount)
	}
	dmgACount = 0
	expect("checking swirl dealt 1 dmg tick to target B", 1, dmgBCount)
	if dmgBCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgBCount)
	}
	dmgBCount = 0
	//check target B should now have 55 cryo
	expect("target 2 should have cryo from swirl (tolerance 0.001)", 0.8*55, targetB.aura.Durability())
	if !durApproxEqual(0.8*55, targetB.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid targetB cryo durability")
		t.FailNow()
	}
	//tick some and then check decay
	for i := 0; i < 60; i++ {
		sim.F++
		targetA.AuraTick()
		targetB.AuraTick()
		targetA.Tick()
		targetB.Tick()
	}
	expect("check target A cryo durability after 60 frames", (40.0*(1-def.Durability(61)/720.0))-12.5, targetA.aura.Durability())
	if !durApproxEqual((40.0*(1-def.Durability(61)/720.0))-12.5, targetA.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid targetA cryo durability")
		t.FailNow()
	}
	expect("check target B cryo durability after 60 frames", (0.8 * 55.0 * (1 - def.Durability(61)/(6*55+420))), targetB.aura.Durability())
	if !durApproxEqual((0.8 * 55.0 * (1 - def.Durability(60)/(6*55+420))), targetB.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid targetB cryo durability")
		t.FailNow()
	}

}
