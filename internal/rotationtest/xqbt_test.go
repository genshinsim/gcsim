package rotationtest

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gsim/internal/characters/bennett"
	"github.com/genshinsim/gsim/internal/characters/xingqiu"
	"github.com/genshinsim/gsim/internal/dummy"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/monster"
	"github.com/genshinsim/gsim/pkg/parse"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestC6XingqiuBennett(t *testing.T) {
	var logger *zap.SugaredLogger
	var sim *dummy.Sim
	var target def.Target
	var xq def.Character
	var bt def.Character

	os.Remove("./out.log")
	// call flag.Parse() here if TestMain uses flags
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	config.EncoderConfig.TimeKey = ""
	config.OutputPaths = []string{"out.log"}
	log, _ := config.Build(zap.AddCallerSkip(1))
	logger = log.Sugar()
	//set up sim
	sim = dummy.NewSim(func(s *dummy.Sim) {

		s.R = rand.New(rand.NewSource(time.Now().Unix()))

		target = monster.New(0, s, logger, 0, def.EnemyProfile{
			Level:  88,
			Resist: defaultResMap(),
		})

	})

	str := `
	char+=xingqiu ele=hydro lvl=80 hp=9059.529 atk=178.823 def=671.415 atk%=0.180 cr=.05 cd=0.5 cons=6 talent=1,8,8;
	weapon+=xingqiu label="sacrificial sword" atk=454.363 er=0.613 refine=4;

	char+=bennett ele=pyro lvl=90 hp=12397.403 atk=191.157 def=771.249 er=0.267 cr=.05 cd=0.5 cons=6 talent=7,10,10;
	weapon+=bennett label="generic sword" atk=23 refine=1;
	`

	p := parse.New("test", str)
	cfg, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error parsing initial config")
		return
	}

	xq, err = xingqiu.NewChar(sim, logger, cfg.Characters.Profile[0])
	if err != nil {
		fmt.Println(err)
		fmt.Println("error parsing initial config")
		return
	}

	bt, err = bennett.NewChar(sim, logger, cfg.Characters.Profile[2])
	if err != nil {
		fmt.Println(err)
		fmt.Println("error parsing initial config")
		return
	}

	sim.Chars = append(sim.Chars, xq)
	sim.Chars = append(sim.Chars, bt)

	delay := 0
	param := make(map[string]int)
	atkCounts := make(map[def.AttackTag]int)
	particleCount := 0
	var totalDmg float64
	//on damage to track what's happening
	sim.OnDamage = func(ds *def.Snapshot) {
		atkCounts[ds.AttackTag]++
		dmg, _ := target.Attack(ds)
		logger.Infow("attack", "abil", ds.Abil, "dmg", dmg)
		totalDmg += dmg
	}
	sim.OnParticle = func(p def.Particle) {
		xq.ReceiveParticle(p, true, 4)
		particleCount += p.Num
	}

	//cast e at frame 0

	delay = bt.Skill(param)
	//damage on frame 9?
	sim.Skip(delay)

}

func expect(key string, a interface{}, b interface{}) bool {
	fmt.Printf("%v: expecting %v, got %v\n", key, a, b)
	return a == b
}

func defaultResMap() map[def.EleType]float64 {
	res := make(map[def.EleType]float64)

	res[def.Electro] = 0.1
	res[def.Pyro] = 0.1
	res[def.Anemo] = 0.1
	res[def.Cryo] = 0.1
	res[def.Frozen] = 0.1
	res[def.Hydro] = 0.1
	res[def.Dendro] = 0.1
	res[def.Geo] = 0.1
	res[def.Physical] = 0.3

	return res
}

func floatApproxEqual(expect, result, tol float64) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}
