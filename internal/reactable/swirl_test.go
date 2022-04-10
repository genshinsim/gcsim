package reactable

import (
	"fmt"
	"math"
	"testing"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestSwirl50to25(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~20")
	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	c.Tick()
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 50,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)
	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}

func TestSwirl25to25(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~20")
	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	c.Tick()
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := core.Durability(25)*1.25 + 23.75
	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)
	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}

func TestSwirl25to50(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~40")
	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 50,
		},
	})
	//1 tick
	c.Tick()
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := core.Durability(25)*1.25 + 23.75
	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)
	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}

func TestSwirl50to50(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~40")
	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 50,
		},
	})
	//1 tick
	c.Tick()
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 50,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := core.Durability(50)*1.25 + 23.75

	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)
	//no durability
	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}

func TestSwirl25to10(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~10")

	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 25,
		},
	})
	//tick 285
	for i := 0; i < 285; i++ {
		c.Tick()
	}
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)

	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}

func TestSwirl50to10(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~10")

	c := testCore()

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	trg2 := &testTarget{src: 1}
	trg2.Reactable = &Reactable{}
	trg2.Init(trg, c)
	c.Targets = append(c.Targets, trg2)

	c.Init()

	var src *core.AttackEvent
	trg2.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 25,
		},
	})
	//tick 285
	for i := 0; i < 285; i++ {
		c.Tick()
	}
	//check durability after 1 tick
	dur := trg.Durability[core.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	c.Tick()
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)

	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}
