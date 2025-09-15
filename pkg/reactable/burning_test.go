package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/model"
)

func TestBurning(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)

	c.Init()

	//TODO: write tests for burning (this is copypasted from quicken for now)
	trg.AttachOrRefill(&model.AttackEvent{
		Info: model.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	})
	trg.React(&model.AttackEvent{
		Info: model.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})
	// dendro electro gone; 20 quicken
	if !durApproxEqual(20, trg.Durability[Quicken], 0.00001) {
		t.Errorf("expecting 20 cryo attached, got %v", trg.Durability[Quicken])
	}
	if trg.AuraContains(attributes.Dendro, attributes.Electro) {
		t.Error("expecting target to not contain any remaining dendro or electro aura")
	}
}
