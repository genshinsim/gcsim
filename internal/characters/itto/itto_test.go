package itto

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Itto, core.Geo, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	c.CharPos[prof.Base.Key] = 0
	//add targets to test with
	eProf := testhelper.EnemyProfile()
	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, eProf))
	p := make(map[string]int)

	var f int

	f, _ = x.Skill(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	f, _ = x.Burst(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	//bunch of attacks
	for j := 0; j < 10; j++ {
		f, _ = x.Attack(p)
		for i := 0; i < f; i++ {
			c.Tick()
		}
	}
	//charge attack
	f, _ = x.ChargeAttack(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	//tick a bunch of times after
	for i := 0; i < 1200; i++ {
		c.Tick()
	}

}
