package characters

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/internal/characters/travelerdendro"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
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
	gr, ok := g.(*travelerdendro.LeaLotus)
	if !ok {
		t.Errorf("expecting gadget to be lea lotus. failed")
		t.FailNow()
	}
	log.Println("initial aura string: ", gr.ActiveAuraString())
	if gr.Durability[reactable.ModifierDendro] != 20 {
		t.Errorf("expecting initial 20 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierDendro])
	}

	//pattern only hit gadet
	pattern := combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100)
	pattern.SkipTargets[combat.TargettableEnemy] = true

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
	if gr.Durability[reactable.ModifierDendro] != 20 {
		t.Errorf("expecting 20 dendro on traveler lea lotus, got %v", gr.Durability[reactable.ModifierDendro])
	}

}
