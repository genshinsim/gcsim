package sucrose

import (
	"os"
	"testing"

	"github.com/genshinsim/gcsim/internal/tests"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/player"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	config := zap.NewDevelopmentConfig()
	debug := os.Getenv("GCSIM_VERBOSE_TEST")
	level := zapcore.InfoLevel
	if debug != "" {
		level = zapcore.DebugLevel
	}
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = ""
	log, _ := config.Build(zap.AddCallerSkip(1))
	logger = log.Sugar()
	os.Exit(m.Run())
}

func TestSkillCDWithC4(t *testing.T) {
	c, err := core.New(func(c *core.Core) error {
		c.Log = logger
		return nil
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	prof := tests.CharProfile(keys.Sucrose, core.Anemo, 6)
	x, err := NewChar(c, prof)
	//cast it to *char so we can access private members
	sucrose := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	c.CharPos[prof.Base.Key] = 0
	//add targets to test with
	eProf := tests.EnemeyProfile()
	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, eProf))

	//check skill is ready
	p := make(map[string]int)
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be ready at start")
		t.FailNow()
	}

	//use skill, wait out animation, and check if ready for another use (2 stacks at c4)
	a, _ := x.Skill(p)
	for i := 0; i < a; i++ {
		c.Tick()
	}
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected second skill charge to be ready at start. At frame %v", c.F)
		t.FailNow()
	}

	//use second charge, next charge should be ready at 900-a
	x.Skill(p)
	for i := 0; i < 900-a-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}

	//tick once to get to 900-a
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v; charges: %v", c.F, x.Cooldown(core.ActionSkill), sucrose.eCharges)
		t.FailNow()
	}

	//next charge should be ready at 900 from now
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}

	for i := 0; i < 900-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

	//next charge should be ready at 900 from now
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}

	for i := 0; i < 900-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

	//use skill and then trigger flat cd reduction
	x.Skill(p)
	if x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected skill to be on cd at frame: %v", c.F)
		t.FailNow()
	}
	//next charge should be ready by 900 - flat cd reduction
	x.ReduceActionCooldown(core.ActionSkill, 100)
	for i := 0; i < 800-1; i++ {
		c.Tick()
		if x.ActionReady(core.ActionSkill, p) {
			t.Errorf("expected skill to be on cd at frame: %v", c.F)
			t.FailNow()
		}
	}
	c.Tick()
	if !x.ActionReady(core.ActionSkill, p) {
		t.Errorf("expected one charge of skill to be ready now. At frame %v; CD left: %v", c.F, x.Cooldown(core.ActionSkill))
		t.FailNow()
	}

}
