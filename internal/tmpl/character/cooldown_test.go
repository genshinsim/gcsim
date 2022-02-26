package character

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/testhelper"
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func TestSkillSetCDSingleCharge(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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

	//set cd for 600 frames
	x.SetCD(core.ActionSkill, 600)
	//expecting skill to come up in f+600
	for c.F < 599 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
}

func TestBurstSetCDSingleCharge(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
	x.Energy = 60
	x.EnergyMax = 60
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

	//set cd for 600 frames
	x.SetCD(core.ActionBurst, 600)
	x.Energy = 0
	//expecting burst to come up in f+600
	for c.F < 599 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionBurst] > 0 {
		t.Errorf("expecting 0 burst stacks got %v", x.AvailableCDCharge[core.ActionBurst])
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
	if x.AvailableCDCharge[core.ActionBurst] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionBurst])
	}
	//cooldown should be 0 since at least one avail?
	if x.Cooldown(core.ActionBurst) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionBurst))
	}
	//action not ready since no energy
	if x.ActionReady(core.ActionBurst, p) {
		t.Error("burst should not be ready since no energy")
	}
	x.Energy = 60
	//action should be ready now
	if !x.ActionReady(core.ActionBurst, p) {
		t.Error("burst should be ready since stack >1")
	}
}

func TestSkillSetCDDoubleCharge(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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
	x.SetNumCharges(core.ActionSkill, 2)

	//set cd for 600 frames
	x.SetCD(core.ActionSkill, 600)
	//expecting skill to come up in f+600
	for c.F < 599 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}

	//try longer cd
	x.SetCD(core.ActionSkill, 15*60)
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}

	//try mixed cd
	x.SetCD(core.ActionSkill, 15*60)
	x.SetCD(core.ActionSkill, 600)
	//expecting 1 stack to be ready in +15, and another to be ready in +10
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	for i := 0; i < 10*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}

	//use press then hold
	x.SetCD(core.ActionSkill, 600)
	x.SetCD(core.ActionSkill, 15*60)
	//expecting 1 stack to be ready in +15, and another to be ready in +10
	for i := 0; i < 10*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	for i := 0; i < 15*60-1; i++ {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}

}

func TestFlatCDReduction(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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

	x.SetCD(core.ActionSkill, 600)
	//reduce cd by 10 frames, should come up at 590
	x.ReduceActionCooldown(core.ActionSkill, 10)
	for c.F < 589 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

	x.SetCD(core.ActionSkill, 600)
	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//at this point there shoudl be 500 frames left
	//we're going to reduce by 1000 and make sure stacks and everything are correct
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 500 {
		t.Errorf("expecting cooldown to be 500, got %v", x.Cooldown(core.ActionSkill))
	}
	x.ReduceActionCooldown(core.ActionSkill, 1000)
	//check now
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}

}

func TestFlatCDReductionDoubleCharge(t *testing.T) {
	c := testhelper.NewTestCore()
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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
	x.SetNumCharges(core.ActionSkill, 2)

	x.SetCD(core.ActionSkill, 600)
	x.SetCD(core.ActionSkill, 15*60)
	//reduce cd by 10 frames, should come up at 590
	x.ReduceActionCooldown(core.ActionSkill, 10)
	for c.F < 589 {
		c.Tick()
	}
	//stack shouldn't be ready yet
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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
	x.SetNumCharges(core.ActionSkill, 2)

	x.SetCD(core.ActionSkill, 600)
	x.SetCD(core.ActionSkill, 15*60)

	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//first charge is 500 from recharging
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 500 {
		t.Errorf("expecting cooldown to be 500, got %v", x.Cooldown(core.ActionSkill))
	}
	x.ResetActionCooldown(core.ActionSkill)
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	prof := testhelper.CharProfile(core.TravelerAnemo, core.Cryo, 0)
	x, err := NewTemplateChar(c, prof)
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
	x.SetNumCharges(core.ActionSkill, 2)

	x.SetCD(core.ActionSkill, 600)
	x.SetCD(core.ActionSkill, 15*60)

	for i := 0; i < 100; i++ {
		c.Tick()
	}
	//first charge is 500 from recharging
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] > 0 {
		t.Errorf("expecting 0 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if x.Cooldown(core.ActionSkill) != 1 {
		t.Errorf("expecting cooldown to be 1, got %v", x.Cooldown(core.ActionSkill))
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 1 {
		t.Errorf("expecting 1 burst stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] > 1 {
		t.Errorf("expecting 1 e stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	//1 more tick to be ready
	c.Tick()
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 2 skill stacks got %v", x.AvailableCDCharge[core.ActionSkill])
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
	if x.AvailableCDCharge[core.ActionSkill] != 2 {
		t.Errorf("expecting 1 skill stacks got %v", x.AvailableCDCharge[core.ActionSkill])
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Error("skill should be ready since stack >1")
	}
	if x.Cooldown(core.ActionSkill) != 0 {
		t.Errorf("expecting cooldown to be 0, got %v", x.Cooldown(core.ActionSkill))
	}
}
