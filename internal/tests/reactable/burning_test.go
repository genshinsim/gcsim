package reactable_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestBurningTicks(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	// expecting 8 ticks: https://www.youtube.com/watch?v=PdZ6Qxo7pSY
	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.AttackTag == attacks.AttackTagBurningDamage {
			count++
		}
		return false
	}, "burning-ticks")

	// yanfei auto at 80
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 70)
	// tighnari skill at 200
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 200)
	// lisa 250
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 250)

	// burning starts ticking at 200 and ticks every 15 frames
	for c.F = 0; c.F < 200; c.F++ {
		c.Tick()
	}
	// burning got queued at f = 200 but first tick actually happens at beginning of
	// 216; so we advanced 1 extra here??
	//TODO: does this need to be adjusted somehow? i think this has to do with the fact
	// that the task got added AFTER the run at f 200 so that's why it doesn't get
	// executed until 201, then delay 15 so we end up at 216 first tick instead of 215
	advanceCoreFrame(c)

	// log.Printf("count should be 0 right now, got %v", count)
	for i := 0; i < 8; i++ {
		for j := 0; j < 15; j++ {
			advanceCoreFrame(c)
		}
		// log.Printf("count should be %v right now, got %v", i+1, count)
	}

	// extra 200 frames to make sure it doesn't go past 8
	for i := 0; i < 200; i++ {
		advanceCoreFrame(c)
	}

	if count != 8 {
		t.Errorf("expecting 8 burning ticks, got %v", count)
	}
}

func TestBurningQuickenFuel(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	//https://www.youtube.com/watch?v=En3Ki_vVgR0
	count := 0
	countByActor := []int{0, 0}
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.AttackTag == attacks.AttackTagBurningDamage {
			count++
			countByActor[ae.Info.ActorIndex]++
		}
		return false
	}, "burning-ticks")

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 290)
	// beidou e should apply hitlag here
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 327)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
			ActorIndex: 0,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 396)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
			ActorIndex: 1,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 462)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 100,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 536)

	f := make(map[event.Event]int)
	cb := func(evt event.Event) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			f[evt] = c.F
			return false
		}
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, cb(i), fmt.Sprintf("event-%v", i))
	}
	i := 0
	// quicken reaction at 327
	for ; i < 327; i++ {
		advanceCoreFrame(c)
	}
	log.Printf("quicken at %v\n", f[event.OnQuicken])

	// burning reaction at 396
	for ; i < 396; i++ {
		advanceCoreFrame(c)
	}
	log.Printf("burning at %v\n", f[event.OnBurning])

	// spread at 462
	for ; i < 462; i++ {
		advanceCoreFrame(c)
	}
	log.Printf("spread at %v\n", f[event.OnSpread])

	// overload, quicken, aggrvate at 536
	for ; i < 535; i++ {
		advanceCoreFrame(c)
	}
	advanceCoreFrame(c)
	log.Printf("overload at %v\n", f[event.OnOverload])
	log.Printf("quicken at %v\n", f[event.OnQuicken])
	log.Printf("aggravate at %v\n", f[event.OnAggravate])

	// 4 burning ticks at yanfei's em, 4 ticks at tighnari em (applied at 462), last tick roughly 523
	log.Printf("number of burning ticks - actor 0: %v", countByActor[0])
	log.Printf("number of burning ticks - actor 1: %v", countByActor[1])

	// dendro or quicken last frame 796
	for ; i < 2000; i++ {
		advanceCoreFrame(c)
		if trg[0].Durability[reactable.ModifierQuicken] == 0 {
			log.Printf("quicken gone at f: %v\n", c.F)
			break
		}
	}
}

func TestPyroDendroCoexist(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	//https://www.youtube.com/watch?v=dXzQTNCYfeU&list=PL10DrkffqpyuwG8i0JOq-TgcqPES6bsja&index=16
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 121)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 133)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 195)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 202)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 249)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 327)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 344)
	// pyro ended 546, dendro ended 689

	f := make(map[event.Event]int)
	cb := func(evt event.Event) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			f[evt] = c.F
			return false
		}
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, cb(i), fmt.Sprintf("event-%v", i))
	}
	i := 0

	for ; i < 1000; i++ {
		advanceCoreFrame(c)
		fmt.Printf("%v: %v\n", i, trg[0].ActiveAuraString())
	}

}

func TestDendroDecayTry1(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 155)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 168)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 230)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 263)

	f := make(map[event.Event]int)
	cb := func(evt event.Event) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			f[evt] = c.F
			return false
		}
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, cb(i), fmt.Sprintf("event-%v", i))
	}
	i := 0

	for ; i < 1900; i++ {
		if i == 263 {
			fmt.Println("hi")
		}
		advanceCoreFrame(c)
		fmt.Printf("%v: %v\n", i, trg[0].ActiveAuraString())
	}

}

func TestDendroDecayTry2(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 80)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 440)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 453)

	f := make(map[event.Event]int)
	cb := func(evt event.Event) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			f[evt] = c.F
			return false
		}
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, cb(i), fmt.Sprintf("event-%v", i))
	}
	i := 0

	for ; i < 1900; i++ {
		if i == 460 {
			fmt.Println("hi")
		}
		advanceCoreFrame(c)
		fmt.Printf("%v: %v\n", i, trg[0].ActiveAuraString())
	}

}

func TestQuickenBurningDecay(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 61)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 128)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 188)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 206)

	f := make(map[event.Event]int)
	cb := func(evt event.Event) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			f[evt] = c.F
			return false
		}
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, cb(i), fmt.Sprintf("event-%v", i))
	}
	i := 0

	for ; i < 1900; i++ {
		if i == 460 {
			fmt.Println("hi")
		}
		advanceCoreFrame(c)
		fmt.Printf("%v: %v\n", i, trg[0].ActiveAuraString())
	}

}
