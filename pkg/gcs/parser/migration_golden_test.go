package parser

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

// This test guards the migration from the hand-written parser to the
// pigeon (PEG) generated parser. The golden file was generated with the
// original hand-written parser (run `go test ./pkg/gcs/parser -run TestMigrationGolden -update-golden`
// to regenerate). The new parser must produce byte-identical fingerprints
// (AST structure incl. positions + ActionList contents) for all corpus entries.

var updateGolden = flag.Bool("update-golden", false, "rewrite testdata/migration_golden.txt")

var goldenSpew = spew.ConfigState{
	SortKeys:                true,
	SpewKeys:                true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	// Do not use String() methods: ast.MapExpr.String() iterates a Go map,
	// which is non-deterministic. Structural dumps also capture Pos values.
	DisableMethods: true,
	Indent:         " ",
}

func goldenFingerprint(src string) (string, error) {
	file := ast.NewFile()
	p := New(file, src)
	res, prog, err := p.Parse()
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("=== prog ===\n")
	sb.WriteString(goldenSpew.Sdump(prog))
	sb.WriteString("=== res ===\n")
	sb.WriteString(goldenSpew.Sdump(res))
	return sb.String(), nil
}

var migrationCorpus = []struct {
	name string
	src  string // if starts with "@", read from file path instead
}{
	{"sample_gcsim", "@testdata/corpus/sample.gcsim"},
	{"basic_cfg", cfg},
	{"char_action", charaction},
	{"char_stats", charstats},
	{"precedence", `
let a = 1+2*3;
let b = (1+2)*3;
let c = 1==2 && 3!=4;
let d = 1 && 0 || 1+2 == 3;
let e = -a * b;
let f = !-a;
let g = a - b;
let h = 1 < 2 == 3 > 4;
let i = 1 <= 2 != 3 >= 4;
let j = 1 <> 2;
`},
	{"literals", `
let i = 123;
let f = 1.5;
let f2 = .5;
let f3 = 1.;
let t = true;
let fl = false;
let s = "hello";
let s2 = "with \"escaped\" quotes";
let m = [a=1, b=2+3, c="str"];
let e = [];
`},
	{"unary_paren", `
let a = -(1 + 2);
let b = - -1;
let c = !!true;
let d = !(1 == 2);
`},
	{"fields", `
if .status.field > 0 {
	print("hi");
}
let x = .a.b.c;
`},
	{"fns", `
fn y(a, b) {
	let c = a + b;
	return c;
}
fn z(a number, b number) string {
	return "x";
}
fn h(cb fn(number) number) number {
	return cb(1);
}
fn empty() {
}
let anon = fn(x) { return x + 1; };
let typed fn(number) string = fn(a number) string { return "q"; };
`},
	{"control_flow", `
let x = 0;
while x < 10 {
	x = x + 1;
	if x == 5 {
		continue;
	} else if x == 8 {
		break;
	} else {
		x = x + 0;
	}
}
for let i = 0; i < 3; i = i + 1 {
	x = x + i;
}
for x = 0; x < 5; x = x + 1 {
	let j = x;
}
for x < 20 {
	x = x + 1;
}
for {
	break;
}
`},
	{"switch", `
switch a {
case 1:
	1 + 1;
	fallthrough;
case 2:
	2 + 2;
	break;
default:
	3 + 3;
}
switch {
case true:
	1;
}
`},
	{"options", `
options debug=true defhalt=false hitlag=false iteration=500 duration=90.5 workers=24 mode=sl;
options swap_delay=1 attack_delay=2 charge_delay=3 skill_delay=4 burst_delay=5 jump_delay=6 dash_delay=7 aim_delay=8;
options frame_defaults=human ignore_burst_energy=true;
`},
	{"target_full", `
target lvl=100 resist=0.1 pyro=0.2 hydro=0.3 hp=1000 pos=1,2.5 radius=2 particle_threshold=1 particle_drop_count=2 particle_element=pyro freeze_resist=0.5;
target lvl=90;
`},
	{"energy_hurt", `
energy once interval=300 amount=1;
energy every interval=300,600 amount=2;
hurt once interval=300 amount=1,300 element=physical;
hurt every interval=480,720 amount=1,300 element=pyro;
energy;
hurt;
`},
	{"random_stats", `
raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="favoniuslance" refine=3 lvl=90/90 +params=[stacks=5];
raiden add set="tenacityofthemillelith" count=4 +params=[x=1];
raiden add stats random rarity=5 sand=hp% goblet=pyro% circlet=cr;
raiden add stats hp=4780 atk=311.0 er=0.5180 label=main;
active raiden;
`},
	{"char_actions", `
xingqiu attack[randomparam=2]:4,skill;
xingqiu burst[orbital=0];
xingqiu attack:4;
xingqiu dash,jump,walk:2,swap;
active xingqiu;
`},
	{"comments", `
# hash comment
// slash comment
let x = 1; # trailing
let y = 2; // trailing
`},
	{"blocks", `
{
	let x = 1;
	{
		let y = 2;
	}
}
`},
}

func TestMigrationGolden(t *testing.T) {
	const goldenPath = "testdata/migration_golden.txt"
	// auto-discover real-world configs dropped in testdata/corpus
	corpus := migrationCorpus
	if ents, err := os.ReadDir("testdata/corpus"); err == nil {
		for _, e := range ents {
			if !e.IsDir() {
				corpus = append(corpus, struct {
					name string
					src  string
				}{"corpus_" + e.Name(), "@testdata/corpus/" + e.Name()})
			}
		}
	}
	var sb strings.Builder
	for _, c := range corpus {
		src := c.src
		if strings.HasPrefix(src, "@") {
			b, err := os.ReadFile(strings.TrimPrefix(src, "@"))
			if err != nil {
				t.Fatalf("%v: %v", c.name, err)
			}
			src = string(b)
		}
		fp, err := goldenFingerprint(src)
		if err != nil {
			t.Fatalf("%v: parse error: %v", c.name, err)
		}
		fmt.Fprintf(&sb, "######## %v ########\n%v\n", c.name, fp)
	}
	got := sb.String()
	if *updateGolden {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("golden updated")
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v (run with -update-golden)", err)
	}
	if string(want) != got {
		wl := strings.Split(string(want), "\n")
		gl := strings.Split(got, "\n")
		for i := 0; i < len(wl) || i < len(gl); i++ {
			var w, g string
			if i < len(wl) {
				w = wl[i]
			}
			if i < len(gl) {
				g = gl[i]
			}
			if w != g {
				t.Errorf("first diff at line %d:\n want: %q\n  got: %q", i, w, g)
				for j := i - 4; j < i+2; j++ {
					if j >= 0 && j < len(wl) && j < len(gl) {
						t.Logf("  ctx %d:\n    want %q\n    got  %q", j, wl[j], gl[j])
					}
				}
				break
			}
		}
		t.Errorf("golden mismatch; run with -update-golden to regenerate (only valid with old parser)")
	}
}
