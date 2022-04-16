package diona

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestC6Below50(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Diona, core.Cryo, 6)
	//make her half dead?
	//default profile has 99em
	prof.Base.StartHP = 1
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	c.CharPos[prof.Base.Key] = 0
	//make sure to set her base stats
	err = x.CalcBaseStats()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//add targets to test with
	eProf := testhelper.EnemyProfile()
	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, eProf))
	p := make(map[string]int)

	c.Init()

	x.Burst(p)
	//c6 should activate 30 frames after start of burst
	if x.ModIsActive("diona-c6") {
		t.Error("diona c6 should not be active")
	}
	//check for bonus tag
	if x.Tag("c6bonus-diona") != 0 {
		t.Errorf("expecting 0 for c6bonus-diona tag, got %v", x.Tag("c6bonus-diona"))
	}
	testhelper.SkipFrames(c, 30)

	if x.Tag("c6bonus-diona") != 30+120 {
		t.Errorf("expecting 150 for c6bonus-diona tag, got %v", x.Tag("c6bonus-diona"))
	}
}

func TestC6Above50(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Diona, core.Cryo, 6)
	//default profile has 99em
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	diona := x.(*char)
	c.Chars = append(c.Chars, x)
	c.CharPos[prof.Base.Key] = 0
	//make sure to set her base stats
	err = x.CalcBaseStats()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//add targets to test with
	eProf := testhelper.EnemyProfile()
	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, eProf))
	c.Init()

	p := make(map[string]int)

	// log.Println(diona.Tmpl.HPMax)
	// log.Println(x.HP())

	x.Burst(p)
	//c6 should activate 30 frames after start of burst
	if x.ModIsActive("diona-c6") {
		t.Error("diona c6 should not be active")
	}
	//check for bonus tag
	if x.Tag("c6bonus-diona") != 0 {
		t.Errorf("expecting 0 for c6bonus-diona tag, got %v", x.Tag("c6bonus-diona"))
	}
	testhelper.SkipFrames(c, 30)
	if x.Tag("c6bonus-diona") != 0 {
		t.Errorf("expecting 0 for c6bonus-diona tag, got %v", x.Tag("c6bonus-diona"))
	}
	if !x.ModIsActive("diona-c6") {
		t.Error("diona c6 should be active")
	}
	//check mod expiry; gotta find it in the array
	idx := -1
	for i, v := range diona.Mods {
		if v.Key == "diona-c6" {
			idx = i
			break
		}
	}
	if idx == -1 {
		t.Error("diona-c6 mod not found??")
	}
	if diona.Mods[idx].Expiry != 120+30 {
		t.Errorf("expecting 150 for diona-c6 mod expiry,  got %v", diona.Mods[idx].Expiry)
	}

}
