package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
)

type AvatarPromotes []struct {
	AvatarPromoteID int       `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	PromoteLevel    int       `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	UnlockMaxLevel int `json:"unlockMaxLevel"`
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
		partialPromoData.MaxLevel = v.UnlockMaxLevel
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

func determineStat(specializedStat string) (attributes.Stat, error) {
	switch specializedStat {
	case "FIGHT_PROP_CRITICAL_HURT":
		return attributes.CD, nil
	case "FIGHT_PROP_HEAL_ADD":
		return attributes.Heal, nil
	case "FIGHT_PROP_ATTACK_PERCENT":
		return attributes.ATKP, nil
	case "FIGHT_PROP_ELEMENT_MASTERY":
		return attributes.EM, nil
	case "FIGHT_PROP_HP_PERCENT":
		return attributes.HPP, nil
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
		return attributes.ER, nil
	case "FIGHT_PROP_CRITICAL":
		return attributes.CR, nil
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
		return attributes.PhyP, nil
	case "FIGHT_PROP_ELEC_ADD_HURT":
		return attributes.ElectroP, nil
	case "FIGHT_PROP_ROCK_ADD_HURT":
		return attributes.GeoP, nil
	case "FIGHT_PROP_FIRE_ADD_HURT":
		return attributes.PyroP, nil
	case "FIGHT_PROP_WATER_ADD_HURT":
		return attributes.HydroP, nil
	case "FIGHT_PROP_DEFENSE_PERCENT":
		return attributes.DEFP, nil
	case "FIGHT_PROP_ICE_ADD_HURT":
		return attributes.CryoP, nil
	case "FIGHT_PROP_WIND_ADD_HURT":
		return attributes.AnemoP, nil
	case "FIGHT_PROP_GRASS_ADD_HURT":
		return attributes.DendroP, nil
	case "":
		return attributes.EndStatType, nil
	default:
		return attributes.EndStatType, errors.New("unknown stat type")

	}
}
