package yanfei

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

func TestBasicAbilUsage(t *testing.T) {
	c, err := core.New(func(c *core.Core) error {
		c.Log = logger
		return nil
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	prof := tests.CharProfile(keys.Yanfei, core.Pyro, 6)
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
	eProf := tests.EnemeyProfile()
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
