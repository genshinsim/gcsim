package main

// import fyne
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var bows = []string{"prototypecrescent", "rust", "sacrificialbow", "skywardharp",
	"thestringless", "thunderingpulse", "theviridescenthunt",
	"windblumeode", "alleyhunter", "amosbow",
	"blackcliffwarbow", "elegy", "favoniuswarbow",
	"hamayumi", "mitternachtswaltz", "mouunsmoon", "polarstar"}

var catalysts = []string{"frostbearer", "hakushinring", "kagura", "mappamare",
	"memoryofdust", "oathsworneye", "eyeofperception", "lostprayertothesacredwinds",
	"prototypeamber", "skywardatlas", "solarpearl", "ttds", "thewidsith", "wineandsong",
	"blackcliffagate", "dodocotales", "favoniuscodex"}
var claymores = []string{"songofbrokenpines", "prototypearchaic", "rainslasher",
	"redhornstonethresher", "lithicblade", "luxurioussealord", "skywardpride", "serpentspine", "snowtombedstarsilver",
	"theunforged", "whiteblind", "wolfsgravestone", "akuoumaru", "blackcliffslasher",
	"favoniusgreatsword", "katsuragikirinagamasa"}
var polearms = []string{"calamityqueller", "thecatch", "crescentpike", "deathmatch", "dragonsbane",
	"dragonspinespear", "favoniuslance", "engulfinglightning", "staffofhoma", "kitaincrossspear", "lithicspear",
	"primordialjadewingedspear", "prototypestarglitter", "skywardspine", "vortexvanquisher", "wavebreakersfin",
	"blackcliffpole"}
var swords = []string{"aquilafavonia", "theblacksword", "blackclifflongsword",
	"favoniussword", "festeringdesire", "theflute", "freedomsworn", "haran",
	"harbingerofdawn", "ironsting", "lionsroar", "mistsplitterreforged",
	"jadecutter", "prototyperancour", "sacrificialsword", "skywardblade", "summitshaper",
	"thealleyflash", "amenomakageuchi"}

var bowusers = []string{
	"aloy", "amber",
	"diona", "fischl",
	"ganyu", "gorou",
	"sara", "tartaglia",
	"venti", "yoimiya",
}

var catalystusers = []string{
	"barbara",
	"klee",
	"lisa",
	"mona",
	"ningguang",
	"kokomi",
	"sucrose",
	"yae",
	"yanfei",
}

var claymoreusers = []string{
	"itto",
	"beidou",
	"chongyun",
	"diluc",
	"eula",
	"noelle",
}

var polearmusers = []string{
	"hutao",
	"raiden",
	"rosaria",
	"shenhe",
	"thoma",
	"xiangling",
	"xiao",
	"yunjin",
	"zhongli",
}

var swordsusers = []string{
	"albedo",
	"bennett",
	"jean",
	"kazuha",
	"kaeya",
	"ayaka",
	"ayato",
	"keqing",
	"qiqi",
	"xingqiu",
}

var weapon_types = []string{"Bows", "Catalysts", "Claymores", "Polearms", "Swords"}

var refines = []string{"1", "2", "3", "4", "5"}

//For main stats
var sands = []string{"hp%", "atk%", "def%", "er", "em"}

var sandsComplete = map[string]string{
	"hp%":  "=0.466",
	"atk%": "=0.466",
	"def%": "=0.466",
	"er":   "=0.518",
	"em":   "=186.5",
}

var goblets = []string{"hp%", "atk%", "def%", "em", "anemo%",
	"cryo%", "electro%", "geo%", "hydro%", "pyro%", "phys%"}

var gobletsComplete = map[string]string{
	"hp%":      "=0.466",
	"atk%":     "=0.466",
	"def%":     "=0.466",
	"em":       "=186.5",
	"anemo%":   "=0.466",
	"cryo%":    "=0.466",
	"electro%": "=0.466",
	"geo%":     "=0.466",
	"hydro%":   "=0.466",
	"pyro%":    "=0.466",
	"phys%":    "=0.466",
}

var circlets = []string{"hp%", "atk%", "def%", "em", "cr", "cd", "heal"}

var circletsComplete = map[string]string{
	"hp%":  "=0.466",
	"atk%": "=0.466",
	"def%": "=0.466",
	"em":   "=186.5",
	"cr":   "=0.311",
	"cd":   "=0.622",
	"heal": "=0.359",
}

func writeconfig(filename string, character string, weapons []string, refines []string, sands []string, goblets []string, circlets []string) {
	fileswritten := 0
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	filename = strings.Replace(filename, "\\", "/", -1)
	//fmt.Println(string(config))
	fmt.Printf("Values on function: %v %v %v %v %v %v\n", character, weapons, refines, sands, goblets, circlets)
	for _, weapon := range weapons {
		for _, refine := range refines {
			for _, sand := range sands {
				for _, goblet := range goblets {
					for _, circlet := range circlets {
						var replaceWeapon = regexp.MustCompile(character + ` add weapon="([a-z]+)" refine=([1-9])`)
						cfg_replaced_weapon := replaceWeapon.ReplaceAllString(string(config), character+" add weapon=\""+weapon+"\" refine="+refine)
						var replaceMainStats = regexp.MustCompile(character + ` add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
						cfg_replaced_total := replaceMainStats.ReplaceAllString(string(cfg_replaced_weapon), character+" add stats hp=4780 atk=311 "+sand+sandsComplete[sand]+" "+goblet+gobletsComplete[goblet]+" "+circlet+circletsComplete[circlet]+";")
						f, err := os.Create(path.Dir(filename) + "/_" + character + " " + weapon + " R" + refine + " " + sand + "-" + goblet + "-" + circlet + ".txt") //Write the .txt files

						if err != nil {
							log.Fatal(err)
						}

						defer f.Close()

						_, err2 := f.WriteString(cfg_replaced_total)

						if err2 != nil {
							log.Fatal(err2)
						}
						fileswritten = +1
						fmt.Println("done")
					}
				}
			}
		}

	}
	fmt.Printf("\nDone writing: %v files\n", fileswritten)

}

func getCharNames(filename string) []string {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	charNames := make([]string, 4, 4)
	//get the characters names to store them later in array
	var reGetCharNames = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	for i, match := range reGetCharNames.FindAllStringSubmatch(string(config), -1) {
		charNames[i] = string(match[1])

	}
	if err != nil {
		log.Println(err)

	}
	return charNames
}

func main() {
	// New app
	a := app.New()
	// New Window & title
	w := a.NewWindow("Multi-File Simulator and Optimizer")
	//Resize main/parent window
	w.Resize(fyne.NewSize(800, 400))
	//selects

	weapons_list := widget.NewCheckGroup((bows), func(selected []string) {})
	weapon_type_select := widget.NewSelect(weapon_types, func(changed string) {
		defer weapons_list.Refresh()
		switch changed {
		case "Bows":
			weapons_list.Options = bows
		case "Catalysts":
			weapons_list.Options = catalysts
		case "Claymores":
			weapons_list.Options = claymores
		case "Polearms":
			weapons_list.Options = polearms
		case "Swords":
			weapons_list.Options = swords
		}

	})
	char_select := widget.NewSelect(nil, func(changed string) {
		if contains(bowusers, changed) {
			weapon_type_select.SetSelected("Bows")
		} else if contains(catalystusers, changed) {
			weapon_type_select.SetSelected("Catalysts")
		} else if contains(claymoreusers, changed) {
			weapon_type_select.SetSelected("Claymores")
		} else if contains(polearmusers, changed) {
			weapon_type_select.SetSelected("Polearms")
		} else if contains(swordsusers, changed) {
			weapon_type_select.SetSelected("Swords")
		}
		weapon_type_select.Refresh()

	})
	weapon_type_select.SetSelected("Bows")

	//lists

	refines_list := widget.NewCheckGroup((refines), func(selected []string) {})

	sands_list := widget.NewCheckGroup((sands), func(selected []string) {})

	goblets_list := widget.NewCheckGroup((goblets), func(selected []string) {})

	circlets_list := widget.NewCheckGroup((circlets), func(selected []string) {})
	//check

	replace := false
	gz := false
	check_replace := widget.NewCheck("Replace configs", func(b bool) {
		if b {
			replace = true
		} else {
			replace = false
		}
	})
	check_replace.SetChecked(true)

	check_gz := widget.NewCheck("Create .gz's", func(b bool) {
		if b {
			gz = true
		} else {
			gz = false
		}
	})
	check_gz.SetChecked(true)
	// buttons
	var file_loaded string

	load := widget.NewButton("Load", func() {
		file_loaded = selectFile()
		charNames := getCharNames(file_loaded)
		char_select.Options = charNames
		char_select.SetSelected(charNames[0])
		char_select.Refresh()
	})

	write := widget.NewButton("Write", func() {
		if len(char_select.Selected) == 0 {
			fmt.Println("Load config to be used as a base first!")
			return
		}
		if len(weapons_list.Selected) == 0 || len(refines_list.Selected) == 0 || len(sands_list.Selected) == 0 || len(goblets_list.Selected) == 0 || len(circlets_list.Selected) == 0 {
			fmt.Println("All lists have to have at least 1 element selected")
			return
		}
		writeconfig(file_loaded, char_select.Selected, weapons_list.Selected, refines_list.Selected, sands_list.Selected, goblets_list.Selected, circlets_list.Selected)
	})
	select_all := widget.NewButton("Select all", func() {
		weapons_list.Selected = weapons_list.Options
		weapons_list.Refresh()
	})
	unselect_all := widget.NewButton("Unselect all", func() {
		weapons_list.Selected = nil
		weapons_list.Refresh()
	})

	btn_run := widget.NewButton("Just run", func() { OptnRunFunc(false, replace, gz) })
	btn_optrun := widget.NewButton("Opt n' run", func() { OptnRunFunc(true, replace, gz) })

	//-------------------------------------------------------Load Menu-------------------------------------------------//
	select1_label := widget.NewLabel("Load character:")
	select2_label := widget.NewLabel("Weapons:")
	col0 := container.New(layout.NewVBoxLayout(), select1_label, load, select2_label, char_select, weapon_type_select, layout.NewSpacer())

	//-------------------------------------------------------Weapons-------------------------------------------------//
	weapons_container := container.New(layout.NewGridLayout(1), weapons_list)
	weapons_scroll := container.NewVScroll(weapons_container)
	weapons_label := widget.NewLabel("Weapons")
	weapons_btns_container := container.New(layout.NewVBoxLayout(), weapons_label, select_all, unselect_all)
	col1 := container.NewVSplit(weapons_scroll, weapons_btns_container)
	col1.SetOffset(1.0)

	//-------------------------------------------------------Refines-------------------------------------------------//
	refines_container := container.New(layout.NewGridLayout(1), refines_list)
	refines_scroll := container.NewVScroll(refines_container)
	refines_label := widget.NewLabel("Refines")
	col2 := container.NewVSplit(refines_scroll, refines_label)
	col2.SetOffset(1.0)

	//-------------------------------------------------------Sands-------------------------------------------------//
	sands_container := container.New(layout.NewGridLayout(1), sands_list)
	sands_scroll := container.NewVScroll(sands_container)
	sands_label := widget.NewLabel("Sands")
	col3 := container.NewVSplit(sands_scroll, sands_label)
	col3.SetOffset(1.0)

	//-------------------------------------------------------Goblets-------------------------------------------------//
	goblets_container := container.New(layout.NewGridLayout(1), goblets_list)
	goblets_scroll := container.NewVScroll(goblets_container)
	goblets_label := widget.NewLabel("Goblets")
	col4 := container.NewVSplit(goblets_scroll, goblets_label)
	col4.SetOffset(1.0)

	//-------------------------------------------------------Circlets-------------------------------------------------//
	circlets_container := container.New(layout.NewGridLayout(1), circlets_list)
	circlets_scroll := container.NewVScroll(circlets_container)
	circlets_label := widget.NewLabel("Circlets")
	col5 := container.NewVSplit(circlets_scroll, circlets_label)
	col5.SetOffset(1.0)

	col6 := container.New(layout.NewVBoxLayout(), write, btn_run, btn_optrun, check_replace, check_gz)

	content := container.New(layout.NewHBoxLayout(), col0, col1, col2, col3, col4, col5, col6)

	w.SetContent(content)
	//show and run
	w.ShowAndRun()
}
