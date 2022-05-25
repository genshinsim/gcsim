package bennett

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/internal/reactable"
	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Bennett, core.Pyro, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testhelper.TestSwordCharacter(c, x)
}

func TestBurstAura(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Bennett, core.Pyro, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	c.CharPos[x.Key()] = 0
	c.Init()

	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, testhelper.EnemyProfile()))

	p := make(map[string]int)
	//pull out the player
	player := c.Targets[0].(*player.Player)

	//player should have no aura
	if player.AuraType() != core.NoElement {
		t.FailNow()
		log.Printf("at frame %v, aura is: %v\n", c.F, c.Targets[0].AuraType())
		t.Errorf("expected no aura, got %v", player.AuraType())
	}

	x.Burst(p)

	//this has pyro because we coded it to apply field right away... i think
	//this has to do with something zajef tested and said that bennett's Q
	//damage should benefit from his own buff?
	if player.AuraType() != core.Pyro {
		t.FailNow()
		log.Printf("at frame %v, aura is: %v, durability is: %v\n", c.F, c.Targets[0].AuraType(), player.Reactable.Durability[core.Pyro])
		t.Errorf("expected pyro aura, got %v", player.AuraType())
	}

	//first tick starts after a hard coded value
	testhelper.SkipFrames(c, burstStartFrame)

	//check player aura
	if player.AuraType() != core.Pyro {
		t.FailNow()
		log.Printf("at frame %v, aura is: %v, durability is: %v\n", c.F, c.Targets[0].AuraType(), player.Reactable.Durability[core.Pyro])
		t.Errorf("expected pyro aura, got %v", player.AuraType())
	}

	//for i := burstStartFrame; i <= 720+burstStartFrame; i += 60
	//the last application of aura should be 720+burstStartFrame
	testhelper.SkipFrames(c, 720)

	//this should be frame 720 + 31 = 751
	if player.AuraType() != core.Pyro {
		t.FailNow()
		log.Printf("at frame %v, aura is: %v, durability is: %v\n", c.F, c.Targets[0].AuraType(), player.Reactable.Durability[core.Pyro])
		t.Errorf("expected pyro aura, got %v", player.AuraType())
	}

	//player.ApplySelfInfusion(core.Pyro, 25, 126)
	//application is hard coded to last 126 frames with 25 durability
	testhelper.SkipFrames(c, bennettSelfInfusionDurationInFrames-1)

	if player.AuraType() != core.Pyro || player.Reactable.Durability[core.Pyro] <= reactable.ZeroDur {
		t.FailNow()
		t.Errorf("expected pyro aura, got %v", player.AuraType())
		log.Printf("at frame %v (should be 1 frame before aura expires), aura is: %v, durability is: %v\n", c.F, c.Targets[0].AuraType(), player.Reactable.Durability[core.Pyro])

	}

	testhelper.SkipFrames(c, 1)
	if player.AuraType() != core.NoElement || player.Reactable.Durability[core.Pyro] > reactable.ZeroDur {
		t.FailNow()
		t.Errorf("expected no aura, got %v", player.AuraType())
		log.Printf("at frame %v (aura should have expired), aura is: %v, durability is: %v\n", c.F, c.Targets[0].AuraType(), player.Reactable.Durability[core.Pyro])

	}

}

func TestCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Bennett, core.Pyro, 6)
	x, err := NewChar(c, prof)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = testhelper.TestSkillCDSingleCharge(c, x, 300-60+14)
	if err != nil {
		t.Error(err)
	}
}
