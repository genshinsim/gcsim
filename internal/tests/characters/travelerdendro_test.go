package characters

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/internal/characters/traveler/common/dendro"
	_ "github.com/genshinsim/gcsim/internal/characters/traveler/dendro/aether"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestTravelerDendroBurstAttach(t *testing.T) {

	c, trg := makeCore(2)
	prof := defProfile(keys.AetherDendro)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.Combat.DefaultTarget = trg[0].Key()
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		log.Printf("hit by %v attack, dur %v", atk.Info.Element, atk.Info.Durability)
		return false
	}, "hit-check")
	advanceCoreFrame(c)

	// use burst to create a ball
	p := make(map[string]int)
	c.Player.Exec(action.ActionBurst, keys.AetherDendro, p)
	for !c.Player.CanQueueNextAction() {
		advanceCoreFrame(c)
	}
	//wait until dendro gadget is created
	for c.Combat.GadgetCount() < 1 {
		advanceCoreFrame(c)
	}
	//skip an additional frame to be safe
	advanceCoreFrame(c)

	//check that gadget has dendro on it
	g := c.Combat.Gadget(0)
	gr, ok := g.(*dendro.LeaLotus)
	if !ok {
		t.Errorf("expecting gadget to be lea lotus. failed")
		t.FailNow()
	}
	log.Println("initial aura string: ", gr.ActiveAuraString())
	if gr.Durability[reactable.ModifierDendro] != 10 {
		t.Errorf("expecting initial 10 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierDendro])
	}

	//pattern only hit gadet
	pattern := combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100)
	pattern.SkipTargets[targets.TargettableEnemy] = true

	// check the cryo attaches
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 100,
		},
		Pattern: pattern,
	}, 0)
	advanceCoreFrame(c)

	log.Println("after applying 100 cyro: ", gr.ActiveAuraString())
	if gr.Durability[reactable.ModifierCryo] != 80 {
		t.Errorf("expecting 80 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierCryo])
	}
	if gr.Durability[reactable.ModifierDendro] != 10 {
		t.Errorf("expecting 10 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierDendro])
	}

}

func TestTravelerDendroBurstPyro(t *testing.T) {

	c, trg := makeCore(1)
	prof := defProfile(keys.AetherDendro)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.Combat.DefaultTarget = trg[0].Key()
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		log.Printf("gadget hit by %v attack, dur %v", atk.Info.Element, atk.Info.Durability)
		return false
	}, "hit-check")
	dmgCount := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Abil == "Lea Lotus Lamp Explosion" {
			dmgCount++
			log.Println("big boom at: ", c.F)
		}
		return false
	}, "hit-check")
	advanceCoreFrame(c)

	// use burst to create a ball
	p := make(map[string]int)
	c.Player.Exec(action.ActionBurst, keys.AetherDendro, p)
	for !c.Player.CanQueueNextAction() {
		advanceCoreFrame(c)
	}
	//wait until dendro gadget is created
	for c.Combat.GadgetCount() < 1 {
		advanceCoreFrame(c)
	}
	//skip an additional frame to be safe
	advanceCoreFrame(c)

	//check that gadget has dendro on it
	g := c.Combat.Gadget(0)
	gr, ok := g.(*dendro.LeaLotus)
	if !ok {
		t.Errorf("expecting gadget to be lea lotus. failed")
		t.FailNow()
	}
	log.Println("initial aura string: ", gr.ActiveAuraString())
	if gr.Durability[reactable.ModifierDendro] != 10 {
		t.Errorf("expecting initial 10 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierDendro])
	}

	//pattern only hit gadet
	pattern := combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100)
	pattern.SkipTargets[targets.TargettableEnemy] = true

	// check the cryo attaches
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 100,
		},
		Pattern: pattern,
	}, 0)
	advanceCoreFrame(c)

	log.Printf("at f %v after applying 100 pyro: %v\n", c.F, gr.ActiveAuraString())
	if gr.Durability[reactable.ModifierPyro] != 0 {
		t.Errorf("expecting 0 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierPyro])
	}

	//should get an explosion 60 frfames later
	for i := 0; i < 100; i++ {
		advanceCoreFrame(c)
	}

	if dmgCount != 1 {
		t.Errorf("expected 1 dmg count, got %v", dmgCount)
	}

}

// lotus is expected to tick at frame 37 after appearing, which is 54+37 after cast
// and then tick every 90 frames after that for the duration
// duration is either 12s at c0 or 15s at c2
func TestTravelerDendroBurstTicks(t *testing.T) {
	c, trg := makeCore(1)
	prof := defProfile(keys.AetherDendro)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.Combat.DefaultTarget = trg[0].Key()
	dmgCount := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Abil == "Lea Lotus Lamp" {
			dmgCount++
			log.Println("boom at (adjusted): ", c.F-54-1)
		}
		return false
	}, "hit-check")
	advanceCoreFrame(c)

	// use burst to create a ball
	p := make(map[string]int)
	log.Println("casting burst: ", c.F)
	c.Player.Exec(action.ActionBurst, keys.AetherDendro, p)

	//expecting to take a total of 54 frames to appear + 15s duration
	totalDuration := 15 * 60
	expectedCount := 1 + (totalDuration-37)/90

	//add 100 for good measures in case bugs from extra ticks
	for i := 0; i < 54+totalDuration+100; i++ {
		advanceCoreFrame(c)
	}

	if dmgCount != expectedCount {
		t.Errorf("expecting %v ticks, got %v", expectedCount, dmgCount)
	}

}

func TestTravelerDendroBurstElectroTicks(t *testing.T) {
	c, trg := makeCore(1)
	prof := defProfile(keys.AetherDendro)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	c.Combat.DefaultTarget = trg[0].Key()
	dmgCount := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Abil == "Lea Lotus Lamp" {
			dmgCount++
			log.Println("boom at (adjusted): ", c.F-54-1)
		}
		return false
	}, "hit-check")
	advanceCoreFrame(c)

	// use burst to create a ball
	p := make(map[string]int)
	log.Println("casting burst: ", c.F)
	c.Player.Exec(action.ActionBurst, keys.AetherDendro, p)
	//wait until dendro gadget is created
	for c.Combat.GadgetCount() < 1 {
		advanceCoreFrame(c)
	}

	//pattern only hit gadet
	pattern := combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100)
	pattern.SkipTargets[targets.TargettableEnemy] = true

	// check the cryo attaches
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 100,
		},
		Pattern: pattern,
	}, 0)

	//first tick at 15, then tick every 54 after that
	totalDuration := 15 * 60
	expectedCount := 1 + (totalDuration-15)/54

	//add 100 for good measures in case bugs from extra ticks
	for i := 0; i < totalDuration+100; i++ {
		advanceCoreFrame(c)
	}

	if dmgCount != expectedCount {
		t.Errorf("expecting %v ticks, got %v", expectedCount, dmgCount)
	}

}
