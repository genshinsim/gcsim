package shield_test

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/eval"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func TestShieldDamage(t *testing.T) {
	cfg := `
options swap_delay=12 iteration=1000;

# Build assumptions


nicole char lvl=90/90 cons=0 talent=9,9,9;
nicole add weapon="oathsworn" refine=5 lvl=90/90;
nicole add set="cg" count=4;
nicole add stats hp=4780 atk=311 atk%=0.466 atk%=0.466 atk%=0.466;
nicole add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=165.4 atk%=0.2976 er=0.3857 em=39.64 cr=0.1655 cd=0.1324;

# Target and Energy assumptions
target lvl=100 resist=0.1 radius=2 pos=0,2.4 hp=999999999; # standard ST target
energy every interval=480,720 amount=1; # standard Energy drops
hurt once interval=300 amount=10000,10000 element=physical;
# Rotation assumptions
active nicole;

nicole skill;

wait(40*60);
  `

	file := ast.NewFile()
	simcfg, gcsl, err := simulator.Parse(file, cfg)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c, err := simulation.NewCore(1923, true, simcfg)
	if err != nil {
		t.Fatalf("NewCore() error = %v", err)
	}
	eval, err := eval.NewEvaluator(file, gcsl, c)
	if err != nil {
		t.Fatalf("NewEvaluator() error = %v", err)
	}
	// create a new simulation and run
	s, err := simulation.New(simcfg, eval, c)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	res, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// t.Fatalf("Shield results: %+v", res.ShieldResults.EffectiveShield)

	expected := []stats.ShieldSingleInterval{{Start: 9, End: 301, HP: 15724.17267882571}, {Start: 301, End: 1209, HP: 5724.172678825709}}
	for i := range res.ShieldResults.EffectiveShield["anemo"] {
		if res.ShieldResults.EffectiveShield["anemo"][i] != expected[i] {
			t.Fatalf("expected shield interval %d to be %+v, got %+v", i, expected[i], res.ShieldResults.EffectiveShield["anemo"][i])
		}
	}
}

func TestShieldBroken(t *testing.T) {
	cfg := `
options swap_delay=12 iteration=1000;

# Build assumptions


nicole char lvl=90/90 cons=0 talent=9,9,9;
nicole add weapon="oathsworn" refine=5 lvl=90/90;
nicole add set="cg" count=4;
nicole add stats hp=4780 atk=311 atk%=0.466 atk%=0.466 atk%=0.466;
nicole add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=165.4 atk%=0.2976 er=0.3857 em=39.64 cr=0.1655 cd=0.1324;

# Target and Energy assumptions
target lvl=100 resist=0.1 radius=2 pos=0,2.4 hp=999999999; # standard ST target
energy every interval=480,720 amount=1; # standard Energy drops
hurt once interval=300 amount=30000,30000 element=physical;
# Rotation assumptions
active nicole;

nicole skill;

wait(40*60);
  `

	file := ast.NewFile()
	simcfg, gcsl, err := simulator.Parse(file, cfg)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c, err := simulation.NewCore(1923, true, simcfg)
	if err != nil {
		t.Fatalf("NewCore() error = %v", err)
	}
	eval, err := eval.NewEvaluator(file, gcsl, c)
	if err != nil {
		t.Fatalf("NewEvaluator() error = %v", err)
	}
	// create a new simulation and run
	s, err := simulation.New(simcfg, eval, c)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	res, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// t.Fatalf("Shield results: %+v", res.ShieldResults.EffectiveShield)

	expected := []stats.ShieldSingleInterval{{Start: 9, End: 301, HP: 15724.17267882571}}
	for i := range res.ShieldResults.EffectiveShield["anemo"] {
		if res.ShieldResults.EffectiveShield["anemo"][i] != expected[i] {
			t.Fatalf("expected shield interval %d to be %+v, got %+v", i, expected[i], res.ShieldResults.EffectiveShield["anemo"][i])
		}
	}
}
