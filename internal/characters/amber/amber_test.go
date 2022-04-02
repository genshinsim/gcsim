package amber

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Amber, core.Pyro, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testhelper.TestBowCharacter(c, x)

	p := make(map[string]int)
	var f int

	f, _ = x.Skill(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	p["bunny"] = 1
	f, _ = x.Aimed(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
}

func TestCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Amber, core.Pyro, 0)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCDSingleCharge(c, x, 900)
	if err != nil {
		t.Error(err)
	}
}

func TestC4CD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Amber, core.Pyro, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCDDoubleCharge(c, x, []int{720, 720})
	if err != nil {
		t.Error(err)
	}
}
