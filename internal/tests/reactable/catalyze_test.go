package reactable_test

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

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
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)

	trg := c.Combat.Enemies()[0].(*enemy.Enemy)
	if math.Abs(float64(trg.Durability[reactable.ModifierQuicken])-19.933333) > 0.000001 { // 2f decay
		t.Errorf(
			"expected electro attack to not consume quicken aura, got quicken=%v",
			trg.Durability[reactable.ModifierQuicken],
		)
	}
	if trg.Durability[reactable.ModifierElectro] != 20 {
		t.Errorf(
			"expected electro attack to not reduce, got electro=%v", trg.Durability[reactable.ModifierElectro])
	}
}
