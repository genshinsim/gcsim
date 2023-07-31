package main

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
	"github.com/genshinsim/gcsim/pkg/worker"
)

const cfg = `options swap_delay=12 debug=true iteration=100 duration=100.58 workers=50 mode=sl;
ayaka char lvl=90/90 cons=0 talent=9,9,9; 
ayaka add weapon="amenomakageuchi" refine=5 lvl=90/90;
ayaka add set="blizzardstrayer" count=4;
ayaka add stats hp=4780 atk=311 atk%=0.466 cryo%=0.466 cd=0.622; #main;
ayaka add stats hp%=0.0992 hp=507.88 atk%=0.1488 atk=33.08 def%=0.124 def=39.36 er=0.1653 em=39.64 cr=0.3972 cd=0.662;

kazuha char lvl=90/90 cons=0 talent=9,9,9; 
kazuha add weapon="favoniussword" refine=3 lvl=90/90;
kazuha add set="viridescentvenerer" count=5;
kazuha add stats hp=4780 atk=311 em=187 em=187 em=187; #main
kazuha add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=118.92 cr=0.3972 cd=0.5296;

mona char lvl=90/90 cons=0 talent=9,9,9; 
mona add weapon="thrillingtalesofdragonslayers" refine=5 lvl=90/90;
#mona add weapon="prototypeamber" refine=5 lvl=90/90;
mona add set="tenacityofthemillelith" count=4;
mona add stats hp=4780 atk=311 er=0.518 hydro%=0.466 cr=0.311; #main
#mona add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311; #main  Prototype Amber
mona add stats hp%=0.0992 hp=507.88 atk%=0.0992 atk=33.08 def%=0.124 def=39.36 er=0.551 em=39.64 cr=0.1324 cd=0.7944;
#mona add stats hp%=0.0992 hp=507.88 atk%=0.0992 atk=33.08 def%=0.124 def=39.36 er=0.4408 em=39.64 cr=0.1986 cd=0.7944;  # Prototype Amber

shenhe char lvl=90/90 cons=0 talent=9,9,9; 
shenhe add weapon="favoniuslance" refine=3 lvl=90/90;
shenhe add set="noblesseoblige" count=4;
shenhe add stats hp=4780 atk=311 atk%=0.466 atk%=0.466 cr=0.311; #main
shenhe add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.6612 em=39.64 cr=0.331 cd=0.1324;

# ----
energy every interval=480,720 amount=1;
target lvl=100 resist=.1;

# ----
active ayaka;

while 1 {
ayaka dash, attack, skill, attack:2;
kazuha skill, high_plunge;
shenhe burst, skill;
kazuha burst;
mona skill,burst;
ayaka skill, dash, burst, attack, charge;
shenhe skill, attack;
kazuha skill, high_plunge;
ayaka dash, attack, charge;
print("Rotation Done");
}

`

func main() {
	p := ast.New(cfg)
	cfg, gcsl, err := p.Parse()
	if err != nil {
		panic(err)
	}

	os.Remove("data.json")
	f, err := os.OpenFile("data.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	respCh := make(chan stats.Result)
	errCh := make(chan error)
	pool := worker.New(cfg.Settings.NumberOfWorkers, respCh, errCh)
	pool.StopCh = make(chan bool)

	go func() {
		//make all the seeds
		wip := 0
		for wip < cfg.Settings.Iterations {
			pool.QueueCh <- worker.Job{
				Cfg:     cfg.Copy(),
				Actions: gcsl.Copy(),
				Seed:    rand.Int63(),
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	count := cfg.Settings.Iterations
	for count > 0 {
		select {
		case r := <-respCh:
			// fmt.Printf("%d\n", count)
			out, _ := json.Marshal(r)
			f.Write(out)
			f.WriteString("\n")
			count--
		case err := <-errCh:
			panic(err)
		}
	}
}
