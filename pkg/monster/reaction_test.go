package monster

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func TestMelt(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	ampCount := 0
	ampMult := 0.0
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
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	c.Events.Subscribe(core.OnAmpReaction, func(args ...interface{}) bool {
		snap := args[1].(*core.Snapshot)
		if snap.ReactionType == core.Melt {
			ampCount++
			ampMult = snap.ReactMult
		}
		return false
	}, "amp-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----melt testing----")

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	ds := &core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	}
	c.Combat.ApplyDamage(ds)
	expect("apply 25 cryo to 50 pyro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("melt test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	expect("checking melt count", 1, ampCount)
	if ampCount != 1 {
		t.Errorf("melt test: expecting 1 melt, got %v", ampCount)
	}
	expect("checking melt multiplier", 1.5, ampMult)
	if !floatApproxEqual(1.5, ampMult, 0.0000001) {
		t.Errorf("melt test: expecting 1.5 multiplier, got %v", ampMult)
	}
	ampCount = 0
	ampMult = 0

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 100,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	ds = &core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	}
	c.Combat.ApplyDamage(ds)

	expect("apply 25 pyro to 100 cryo (tolerance 0.01)", 30, target.aura.Durability())
	if !durApproxEqual(30, target.aura.Durability(), 0.01) {
		t.Error("melt test: invalid durability")
		t.FailNow()
	}
	//check our snapshot, should have been modified
	if ampCount != 1 {
		t.Errorf("melt test: expecting 1 melt, got %v", ampCount)
	}
	expect("checking melt multiplier", 2, ampMult)
	if !floatApproxEqual(2, ampMult, 0.0000001) {
		t.Errorf("melt test: expecting 1.5 multiplier, got %v", ampMult)
	}

}

func TestSuperconduct(t *testing.T) {

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
		snap := args[1].(*core.Snapshot)
		if snap.AttackTag == core.AttackTagSuperconductDamage {
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

	fmt.Println("----superconduct testing----")

	//TEST SUPERCONDUCT
	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Cryo,
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
	expect("apply 25 electro to 25 cryo (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("superconduct test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking superconduct dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("superconduct test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
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
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 cryo to 25 electro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("superconduct test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
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
		snap := args[1].(*core.Snapshot)
		if snap.AttackTag == core.AttackTagOverloadDamage {
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

	fmt.Println("----overload testing----")

	target.aura = nil
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
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 pyro to 25 electro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("overload test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking overload dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Pyro,
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
	expect("apply 25 electro to 25 pyro (tolerance 0.01)", 0, target.aura.Durability())
	if !durApproxEqual(0, target.aura.Durability(), 0.01) {
		t.Error("overload test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking overload dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

}

func TestMultiOverload(t *testing.T) {

	aCount := 0
	bCount := 0
	var targetA *Target
	var targetB *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	targetA = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, targetA)

	targetB = New(1, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, targetB)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		snap := args[1].(*core.Snapshot)
		dmg := args[2].(float64)
		if snap.AttackTag == core.AttackTagOverloadDamage && dmg > 0 {
			if t.Index() == 0 {
				aCount++
			} else {
				bCount++
			}
		}

		return false
	}, "atk-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----multi target overload testing----")

	targetA.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		ActorIndex: 0,
		Durability: 100,
		Element:    core.Electro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    0,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		ActorIndex: 0,
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    0,
		DamageSrc:  -1,
	})
	//we should get 2 ticks of damage here one of each target
	c.Skip(2)

	expect("expecting 2 overload ticks, one on each target", 2, aCount+bCount)
	if aCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage on target A, got %v", aCount)
	}
	if bCount != 1 {
		t.Errorf("overload test: expecting 1 tick of damage on target A, got %v", bCount)
	}
	aCount = 0
	bCount = 0

	//there should be electro left still = 80-25-1 tick delay
	//trigger overload again, should be no dmg this time
	c.Combat.ApplyDamage(&core.Snapshot{
		ActorIndex: 0,
		CharLvl:    90,
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    0,
		DamageSrc:  -1,
	})
	c.Skip(2)
	expect("expecting 0 overload ticks, one on each target", 0, aCount+bCount)
	if aCount != 0 {
		t.Errorf("overload test: expecting 0 tick of damage on target A, got %v", aCount)
	}
	if bCount != 0 {
		t.Errorf("overload test: expecting 0 tick of damage on target B, got %v", bCount)
	}
	aCount = 0
	bCount = 0

}

func TestVaporize(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target
	ampCount := 0
	ampMult := 0.0

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
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	c.Events.Subscribe(core.OnAmpReaction, func(args ...interface{}) bool {
		snap := args[1].(*core.Snapshot)
		if snap.ReactionType == core.Vaporize {
			ampCount++
			ampMult = snap.ReactMult
		}
		return false
	}, "amp-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----vaporize testing----")

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 100,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	ds := &core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	}
	c.Combat.ApplyDamage(ds)
	expect("apply 25 hydro to 100 pyro (tolerance 0.01)", 30, target.aura.Durability())
	if !durApproxEqual(30, target.aura.Durability(), 0.01) {
		t.Error("vaporize test: invalid durability")
		t.FailNow()
	}
	expect("checking vape count", 1, ampCount)
	if ampCount != 1 {
		t.Errorf("vaporize test: expecting 1 vape, got %v", ampCount)
	}
	expect("checking melt multiplier", 2, ampMult)
	if !floatApproxEqual(2, ampMult, 0.0000001) {
		t.Errorf("vaporize test: expecting 2 multiplier, got %v", ampMult)
	}
	ampCount = 0
	ampMult = 0

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	ds = &core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	}
	c.Combat.ApplyDamage(ds)
	expect("apply 25 pyro to 50 hydro (tolerance 0.01)", 27.5, target.aura.Durability())
	if !durApproxEqual(27.5, target.aura.Durability(), 0.01) {
		t.Error("vaporize test: invalid durability")
		t.FailNow()
	}
	expect("checking vape count", 1, ampCount)
	if ampCount != 1 {
		t.Errorf("vaporize test: expecting 1 vape, got %v", ampCount)
	}
	expect("checking melt multiplier", 1.5, ampMult)
	if !floatApproxEqual(1.5, ampMult, 0.0000001) {
		t.Errorf("vaporize test: expecting 2 multiplier, got %v", ampMult)
	}
	ampCount = 0
	ampMult = 0

}

func TestCrystallize(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target
	var shdHP float64
	var shdEle core.EleType

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
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shd := args[0].(core.Shield)
		shdHP = shd.CurrentHP()
		shdEle = shd.Element()
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

	fmt.Println("----crystallize testing----")

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Geo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
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
	expect("checking crystallize shield ele type", core.Pyro, shdEle)
	if shdEle != core.Pyro {
		t.Errorf("crystallize test: expecting pyro shield, got %v", shdEle)
	}

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
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
		Element:    core.Geo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
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
	expect("checking crystallize shield ele type", core.Hydro, shdEle)
	if shdEle != core.Hydro {
		t.Errorf("crystallize test: expecting hydro shield, got %v", shdEle)
	}

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Geo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
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
	expect("checking crystallize shield ele type", core.Cryo, shdEle)
	if shdEle != core.Cryo {
		t.Errorf("crystallize test: expecting cryo shield, got %v", shdEle)
	}

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
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
		Element:    core.Geo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
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
	expect("checking crystallize shield ele type", core.Electro, shdEle)
	if shdEle != core.Electro {
		t.Errorf("crystallize test: expecting electro shield, got %v", shdEle)
	}

}

func TestSwirl(t *testing.T) {

	attackCount := 0
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
		dmg := args[2].(float64)
		snap := args[1].(*core.Snapshot)
		if snap.AttackTag <= core.ReactionAttackDelim {
			return false
		}
		log.Println(args)
		attackCount++
		if dmg > 0 {
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

	fmt.Println("----swirl testing----")

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
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
	expect("apply 25 anemo to 25 cryo (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking swirl triggered 1 attacks", 1, attackCount)
	if attackCount != 1 {
		t.Errorf("swirl test: expecting 1 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	c.Skip(120)

	target.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
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
	expect("apply 25 anemo to 25 pyro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking swirl triggered 1 attacks", 1, attackCount)
	if attackCount != 1 {
		t.Errorf("swirl test: expecting 1 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 1 dmg tick", 1, dmgCount)
	if dmgCount != 1 {
		t.Errorf("swirl test: expecting 1 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	c.Skip(120)

	target.aura = nil
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
		Element:    core.Anemo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 25 hydro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking swirl triggered 1 attacks", 1, attackCount)
	if attackCount != 1 {
		t.Errorf("swirl test: expecting 1 attacks, got %v", attackCount)
	}
	attackCount = 0
	expect("checking swirl dealt 0 dmg tick", 0, dmgCount)
	if dmgCount != 0 {
		t.Errorf("swirl test: expecting 0 tick of damage, got %v", dmgCount)
	}
	dmgCount = 0

	c.Skip(120)

	target.aura = nil
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
		Element:    core.Anemo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 25 electro (tolerance 0.01)", 7.5, target.aura.Durability())
	if !durApproxEqual(7.5, target.aura.Durability(), 0.01) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
	expect("checking swirl triggered 1 attacks", 1, attackCount)
	if attackCount != 1 {
		t.Errorf("swirl test: expecting 1 attacks, got %v", attackCount)
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

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	targetA = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, targetA)

	targetB = New(1, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, targetB)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		dmg := args[2].(float64)
		snap := args[1].(*core.Snapshot)
		if snap.AttackTag <= core.ReactionAttackDelim {
			return false
		}
		attackCount++
		if dmg > 0 {
			if t.Index() == 0 {
				dmgACount++
			} else {
				dmgBCount++
			}

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

	fmt.Println("----multitarget swirl testing----")

	targetA.aura = nil
	targetB.aura = nil
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    0,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Anemo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    0,
		DamageSrc:  -1,
	})
	expect("apply 25 anemo to 50 cryo on target 1 (tolerance 0.001)", 27.5, targetA.aura.Durability())
	if !durApproxEqual(27.5, targetA.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid durability")
		t.FailNow()
	}
	//next tick should deal damage
	c.Tick()
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
	c.Skip(60)
	expect("check target A cryo durability after 60 frames", (40.0*(1-core.Durability(61)/720.0))-12.5, targetA.aura.Durability())
	if !durApproxEqual((40.0*(1-core.Durability(61)/720.0))-12.5, targetA.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid targetA cryo durability")
		t.FailNow()
	}
	expect("check target B cryo durability after 60 frames", (0.8 * 55.0 * (1 - core.Durability(61)/(6*55+420))), targetB.aura.Durability())
	if !durApproxEqual((0.8 * 55.0 * (1 - core.Durability(60)/(6*55+420))), targetB.aura.Durability(), 0.001) {
		t.Error("swirl test: invalid targetB cryo durability")
		t.FailNow()
	}

}
