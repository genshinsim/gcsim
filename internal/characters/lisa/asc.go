package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Hits by Charged Attacks apply Violet Arc's Conductive status to opponents.
func (c *char) makeA1CB() info.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a info.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		count := t.GetTag(conductiveTag)
		if count < 3 {
			t.SetTag(conductiveTag, count+1)
		}
	}
}

// Opponents hit by Lightning Rose have their DEF decreased by 15% for 10s.
func (c *char) makeA4CB() info.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a info.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		t.AddDefMod(info.DefMod{
			Base:  modifier.NewBaseWithHitlag("lisa-a4", 600),
			Value: -0.15,
		})
	}
}
