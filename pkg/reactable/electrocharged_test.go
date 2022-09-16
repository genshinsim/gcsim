package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestEC(t *testing.T) {

	c := testCore()
	trg := addTargetToCore(c)

	c.Init()

	count := 0
	trg.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		if atk.Info.Abil == "electrocharged" {
			count++
		}
		return 0, false
	}

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})
	//tick once every 60 second. we should get 2 ticks total
	for i := 0; i < 121; i++ {
		advanceCoreFrame(c)
	}
	if count != 2 {
		t.Errorf("expecting 2 ticks, got %v", count)
	}
}
