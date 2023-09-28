package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func TestSuperconduct(t *testing.T) {
	c, trg := testCoreWithTrgs(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(makeAOEAttack(attributes.Cryo, 25), 0)
	c.Tick()
	c.QueueAttackEvent(makeAOEAttack(attributes.Electro, 25), 0)
	advanceCoreFrame(c)
	if trg[0].last.Info.Abil != "superconduct" {
		t.Errorf("expecting superconduct, got %v", trg[0].last.Info.Abil)
	}
}

func TestFrozenSuperconduct(t *testing.T) {
	c, trg := testCoreWithTrgs(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// trigger a freeze
	c.QueueAttackEvent(makeAOEAttack(attributes.Cryo, 25), 0)
	c.Tick()
	c.QueueAttackEvent(makeAOEAttack(attributes.Hydro, 25), 0)
	c.Tick()
	c.QueueAttackEvent(makeAOEAttack(attributes.Electro, 25), 0)
	advanceCoreFrame(c)
	if trg[0].last.Info.Abil != "superconduct" {
		t.Errorf("expecting superconduct, got %v", trg[0].last.Info.Abil)
	}
}
