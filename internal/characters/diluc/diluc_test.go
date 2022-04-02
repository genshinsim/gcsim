package diluc

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Diluc, core.Pyro, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testhelper.TestClaymoreCharacter(c, x)
}

func TestSkillCDCon0(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Diluc, core.Pyro, 6)
	prof.Base.Cons = 0
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
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
	p := make(map[string]int)

	// sh := x.(*char)

	x.Skill(p)
	//wait 4s, skill should go into cooldown for 600 - 240
	for i := 0; i < 239; i++ {
		c.Tick()
	}
	//should not be in cd yet
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready still (2nd charge)")
	}
	//1 more tick to be ready
	c.Tick()
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be on cd")
	}
	if x.Cooldown(core.ActionSkill) != 360 {
		t.Errorf("expecting cooldown to be 360, got %v", x.Cooldown(core.ActionSkill))
	}

	//wait out cd, try 2 press
	for i := 0; i < 360; i++ {
		c.Tick()
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready now")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	x.Skill(p)
	x.Skill(p)
	//wait 4s, skill should go into cooldown for 600 - 240
	for i := 0; i < 239; i++ {
		c.Tick()
	}
	//should not be in cd yet
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready still (2nd charge)")
	}
	//1 more tick to be ready
	c.Tick()
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be on cd")
	}
	if x.Cooldown(core.ActionSkill) != 360 {
		t.Errorf("expecting cooldown to be 360, got %v", x.Cooldown(core.ActionSkill))
	}
	for i := 0; i < 360; i++ {
		c.Tick()
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready now")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	x.Skill(p)
	x.Skill(p)
	x.Skill(p)
	//skill should be in cd now for full 600s
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be on cd")
	}
	if x.Cooldown(core.ActionSkill) != 600 {
		t.Errorf("expecting cooldown to be 600, got %v", x.Cooldown(core.ActionSkill))
	}
	for i := 0; i < 599; i++ {
		c.Tick()
	}
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should not be ready")
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be on ready")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
}
