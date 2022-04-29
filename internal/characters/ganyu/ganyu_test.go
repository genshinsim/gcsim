package ganyu

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Ganyu, core.Cryo, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testhelper.TestBowCharacter(c, x)

}

func TestCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Ganyu, core.Cryo, 0)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
		err = testhelper.TestSkillCDSingleCharge(c, x, 600+10)
	if err != nil {
		t.Error(err)
	}
}

func TestC2CD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Ganyu, core.Cryo, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCooldown(c, x, []int{600+10, 600})
	if err != nil {
		t.Error(err)
	}
}
