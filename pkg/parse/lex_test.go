package parse

import (
	"log"
	"testing"
)

var s = `
xiangling char lvl=80/90 cons=4 talent=6,9,9;
xiangling add weapon="staff of homa" lvl=80/90 refine=3;
xiangling add set="seal of insulation" count=4;
xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
xiangling add stats atk%=.0992 cr=.1655 cd=.7282 em=39.64 er=.5510 hp%=.0992 hp=507.88 atk=33.08 def%=.124 def=39.36;

xiangling skill +if=.debuff.res.1.test==1
xiangling burst,skill;
`

func TestLex(t *testing.T) {
	log.Println("testing lex")

	l := lex("test", s)
	// last := "roar"
	// stop := false
	for n := l.nextItem(); n.typ != itemEOF; n = l.nextItem() {
		log.Printf("%v - %v: %v\n", n.line, n.pos, n)
		if n.typ == itemError {
			t.FailNow()
		}
		// last = n.val
	}

}
