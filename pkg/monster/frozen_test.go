package monster

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func TestFrozenDuration(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	fmt.Println("----testing applying 25 cryo on 50 hydro (no delay)----")

	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Hydro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})
	c.Combat.ApplyDamage(&core.Snapshot{
		CharLvl:    90,
		Durability: 25,
		Element:    core.Cryo,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	dur := 0

	for i := 0; i < 220; i++ {

		if target.aura != nil {
			if target.aura.Type() == core.Frozen {
				dur++
			}
		}
		c.Tick()

	}
	expect("checking freeze duration in frames (25 cryo on 50 hydro no delay)", 209, dur)
	if dur != 209 {
		t.Errorf("freeze test: expecting 209 frames in duration, got %v", dur)
		t.FailNow()
	}

	//extending should add to existing durability? capped at?

}
