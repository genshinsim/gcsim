package xingqiu

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
// var target core.Target
// var xq core.Character

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

// 		target = monster.New(0, s, logger, 0, core.EnemyProfile{
// 			Level:  88,
// 			Resist: defaultResMap(),
// 		})

// 	})

// 	str := `
// 	char+=xingqiu ele=hydro lvl=80 hp=9514.469 atk=187.803 def=705.132 atk%=0.240 cr=.05 cd=0.5 cons=6 talent=1,8,10;
// 	weapon+=xingqiu label="sacrificial sword" atk=454.363 er=0.613 refine=4;
// 	art+=xingqiu label="gladiator's finale" count=2;
// 	art+=xingqiu label="noblesse oblige" count=2;
// 	stats+=xingqiu label=flower hp=4780 def=44 er=.065 cr=.097 cd=.124;
// 	stats+=xingqiu label=feather atk=311 cd=.218 def=19 atk%=.117 em=40;
// 	stats+=xingqiu label=sands atk%=0.466 cd=.124 def%=.175 er=.045 hp=478;
// 	stats+=xingqiu label=goblet hydro%=.466 cd=.202 atk%=.14 hp=299 atk=39;
// 	stats+=xingqiu label=circlet cr=.311 cd=0.062 atk%=.192 hp%=.082 atk=39;
// 	`

// 	p := parse.New("test", str)
// 	cfg, err := p.Parse()
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("error parsing initial config")
// 		return
// 	}

// 	xq, err = NewChar(sim, logger, cfg.Characters.Profile[0])
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("error parsing initial config")
// 		return
// 	}

// 	sim.Chars = append(sim.Chars, xq)

// 	//manually add 2 pc glad + 2pc no
// 	gladiator.New(xq, sim, logger, 2)
// 	noblesse.New(xq, sim, logger, 2)

// 	os.Exit(m.Run())
// }

// func TestXingqiuSkill(t *testing.T) {
// 	delay := 0
// 	param := make(map[string]int)
// 	atkCounts := make(map[core.AttackTag]int)
// 	particleCount := 0
// 	var totalDmg float64
// 	//on damage to track what's happening
// 	sim.OnDamage = func(ds *core.Snapshot) {
// 		ds.Stats[core.CR] = 1
// 		atkCounts[ds.AttackTag]++
// 		dmg, _ := target.Attack(ds)
// 		logger.Infow("attack", "abil", ds.Abil, "dmg", dmg)
// 		totalDmg += dmg
// 	}
// 	sim.OnParticle = func(p core.Particle) {
// 		xq.ReceiveParticle(p, true, 4)
// 		particleCount += p.Num
// 	}

// 	fmt.Println("----xingqiu skill testing----")
// 	delay = xq.Skill(param)
// 	sim.Skip(delay + 200)
// 	if !expect("skill attack count", 2, atkCounts[core.AttackTagElementalArt]) {
// 		t.Error("invalid attack count")
// 	}
// 	if !expect("particle count", 5, particleCount) {
// 		t.Error("invalid particle count")
// 	}
// 	//9859+11221

// 	expect("total skill damage", 9859+11221, totalDmg)
// 	if !floatApproxEqual(9859+11221, totalDmg, 10) {
// 		t.Error("invalid total damage")
// 	}

// }

// func TestXingqiuAttack(t *testing.T) {
// 	sim.OnDamage = nil
// 	delay := 0
// 	e := 0
// 	param := make(map[string]int)

// 	xq.ResetNormalCounter()
// 	e = xq.ActionFrames(core.ActionAttack, param)
// 	delay = xq.Attack(param)
// 	if !expect("normal attack delay", e, delay) {
// 		t.Error("invalid normal attack delay")
// 	}
// 	sim.Skip(delay)
// 	e = xq.ActionFrames(core.ActionAttack, param)
// 	delay = xq.Attack(param)
// 	if !expect("normal attack delay", e, delay) {
// 		t.Error("invalid normal attack delay")
// 	}
// 	sim.Skip(delay)
// 	e = xq.ActionFrames(core.ActionAttack, param)
// 	delay = xq.Attack(param)
// 	if !expect("normal attack delay", e, delay) {
// 		t.Error("invalid normal attack delay")
// 	}
// 	sim.Skip(delay)
// 	e = xq.ActionFrames(core.ActionAttack, param)
// 	delay = xq.Attack(param)
// 	if !expect("normal attack delay", e, delay) {
// 		t.Error("invalid normal attack delay")
// 	}
// 	sim.Skip(delay)
// 	e = xq.ActionFrames(core.ActionAttack, param)
// 	delay = xq.Attack(param)
// 	if !expect("normal attack delay", e, delay) {
// 		t.Error("invalid normal attack delay")
// 	}
// 	sim.Skip(delay + 100)

// }

// func TestXingqiuBurst(t *testing.T) {

// 	delay := 0
// 	param := make(map[string]int)

// 	delay = xq.Burst(param)

// 	sim.Skip(delay)

// 	atkCounts := make(map[core.AttackTag]int)
// 	//on damage to track what's happening
// 	sim.OnDamage = func(ds *core.Snapshot) {
// 		atkCounts[ds.AttackTag]++
// 		dmg, _ := target.Attack(ds)
// 		logger.Debugw("attack", "abil", ds.Abil, "dmg", dmg)
// 	}
// 	xq.ResetNormalCounter()
// 	//attack 3 times, should trigger 3 different waves of burst
// 	fmt.Println("----xingqiu burst testing----")
// 	fmt.Println("checking burst first wave")
// 	xq.Attack(param)
// 	sim.ExecuteEventHook(core.PostAttackHook)
// 	sim.Skip(200)
// 	if !expect("normal attack count", 1, atkCounts[core.AttackTagNormal]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagNormal] = 0
// 	if !expect("burst attack count", 2, atkCounts[core.AttackTagElementalBurst]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagElementalBurst] = 0
// 	xq.ResetNormalCounter()

// 	fmt.Println("checking burst second wave")
// 	xq.Attack(param)
// 	sim.ExecuteEventHook(core.PostAttackHook)
// 	sim.Skip(200)
// 	if !expect("normal attack count", 1, atkCounts[core.AttackTagNormal]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagNormal] = 0
// 	if !expect("burst attack count", 3, atkCounts[core.AttackTagElementalBurst]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagElementalBurst] = 0
// 	xq.ResetNormalCounter()

// 	fmt.Println("checking burst third wave")
// 	xq.Attack(param)
// 	sim.ExecuteEventHook(core.PostAttackHook)
// 	sim.Skip(200)
// 	if !expect("normal attack count", 1, atkCounts[core.AttackTagNormal]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagNormal] = 0
// 	if !expect("burst attack count", 5, atkCounts[core.AttackTagElementalBurst]) {
// 		t.Error("invalid attack count")
// 	}
// 	atkCounts[core.AttackTagElementalBurst] = 0
// 	xq.ResetNormalCounter()

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
