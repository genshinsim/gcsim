package yaemiko

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.YaeMiko, core.Electro, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	yaemiko := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	c.CharPos[prof.Base.Key] = 0
	c.Init()
	//add targets to test with
	eProf := testhelper.EnemyProfile()
	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, eProf))

	//check skill is ready
	p := make(map[string]int)
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be ready at start")
		t.FailNow()
	}

	//use skill, wait out animation, and check if ready for another use (3 stacks at c0)
	a, _ := x.Skill(p)
	for i := 0; i < a; i++ {
		c.Tick()
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected second skill charge to be ready at start. At frame %v", c.F)
		t.FailNow()
	}

	//use skill, wait out animation, and check if ready for another use (3 stacks at c0)
	a, _ = x.Skill(p)
	for i := 0; i < a; i++ {
		c.Tick()
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected third skill charge to be ready at start. At frame %v", c.F)
		t.FailNow()
	}
	//use third charge, next charge should be ready at 9*60-a-1
	x.Skill(p)
	for i := 0; i < 4*60-2*a-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}

	//tick once to get to 9*60-a
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v; charges: %v", c.F, x.Cooldown(core.ActionSkill), yaemiko.availableCDCharge)
		t.FailNow()
	}

	//next charge should be ready at 9*60-a-1 from now
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}

	for i := 0; i < 4*60-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

	//next charge should be ready at 900 from now
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}

	for i := 0; i < 4*60-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

	//use skill and then trigger flat cd reduction
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}
	//next charge should be ready by 900 - flat cd reduction
	x.ReduceActionCooldown(core.ActionSkill, 100)
	for i := 0; i < 4*60-101; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

}
