package monster

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func TestElectroOnHydro(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New(func(c *core.Core) error {
		c.Log, _ = core.NewDefaultLogger(c, false, false, nil)
		return nil
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.AttackTag != core.AttackTagECDamage {
			return false
		}
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	//TEST SWIRL
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	// expect("apply 25 hydro to 25 electro (tolerance 0.01)", 20, target.aura.Durability())
	// if !durApproxEqual(20, target.aura.Durability(), 0.01) {
	// 	t.Error("ec test: invalid durability")
	// 	t.FailNow()
	// }
	//next tick should deal damage
	c.Tick()

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
	c.Skip(6)

	expect("check electro durability after 0.1s, t=0s (frame 1)", (20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//expecting another tick at frame 61
	c.Skip(60 - 6)
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
	c.Skip(6)
	expect("expecting ec to be gone now at t=1.1s (frame 67), nothing left", nil, target.aura)
	if target.aura != nil {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}
}

func TestHydroOnElectro(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.AttackTag != core.AttackTagECDamage {
			return false
		}
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----testing applying 25 electro on 25 hydro (no delay)----")

	//TEST SWIRL
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	// expect("apply 25 hydro to 25 electro (tolerance 0.01)", 20, target.aura.Durability())
	// if !durApproxEqual(20, target.aura.Durability(), 0.01) {
	// 	t.Error("ec test: invalid durability")
	// 	t.FailNow()
	// }
	//next tick should deal damage
	c.Tick()
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
	c.Skip(6)
	expect("check electro durability after 0.1s, t=0s (frame 1)", (20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-7.0/570.0))-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid durability")
		t.FailNow()
	}
	//expecting another tick at frame 61
	c.Skip(60 - 6)
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
	c.Skip(6)
	expect("expecting ec to be gone now at t=1.1s (frame 67), nothing left", nil, target.aura)
	if target.aura != nil {
		t.Error("ec test: invalid aura")
		t.FailNow()
	}
}

func TestECChain(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.AttackTag != core.AttackTagECDamage {
			return false
		}
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----testing 25 hydro + 25 electro, wait 1 sec, + 25 electro----")
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	tickCount := 61

	//tick to frame 61, then refresh
	c.Skip(61)
	//after 61 frames should have 1 tick wane + decay; 2nd wane not in yet
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, t=1s (frame 61)", (20.0*(1-core.Durability(tickCount)/570.0))-10, ec.electro.CurrentDurability)
	if !durApproxEqual((20.0*(1-core.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
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
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("check electro durability after + 25 electro, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10+25, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after + 25 electro, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should tick once next frame
	c.Tick()
	tickCount++
	expect("checking ec dealt 1 dmg tick after reapply electro (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//wait 5 more frames, we should get 1 wane from the initial tick
	tickCount += 5
	c.Skip(5)
	expect("check electro durability after wane #2 (frame 67)", 20.0*(1-core.Durability(tickCount)/570.0)-20+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-20+25, ec.electro.CurrentDurability, 0.001) {
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
	c.Tick()
	expect("check electro durability after wane #3; even though hydro gone already (frame 68)", 20.0*(1-core.Durability(tickCount)/570.0)-30+25, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-30+25, ec.electro.CurrentDurability, 0.01) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	//tick a bunch more times to make sure no more damage and wane
	c.Skip(120)
	expect("expecting no more damage ticks", 0, dmgCount)
	if dmgCount != 0 {
		t.Errorf("ec test: expecting 0 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

}

func TestECHydroChain(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.AttackTag != core.AttackTagECDamage {
			return false
		}
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----testing 25 hydro + 25 electro, wait 1 sec, + 25 hydro----")
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	tickCount := 61

	//tick to frame 61, then refresh
	c.Skip(61)
	//after 61 frames should have 1 tick wane + decay; 2nd wane not in yet
	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, t=1s (frame 61)", (20.0*(1-core.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability)
	if !durApproxEqual((20.0*(1-core.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
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
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("check electro durability after + 25 hydro, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after + 25 hydro, t=1s (frame 61)", 20.0*(1-core.Durability(tickCount)/570.0)-10+25, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10+25, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}
	//should tick once next frame
	c.Tick()
	tickCount++
	expect("checking ec dealt 1 dmg tick after reapply electro (frame 61)", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("ec test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0
	//wait 5 more frames, we should get 1 wane from the initial tick
	tickCount += 5
	c.Skip(5)
	expect("check hydro durability after wane #2 (frame 67)", 20.0*(1-core.Durability(tickCount)/570.0)-20+25, ec.hydro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-20+25, ec.hydro.CurrentDurability, 0.001) {
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
	c.Tick()
	expect("check hydro durability after wane #3; even though hydro gone already (frame 68)", 20.0*(1-core.Durability(tickCount)/570.0)-30+25, target.aura.Durability())
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-30+25, target.aura.Durability(), 0.01) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}

}

func TestECSwirl(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.AttackTag == core.AttackTagECDamage || ds.AttackTag == core.AttackTagSwirlElectro || ds.AttackTag == core.AttackTagSwirlHydro {
			dmgCount++
		}
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----testing 25/25 ec +  25 anemo swirl----")

	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	tickCount := 0 //enough for 1 tick at f = 1, and wane at f = 1 + 6
	for i := 0; i < 7; i++ {
		c.Tick()
		tickCount++
	}

	ec, ok := target.aura.(*AuraEC)
	if !ok {
		t.Errorf("expecting aura to cast to EC but failed; got %v", target.aura.Type())
		t.FailNow()
	}
	expect("check electro durability after tick, f = 7", 20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability)
	if !durApproxEqual(20.0*(1-core.Durability(tickCount)/570.0)-10, ec.electro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid electro durability")
		t.FailNow()
	}
	expect("check hydro durability after 0.1s, f = 7", (20.0*(1-core.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability)
	if !durApproxEqual((20.0*(1-core.Durability(tickCount)/570.0))-10, ec.hydro.CurrentDurability, 0.001) {
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
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Anemo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	//expecting only hydro left
	hleft := 25 - 2*((20.0*(1-core.Durability(tickCount)/570.0))-10) //amount anemo left after reducing electro
	hleft = ((20.0 * (1 - core.Durability(tickCount)/570.0)) - 10) - 0.5*hleft
	expect("check hydro durability after swirl, f = 7", hleft, ec.hydro.CurrentDurability)
	if !durApproxEqual(hleft, ec.hydro.CurrentDurability, 0.001) {
		t.Error("ec test: invalid hydro durability")
		t.FailNow()
	}

	//tick once, expecting 2 swirls
	for i := 0; i < 1; i++ {
		c.Tick()
		tickCount++
	}
	expect("checking swirl dealt 2 dmg tick at f = 8", 2, dmgCount)
	if dmgCount != 2 {
		t.Errorf("ec test: expecting 2 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	log.Println(target.aura.Type())
	//expecting hydro?
	expect("check if ending aura after swirl dmg is hydro", core.Hydro, target.aura.Type())
	if target.aura.Type() != core.Hydro {
		t.Error("ec test: expecting residual hydro aura")
		t.FailNow()
	}
}
