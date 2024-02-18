package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/scripts/extract-talents/excel"
)

var (
	avatarId  = flag.Int("avatar", 0, "extract using avatar id")
	depotId   = flag.Int("depot", 0, "extract using depot id (for traveler use)")
	excelPath = flag.String("excels", "./pipeline/data", "folder to look for excel data dump")
	level     = flag.Int("level", 1, "talent level to extract from")
)

func main() {
	flag.Parse()

	err := excel.LoadResources(func(name string, v any) error {
		d, err := os.ReadFile(filepath.Join(*excelPath, name))
		if err != nil {
			return err
		}
		return json.Unmarshal(d, v)
	})
	if err != nil {
		panic(err)
	}

	var depot *excel.AvatarSkillDepot
	if *depotId == 0 {
		id := uint32(*avatarId)
		avatar := excel.FindAvatar(id)
		if avatar == nil {
			log.Fatalf("no avatar found for id: %d", id)
		}
		depot = avatar.SkillDepot()
	} else {
		id := uint32(*depotId)
		depot = excel.FindSkillDepot(id)
		if depot == nil {
			log.Fatalf("no depot found for id: %d", id)
		}
	}

	dump("attack", excel.FindSkill(depot.Skills[0]))
	dump("skill", excel.FindSkill(depot.Skills[1]))
	dump("burst", excel.FindSkill(depot.EnergySkill))
}

func dump(prefix string, s *excel.AvatarSkill) {
	indent := strings.Repeat(" ", 4)
	fmt.Printf("%s%s: # %s\n", indent, prefix, s.Name())

	indent = strings.Repeat(indent, 3)
	for _, text := range s.ProudSkill(uint32(*level)).ParamDescList {
		text, params := parseParam(text.String())
		if text == "" {
			continue
		}

		for _, index := range params {
			fmt.Printf("%s- %d # %s\n", indent, index, text)
		}
	}
}

func parseParam(s string) (string, []int) {
	const paramPrefix = "{param"
	var i int
	var params []int
	for {
		start := strings.Index(s[i:], paramPrefix)
		if start == -1 {
			break
		}
		start += i + len(paramPrefix)

		end := strings.Index(s[start:], ":")
		if end == -1 {
			break
		}

		paramIndex, err := strconv.Atoi(s[start : start+end])
		if err != nil {
			panic(err)
		}

		paramIndex--
		params = append(params, paramIndex)

		s = s[:start] + strconv.Itoa(paramIndex) + s[start+end:]
		i = start
	}
	return s, params
}
