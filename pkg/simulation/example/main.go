package main

import (
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

const cfg = `target lvl=100 resist=0.1;
options swap_delay=12 debug=true iteration=1000 duration=109 workers=30 mode=sl;
energy every interval=480,720 amount=1;

#Chars Builds

yelan char lvl=90/90 cons=0 talent=9,9,9; 
yelan add weapon="favoniuswarbow" refine=3 lvl=90/90;
yelan add set="emblemofseveredfate" count=4;
# yelan add set="noblesseoblige" count=4;
yelan add stats hp=4780 atk=311 hp%=0.466 hydro%=0.466 cr=0.311; #main
yelan add stats def%=0.124 def=39.36 hp=507.88 hp%=0.1984 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

xingqiu char lvl=90/90 cons=6 talent=9,9,9; 
xingqiu add weapon="harbingerofdawn" refine=5 lvl=90/90;
xingqiu add set="emblemofseveredfate" count=4;
xingqiu add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311; #main
xingqiu add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

albedo char lvl=90/90 cons=0 talent=9,9,9;
albedo add weapon="cinnabarspindle" lvl=90/90 refine=5;
albedo add set="huskofopulentdreams" count=4 +params=[stacks=4];
albedo add stats hp=4780 atk=311 def%=0.583 geo%=0.466 cr=0.311;
albedo add stats def=39.36 def%=0.248 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.331 cd=0.7944 ;

noelle char lvl=90/90 cons=6 talent=9,9,9;
noelle add weapon="favoniusgreatsword" refine=3 lvl=90/90;
noelle add set="archaicpetra" count=4;
noelle add stats hp=4780 atk=311 def%=0.583 geo%=0.466 cr=0.311; #main 5* set
noelle add stats def%=0.248 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

active albedo;
while 1 {
albedo skill;
yelan burst, attack, skill, attack;
xingqiu burst, attack[delay=3];
noelle burst, attack, skill, attack:3, dash, attack:3, dash, attack:3;
yelan skill, attack;
xingqiu skill[delay=6], attack:2;
noelle attack:3, dash, attack:3, dash, attack;
}



`

func main() {
	// parse cfg
	p := ast.New(cfg)
	cfg, gcsl, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(gcsl.String())

	// create core
	c, err := core.New(core.Opt{
		Seed:       0,
		Debug:      true,
		DefHalt:    cfg.Settings.DefHalt,
		DamageMode: cfg.Settings.DamageMode,
	})
	if err != nil {
		panic(err)
	}

	// create new eval
	eval, err := gcs.NewEvaluator(gcsl, c)
	if err != nil {
		panic(err)
	}

	// create simulation
	sim, err := simulation.New(cfg, eval, c)
	if err != nil {
		panic(err)
	}
	// run simulatin
	_, err = sim.Run()

	if err != nil {
		panic(err)
	}

	// fmt.Println(res)

	logs, err := c.Log.Dump()
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(logs))

	os.Remove("logs.json")
	os.WriteFile("logs.json", logs, 0o600)

	// do stuff
}
