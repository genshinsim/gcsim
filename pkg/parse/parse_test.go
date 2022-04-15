package parse

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
)

func TestParse(t *testing.T) {
	p := New("test", pteststring)
	a, err := p.Parse()
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

	p = New("test", raiden)
	_, err = p.Parse()
	if err != nil {
		t.Error(err)
	}

	p = New("test", s2)
	_, err = p.Parse()
	if err != nil {
		t.Error(err)
	}

	p = New("test", check)
	_, err = p.Parse()
	if err != nil {
		t.Error(err)
	}

}

func TestConfigClone(t *testing.T) {
	parser := New("test", pteststring)
	c, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}
	next := c.Clone()
	// cjson, _ := json.Marshal(c)
	// fmt.Println(string(cjson))
	// fmt.Println("")
	// njson, _ := json.Marshal(next)
	// fmt.Println(string(njson))
	if !reflect.DeepEqual(c, next) {
		t.Fail()
	}

}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parser := New("test", raiden)
		parser.Parse()
	}
}

func BenchmarkClone(b *testing.B) {
	parser := New("test", raiden)
	c, err := parser.Parse()
	if err != nil {
		b.Error(err)
	}
	// run clone b.N times
	for n := 0; n < b.N; n++ {
		c.Clone()
	}
}

func BenchmarkSerialization(b *testing.B) {
	parser := New("test", raiden)
	c, err := parser.Parse()
	if err != nil {
		b.Error(err)
	}
	j, _ := json.Marshal(c)
	for n := 0; n < b.N; n++ {
		var out core.SimulationConfig
		json.Unmarshal(j, &out)
	}
}

var check = `
target lvl=100 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;
target lvl=100 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;

##Actions List
active bennett;

# HP particle simulation. Per srl:
# it adds 1 particle randomly, uniformly distributed between 200 to 300 frames after the last time an energy drops
# so in the case above, it adds on avg one particle every 250 frames in effect
# so over 90s of combat that's 90 * 60 / 250 = 21.6 on avg
# From elijam assumptions: https://discord.com/channels/763583452762734592/851428030094114847/884832120650993805
energy every interval=200,300 amount=1;

# xiao attack,charge,high_plunge[plunge_hits=1]  +if=.status.xiaoburst>60;
xiao high_plunge[plunge_hits=1]  +if=.status.xiaoburst>60;

zhongli skill[hold_nostele=1] +limit=1;
zhongli skill[hold_nostele=1] +if=.cd.xiao.burst<500;

bennett skill,burst  +swap_lock=100;
bennett skill  +if=.energy.bennett>40 +swap_lock=100;
bennett burst; 

sucrose skill,skill +swap_to=xiao +if=.cd.xiao.burst==0&&.energy.xiao<60;
sucrose skill +swap_to=xiao +if=.cd.xiao.burst==0;

# xiao skill,skill,burst  +if=.tags.xiao.eCharge>1;
xiao skill,skill,burst; 
# xiao skill,burst; 
# xiao skill  +if=.energy.xiao>60&&.cd.xiao.burst<120 +swap_lock=120;

# Funneling
sucrose skill +swap_to=xiao +if=.cd.xiao.burst==0;

xiao attack  +swap_lock=100;
zhongli attack  +is_onfield;

restart;
`

var pteststring = `
options debug=true iteration=3000 duration=41 workers=24;

xiangling char lvl=80/90 cons=4 talent=6,9,9 start_hp=100 +params=[a=1,b=2];
xiangling add weapon="staff of homa" lvl=80/90 refine=3 +params=[a=1,b=2];
xiangling add set="seal of insulation" count=4 +params=[a=1,b=2];
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

bennett skill;

# chain; macros have to be defined first
chain a,b,c +if=.field1.field2.field3>1 +swap_to=xiangling +limit=1 +try=wait;

# reset
reset_limit;

# wait
wait_for mods value=.xiangling.bennettbuff==1 max=10;
wait_for time max=10;
wait_for time max=100 +filler=attack[param=1];

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
xiangling attack +try=wait;
xiangling attack +try=drop;

# calc mode wait
wait 10;
wait until 1000;
`

var s2 = `
options debug=true iteration=3000 duration=41 workers=24 mode=apl;

xiangling char lvl=80/90 cons=4 talent=6,9,9;
xiangling add weapon="staff of homa" lvl=90/90 refine=1;

bennett char lvl=80/90 cons=1 talent=6,6,6;
bennett add weapon="favoniussword" lvl=90/90 refine=5;

target lvl=80 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.3;

active xiangling;

energy once interval=300 amount=1; #once at frame 300
hurt every interval=300,600 amount=100,200 ele=pyro; #randomly 100 to 200 dmg every 300 to 600 frames

# macros
a:xiangling skill +label=a;
b:wait_for particles value=xiangling max=100;


# list

bennett skill +swap_to=xiangling +label=battery;
chain a,b +label=xlcollect;

xiangling attack +is_onfield +label=fill;`

var raiden = `
options debug=true iteration=300 duration=60 workers=24;

bennett char lvl=70/80 cons=2 talent=6,8,8;
bennett add weapon="favoniussword" lvl=90/90 refine=1;
bennett add set="noblesseoblige" count=4;
bennett add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311 label=main ;
bennett add stats hp=717 hp%=0.057999999999999996 atk=78 atk%=0.663 def=118 em=42 er=0.221 cr=0.039 cd=0.475 label=subs ;

raiden char lvl=90/90 cons=1 talent=10,10,10;
raiden add weapon="engulfinglightning" lvl=90/90 refine=1;
raiden add set="emblemofseveredfate" count=4;
raiden add stats hp=4780 atk=311 er=0.518 electro%=0.466 cr=0.311 ;
raiden add stats hp=538 hp%=0.040999999999999995 atk=68 atk%=0.134 def%=0.073 em=89 er=0.057999999999999996 cr=0.32699999999999996 cd=0.948 ;

xiangling char lvl=80/90 cons=6 talent=6,9,10;
xiangling add weapon="staffofhoma" refine=1 lvl=90/90;
xiangling add set="crimsonwitchofflames" count=2;
xiangling add set="gladiatorsfinale" count=2;
xiangling add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311 ;
xiangling add stats hp=478 hp%=0.047 atk=65 atk%=0.152 def=76 def%=0.051 em=63 er=0.16199999999999998 cr=0.264 cd=0.9940000000000001 ;

xingqiu char lvl=80/90 cons=6 talent=1,9,10;
xingqiu add weapon="sacrificialsword" refine=5 lvl=90/90;
xingqiu add set="noblesseoblige" count=2;
xingqiu add set="heartofdepth" count=2;
xingqiu add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311 ;
xingqiu add stats hp=598 atk=58 atk%=0.391 def=113 em=23 er=0.279 cr=0.23299999999999998 cd=0.575 ;


##Default Enemy
target lvl=100 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;
# target lvl=100 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;

##Actions List
active raidenshogun;

# HP particle simulation. Per srl:
# it adds 1 particle randomly, uniformly distributed between 200 to 300 frames after the last time an energy drops
# so in the case above, it adds on avg one particle every 250 frames in effect
# so over 90s of combat that's 90 * 60 / 250 = 21.6 on avg
energy every interval=200,300 amount=1;

raidenshogun attack:4,dash,attack:4,dash,attack:4,dash,attack:2,charge +if=.status.raidenburst>0;

# Additional check to reset at the start of the next rotation
raidenshogun skill +if=.status.xianglingburst==0&&.energy.xingqiu>70&&.energy.xiangling>70;
raidenshogun skill +if=.status.raidenskill==0;

# Skill is required before burst to activate Kageuchi. Otherwise ER is barely not enough
# For rotations #2 and beyond, need to ensure that Guoba is ready to go. Guoba timing is about 300 frames after XQ fires his skill
xingqiu skill[orbital=1],burst[orbital=1],attack +if=.cd.xiangling.skill<300;

# Bennett burst goes after XQ burst for uptime alignment. Attack to proc swords
bennett burst,attack,skill +if=.status.xqburst>0&&.cd.xiangling.burst<180;

# Only ever want to XL burst in Bennett buff and after XQ burst for uptime alignment
xiangling burst,attack,skill,attack,attack +if=.status.xqburst>0&&.status.btburst>0;
# Second set of actions needed in case Guoba CD comes off while pyronado is spinning
xiangling burst,attack +if=.status.xqburst>0&&.status.btburst>0;
xiangling skill ;

# Raiden must burst after all others. Requires an attack to allow Bennett buff to apply
raidenshogun burst +if=.status.xqburst>0&&.status.xianglingburst>0&&.status.btburst>0;

# Funnelling
bennett attack,skill +if=.status.xqburst>0&&.energy.xiangling<70 +swap_to=xiangling;
bennett skill +if=.energy.xiangling<70 +swap_to=xiangling;
bennett skill +if=.energy.xingqiu<80 +swap_to=xingqiu;
bennett attack,skill +if=.status.xqburst>0 +if=.energy.raidenshogun<90 +swap_to=raidenshogun;

xingqiu attack +if=.status.xqburst>0;
xiangling attack +is_onfield;
bennett attack +is_onfield;
xingqiu attack +is_onfield;
raidenshogun attack +is_onfield;

`

func TestConfigOrdering(t *testing.T) {
	parser := New("test", orderstr)
	c, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}

	//expecting raiden, ayaka
	var order = []core.CharKey{core.Raiden, core.Ayaka}
	for i, v := range c.Characters.Profile {
		if v.Base.Key != order[i] {
			t.Errorf("expecting %v, got %v", order[i], v)
		}
	}

}

var orderstr = `
raiden char lvl=90/90 cons=2 talent=10,10,10;
raiden add weapon="engulfinglightning" refine=1 lvl=90/90;
raiden add set="emblemofseveredfate" count=4;
raiden add stats def%=0.07300 hp=5079 hp%=0.05300 atk=412.0 atk%=0.1920 er=0.6660 em=42.00 cr=0.5720 cd=1.119 electro%=0.4660;

ayaka char lvl=90/90 cons=0 talent=9,9,9;
ayaka add weapon="amenomakageuchi" refine=5 lvl=90/90;
ayaka add set="blizzardstrayer" count=4;
ayaka add stats def=95.00 hp=4780 hp%=0.04100 atk=503.0 atk%=0.7630 er=0.05800 em=42.00 cr=0.1590 cd=1.243 cryo%=0.4660;

options mode=sl;

target lvl=100 resist=0.1;
active raiden;
 
raiden skill;
`

var statLabelTest = `
options debug=true iteration=300 duration=60 workers=24;

bennett char lvl=70/80 cons=2 talent=6,8,8;
bennett add weapon="favoniussword" lvl=90/90 refine=1;
bennett add set="noblesseoblige" count=4;
bennett add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311 label=main ;
bennett add stats hp=717 hp%=0.057999999999999996 atk=78 atk%=0.663 def=118 em=42 er=0.221 cr=0.039 cd=0.475 label=subs ;
bennett add stats roll=avg rarity=5 label=rolltest
def%=1
def=1
hp=1
hp%=1
atk=1
atk%=1
er=1
em=1
cr=1
cd=1
heal=1
pyro%=1
hydro%=1
cryo%=1
electro%=1
anemo%=1
geo%=1
phys%=1
;

active bennett;
target lvl=100 resist=0.1;


`

func TestStatsLabel(t *testing.T) {
	p := New("test", statLabelTest)
	a, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
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
		//print stat map
		for k, vv := range v.StatsByLabel {
			fmt.Printf("\t stats for %v: %v\n", k, vv)
		}
	}

}
