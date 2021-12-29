package parse

import (
	"fmt"
	"testing"
)

var pteststring = `
options debug=true iteration=3000 duration=41 workers=24;

xiangling char lvl=80/90 cons=4 talent=6,9,9 start_hp=100 params=[a=1,b=2];
xiangling add weapon="staff of homa" lvl=80/90 refine=3 params=[a=1,b=2];
xiangling add set="seal of insulation" count=4 params=[a=1,b=2];
xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
xiangling add stats atk%=.0992 cr=.1655 cd=.7282 em=39.64 er=.5510 hp%=.0992 hp=507.88 atk=33.08 def%=.124 def=39.36;

target lvl=80 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.3;
target lvl=80 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.3;

energy once interval=300 amount=1; #once at frame 300
hurt every interval=300,600 amount=100,200 ele=pyro; #randomly 100 to 200 dmg every 300 to 600 frames

# macros
a:xiangling skill;
b:wait_for particles value=xiangling max=100;
c:reset_limit;

# chain; macros have to be defined first
chain a,b,c +if=.field1.field2.field3>1 +swap_to=xiangling +limit=1 +try=1;

# reset
reset_limit;

# wait
wait_for mods value=.xiangling.bennettbuff==1 max=10;
wait_for time max=10;

# basic char abil
xiangling burst,skill;

# repeater
xiangling attack:4,charge,attack:4;

# conditions
xiangling attack +if=.debuff.res.t1.cryo>1;
xiangling attack +swap_to=xiangling;
xiangling attack +swap_lock=100;
xiangling attack +is_onfield;
xiangling attack +label=hi;
xiangling attack +needs=hi;
xiangling attack +limit=2;
xiangling attack +timeout=100;
xiangling attack +try;
xiangling attack +try=1;
xiangling attack +try=0;

# calc mode wait
wait 10;
wait until 1000;
`

func TestParse(t *testing.T) {
	p := New("test", pteststring)
	a, _, err := p.Parse()
	fmt.Println("characters:")
	for _, v := range a.Characters.Profile {
		fmt.Println(v.Base.Key.String())
		//basic stats:
		fmt.Println("\t basics", v.Base)
		fmt.Println("\t char params", v.Params)
		fmt.Println("\t weapons", v.Weapon)
		fmt.Println("\t talents", v.Talents)
		fmt.Println("\t sets", v.Sets)
		fmt.Println("\t set params", v.SetParams)
		//pretty print stats
		fmt.Println("\t stats", v.Stats)
	}
	fmt.Println("rotations:")
	for _, v := range a.Rotation {
		fmt.Println(v)
	}
	fmt.Println("targets:")
	fmt.Println("\t", a.Targets)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("hurt event:")
	fmt.Println("\t", a.Hurt)
	fmt.Println("energy:")
	fmt.Println("\t", a.Energy)

}
