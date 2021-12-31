package reactable

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestOverload(t *testing.T) {
	ctrl := &combatTestCtrl{}
	c, err := core.New(func(c *core.Core) error {
		ctrl.core = c
		c.Combat = ctrl
		c.Log = logger
		ctrl.Init(c)
		return nil
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)
	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)
	c.Targets = append(c.Targets, trg)

	c.Init()

	var src *core.AttackEvent
	trg.onDmgCallBack = func(atk *core.AttackEvent) (float64, bool) {
		src = atk
		log.Println(atk.Info.Abil)
		log.Println(atk.Info.Element)
		log.Println(atk)
		return 0, false
	}

	//electro into pyro
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	trg.Tick()
	next := &core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Electro,
			Durability: 25,
		},
	}
	trg.React(next)
	//dmg should trigger next tick
	trg.Tick()
	c.Tick()
	if src == nil || src.Info.Abil != "overload" {
		t.Errorf("expecting overload, got %v", src)
	}
	//no durability
	if next.Info.Durability > ZeroDur {
		t.Errorf("expected durability to be 0, got %v", next.Info.Durability)
	}
}
