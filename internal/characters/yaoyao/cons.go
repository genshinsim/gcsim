package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(a combat.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(enemy.ResistMod{
		Base:  modifier.NewBaseWithHitlag("yaoyao-c1", 6*60),
		Ele:   attributes.Pyro,
		Value: -0.15,
	})
}

func (c *char) c2(done bool) combat.AttackCBFunc {
	return func(atk combat.AttackCB) {
		if done {
			return
		}
		trg, ok := atk.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		if !trg.StatusIsActive(c2Debuff) {
			trg.QueueEnemyTask(c.c2Explode(c.Core.F, trg), 120)
			trg.AddStatus(c2Debuff, 120, true)
		}
		done = true
	}
}
func (c *char) c2Explode(src int, trg *enemy.Enemy) func() {
	return func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Oil Meets Fire (C2)",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       .75,
		}

		c.Core.QueueAttack(ai, combat.NewCircleHit(trg, 2), 0, 0)

		c.Core.Log.NewEvent("Triggered yaoyao C2 explosion", glog.LogCharacterEvent, c.Index).
			Write("src", src)
	}
}

func (c *char) c6(dur int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.15

	c.Core.Status.Add("xlc6", dur)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("yaoyao-c6", dur),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
