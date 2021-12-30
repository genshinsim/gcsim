package reactable

import (
	"os"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testTarget struct {
	*Reactable
	src           int
	onDmgCallBack func(*core.AttackEvent) (float64, bool)
}

func (t *testTarget) Type() core.TargettableType                 { return core.TargettableEnemy }
func (t *testTarget) Index() int                                 { return 0 }
func (t *testTarget) SetIndex(ind int)                           {}
func (t *testTarget) MaxHP() float64                             { return 1 }
func (t *testTarget) HP() float64                                { return 1 }
func (t *testTarget) Shape() core.Shape                          { return core.NewCircle(0, 0, 1) }
func (t *testTarget) AddDefMod(key string, val float64, dur int) {}
func (t *testTarget) AddResMod(key string, val core.ResistMod)   {}
func (t *testTarget) RemoveResMod(key string)                    {}
func (t *testTarget) RemoveDefMod(key string)                    {}
func (t *testTarget) HasDefMod(key string) bool                  { return false }
func (t *testTarget) HasResMod(key string) bool                  { return false }
func (t *testTarget) AddReactBonusMod(mod core.ReactionBonusMod) {}
func (t *testTarget) ReactBonus(atk core.AttackInfo) float64     { return 0 }
func (t *testTarget) Kill()                                      {}
func (t *testTarget) SetTag(key string, val int)                 {}
func (t *testTarget) GetTag(key string) int                      { return 0 }
func (t *testTarget) RemoveTag(key string)                       {}

type combatTestCtrl struct {
	core *core.Core
	*core.CombatCtrl
}

func (t *testTarget) Attack(atk *core.AttackEvent) (float64, bool) {
	if t.onDmgCallBack != nil {
		return t.onDmgCallBack(atk)
	}
	return 0, false
}

func (c *combatTestCtrl) Init(x *core.Core) {
	c.CombatCtrl = core.NewCombatCtrl(x)
}

var logger *zap.SugaredLogger

var testChar core.CharacterProfile

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

	testChar.Stats = make([]float64, core.EndStatType)
	testChar.Talents.Attack = 1
	testChar.Talents.Burst = 1
	testChar.Talents.Skill = 1
	testChar.Base.Level = 90
	testChar.Stats[core.EM] = 100

	os.Exit(m.Run())
}

func TestReduce(t *testing.T) {
	r := &Reactable{}
	r.Durability = make([]core.Durability, core.ElementDelimAttachable)
	r.DecayRate = make([]core.Durability, core.ElementDelimAttachable)
	r.Durability[core.Electro] = 20
	r.reduce(core.Electro, 20, 1)
	if r.Durability[core.Electro] != 0 {
		t.Errorf("expecting nil electro balance, got %v", r.Durability[core.Electro])
	}
}

func TestTick(t *testing.T) {
	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	trg := &testTarget{src: 1}
	trg.Reactable = &Reactable{}
	trg.Init(trg, c)

	//test electro
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Electro,
			Durability: 25,
		},
	})

	if trg.Durability[core.Electro] != 0.8*25 {
		t.Errorf("expecting 20 electro, got %v", trg.Durability[core.Electro])
	}
	if trg.DecayRate[core.Electro] != 20.0/(6*25+420) {
		t.Errorf("expecting %v decay rate, got %v", 1.0/(6*25+420), trg.DecayRate[core.Electro])
	}

	//should deplete fully in 570 ticks
	for i := 0; i < 6*50+420; i++ {
		trg.Tick()
		// log.Println(target.Durability)
	}
	//check all durability should be nil
	ok := trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test multiple aura
	trg.attach(core.Electro, 50, 0.8)
	trg.attach(core.Hydro, 50, 0.8)
	trg.attach(core.Cryo, 50, 0.8)
	trg.attach(core.Pyro, 50, 0.8)
	for i := 0; i < 6*50+420; i++ {
		trg.Tick()
	}
	ok = trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test refilling
	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Electro,
			Durability: 25,
		},
	})
	for i := 0; i < 100; i++ {
		trg.Tick()
	}
	//calculate expected duration
	decay := core.Durability(20.0 / (6*25 + 420))
	left := 20 - 100*decay
	life := int((left + 40) / decay)
	// log.Println(decay, left, life)

	trg.React(&core.AttackEvent{
		Info: core.AttackInfo{
			Element:    core.Electro,
			Durability: 50,
		},
	})
	for i := 0; i < life-1; i++ {
		trg.Tick()
	}
	//make sure > 0
	if trg.Durability[core.Electro] < 0 {
		t.Errorf("expecting electro not to be 0 yet, got %v", trg.Durability[core.Electro])
	}
	//1 more tick and should be gone
	trg.Tick()
	ok = trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test frozen
	//50 frozen should last just over 208 frames (i.e. 0 by 209)
	trg.Durability[core.Frozen] = 50
	for i := 0; i < 208; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//should be > 0 still
	if trg.Durability[core.Frozen] < 0 {
		t.Errorf("expecting frozen not to be 0 yet, got %v", trg.Durability[core.Frozen])
	}
	//1 more tick and should be gone
	trg.Tick()
	if trg.Durability[core.Frozen] > 0 {
		t.Errorf("expecting frozen to be gone, got %v", trg.Durability[core.Frozen])
	}
	//105 more frames to full recover
	for i := 0; i < 104; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//decay should be > 0 still
	if trg.DecayRate[core.Frozen] < frzDecayCap {
		t.Errorf("expecting frozen decay to > cap, got %v", trg.Durability[core.Frozen])
	}
	//1 more tick to reset decay
	trg.Tick()
	if trg.DecayRate[core.Frozen] > frzDecayCap {
		t.Errorf("expecting frozen decay to reset, got %v", trg.Durability[core.Frozen])
	}

}

func (target *testTarget) allNil(t *testing.T) bool {
	ok := true
	for i, v := range target.Durability {
		ele := core.EleType(i)
		if !durApproxEqual(0, v, 0.00001) {
			t.Errorf("ele %v expected 0 durability got %v", ele, v)
			ok = false
		}
		if !durApproxEqual(0, target.DecayRate[i], 0.00001) && ele != core.Frozen {
			t.Errorf("ele %v expected 0 decay got %v", ele, target.DecayRate[i])
			ok = false
		} else if !durApproxEqual(frzDecayCap, target.DecayRate[i], 0.00001) && ele == core.Frozen {
			t.Errorf("frozen decay expected %v got %v", frzDecayCap, target.DecayRate[i])
			ok = false
		}
	}
	return ok
}

func durApproxEqual(expect, result, tol core.Durability) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}

// func floatApproxEqual(expect, result, tol float64) bool {
// 	if expect > result {
// 		return expect-result < tol
// 	}
// 	return result-expect < tol
// }
