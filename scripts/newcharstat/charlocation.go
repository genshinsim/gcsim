package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

type FetterInfo []struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int    `json:"avatarId"`
}

func getCharLocationMap() map[int]string {
	fetterInfoJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/FetterInfoExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	var fetterInfo FetterInfo
	if err := json.Unmarshal([]byte(fetterInfoJson), &fetterInfo); err != nil {
		log.Fatal(err)
	}

	// reshape fetterInfo to map of avatarId to AvatarAssocType
	locationMap := make(map[int]string)
	for _, v := range fetterInfo {
		locationMap[v.AvatarId] = v.AvatarAssocType
	}
	return locationMap
}

func determineCharRegion(location string) (profile.ZoneType, error) {
	switch location {
	case "ASSOC_TYPE_INAZUMA":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_LIYUE":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_MONDSTADT":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_SUMERU":
		return profile.ZoneSumeru, nil
	case "ASSOC_TYPE_MAINACTOR", "ASSOC_TYPE_RANGER":
		return profile.ZoneUnknown, nil
	case "ASSOC_TYPE_FATUI":
		return profile.ZoneSnezhnaya, nil
	default:
		return profile.ZoneUnknown, errors.New("unknown location")
	}

}
