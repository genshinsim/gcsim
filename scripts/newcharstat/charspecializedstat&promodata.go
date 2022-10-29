package main

import (
	"encoding/json"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/curves"
)

type AvatarPromotes []struct {
	AvatarPromoteID int       `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	PromoteLevel    int       `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	// UnlockMaxLevel int `json:"unlockMaxLevel"`
	// ScoinCost           int `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
}
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

func getCharSpecializedStatandPromoData() (map[int]string, map[int][]curves.PromoData) {

	avatarPromoteJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/AvatarPromoteExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	var avatarPromotes AvatarPromotes
	if err := json.Unmarshal([]byte(avatarPromoteJson), &avatarPromotes); err != nil {
		log.Fatal(err)
	}
	//reshape avatarPromotes and remove duplicates to map of avatarPromoteId to AddProps
	specializedStatMap := make(map[int]string)
	for _, v := range avatarPromotes {
		specializedStatMap[v.AvatarPromoteID] = v.AddProps[3].PropType
	}
	// reshape avatarPromotes to map of avatarPromoteId to PromoData array of 7 items(1 for each ascension)
	promoDataMap := make(map[int][]curves.PromoData)
	for _, v := range avatarPromotes {
		partialPromoData := addPropArraytoPromoData(v.AddProps)

		partialPromoData.MaxLevel = promoteLevelToMaxLevel(v.PromoteLevel)
		promoDataMap[v.AvatarPromoteID] = append(promoDataMap[v.AvatarPromoteID], partialPromoData)
	}

	return specializedStatMap, promoDataMap
}

func addPropArraytoPromoData(addProps []AddProp) curves.PromoData {
	var out curves.PromoData
	for _, prop := range addProps {

		switch prop.PropType {
		case "FIGHT_PROP_BASE_HP":
			out.HP = prop.Value
		case "FIGHT_PROP_BASE_ATTACK":
			out.Atk = prop.Value
		case "FIGHT_PROP_BASE_DEFENSE":
			out.Def = prop.Value
		default:
			out.Special = prop.Value
		}
	}
	return out
}

func promoteLevelToMaxLevel(level int) int {
	switch level {
	case 0:
		return 20
	case 1:
		return 40
	case 2:
		return 50
	case 3:
		return 60
	case 4:
		return 70
	case 5:
		return 80
	case 6:
		return 90
	default:
		return 0
	}
}
