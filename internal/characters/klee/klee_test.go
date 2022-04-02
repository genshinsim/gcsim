package klee

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Klee, core.Pyro, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testhelper.TestCatalystCharacter(c, x)
}

func TestCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Klee, core.Pyro, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCDDoubleCharge(c, x, []int{1200, 1200})
	if err != nil {
		t.Error(err)
	}
}
