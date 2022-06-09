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
raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="dullblade" refine=3 lvl=90/90;
raiden add stats hp=4780 atk=311.0 er=0.518 cr=0.3110 electro%=0.4660;
raiden add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

target lvl=100 resist=0.1;
active raiden;

let x = 0;
// while x < 5 {
// 	x = x + 1;
// 	raiden attack;
// }
raiden attack;
// raiden skill;
raiden attack:4;
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
	c, err := core.New(0, true)
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
