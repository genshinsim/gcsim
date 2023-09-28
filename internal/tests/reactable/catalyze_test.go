package reactable_test

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestQuicken(t *testing.T) {
	c, _ := makeCore(1)
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
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	trg := c.Combat.Enemies()[0].(*enemy.Enemy)
	if math.Abs(float64(trg.Durability[reactable.Quicken])-19.966666) > 0.000001 {
		t.Errorf(
			"expected quicken=%v, got quicken=%v",
			19.666666,
			trg.Durability[reactable.Quicken],
		)
	}
	if trg.AuraContains(attributes.Dendro, attributes.Electro) {
		t.Error("expecting target to not contain any remaining dendro or electro aura")
	}
}

func TestElectroDoesNotReduceQuicken(t *testing.T) {
	c, _ := makeCore(1)
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
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	trg := c.Combat.Enemies()[0].(*enemy.Enemy)
	if math.Abs(float64(trg.Durability[reactable.Quicken])-19.933333) > 0.000001 { // 2f decay
		t.Errorf(
			"expected electro attack to not consume quicken aura, got quicken=%v",
			trg.Durability[reactable.Quicken],
		)
	}
	if trg.Durability[reactable.Electro] != 20 {
		t.Errorf(
			"expected electro attack to not reduce, got electro=%v", trg.Durability[reactable.Electro])
	}
}
