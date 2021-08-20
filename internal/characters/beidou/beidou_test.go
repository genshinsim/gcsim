package beidou

// import (
// 	"fmt"
// 	"math/rand"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/genshinsim/gsim/internal/artifacts/gladiator"
// 	"github.com/genshinsim/gsim/internal/artifacts/noblesse"
// 	"github.com/genshinsim/gsim/internal/dummy"
// 	"github.com/genshinsim/gsim/pkg/core"
// 	"github.com/genshinsim/gsim/pkg/monster"
// 	"github.com/genshinsim/gsim/pkg/parse"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// var logger *zap.SugaredLogger
// var sim *dummy.Sim
// var targetA core.Target
// var targetB core.Target
// var bd core.Character

// func TestMain(m *testing.M) {
// 	os.Remove("./out.log")
// 	// call flag.Parse() here if TestMain uses flags
// 	config := zap.NewDevelopmentConfig()
// 	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
// 	config.EncoderConfig.TimeKey = ""
// 	config.OutputPaths = []string{"out.log"}
// 	log, _ := config.Build(zap.AddCallerSkip(1))
// 	logger = log.Sugar()
// 	//set up sim
// 	sim = dummy.NewSim(func(s *dummy.Sim) {

// 		s.R = rand.New(rand.NewSource(time.Now().Unix()))

// 		targetA = monster.New(0, s, logger, 0, core.EnemyProfile{
// 			Level:  88,
// 			Resist: defaultResMap(),
// 		})
// 		targetB = monster.New(1, s, logger, 0, core.EnemyProfile{
// 			Level:  88,
// 			Resist: defaultResMap(),
// 		})

// 	})

// 	str := `
// 	char+=beidou ele=electro lvl=80 hp=11565 atk=200 def=575 cr=0.05 cd=0.50 electro%=.18 cons=6 talent=6,8,8;
// 	weapon+=beidou label="serpent spine" atk=510 refine=5 cr=.276;
// 	stats+=beidou label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
// 	stats+=beidou label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
// 	`

// 	p := parse.New("test", str)
// 	cfg, err := p.Parse()
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("error parsing initial config")
// 		return
// 	}

// 	bd, err = NewChar(sim, logger, cfg.Characters.Profile[0])
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("error parsing initial config")
// 		return
// 	}

// 	sim.Chars = append(sim.Chars, bd)
// 	sim.Targs = append(sim.Targs, targetA)
// 	sim.Targs = append(sim.Targs, targetB)

// 	//manually add 2 pc glad + 2pc no
// 	gladiator.New(bd, sim, logger, 2)
// 	noblesse.New(bd, sim, logger, 4)

// 	os.Exit(m.Run())
// }

// func TestXingqiuSkill(t *testing.T) {
// 	delay := 0
// 	param := make(map[string]int)
// 	var aCount, bCount int
// 	//on damage to track what's happening
// 	sim.OnDamage = func(ds *core.Snapshot) {

// 		a, _ := targetA.Attack(ds)
// 		if a > 0 {
// 			// log.Printf("a got hit by %v, src %v\n", ds.Abil, ds.SourceFrame)
// 			aCount++
// 		}
// 		b, _ := targetB.Attack(ds)
// 		if b > 0 {
// 			// log.Printf("b got hit by %v, src %v\n", ds.Abil, ds.SourceFrame)
// 			bCount++
// 		}
// 	}

// 	fmt.Println("----beidou Q bounce testing----")
// 	delay = bd.Burst(param)
// 	sim.Skip(delay + 100)
// 	if !expect("attack count on target A and B", 2, aCount+bCount) {
// 		t.Error("invalid attack count")
// 	}
// 	aCount = 0
// 	bCount = 0
// 	delay = bd.Attack(param)
// 	sim.Skip(delay + 200)
// 	//expecting 4 hit from attack + 3 Q
// 	if !expect("attack count on target A", 4, aCount) {
// 		t.Error("invalid attack count")
// 	}
// 	if !expect("attack count on target B", 2, bCount) {
// 		t.Error("invalid attack count")
// 	}

// }

// func expect(key string, a interface{}, b interface{}) bool {
// 	fmt.Printf("%v: expecting %v, got %v\n", key, a, b)
// 	return a == b
// }

// func defaultResMap() map[core.EleType]float64 {
// 	res := make(map[core.EleType]float64)

// 	res[core.Electro] = 0.1
// 	res[core.Pyro] = 0.1
// 	res[core.Anemo] = 0.1
// 	res[core.Cryo] = 0.1
// 	res[core.Frozen] = 0.1
// 	res[core.Hydro] = 0.1
// 	res[core.Dendro] = 0.1
// 	res[core.Geo] = 0.1
// 	res[core.Physical] = 0.3

// 	return res
// }

// func floatApproxEqual(expect, result, tol float64) bool {
// 	if expect > result {
// 		return expect-result < tol
// 	}
// 	return result-expect < tol
// }
