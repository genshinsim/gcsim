package shenhe

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestBasicAbilUsage(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
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

func TestSkillCDCon0(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
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

	sh := x.(*char)

	var f int

	f, _ = x.Skill(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	//expecting skill to come up in f+10*60
	for c.F < 599 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//check action ready
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill shouldn't be ready yet")
	}
	//cooldown should be 1
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

	//use skill hold
	p["hold"] = 1
	x.Skill(p)
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//check action ready
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill shouldn't be ready yet")
	}
	//cooldown should be 1
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

}

func TestBurstCDBasic(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
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

	sh := x.(*char)

	_, _ = x.Burst(p)
	// Burst should be ready at 1200+11 frames after usage
	for i := 0; i < 1200+11-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionBurst] > 0 {
		t.Errorf("expecting 0 burst stacks got %v", sh.availableCDCharge[core.ActionBurst])
	}
	//check action ready
	if x.ActionReady(core.ActionBurst, p) {
		t.Error("burst shouldn't be ready yet")
	}
	//cooldown should be 1
	if x.Cooldown(core.ActionBurst) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionBurst))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionBurst] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionBurst])
	}
	//cooldown should be 0 since at least one avail?
	if x.Cooldown(core.ActionBurst) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionBurst))
	}
	//action not ready since no energy
	if x.ActionReady(core.ActionBurst, p) {
		t.Error("burst should be ready since no energy")
	}
	sh.Energy = sh.MaxEnergy()
	//action should be ready now
	if !x.ActionReady(core.ActionBurst, p) {
		t.Error("burst should be ready since stack >1")
	}
}

func TestSkillCDCon1(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
	prof.Base.Cons = 1
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

	sh := x.(*char)

	var f int

	f, _ = x.Skill(p)
	for i := 0; i < f; i++ {
		c.Tick()
	}
	//expecting skill to come up in f+10*60
	for c.F < 599 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

	//use skill hold
	p["hold"] = 1
	x.Skill(p)
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//check action ready (we have another stack)
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill shouldn't be ready yet")
	}
	//cooldown should be 0
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

	//use hold then press
	p["hold"] = 1
	x.Skill(p)
	p["hold"] = 0
	x.Skill(p)
	//expecting 1 stack to be ready in +15, and another to be ready in +10
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//check action ready
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill shouldn't be ready yet")
	}
	//cooldown should be 1
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	for i := 0; i < 10*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//action should be ready since we have 1 stack
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	//cooldown should be 0 since at least one avail?
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

	//use press then hold
	p["hold"] = 0
	x.Skill(p)
	p["hold"] = 1
	x.Skill(p)
	//expecting 1 stack to be ready in +15, and another to be ready in +10
	for i := 0; i < 10*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//check action ready
	if x.ActionReady(core.ActionSkill, p) {
		t.Error("skill shouldn't be ready yet")
	}
	//cooldown should be 1
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//action should be ready since we have 1 stack
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	//cooldown should be 0 since at least one avail?
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}

}

func TestFlatCDReduction(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
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

	sh := x.(*char)

	x.Skill(p)
	//reduce cd by 10 frames, should come up at 590
	x.ReduceActionCooldown(core.ActionSkill, 10)
	for c.F < 589 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	x.Skill(p)
	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//at this point there shoudl be 500 frames left
	//we're going to reduce by 1000 and make sure stacks and everything are correct
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 500 {
		t.Errorf("expecting cooldown to be 500, got %v", x.Cooldown(core.ActionSkill))
	}
	x.ReduceActionCooldown(core.ActionSkill, 1000)
	//check now
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//and then keep looping for a while to make sure it doesn't get weird
	for i := 0; i < 5000; i++ {
		c.Tick()
	}
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

}

func TestFlatCDReductionCon1(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
	prof.Base.Cons = 1
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

	sh := x.(*char)

	// var f int

	x.Skill(p)
	p["hold"] = 1
	x.Skill(p)
	//reduce cd by 10 frames, should come up at 590
	x.ReduceActionCooldown(core.ActionSkill, 10)
	for c.F < 589 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	//next one should be ready at 590 + 900
	for c.F < 590+900-1 {
		c.Tick()
	}
	//should be at 1 only
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
}

func TestResetSkillCD(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
	prof.Base.Cons = 1
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

	sh := x.(*char)

	x.Skill(p)
	p["hold"] = 1
	x.Skill(p)

	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//first charge is 500 from recharging
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 500 {
		t.Errorf("expecting cooldown to be 500, got %v", x.Cooldown(core.ActionSkill))
	}
	x.ResetActionCooldown(core.ActionSkill)
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//next should be ready in 900 more frames
	for i := 0; i < 899; i++ {
		c.Tick()
	}
	//should be at 1 only
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//and then keep looping for a while to make sure it doesn't get weird
	for i := 0; i < 5000; i++ {
		c.Tick()
	}
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
}

func TestResetSkillCooldownReduction(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.Shenhe, core.Cryo, 6)
	prof.Base.Cons = 1
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

	sh := x.(*char)

	x.Skill(p)
	p["hold"] = 1
	x.Skill(p)

	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//first charge is 500 from recharging
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 500 {
		t.Errorf("expecting cooldown to be 500, got %v", x.Cooldown(core.ActionSkill))
	}
	//add cooldown reduction here
	//this shouldn't affect first charge duration
	x.AddCDAdjustFunc(core.CDAdjust{
		Key: "test",
		Amount: func(a core.ActionType) float64 {
			return -0.1
		},
		Expiry: -1,
	})
	for i := 0; i < 499; i++ {
		c.Tick()
	}
	if sh.availableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	//next should be ready in 810 more frames (900 less 10%)
	for i := 0; i < 809; i++ {
		c.Tick()
	}
	//should be at 1 only
	if sh.availableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	//1 more tick to be ready
	c.Tick()
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack == 2")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//and then keep looping for a while to make sure it doesn't get weird
	for i := 0; i < 5000; i++ {
		c.Tick()
	}
	if sh.availableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 skill stacks got %v", sh.availableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
}
