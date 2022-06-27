package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

const cfg = `
options swap_delay=12 debug=true iteration=1000 duration=105.35 workers=30 mode=sl;
target lvl=100 resist=.1;
# this translate to a lambda of 1/10 per second
energy every=10 amount=1;

raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="favoniuslance" refine=3 lvl=90/90;
raiden add set="tenacityofthemillelith" count=4;
raiden add stats hp=4780 atk=311.0 er=0.518 cr=0.3110 electro%=0.4660;
raiden add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1488 er=0.1653 em=39.64 cr=0.331 cd=0.7944 ;																																																																											

xingqiu char lvl=90/90 cons=6 talent=9,9,9;
xingqiu add weapon="harbingerofdawn" refine=5 lvl=90/90;
xingqiu add set="noblesseoblige" count=4;
xingqiu add stats hp=4780 atk=311.0 atk%=0.4660 cr=0.3110 hydro%=0.4660;
xingqiu add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.3306 em=39.64 cr=0.2648 cd=0.7944 ;																																																																																																																																						

bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="thealleyflash" refine=1 lvl=90/90;
bennett add set="instructor" count=4;
bennett add stats hp=3571 atk=232.0 em=187.0 cr=0.2320 pyro%=0.3480;
bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.2979 cd=0.4634 ;																						

xiangling char lvl=90/90 cons=6 talent=9,9,9;
xiangling add weapon="thecatch" refine=5 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;
xiangling add stats hp=4780 atk=311.0 em=187.0 cr=0.3110 pyro%=0.4660;
xiangling add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944 ;																					

active raiden;
while 1 {
	print("raiden start");
    raiden attack, skill;

	print("xq burst");
    xingqiu burst, attack;

	print("bennett's turn now?");
    bennett burst, attack, skill, attack;

	print("xiangling turn");
    xiangling burst, attack, skill;

	print("back on xingqiu");
    xingqiu attack:3, skill, dash, attack:3;

	print("more bennett");
    bennett skill, attack;

	print("raiden combo");
    raiden burst,attack:4, dash, attack:4, dash, attack:4, dash, attack:3;

	print("bennett battery");
    bennett attack, skill;

    print("waiting");
    wait(7);

	print("restarting at frame: ", f());
}
`

func main() {
	//parse cfg
	p := ast.New(cfg)
	cfg, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Program.String())

	//create core
	c, err := core.New(core.CoreOpt{
		Seed:       0,
		Debug:      true,
		Delays:     cfg.Settings.Delays,
		DefHalt:    cfg.Settings.DefHalt,
		DamageMode: cfg.Settings.DamageMode,
	})
	if err != nil {
		panic(err)
	}

	//create simulation
	sim, err := simulation.New(cfg, c)
	if err != nil {
		panic(err)
	}
	//run simulatin
	res, err := sim.Run()

	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	logs, err := c.Log.Dump()
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(logs))

	os.Remove("logs.json")
	ioutil.WriteFile("logs.json", logs, 0600)

	//do stuff
}
