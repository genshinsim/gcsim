package sayu

import (
	"errors"
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Sayu, core.Anemo, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testhelper.TestClaymoreCharacter(c, x)
}

func TestCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Sayu, core.Anemo, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// skill press
	cd := 6 * 60
	p := make(map[string]int)
	p["hold"] = 0
	x.Skill(p)

	testhelper.SkipFrames(c, 15+cd-1)
	if x.ActionReady(core.ActionSkill, p) {
		t.Error(errors.New("skill shouldn't be ready yet"))
	}
	v := x.Cooldown(core.ActionSkill)
	if v != 1 {
		t.Error(fmt.Errorf("expecting cooldown to be 1, got %v", v))
	}
	testhelper.SkipFrames(c, 1)
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error(errors.New("skill should be ready now"))
	}
	v = x.Cooldown(core.ActionSkill)
	if v != 0 {
		t.Error(fmt.Errorf("expecting cooldown to be 0, got %v", v))
	}

	// skill hold
	p["hold"] = 300
	cd = 6*60 + 300 + 150 // 150 = 300 * 0.5
	x.Skill(p)

	// 1 = -1 + 2 (look at abil.go)
	testhelper.SkipFrames(c, 18+cd+1)
	if x.ActionReady(core.ActionSkill, p) {
		t.Error(errors.New("skill shouldn't be ready yet"))
	}
	v = x.Cooldown(core.ActionSkill)
	if v != 1 {
		t.Error(fmt.Errorf("expecting cooldown to be 1, got %v", v))
	}
	testhelper.SkipFrames(c, 1)
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error(errors.New("skill should be ready now"))
	}
	v = x.Cooldown(core.ActionSkill)
	if v != 0 {
		t.Error(fmt.Errorf("expecting cooldown to be 0, got %v", v))
	}
}
