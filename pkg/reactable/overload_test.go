package reactable

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestOverload(t *testing.T) {

	c := testCore()

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

	var src *combat.AttackEvent
	trg.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		src = atk
		log.Println(atk.Info.Abil)
		log.Println(atk.Info.Element)
		log.Println(atk)
		return 0, false
	}

	//electro into pyro
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	//1 tick
	trg.Tick()
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
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
