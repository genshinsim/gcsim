package parse

import (
	"fmt"
	"log"
	"testing"
)

var s = `
## sucrose

char+=sucrose ele=anemo lvl=60 hp=6501 atk=120 def=494 cr=0.05 cd=0.50 anemo%=.12 cons=2 talent=1,1,1;
weapon+=sucrose label="sacrificial fragments" atk=99 refine=1 em=85;
art+=sucrose label="viridescent venerer" count=4;

## xingqiu

char+=xingqiu ele=hydro lvl=70 hp=8352 atk=165 def=619 cr=0.05 cd=0.50 atk%=.18 cons=6 talent=1,8,8;
weapon+=xingqiu label="sacrificial sword" atk=401 refine=3 er=.559;
art+=xingqiu label="gladiator's finale" count=2;
art+=xingqiu label="noblesse oblige" count=2;
stats+=xingqiu label=flower hp=4780 def=44 er=.065 cr=.097 cd=.124;
stats+=xingqiu label=feather atk=311 cd=.218 def=19 atk=.117 em=40;
stats+=xingqiu label=sands atk%=0.466 cd=.124 def%=.175 er=.045 hp=478;
stats+=xingqiu label=goblet hydro%=.466 cd=.202 atk=.14 hp=299 atk=39;
stats+=xingqiu label=circlet cr=.311 cd=0.062 atk%=.192 hp%=.082 atk=39;

## diona

char+=diona ele=cryo lvl=70 hp=10129 atk=156 def=630 cr=0.05 cd=0.50 er=.2 cons=0 talent=1,1,1;
weapon+=diona label="favonius warbow" atk=401 refine=5;

## ganyu

char+=ganyu ele=cryo lvl=90 hp=9797 atk=335 def=630 cr=0.05 cd=0.884 cons=0 talent=10,6,6;
weapon+=ganyu label="amos bow" atk=608 atk%=0.496 refine=1 param=[stacks=3];
#weapon+=ganyu label="prototype crescent" atk=510 atk%=0.413 refine=5;
art+=ganyu label="blizzard strayer" count=4;
stats+=ganyu hp=4780 em=21 atk=47 cd=.179 def=19;
stats+=ganyu atk=311 cd=.062 em=35 atk%=.157 cr=.07;
stats+=ganyu atk%=.466 atk=31 cd=.225 hp%=.047 er=.168;
stats+=ganyu cryo%=.466 cd=.07 cr=.093 hp=717 def=16;
stats+=ganyu cd=.622 cr=.097 def=21 atk%=.14 def%=.066;


target+="blazing axe mitachurl" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.3;
target+="primo geovishap" lvl=95 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.5 anemo=0.1 physical=.3;
active+=ganyu;

actions+=skill target=sucrose if=.element.cryo==1&&.debuff.res.vvcryo==0 label=vv;
actions+=burst target=ganyu;
actions+=skill target=ganyu;
actions+=aim target=ganyu;
`

func TestLex(t *testing.T) {
	log.Println("testing lex")

	l := lex("test", s)
	last := "roar"
	// stop := false
	for n := l.nextItem(); n.typ != itemEOF; n = l.nextItem() {
		if n.typ == itemAction {
			log.Printf("action start at line %v\n", n.line)
		}
		log.Println(n)
		if n.val == last {
			t.FailNow()
		}
		last = n.val
	}

}

func TestParse(t *testing.T) {
	p := New("test", s)
	a, _, err := p.Parse()
	fmt.Println("characters:")
	for _, v := range a.Characters.Profile {
		fmt.Println(v.Base.Key.String())
		//basic stats:
		fmt.Println("\t basics", v.Base)
		fmt.Println("\t weapons", v.Weapon)
		fmt.Println("\t talents", v.Talents)
		fmt.Println("\t sets", v.Sets)
		//pretty print stats
		fmt.Println("\t stats", v.Stats)
	}
	fmt.Println("rotations:")
	for _, v := range a.Rotation {
		fmt.Println(v)
	}
	fmt.Println(a.Targets)
	if err != nil {
		t.Error(err)
	}

}
