package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func TestOverload(t *testing.T) {
	c, trg := testCoreWithTrgs(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(makeAOEAttack(attributes.Pyro, 25), 0)
	c.Tick()
	c.QueueAttackEvent(makeAOEAttack(attributes.Electro, 25), 0)
	advanceCoreFrame(c)
	if trg[0].last.Info.Abil != "overload" {
		t.Errorf("expecting overload, got %v", trg[0].last.Info.Abil)
	}
}
