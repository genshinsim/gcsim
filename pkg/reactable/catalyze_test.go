package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestQuicken(t *testing.T) {

	c := testCore()
	trg := addTargetToCore(c)

	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})
	//dendro electro gone; 20 quicken
	if !durApproxEqual(20, trg.Durability[ModifierQuicken], 0.00001) {
		t.Errorf("expecting 20 cryo attached, got %v", trg.Durability[ModifierQuicken])
	}
	if trg.AuraContains(attributes.Dendro, attributes.Electro) {
		t.Error("expecting target to not contain any remaining dendro or electro aura")
	}
}
