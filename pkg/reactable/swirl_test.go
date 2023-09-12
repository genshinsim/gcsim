package reactable

import (
	"fmt"
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

func TestSwirl50to25(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~20")
	c, trg := testCoreWithTrgs(2)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// apply 25 pyro first
	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 25, 1), 0)
	// 1 tick
	advanceCoreFrame(c)
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 50, 1), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
		t.FailNow()
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}

func TestSwirl25to25(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~20")
	c, trg := testCoreWithTrgs(2)

	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 25, 0), 0)
	// 1 tick
	advanceCoreFrame(c)
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 25, 0), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := reactions.Durability(25)*1.25 + 23.75
	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}

func TestSwirl25to50(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~40")
	c, trg := testCoreWithTrgs(2)

	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 50, 0), 0)
	// 1 tick
	advanceCoreFrame(c)
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 25, 0), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := reactions.Durability(25)*1.25 + 23.75
	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}

func TestSwirl50to50(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~40")
	c, trg := testCoreWithTrgs(2)

	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// apply 25 pyro first
	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 50, 0), 0)
	// 1 tick
	advanceCoreFrame(c)
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 50, 0), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := reactions.Durability(50)*1.25 + 23.75

	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}

func TestSwirl25to10(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 25 applied to ~10")
	c, trg := testCoreWithTrgs(2)
	c.Init()

	// apply 25 pyro first
	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 25, 1), 0)
	// tick 285
	for i := 0; i < 285; i++ {
		advanceCoreFrame(c)
	}
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 25, 1), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}

func TestSwirl50to10(t *testing.T) {
	fmt.Println("------------------------------\ntesting swirl 50 applied to ~10")

	c, trg := testCoreWithTrgs(2)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// apply 25 pyro first
	c.QueueAttackEvent(makeSTAttack(attributes.Pyro, 25, 1), 0)
	// tick 285
	for i := 0; i < 285; i++ {
		advanceCoreFrame(c)
	}
	// check durability after 1 tick
	dur := trg[0].Durability[Pyro]
	fmt.Printf("pyro left: %v\n", dur)
	c.QueueAttackEvent(makeSTAttack(attributes.Anemo, 25, 1), 0)
	// dmg should trigger next tick
	// i'm expecting an aoe swirl with durability = dur * 1.25 + 23.75
	expected := dur*1.25 + 23.75
	advanceCoreFrame(c)
	if trg[1].last.Info.Abil != "swirl-pyro (aoe)" {
		t.Errorf("expecting swirl, got %v", trg[1].last.Info.Abil)
	}
	// no durability
	if math.Abs(float64(trg[1].last.Info.Durability-expected)) > float64(ZeroDur) {
		t.Errorf("expected durability to be %v, got %v", expected, trg[1].last.Info.Durability)
	}
}
