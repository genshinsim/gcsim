package reactable

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestSwirl50to25(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~20")
	c := testCore()
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	advanceCoreFrame(c)
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 50,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
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
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	advanceCoreFrame(c)
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := combat.Durability(25)*1.25 + 23.75
	advanceCoreFrame(c)
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
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
	})
	//1 tick
	advanceCoreFrame(c)
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := combat.Durability(25)*1.25 + 23.75
	advanceCoreFrame(c)
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
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
	})
	//1 tick
	advanceCoreFrame(c)
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 50,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := combat.Durability(50)*1.25 + 23.75

	advanceCoreFrame(c)
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
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	//tick 285
	for i := 0; i < 285; i++ {
		advanceCoreFrame(c)
	}
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
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
	trg := addTargetToCore(c)
	trg2 := addTargetToCore(c)

	err := c.Init()
	if err != nil {
		fmt.Println("error initialize core: ", err)
		t.FailNow()
	}

	var src *combat.AttackEvent
	trg2.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		// log.Println(atk.Info.Abil)
		// log.Println(atk.Info.Element)
		// log.Println(atk)
		return 0, false
	}

	//apply 25 pyro first
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	//tick 285
	for i := 0; i < 285; i++ {
		advanceCoreFrame(c)
	}
	//check durability after 1 tick
	dur := trg.Durability[attributes.Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Anemo,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	//i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
	if src == nil || src.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", src)
	}
	//no durability
	fmt.Printf("expected durability to be %v, got %v\n", expected, src.Info.Durability)

	if math.Abs(float64(src.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, src.Info.Durability)
	}
}
