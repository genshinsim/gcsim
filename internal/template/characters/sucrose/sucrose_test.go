package sucrose

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Sucrose, core.Anemo, 6)
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
	prof := testhelper.CharProfile(core.Sucrose, core.Anemo, 0)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCDSingleCharge(c, x, 15*60+9)
	if err != nil {
		t.Error(err)
	}
}

func TestC4CD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Sucrose, core.Anemo, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//second charge shouldn't have delay because it will only start
	//recharging after the first for the purpose for this test
	err = testhelper.TestSkillCooldown(c, x, []int{15*60 + 9, 15 * 60})
	if err != nil {
		t.Error(err)
	}
}
