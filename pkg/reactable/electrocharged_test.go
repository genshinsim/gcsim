package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func TestEC(t *testing.T) {
	c, _ := testCoreWithTrgs(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if ae, ok := args[1].(*combat.AttackEvent); ok {
			if ae.Info.Abil == "electrocharged" {
				count++
			}
		}
		return false
	}, "ec-dmg")

	c.QueueAttackEvent(makeAOEAttack(attributes.Hydro, 25), 0)
	c.Tick()
	c.QueueAttackEvent(makeAOEAttack(attributes.Electro, 25), 0)
	// tick once every 60 second. we should get 2 ticks total
	for i := 0; i < 121; i++ {
		advanceCoreFrame(c)
	}
	if count != 2 {
		t.Errorf("expecting 2 ticks, got %v", count)
	}
}
